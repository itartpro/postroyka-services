# For dev
FROM golang:alpine
WORKDIR /appservices/goservices/shared
COPY shared .
WORKDIR /appservices/goservices/gowebbackend
COPY gowebbackend .
RUN go mod download
ENTRYPOINT go run main.go

#### For build
#FROM golang:alpine
#WORKDIR /appservices/goservices/shared
#COPY shared .
#WORKDIR /appservices/goservices/gowebbackend
#COPY gowebbackend .
#RUN go mod download
#CMD go build -ldflags="-s -w" -o gowebbackend && ./gowebbackend

#### if i want to zip it
# RUN apk add upx
# CMD go build -ldflags="-s -w" -o gowebbackend && upx -9 -k gowebbackend && ./gowebbackend

#### For production
#FROM alpine
#WORKDIR /appservices/goservices/shared
#COPY shared .
#WORKDIR /appservices/goservices/gowebbackend
#COPY gowebbackend .
#CMD chmod +x ./gowebbackend && ./gowebbackend