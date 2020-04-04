package main

import (
	"flag"
	"fmt"
	"github.com/maesoser/ping_exporter/pkg/icmping"
	"github.com/maesoser/ping_exporter/pkg/tcping"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func GetEnvStr(name, value string) string {
	if os.Getenv(name) != "" {
		return os.Getenv(name)
	}
	return value
}

func getInt(text string, defval int) int {
	if text == "" {
		return defval
	}
	value, err := strconv.Atoi(text)
	if err != nil {
		return defval
	}
	return value
}

type pingHandler struct{}

func (h pingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m, _ := url.ParseQuery(r.URL.RawQuery)
	address := m.Get("addr")
	if address != "" {
		count := getInt(m.Get("count"), 1)
		ttl := getInt(m.Get("ttl"), 128)
		tos := getInt(m.Get("tos"), 0)
		size := getInt(m.Get("size"), 32)
		port := m.Get("port")
		//log.Printf("Received request to send %d pings to %s\n", count, addressi)
		if port != "" {
			pinger, err := tcping.NewTCPinger(address, port)
			if err != nil {
				log.Printf("Error executing TCP ping to %s: %s", address, err)
			} else {
				pinger.Count = count
				pinger.Run()
				stats := pinger.Statistics()
				fmt.Fprintf(w, "ping_pkt_tx_count{host=\"%s\",addr=\"%s:%s\"} %d\n", pinger.Addr(), pinger.IPAddr(), pinger.Port, stats.PacketsSent)
				fmt.Fprintf(w, "ping_pkt_rx_count{host=\"%s\",addr=\"%s:%s\"} %d\n", pinger.Addr(), pinger.IPAddr(), pinger.Port, stats.PacketsRecv)
				fmt.Fprintf(w, "ping_pkt_loss_pcnt{host=\"%s\",addr=\"%s:%s\"} %f\n", pinger.Addr(), pinger.IPAddr(), pinger.Port, stats.PacketLoss)
				fmt.Fprintf(w, "ping_rtt_min_ns{host=\"%s\",addr=\"%s:%s\"} %d\n", pinger.Addr(), pinger.IPAddr(), pinger.Port, int64(stats.MinRtt))
				fmt.Fprintf(w, "ping_rtt_max_ns{host=\"%s\",addr=\"%s:%s\"} %d\n", pinger.Addr(), pinger.IPAddr(), pinger.Port, int64(stats.MaxRtt))
				fmt.Fprintf(w, "ping_rtt_avg_ns{host=\"%s\",addr=\"%s:%s\"} %d\n", pinger.Addr(), pinger.IPAddr(), pinger.Port, int64(stats.AvgRtt))
				fmt.Fprintf(w, "ping_rtt_mdev_ns{host=\"%s\",addr=\"%s:%s\"} %d\n", pinger.Addr(), pinger.IPAddr(), pinger.Port, int64(stats.StdDevRtt))
			}
		} else {
			pinger, err := icmping.NewPinger(address)
			if err != nil {
				log.Printf("Error executing ping to %s: %s", address, err)
			} else {
				pinger.TTL = ttl
				pinger.TOS = tos
				pinger.Size = size
				pinger.Count = count
				pinger.Run()
				stats := pinger.Statistics()
				fmt.Fprintf(w, "ping_ttl_min{host=\"%s\",addr=\"%s\"} %d\n", pinger.Addr(), pinger.IPAddr(), int64(stats.MaxTTL))
				fmt.Fprintf(w, "ping_ttl_max{host=\"%s\",addr=\"%s\"} %d\n", pinger.Addr(), pinger.IPAddr(), int64(stats.MinTTL))
				fmt.Fprintf(w, "ping_pkt_tx_count{host=\"%s\",addr=\"%s\"} %d\n", pinger.Addr(), pinger.IPAddr(), stats.PacketsSent)
				fmt.Fprintf(w, "ping_pkt_rx_count{host=\"%s\",addr=\"%s\"} %d\n", pinger.Addr(), pinger.IPAddr(), stats.PacketsRecv)
				fmt.Fprintf(w, "ping_pkt_loss_pcnt{host=\"%s\",addr=\"%s\"} %f\n", pinger.Addr(), pinger.IPAddr(), stats.PacketLoss)
				fmt.Fprintf(w, "ping_rtt_min_ns{host=\"%s\",addr=\"%s\"} %d\n", pinger.Addr(), pinger.IPAddr(), int64(stats.MinRtt))
				fmt.Fprintf(w, "ping_rtt_max_ns{host=\"%s\",addr=\"%s\"} %d\n", pinger.Addr(), pinger.IPAddr(), int64(stats.MaxRtt))
				fmt.Fprintf(w, "ping_rtt_avg_ns{host=\"%s\",addr=\"%s\"} %d\n", pinger.Addr(), pinger.IPAddr(), int64(stats.AvgRtt))
				fmt.Fprintf(w, "ping_rtt_mdev_ns{host=\"%s\",addr=\"%s\"} %d\n", pinger.Addr(), pinger.IPAddr(), int64(stats.StdDevRtt))
			}
		}
	}
}

func main() {

	Port := flag.String("port", GetEnvStr("PROMETHEUS_PORT", "8080"), "Port")
	flag.Parse()

	log.Printf("Serving at port %s\n", *Port)
	err := http.ListenAndServe(":"+*Port, pingHandler{})
	log.Fatal(err)

}
