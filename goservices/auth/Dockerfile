### For dev
FROM golang:alpine
WORKDIR /appservices/goservices/shared
COPY shared .
WORKDIR /appservices/goservices/auth
COPY auth .
RUN go mod download
ENTRYPOINT go run main.go

#### For build
#FROM golang:alpine
#WORKDIR /app/goservices/shared
#COPY shared .
#WORKDIR /app/goservices/auth
#COPY auth .
#RUN go mod download
#CMD go build -ldflags="-s -w" -o auth && ./auth

#### if i want to zip it
# RUN apk add upx
# CMD go build -ldflags="-s -w" -o auth && upx -9 -k auth && ./auth

#### For production
#FROM alpine
#WORKDIR /appservices/goservices/shared
#COPY shared .
#WORKDIR /appservices/goservices/auth
#COPY auth .
#CMD chmod +x ./auth && ./auth