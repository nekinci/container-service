FROM docker:dind

RUN apk add --no-cache git make go curl

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
COPY . .
RUN go build .

EXPOSE 80
CMD ["./paas"]