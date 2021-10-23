# Install libraries needed for certs generation on alpine linux
apk update && apk upgrade &&
apk add --no-cache bash procps drill git coreutils libidn curl openssh-keygen openssl &&
chmod +x newcerts.sh &&
# Remove old certificates
rm *.pem &&
rm *.srl &&
rm *.csr &&
rm *.key &&
rm *.pub &&
rm *.cert &&
# Make simple aes256 pem, pub (for use with  golang-jwt/jwt )
# https://gist.github.com/ygotthilf/baa58da5c3dd1f69fae9
ssh-keygen -t rsa -P "" -b 2048 -m PEM -f jwtkey.pem &&
openssl rsa -in jwtkey.pem -pubout -outform PEM -out jwtkey.pub &&
# Make key, srl, cert, key using config in service.conf.
# In my use case for use with credentials.NewServerTLSFromFile in Go, for gRPC communication
openssl genrsa -out ca.key 2048 &&
openssl req -new -x509 -key ca.key -sha256 -subj "/C=RU/ST=Kineshma/O=postroyka.localhost" -days 3650 -out ca.cert &&
openssl genrsa -out service.key 2048 &&
sleep 1 &&
openssl req -new -key service.key -out service.csr -config certificate.conf &&
sleep 1 &&
openssl x509 -req -in service.csr -CA ca.cert -CAkey ca.key -CAcreateserial -out service.pem -days 3650 -sha256 -extfile certificate.conf -extensions req_ext