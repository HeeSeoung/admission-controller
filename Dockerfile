FROM golang:1.21 as builder
WORKDIR /app
 COPY go.mod .
 COPY go.sum .

 RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o serverd cmd/main.go

FROM gcr.io/distroless/base
COPY --from=builder /app/serverd /
EXPOSE 8443

CMD ["/serverd"]