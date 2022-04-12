package main

import (
	log "github.com/jensneuse/abstractlogger"
	"net"
	Url "net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func getServicesUralsFromEnvironment() []ServiceConfig {
	var servicesUrls []ServiceConfig
	for _, e := range os.Environ() {
		ePair := strings.SplitN(e, ";", 2)
		urlPairs := strings.SplitN(ePair[0], "=", 2)
		if strings.HasPrefix(urlPairs[0], "URL_") {
			u, err := Url.Parse(urlPairs[1])
			if err != nil {
				log.Error(err)
				continue
			}
			host, _, _ := net.SplitHostPort(u.Host)

			if host == "localhost" || host == "127.0.0.1" {
				host = urlPairs[0] + "_" + host
			}

			serviceUrlConf := ServiceConfig{
				Name: host,
				URL:  urlPairs[1],
			}
			if len(ePair) > 1 && ePair[1] == "ws=true" {
				serviceUrlConf.WS = strings.Replace(urlPairs[1], "http", "ws", 1)
			}

			servicesUrls = append(servicesUrls, serviceUrlConf)
		}
	}
	if len(servicesUrls) == 0 {
		logger().Warn("no services found in env vas")
	}
	return servicesUrls
}

func getSchemaPollInterval() time.Duration {
	envValue, _ := os.LookupEnv("SCHEMA_POLL_INTERVAL")
	pollInterval, err := strconv.Atoi(envValue)
	if err != nil {
		pollInterval = 0
	}

	return time.Duration(pollInterval) * time.Second
}
