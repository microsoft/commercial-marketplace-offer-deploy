

```
openssl req -x509 -nodes -new -sha256 -days 1024 -newkey rsa:2048 -keyout ca.key -out ca.pem -subj "/C=US/CN=Localhost-Root-CA"
openssl x509 -outform pem -in ca.pem -out ca.crt
```

```
openssl req -new -nodes -newkey rsa:2048 -keyout localhost.key -out localhost.csr -subj "/C=US/ST=WA/L=Redmond/O=GPS/CN=localhost.local"
openssl x509 -req -sha256 -days 1024 -in localhost.csr -CA ca.pem -CAkey ca.key -CAcreateserial -extfile domains.ext -out localhost.crt
```