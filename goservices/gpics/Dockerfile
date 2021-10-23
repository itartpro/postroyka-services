### For dev
#FROM golang:alpine
#WORKDIR /appservices/goservices/shared
#COPY shared .
#WORKDIR /appservices/goservices/gpics
#COPY gpics .
#RUN go mod download
#ENTRYPOINT go run main.go

#### For build
#FROM golang:alpine
#WORKDIR /appservices/goservices/shared
#COPY shared .
#WORKDIR /appservices/goservices/gpics
#COPY gpics .
#RUN go mod download
#CMD go build -ldflags="-s -w" -o gpics && ./gpics

#### if i want to zip it
# RUN apk add upx
# CMD go build -ldflags="-s -w" -o gpics && upx -9 -k gpics && ./gpics

## For production
FROM alpine

RUN apk add --no-cache tzdata
ENV TZ Europe/Moscow
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apk del tzdata

WORKDIR /home/appuser/appservices/goservices/shared
COPY shared .
WORKDIR /home/appuser/appservices/goservices/gpics
COPY gpics .

ARG UID
ARG GID
RUN apk add runuser
RUN addgroup -g $GID -S appusergroup && \
adduser -S -u $UID appuser -G appusergroup
ENTRYPOINT ["/bin/sh", "entrypoint.sh"]
CMD ["./gpics"]