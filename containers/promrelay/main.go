package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)


// GetEnvStr gets a string env variable
func GetEnvStr(name, value string) string {
	if os.Getenv(name) != "" {
		return os.Getenv(name)
	}
	return value
}

// GetEnvInt gets a int env variable
func GetEnvInt(name string, value int) int {
	if os.Getenv(name) != "" {
		i, err := strconv.Atoi(os.Getenv(name))
		if err != nil {
			return value
		}
		return i
	}
	return value
}

type exporters []string

func (i *exporters) String() string {
    return "my string representation"
}

func (i *exporters) Set(value string) error {
    *i = append(*i, value)
    return nil
}

func main() {

	ServerOpt := flag.Bool("server", false, "Server mode")
	NoReuseOpt := flag.Bool("noreuse", false, "Remove the metrics once has been scraped once.")
	ExposedPort := flag.Int("expose", GetEnvInt("EXPOSED_PORT", 9091), "Port to expose the metrics.")
	ListenPort := flag.Int("listen", GetEnvInt("LISTEN_PORT", 0), "HTTP Port to listen to incoming connections. If it is not configure it will use the ports by default 80/443")
	OutputFile := flag.String("output", GetEnvStr("OUTPUT_FILE", ""), "Output file to write when an updated is received")
	HMACKey := flag.String("hmac", GetEnvStr("HMAC_KEY", ""), "HMAC symtetric key")

	ClientOpt := flag.Bool("client", false, "Server mode")
        var Exporters exporters
	flag.Var(&Exporters, "exporter", "Endpoint in which the prometheus exporter exposes its metrics")
	Target := flag.String("target", GetEnvStr("TARGET_URL", "http://domain.com"), "Target domain to send the collected metrics")
	Interval := flag.String("interval", GetEnvStr("SCRAPE_INTERVAL", "10m"), "Interval between scrapes")
	CompressOpt := flag.Bool("compress", false, "Compress the payload")
        OneshotOpt := flag.Bool("oneshot", false, "Send the data once and then exit")
	flag.Parse()

	if *ServerOpt == true && *ClientOpt == true {
		log.Fatal("PromRelay has to work on either --client or --server mode")
	} else if *ServerOpt {
		var server RelayServer
		server.Reuse = !*NoReuseOpt
		server.ExposedPort = *ExposedPort
		server.ListenPort = *ListenPort
		server.Verbose = true
		server.Output = *OutputFile
		server.HMACKey = *HMACKey
		server.Run()
	} else if *ClientOpt {
		var err error
		var client RelayClient
		client.Exporters = Exporters
		client.Compress = *CompressOpt
		client.Oneshot = *OneshotOpt
		client.HMACKey = *HMACKey
		client.AddTarget(*Target)
		client.Interval, err = time.ParseDuration(*Interval)
		if err != nil {
			log.Fatalf("Error parsing interval: %s\n", err)
		}
		client.Run()
	} else {
		log.Fatal("PromRelay has to work on either --client or --server mode")
	}
}
