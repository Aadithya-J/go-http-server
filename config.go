package main

import "flag"

type Config struct {
	Directory            string
	Port                 string
	UseHttps             bool
	CertFile             string
	KeyFile              string
	MaxRequestsPerSecond int
}

func parseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Directory, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.StringVar(&cfg.Port, "port", "8080", "the port to serve HTTP on")
	flag.BoolVar(&cfg.UseHttps, "https", false, "whether to use HTTPS")
	flag.StringVar(&cfg.CertFile, "cert", "", "the certificate file")
	flag.StringVar(&cfg.KeyFile, "key", "", "the key file")
	flag.IntVar(&cfg.MaxRequestsPerSecond, "max", 10, "the maximum number of requests per second")

	return cfg
}
