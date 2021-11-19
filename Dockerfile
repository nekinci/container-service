FROM docker:dind

RUN apk add --no-cache git make go curl

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
VOLUME /var/run/docker.sock /var/run/docker.sock
COPY . .
RUN go build .

EXPOSE 80
CMD ["./paas"]