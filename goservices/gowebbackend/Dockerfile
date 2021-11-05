### For dev
#FROM golang:alpine
#
#ENV MAINPNAME gowebbackend
#ARG UID
#ARG GID
#ARG MYUSERNAME
#ARG MYUSERGROUP
#ENV MYUSERNAME ${MYUSERNAME}
#ENV MYUSERGROUP ${MYUSERGROUP}
#
#RUN apk add runuser
#RUN addgroup -g $GID -S $MYUSERGROUP && \
#adduser -S -u $UID $MYUSERNAME -G $MYUSERGROUP
#
#WORKDIR /home/$MYUSERNAME/appservices/goservices/shared
#COPY shared .
#WORKDIR /home/$MYUSERNAME/appservices/goservices/$MAINPNAME
#COPY $MAINPNAME .

#ENTRYPOINT chown -R $MYUSERNAME:$MYUSERGROUP /home/$MYUSERNAME/appservices && \
#exec runuser -u $MYUSERNAME go run main.go
## For build
#CMD go build -ldflags="-s -w" -o $MAINPNAME && ./$MAINPNAME
# if i want to zip it
#RUN apk add upx
#ENTRYPOINT go build -ldflags="-s -w" -o $MAINPNAME && upx -9 -k $MAINPNAME && ./$MAINPNAME

## For production
FROM alpine

ENV MAINPNAME gowebbackend
ARG UID
ARG GID
ARG MYUSERNAME
ARG MYUSERGROUP
ENV MYUSERNAME ${MYUSERNAME}
ENV MYUSERGROUP ${MYUSERGROUP}

RUN apk add runuser
RUN addgroup -g $GID -S $MYUSERGROUP && \
adduser -S -u $UID $MYUSERNAME -G $MYUSERGROUP

WORKDIR /home/$MYUSERNAME/appservices/goservices/shared
COPY shared .
WORKDIR /home/$MYUSERNAME/appservices/goservices/$MAINPNAME
COPY $MAINPNAME .

ENTRYPOINT chown -R $MYUSERNAME:$MYUSERGROUP /home/$MYUSERNAME/appservices && \
exec runuser -u $MYUSERNAME ./$MAINPNAME