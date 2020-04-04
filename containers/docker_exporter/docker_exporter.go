package main

import (
	"flag"
	"log"
	"math"
	"net/http"
        "os"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// return rounded percentage
func percent(val float64, total float64) int {
	if total <= 0 {
		return 0
	}
	return round((val / total) * 100)
}

func GetEnvStr(name, value string) string {
	if os.Getenv(name) != "" {
		return os.Getenv(name)
	}
	return value
}

func main() {

	PromPort := flag.String("port",
		GetEnvStr("DOCKER_EXPORTER_PORT","9134"),
		"Prometheus Port"
	)
	DockerEndPoint := flag.String("endpoint",
		GetEnvStr("DOCKER_EXPORTER_SOCKET","/var/run/docker.sock"),
		"Docker Endpoint"
	)
	flag.Parse()

	log.Println("Starting docker runtime")
	docker := &Docker{}
	docker.Start(*DockerEndPoint)

	log.Printf("Serving exporter at port %s\n", *PromPort)
	err := http.ListenAndServe(":"+*PromPort, docker)
	log.Panic(err)

}
