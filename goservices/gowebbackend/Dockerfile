## For dev
#FROM golang:alpine
#WORKDIR /home/appuser/appservices/goservices/shared
#COPY shared .
#WORKDIR /home/appuser/appservices/goservices/gowebbackend
#COPY gowebbackend .
#RUN go mod download
#ENTRYPOINT go run main.go

## For build
#FROM golang:alpine
#WORKDIR /appservices/goservices/shared
#COPY shared .
#WORKDIR /appservices/goservices/gowebbackend
#COPY gowebbackend .
#RUN go mod download
#CMD go build -ldflags="-s -w" -o gowebbackend && ./gowebbackend

## if i want to zip it
# RUN apk add upx
# CMD go build -ldflags="-s -w" -o gowebbackend && upx -9 -k gowebbackend && ./gowebbackend

## For production
FROM alpine

RUN apk add --no-cache tzdata
ENV TZ Europe/Moscow
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apk del tzdata

WORKDIR /home/appuser/appservices/goservices/shared
COPY shared .
WORKDIR /home/appuser/appservices/goservices/gowebbackend
COPY gowebbackend .

ARG UID
ARG GID
RUN apk add runuser
RUN addgroup -g $GID -S appusergroup && \
adduser -S -u $UID appuser -G appusergroup
ENTRYPOINT ["/bin/sh", "entrypoint.sh"]
CMD ["./gowebbackend"]