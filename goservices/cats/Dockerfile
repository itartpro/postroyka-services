## For dev
FROM golang:alpine
WORKDIR /appservices/goservices/shared
COPY shared .
WORKDIR /appservices/goservices/cats
COPY cats .
RUN go mod download
ENTRYPOINT go run main.go

#### For build
#FROM golang:alpine
#WORKDIR /appservices/goservices/shared
#COPY shared .
#WORKDIR /appservices/goservices/cats
#COPY cats .
#RUN go mod download
#CMD go build -ldflags="-s -w" -o cats && ./cats

#### if i want to zip it
# RUN apk add upx
# CMD go build -ldflags="-s -w" -o cats && upx -9 -k cats && ./cats

#### For production
#FROM alpine
#WORKDIR /appservices/goservices/shared
#COPY shared .
#WORKDIR /appservices/goservices/cats
#COPY cats .
#CMD chmod +x ./cats && ./cats