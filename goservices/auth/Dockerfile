### For dev
FROM golang:alpine

RUN apk add --no-cache tzdata
ENV TZ Europe/Moscow
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apk del tzdata

WORKDIR /home/appuser/appservices/goservices/shared
COPY shared .
WORKDIR /home/appuser/appservices/goservices/auth
COPY auth .
RUN go mod download
CMD go run main.go

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

#RUN apk add --no-cache tzdata
#ENV TZ Europe/Moscow
#RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
#RUN apk del tzdata

#WORKDIR /appservices/goservices/shared
#COPY shared .
#WORKDIR /appservices/goservices/auth
#COPY auth .

#ARG UID
#ARG GID
#RUN apk add runuser
#RUN addgroup -g $GID -S appusergroup && \
#adduser -S -u $UID appuser -G appusergroup
#ENTRYPOINT ["/bin/sh", "entrypoint.sh"]

#CMD ["./auth"]