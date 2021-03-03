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

Run the app.

```bash
go run main.go
```

Now, we can list all package dependencies.

```bash
go list -m all
```

We can also build an executable file.

```bash
go build
```

## Build a Docker Image

### Single stage docker build

We can build a single-stage docker image using the `golang:alpine` image. Here’s the Dockerfile.

```bash
FROM golang:alpine
WORKDIR /app
ADD . /app
RUN cd /app && go build -o goapp
ENTRYPOINT ./goapp
```

Now we can build the docker image.

```bash
docker build -t somnidev/go-kubernetes-api:latest -t somnidev/go-kubernetes-api:0.1 -f Dockerfile .
```

When we check the size using `docker images` we get about **408** MB, just for our single little Go binary. That's pretty big.

### Multi stage docker build

With multi-stage builds, you use multiple FROM statements in your Dockerfile. Each FROM instruction can use a different base, and each of them begins a new stage of the build. You can selectively copy artifacts from one stage to another, leaving behind everything you don’t want in the final image - see [Use multi-stage builds](https://docs.docker.com/develop/develop-images/multistage-build/).

By default, the stages are not named, and you refer to them by their integer number, starting with 0 for the first FROM instruction. However, you can name your stages, by adding an `AS <NAME>` to the FROM instruction.

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

Now we can build our _multi stage_ docker image.

```bash
docker build -t somnidev/go-kubernetes-api:latest -t somnidev/go-kubernetes-api:0.1 -f Dockerfile .
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

## Creating a Deployment for Pods

For our first deployment we create a file called `pods.yaml`for our Go app. It creates a Replicaset that brings up one pod.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-kubernetes
  labels:
    app: go-kubernetes
spec:
  replicas: 1
  selector:
    matchLabels:
      app: go-kubernetes
  template:
    metadata:
      labels:
        app: go-kubernetes
    spec:
      containers:
      - name: go-kubernetes
        image: somnidev/go-kubernetes-api:0.1
        ports:
        - containerPort: 80
```

Before you begin, make sure your Kubernetes cluster is up and running.

Create the Deployment by running the following command.

```bash
kubectl apply -f pods.yaml
```

Since our pod defined in the _template_ has the _label_ `app: go-kubernetes` we can list `-l` only pods with that label.

```bash
kubectl get pods -l app=go-kubernetes -o wide
```

Get a shell into the running container.

```bash
kubectl exec --stdin --tty shell-demo -- /bin/bash
```

Now we can delete the deployment with the following command.

```bash
kubectl delete -f pod.yaml
```

## Creating a Service

Since you can't connect to Pods there is a concept called a _Kubernetes Service_. A _Kubernetes Service_ is an abstraction which defines a logical set of Pods running, whereby each Pod provides the same functionality.

When created, each Service is assigned its own unique IP address, also called _clusterIP_. This address is tied to the lifespan of the Service, and will not change while the Service is alive.

Pods can be configured to talk to the Service, and know that communication to the Service will be automatically load-balanced out to some pod that is a member of the Service.

We create a new file `services.yaml` for our service definitions.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: go-kubernetes
spec:
  ports:
  - name: http
    port: 80
    targetPort: 80
    nodePort: 30080
  selector:
    app: go-kubernetes
  type: NodePort
```

We need a _type: NodePort_ to expose our Service to the outside of the cluster. Notice that the NodePort has to be greater than `30000`. Notice we can remove the `NodePort` after we have tested the configuration and replace it with `type: ClusterIP`.

## Creating an Ingress

Assuming you have Docker for Mac installed, follow the next steps to set up the Nginx Ingress Controller on your local Kubernetes cluster - and take a look at the [Installation Guide](https://kubernetes.github.io/ingress-nginx/deploy/#docker-for-mac) for the actual link.

Run the following command to set up the NGINX Ingress.

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.44.0/deploy/static/provider/cloud/deploy.yaml
```

Verify the service was enabled by running the following:

```bash
kubectl get pods -n ingress-nginx
```

Or use the label.

```bash
kubectl get pods --all-namespaces -l app=ingress-nginx
```

## Ingress Configuration

Since `apiVersion: networking.k8s.io/v1beta1`is deprecated we use `apiVersion: networking.k8s.io/v1`. Take a look at Kubernetes Documentation for more information on [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/).

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: nginx-ingress
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/use-regex: 'true'
    nginx.ingress.kubernetes.io/rewrite-target: /$1
spec:
  rules:
  - http:
      paths:
      - path: /api/?(.*)
        pathType: Prefix
        backend:
          service:
            name: go-kubernetes-api
            port:
              number: 80
```

Verify that the configuration is correct and open `http://localhost:80/api/ping`in the browser.

## Securing the Ingress

In order to secure our Ingress Service we have to follow the steps shown in the Kubernetes Documentation [Connecting Applications with Services - Securing the Service](https://kubernetes.io/docs/concepts/services-networking/connect-applications-service/#securing-the-service).

