package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os/exec"
	"strings"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	allowedEndpoints = []string{"/register", "/login", "/refreshToken", "/ws"}
)

func main() {

	addr := "127.0.0.1:8070"
	r := gin.Default()

	r.Use(CORSMiddleware())
	r.Use(wsMiddleware)
	r.Use(authMiddleware)

	r.POST("/login", login)

	r.POST("/register", register)

	r.POST("/refreshToken", refreshToken)

	r.GET("/ws", serveWs)

	_ = r.Run(addr)

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func isAllowed(endpoint string) bool {
	for _, e := range allowedEndpoints {
		if e == endpoint {
			return true
		}
	}
	return false
}

func authMiddleware(context *gin.Context) {
	path := context.Request.URL.Path
	if isAllowed(path) {
		return
	}

	authHeader := context.GetHeader("Authorization")
	if authHeader == "" {
		context.String(401, "Un-authorized")
		context.Abort()
		return
	}

	token := strings.Replace(authHeader, "Bearer ", "", 1)
	token = strings.Replace(token, "bearer ", "", 1)

	payload, tokenErr := validateToken(token)
	if tokenErr != nil {
		context.String(401, "Un-authorized")
		return
	}
	context.Set("CurrentUserEmail", payload.Email)

}

func wsMiddleware(context *gin.Context) {
	path := context.Request.URL.Path
	if path != "/ws" {
		return
	}

	token := context.Request.URL.Query().Get("token")
	if token == "" {
		context.String(401, "Un-authorized")
		context.Abort()
		return
	}

	payload, tokenErr := validateToken(token)
	if tokenErr != nil {
		context.String(401, "Un-authorized")
		return
	}
	context.Set("CurrentUserEmail", payload.Email)

}

func login(ctx *gin.Context) {
	var u User
	decoder := json.NewDecoder(ctx.Request.Body)
	_ = decoder.Decode(&u)
	h := sha256.New()
	h.Write([]byte(u.Password))
	password := hex.EncodeToString(h.Sum(nil))

	if !isEmailValid(u.Email) {
		ctx.String(400, "Email address is invalid.")
		return
	}

	user, err := userRepository{}.GetOne(u.Email, password)
	if err != nil {
		log.Printf("%v", err)
		ctx.String(404, "User not found!")
		return
	}

	if user != nil {
		tokenResponse := *generateToken(*user)
		if user.RefreshToken == "" || (user.RefreshTokenExpire == nil || time.Now().After(*user.RefreshTokenExpire)) {
			updatedAt := time.Now()
			user.UpdatedAt = &updatedAt
			user.RefreshToken = tokenResponse.RefreshToken
			expire := time.Now().Add(time.Hour * 24 * 30) // User's refresh token expires after 30 days.
			user.RefreshTokenExpire = &expire
			_, _ = userRepository{}.Update(*user)
		}
		tokenResponse.RefreshToken = user.RefreshToken
		ctx.JSON(200, tokenResponse)
		return
	}

	ctx.String(401, "Un-authorized")
}

func refreshToken(ctx *gin.Context) {

	var refreshTokenBody = struct {
		RefreshToken string `json:"refresh_token"`
	}{}

	decoder := json.NewDecoder(ctx.Request.Body)
	_ = decoder.Decode(&refreshTokenBody)

	u, err := userRepository{}.GetByRefreshToken(refreshTokenBody.RefreshToken)

	if err != nil {
		ctx.String(401, "Token is invalid.")
		ctx.Abort()
		return
	}

	if time.Now().After(*u.RefreshTokenExpire) {
		ctx.String(401, "Token is invalid.")
		ctx.Abort()
		return
	}
	tokenResponse := *generateToken(*u)
	tokenResponse.RefreshToken = refreshTokenBody.RefreshToken
	ctx.JSON(200, tokenResponse)

}

func register(ctx *gin.Context) {
	userRepository := newUserRepository()
	var u User
	decoder := json.NewDecoder(ctx.Request.Body)
	_ = decoder.Decode(&u)

	if !isEmailValid(u.Email) {
		ctx.String(400, "Email address is invalid.")
		return
	}

	us, _ := userRepository.Get(u.Email)

	if us != nil {
		ctx.String(409, "Account already exist!")
		return
	}

	h := sha256.New()
	h.Write([]byte(u.Password))
	password := hex.EncodeToString(h.Sum(nil))
	u.Password = password
	u.Id = uuid.New().String()
	u.CreatedAt = time.Now()
	_, err := userRepository.Save(u)
	if err != nil {
		ctx.String(400, "An error occurred while creating new account!")
		return
	}

	ctx.String(201, "Account created successfully!")
}

func serveWs(context *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(context.Writer, context.Request, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	cmd := exec.Command("docker", "exec", "-i", "f18b", "/bin/sh")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println(err)
		return
	}

	go func(c1 *websocket.Conn) {
		for {
			_, message, err := c1.NextReader()
			if err != nil {
				c1.Close()
				break
			}
			io.Copy(stdin, message)
		}
	}(ws)

	ws.WriteMessage(1, []byte("Starting...\n"))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return
	}

	go func() {
		outScan := bufio.NewScanner(stdout)
		for outScan.Scan() {
			ws.WriteMessage(1, outScan.Bytes())
		}

	}()

	errScan := bufio.NewScanner(stderr)
	for errScan.Scan() {
		ws.WriteMessage(1, errScan.Bytes())
	}

	if err := cmd.Wait(); err != nil {
		log.Println(err)
		return
	}

	ws.WriteMessage(1, []byte("Ending.."))

}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
