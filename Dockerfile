FROM golang:1.11-alpine as builder
COPY . /go/src/github.com/cjsaylor/contact-api
WORKDIR /go/src/github.com/cjsaylor/contact-api
RUN ls vendor/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s" -v -o web ./cmd/web/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/src/github.com/cjsaylor/contact-api/web web
EXPOSE 8080

CMD ["./web"]