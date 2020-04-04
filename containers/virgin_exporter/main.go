package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/maesoser/virgin_exporter/virgin"
)

func GetEnvStr(name, value string) string {
	if os.Getenv(name) != "" {
		return os.Getenv(name)
	}
	return value
}

func main() {
	Address := flag.String("a", GetEnvStr("VIRGIN_ADDR", "192.168.0.1"), "Router's address")
	Pass := flag.String("w", GetEnvStr("VIRGIN_PASSWD", "admin"), "Router's password")
	User := flag.String("u", GetEnvStr("VIRGIN_USER", "admin"), "Router's username")
	Port := flag.String("p", GetEnvStr("VIRGIN_PORT", "9235"), "Prometheus port")
	Verbose := flag.Bool("v", false, "Verbose output")
	flag.Parse()
	router := virgin.NewRouter(*Address, *User, *Pass)
	router.Verbose = *Verbose

	log.Printf("Serving exporter at port %s\n", *Port)
	err := http.ListenAndServe(":"+*Port, router)
	log.Panic(err)

}
