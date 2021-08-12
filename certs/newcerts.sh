# Install libraries needed for certs generation on alpine linux
apk update && apk upgrade && apk add --no-cache bash procps drill git coreutils libidn curl &&
apk add openssl && chmod +x newcerts.sh &&
# Remove old certificates
rm *.pem &&
rm *.srl &&
rm *.csr &&
rm *.key &&
rm *.pub &&
rm *.cert &&
# Make simple aes256 pem, pub, key (for use with dgrijalva/jwt-go)
# Replace $JWTKEY_PASS with your own password
openssl genrsa -aes256 -passout pass:$JWTKEY_PASS -out jwtkey.pem 2048 &&
openssl rsa -pubout -in jwtkey.pem -passin pass:$JWTKEY_PASS -pubout -out jwtkey.pub &&
openssl pkey -in jwtkey.pem -passin pass:$JWTKEY_PASS -out jwtkey.key &&
# Make key, srl, cert, key using config in service.conf.
# In my use case for use with credentials.NewServerTLSFromFile in Go, for gRPC communication
openssl genrsa -out ca.key 4096 &&
openssl req -new -x509 -key ca.key -sha256 -subj "/C=RU/ST=Kineshma/O=client.local" -days 3650 -out ca.cert &&
openssl genrsa -out service.key 4096 &&
sleep 1 &&
openssl req -new -key service.key -out service.csr -config certificate.conf &&
sleep 1 &&
openssl x509 -req -in service.csr -CA ca.cert -CAkey ca.key -CAcreateserial -out service.pem -days 3650 -sha256 -extfile certificate.conf -extensions req_ext