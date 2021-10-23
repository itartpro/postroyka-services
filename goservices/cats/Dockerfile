## For dev
FROM golang:alpine

RUN apk add --no-cache tzdata
ENV TZ Europe/Moscow
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apk del tzdata

WORKDIR /home/appuser/appservices/goservices/shared
COPY shared .
WORKDIR /home/appuser/appservices/goservices/cats
COPY cats .
RUN go mod download
CMD go run main.go

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

#RUN apk add --no-cache tzdata
#ENV TZ Europe/Moscow
#RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
#RUN apk del tzdata

#WORKDIR /appservices/goservices/shared
#COPY shared .
#WORKDIR /appservices/goservices/cats
#COPY cats .

#ARG UID
#ARG GID
#RUN apk add runuser
#RUN addgroup -g $GID -S appusergroup && \
#adduser -S -u $UID appuser -G appusergroup
#ENTRYPOINT ["/bin/sh", "entrypoint.sh"]

#CMD ["./cats"]