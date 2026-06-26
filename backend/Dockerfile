FROM golang:1.25.10-alpine AS build

WORKDIR /src

ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
ENV TZ=Asia/Shanghai

COPY --from=build /out/server /app/server

EXPOSE 8080
CMD ["/app/server"]
