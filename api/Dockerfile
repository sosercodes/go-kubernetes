# build the app - builder image
FROM golang:1.23.6-alpine3.20 AS builder
RUN mkdir /build
ADD *.go *.mod *.sum /build/
WORKDIR /build
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o golang-app .

# create clean app image
FROM alpine:3.21
COPY --from=builder /build/golang-app .
ENTRYPOINT [ "./golang-app" ]