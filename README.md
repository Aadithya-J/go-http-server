# go-http-server

This project implements HTTP and HTTPS servers using Go's standard libraries (net and crypto/tls). It serves static files, handles basic HTTP GET requests, and logs request details. Configurable via command-line flags for port, directory, and optional HTTPS with custom SSL/TLS certificates. Ideal for educational purposes to explore fundamental server concepts.

To run the http server with default config 
```bash
go run main.go config.go server.go --dir=./public --port=3000
```
To run the https server generate keys using openssl and include the paths in flag
```bash
go run main.go config.go server.go --dir=./public --port=3000 --https --cert=./tlsCert/cert.pem --key=./tlsCert/key.pem
```
server runs in 
```
http://localhost:3000/home
https://localhost:3000/home
```

#### Command-Line Flags

| Flag         | Description                      | Default Value |
|--------------|----------------------------------|---------------|
| `--dir`| Directory to serve files from    | `.`           |
| `--port`     | Port to bind the server to       | `4221`        |
| `--https`    | Enable HTTPS (requires cert and key) | `false`       |
| `--cert`     | Path to SSL certificate file     | `cert.pem`    |
| `--key`      | Path to SSL key file             | `key.pem`     |
