#!/bin/bash
# Inspired from: https://github.com/grpc/grpc-java/tree/master/examples#generating-self-signed-certificates-for-use-with-grpc

# Output files
# ca.key: Certificate Authority private key file (this shouldn't be shared in real-life)
# ca.crt: Certificate Authority trust certificate (this should be shared with users in real-life)
# service.key: Server private key, password protected (this shouldn't be shared)
# service.csr: Server certificate signing request (this should be shared with the CA owner)
# service.pem: Service certificate file (this shouldn't be shared)

# Summary 
# Private files: ca.key, server.key, service.pem
# "Share" files: ca.crt (needed by the client), service.csr (needed by the CA)

openssl genrsa -out ca.key 4096

openssl req -new -x509 -key ca.key -sha256 -days 365 -subj "/C=US/ST=NJ/O=CA, Inc." -out ca.crt

openssl genrsa -out service.key 4096

openssl req -new -key service.key -out service.csr -config certificate.conf

openssl x509 -req -in service.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out service.pem -days 365 -sha256 -extfile certificate.conf -extensions req_ext

openssl x509 -in service.pem -text -noout


# Reference: https://itnext.io/practical-guide-to-securing-grpc-connections-with-go-and-tls-part-1-f63058e9d6d1