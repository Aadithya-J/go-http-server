# go-http-server

A simple ,lightweight , HTTP/HTTPS server written in Go, without using the http package and utilizing only the `net` package. This server supports file serving, directory listing, and basic HTTP GET and POST request handling. It features command-line flag configurations for directory, port, and HTTPS settings, making it versatile and easy to deploy for various use cases.

#### Command-Line Flags

| Flag         | Description                      | Default Value |
|--------------|----------------------------------|---------------|
| `--directory`| Directory to serve files from    | `.`           |
| `--port`     | Port to bind the server to       | `4221`        |
| `--https`    | Enable HTTPS (requires cert and key) | `false`       |
| `--cert`     | Path to SSL certificate file     | `cert.pem`    |
| `--key`      | Path to SSL key file             | `key.pem`     |
