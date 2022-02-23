# Stage 1 - build
FROM golang:alpine AS build
RUN go version

WORKDIR /build/
COPY go.* ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app

# Stage 2 - run
FROM alpine:latest
WORKDIR /root/
COPY --from=build /build/app .
RUN apk --no-cache add curl

ENTRYPOINT ["./app"]
EXPOSE 8080