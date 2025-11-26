FROM golang:1.25-alpine AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY . .
RUN go build -o main main.go

FROM alpine
WORKDIR /app

COPY --from=builder /app/main .
COPY app.env .


EXPOSE 8080
CMD [ "/app/main" ]