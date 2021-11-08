package api

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/nekinci/paas/application"
	"github.com/nekinci/paas/specification"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"sync"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	allowedEndpoints = []string{"/register", "/login", "/refreshToken", "/logs", "/terminal", "/appState"}
	host             = "containerdemo.live"
)

type WebApi struct {
	appCtx *application.Context
	r      *gin.Engine
}

func ListenAndServe(c *application.Context) {
	addr := "127.0.0.1:8070"
	r := gin.Default()
	w := WebApi{
		appCtx: c,
		r:      r,
	}

	//	r.Use(CORSMiddleware())
	r.Use(wsMiddleware)
	r.Use(authMiddleware)

	r.POST("/login", login)

	r.POST("/register", register)

	r.POST("/refreshToken", refreshToken)

	r.GET("/myApps", w.myApps)

	r.GET("/terminal", w.terminalWS)
	r.GET("/logs", w.logsWS)
	r.GET("/ws", w.logsWS)
	r.POST("/run", w.runApplication)
	r.GET("/info/:appName", w.appInfo)
	r.GET("/appState", w.getState)
	_ = w
	_ = r.Run(addr)

}

func (w *WebApi) appInfo(context *gin.Context) {
	appName := (context).Param("appName")

	if appName == "" {
		context.JSON(404, gin.H{
			"message": "Application name not allowed.",
		})
		return
	}

	app := w.appCtx.GetApplication(appName)
	if app != nil {
		info := app.GetApplicationInfo()
		context.JSON(200, gin.H{
			"name":         info.Name,
			"url":          "http://" + info.Name + "." + host,
			"environments": nil,
			"owner":        info.UserEmail,
			"containerId":  info.Id,
			"startTime":    info.StartTime,
			"status":       info.Status,
			"image":        info.Image,
		})
		return
	}

	context.JSON(404, gin.H{
		"message": "Application not found!",
	})

}

func (w *WebApi) runApplication(context *gin.Context) {
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(400, gin.H{
			"code":    400,
			"message": "File request invalid.",
		})
		return
	}
	open, err := file.Open()

	if err != nil {
		context.JSON(400, gin.H{
			"code":    400,
			"message": "Unknown error",
		})
		return
	}

	defer open.Close()
	all, _ := ioutil.ReadAll(open)
	application, appErr := specification.NewApplication(all)

	if appErr != nil {
		context.JSON(400, gin.H{
			"code":    400,
			"message": "File invalid, please upload valid yaml file.",
		})
		return
	}
	currentUserEmail, _ := context.Get("CurrentUserEmail")
	application.Email = currentUserEmail.(string)

	if application.Name == "nginxapp" && application.Image == "nginx" {
		w.appCtx.AddTryItUser(application.Name, application.Email)
		context.JSON(200, gin.H{
			"code":    200,
			"message": "Container started.",
			"appName": application.Name,
		})
		return
	}

	handleErr := w.appCtx.Handle(application, true)

	if handleErr != nil {
		context.JSON(400, gin.H{
			"code":    400,
			"message": handleErr.Error(),
		})
		return
	}

	context.JSON(200, gin.H{
		"code":    200,
		"message": "Container started.",
		"appName": application.Name,
	})
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
	if path != "/logs" && path != "/terminal" && path != "/appState" {
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

	currentApp, isThere := context.GetQuery("currentApp")
	if !isThere {
		currentApp = ""
	}

	context.Set("CurrentApp", currentApp)
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

func (w *WebApi) logsWS(context *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		log.Printf("upgrade: %v\n", err)
		return
	}

	currentApp, _ := context.Get("CurrentApp")
	currentEmail, _ := context.Get("CurrentUserEmail")
	app := w.appCtx.GetApplication(currentApp.(string))

	if app == nil {
		context.Status(404)
		return
	}

	if app.GetApplicationInfo().UserEmail != currentEmail.(string) && app.GetApplicationInfo().UserEmail != "superuser@containerdemo.live" {
		context.Status(401)
		return
	}

	var mu sync.Mutex
	app.LogStream(func(l application.Log) {
		mu.Lock()
		defer mu.Unlock()
		if l.Show {
			ws.WriteMessage(1, []byte(l.Format()))
		}
	})

	go app.ListenLogs()

}

func (w *WebApi) terminalWS(context *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(context.Writer, context.Request, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	currentApp, _ := context.Get("CurrentApp")
	currentEmail, _ := context.Get("CurrentUserEmail")
	app := w.appCtx.GetApplication(currentApp.(string))

	if app == nil {
		context.Status(404)
		return
	}

	if app.GetApplicationInfo().UserEmail != currentEmail.(string) && app.GetApplicationInfo().UserEmail != "superuser@containerdemo.live" {
		context.Status(401)
		return
	}

	terminal, cancel, err := app.OpenTerminal()

	if err != nil {
		ws.WriteMessage(1, []byte(fmt.Sprintf("%v", err)))
		ws.Close()
		return
	}

	defer cancel()

	go func(c1 *websocket.Conn) {
		for {
			_, msg, err := c1.ReadMessage()
			if err != nil {
				c1.Close()
				break
			}

			io.Copy(*terminal.Stdin, bytes.NewReader(msg))

		}
	}(ws)

	go func() {
		outScan := bufio.NewScanner(*terminal.Stdout)
		var mu sync.Mutex
		for outScan.Scan() {
			mu.Lock()
			ws.WriteMessage(1, outScan.Bytes())
			mu.Unlock()
		}
	}()

	errScan := bufio.NewScanner(*terminal.Stderr)
	var mu sync.Mutex
	for errScan.Scan() {
		mu.Lock()
		ws.WriteMessage(1, errScan.Bytes())
		mu.Unlock()
	}

	ws.WriteMessage(1, []byte("Ending.."))

}

func (w *WebApi) myApps(context *gin.Context) {
	email, _ := context.Get("CurrentUserEmail")
	context.JSON(200, w.appCtx.GetApplicationsByUser(email.(string)))
}

func (w *WebApi) getState(context *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	eventListener := func(event application.StateEvent) {
		ws.WriteJSON(event)
	}

	w.appCtx.SendInitEvent(eventListener)

	w.appCtx.AddStateListener(eventListener)
}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
