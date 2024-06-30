package main

import "flag"

type Config struct {
	Directory string
	Port      string
	UseHttps  bool
	CertFile  string
	KeyFile   string
}

func parseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Directory, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.StringVar(&cfg.Port, "port", "8080", "the port to serve HTTP on")
	flag.BoolVar(&cfg.UseHttps, "https", false, "whether to use HTTPS")
	flag.StringVar(&cfg.CertFile, "cert", "", "the certificate file")
	flag.StringVar(&cfg.KeyFile, "key", "", "the key file")

	return cfg
}
