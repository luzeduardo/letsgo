## TLS config for local development
To generate the tls certificate for local development just:

```sh
cd tls
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```
