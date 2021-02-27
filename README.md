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
