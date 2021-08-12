### For dev
FROM golang:alpine
WORKDIR /appservices/goservices/shared
COPY shared .
WORKDIR /appservices/goservices/gpics
COPY gpics .
RUN go mod download
ENTRYPOINT go run main.go

#### For build
#FROM golang:alpine
#WORKDIR /app/goservices/shared
#COPY shared .
#WORKDIR /app/goservices/gpics
#COPY gpics .
#RUN go mod download
#CMD go build -ldflags="-s -w" -o gpics && ./gpics

#### if i want to zip it
# RUN apk add upx
# CMD go build -ldflags="-s -w" -o gpics && upx -9 -k gpics && ./gpics

#### For production
#FROM alpine
#WORKDIR /appservices/goservices/shared
#COPY shared .
#WORKDIR /appservices/goservices/gpics
#COPY gpics .
#CMD chmod +x ./gpics && ./gpics