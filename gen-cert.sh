#!/bin/bash
openssl genrsa -out key.pem 4096
openssl req -new -x509 -key key.pem -out cert.pem -days 365 -subj "/CN=my-proxy-ca"

cp ./cert.pem /usr/local/share/ca-certificates/veleno.crt
update-ca-certificates
