FROM alpine
WORKDIR /home/appuser/appservices/certs
COPY . .
ENTRYPOINT ./checkcerts

#### For build
#FROM golang:alpine
#WORKDIR /appservices/certs
#COPY . .
#RUN go mod download
#CMD go build -ldflags="-s -w" -o checkcerts && ./checkcerts