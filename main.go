package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/bogdanovich/dns_resolver"
	"github.com/gorilla/mux"
	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

var cfg ConfigOptions

type ConfigOptions struct {
	Resolver string `env:"RESOLVER" env-default:127.0.0.1"`
	Lookup   string `env:"LOOKUP" env-default:2.0.0.127.my.domain"`
	HttpPort string `env:"HTTPPORT" env-default:"8080"`
	LogLevel string `env:"LOGLEVEL" env-default:"info"`
}

func httpHealth(w http.ResponseWriter, r *http.Request) {
	lookup_host := r.Header.Get("lookup")

	if lookup_host == "" {
		lookup_host = cfg.Lookup
	} else {
		log.Debug("lookup override: ", lookup_host)
	}

	resolver := dns_resolver.New([]string{cfg.Resolver})
	resolver.RetryTimes = 3
	ip, err := resolver.LookupHost(lookup_host)
	if err != nil {
		writeHttpError(w, fmt.Sprintf("%v", err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	result := "{\"result\": alive}"

	log.Debug(ip)
	io.WriteString(w, result)
}

func writeHttpError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, fmt.Sprintf("{\"level\":\"error\",\"msg\":\"%v\"}", msg))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, "{\"level\":\"error\",\"msg\":\"endpoint not found\"}")
}

func setupLogging(logLevel string) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	if logLevel == "info" {
		log.SetLevel(log.InfoLevel)
	} else if logLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	}
}

func init() {
	// config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Panic(err)
	}
	// logging
	setupLogging(cfg.LogLevel)
}

func main() {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFound)

	r.HandleFunc("/health", httpHealth)

	srv := &http.Server{
		Addr:         "0.0.0.0:" + cfg.HttpPort,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
