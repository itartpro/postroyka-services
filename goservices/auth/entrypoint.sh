#!/bin/sh

chown -R appuser:appusergroup /home/appuser/appservices && \
chown -R appuser:appusergroup /home/appuser/appservices/certs && \
chown -R appuser:appusergroup /home/appuser/appservices/uploads && \
chown -R appuser:appusergroup /home/appuser/appservices/goservices/shared && \
exec runuser -u appuser "$@"