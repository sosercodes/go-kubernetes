# Go React Kubernetes Tutorial

Let's create a new directory `api` and change into it.

```bash
go mod init github.com/somnidev/go-kubernetes
```

Create a new `main.go` file and add the following text.

```go
package main

import (
    "log"
)

func main() {
    log.Println("Hello Go!")
}
```

Once you saved the file return to your shell and run the application.

```bash
% go run main.go 
2021/02/26 18:40:46 Hello Go!
```

## Setting up Gin and your endpoint

We can include Gin just like any other dependency in Go.

```bash
go get -u github.com/gin-gonic/gin
```

This downloads the Gin dependency and adds it to the 'go.mod' file. 

Now we recognize yet another new file in our directory – `go.sum`. That file contains checksums for direct and indirect dependencies of our application and actually some more things not that relevant for our cause.

Now that we have Gin included, let’s set up and start creating the server.

```bash
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.Run()
}
```

Yes, that’s all we need to set up Gin. Using `go build` to build the application and then running the created executable will show us where we can reach Gin (should be curl localhost:8080 on default).

That will give us an HTTP status `404` and some log output indicating that we’re set up and are successfully serving HTTP errors. How to setup [a starting template inside your project](https://gin-gonic.com/docs/quickstart/) and [some sample projcets](https://gin-gonic.com/docs/users/).

```bash
package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    // Listen and Server to 0.0.0.0:8080
    r.Run(":8080")
}
```

## Build a Docker Image

### Single stage docker build

We can build a single-stage docker image using the `golang:alpine` image. Here’s the Dockerfile

```bash
FROM golang:alpine
WORKDIR /app
ADD . /app
RUN cd /app && go build -o goapp
ENTRYPOINT ./goapp
```

When we check the size using `docker images` we get about **408** MB, just for our single little Go binary. That's pretty big.

### Multi stage docker build

Now let’s try a multi-stage build using this new Dockerfile.

```bash
# build the app - builder image
FROM golang:1.16.0-alpine3.13 AS builder
RUN mkdir /build
ADD *.go *.mod *.sum /build/
WORKDIR /build
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o golang-app .

# create clean app image
FROM alpine:3.13
COPY --from=builder /build/golang-app .
ENTRYPOINT [ "./golang-app" ]
```

Let's check the size of the image.

```bash
% docker images | grep somnidev
somnidev/go-kubernetes-api                       latest                                                  f5d04da81e8d   14 minutes ago   15MB
```

Now we get an image that is really small. Only **15MB**. Let's run it.

```bash
docker run --rm -p 8080:8080 somnidev/go-kubernetes-api
```
