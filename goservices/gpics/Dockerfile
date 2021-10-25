## For production
FROM alpine

ENV MAINPNAME gpics
ARG UID
ARG GID
ARG MYUSERNAME
ARG MYUSERGROUP
ENV MYUSERNAME ${MYUSERNAME}
ENV MYUSERGROUP ${MYUSERGROUP}

RUN apk add --no-cache tzdata
ENV TZ Europe/Moscow
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apk del tzdata

RUN apk add runuser
RUN addgroup -g $GID -S $MYUSERGROUP && \
adduser -S -u $UID $MYUSERNAME -G $MYUSERGROUP

WORKDIR /home/$MYUSERNAME/appservices/goservices/shared
COPY shared .
WORKDIR /home/$MYUSERNAME/appservices/goservices/$MAINPNAME
COPY $MAINPNAME .

ENTRYPOINT chown -R $MYUSERNAME:$MYUSERGROUP /home/$MYUSERNAME/appservices && \
exec runuser -u $MYUSERNAME ./$MAINPNAME