Go application that has a single endpoint `/` which shows how many times the page was requested. Docker Compose spins 3 instances of the application, and **nginx**, acting like a load-balancer, in ~~round-robin~~ least-connections manner distributes requests between them. 

## Technologies

- nginx
- jenkins
- docker & docker compose
- redis
- go http server

## Prerequisite

Generate self-signed certificate:

```sh
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout nginx-selfsigned.key -out nginx-selfsigned.crt
```

- `req`: command to create and process certificate requests
- `-x509`: specifies that you want to create a self-signed certificate instead of a certificate request
- `-nodes`: tells to not encrypt the private key for server application to read the key without a passphrase
- `-days 365`: sets the validity of the certificate to 365 days
- `-newkey rsa:2048`: generates a new RSA key of 2048 bits
- `-keyout nginx-selfsigned.key`: specifies the filename for the private key
- `-out nginx-selfsigned.crt`: specifies the filename for the self-signed certificate
