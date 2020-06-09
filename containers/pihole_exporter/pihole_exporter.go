package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type GravityStatus struct {
	FileExist  bool        `json:"file_exists"`
	AbsoluteTS uint64      `json:"absolute"`
	RelativeTS interface{} `json:"relative"`
}
type PiHoleStatus struct {
	BlockedDomains      uint64        `json:"domains_being_blocked"`
	TotalQueriesToday   uint64        `json:"dns_queries_today"`
	BlockedQueriesToday uint64        `json:"ads_blocked_today"`
	BlockedPcntToday    float64       `json:"ads_percentage_today"`
	UniqueDomains       uint64        `json:"unique_domains"`
	ForwardedQueries    uint64        `json:"queries_forwarded"`
	CachedQueries       uint64        `json:"queries_cached"`
	ClientsEverSeen     uint64        `json:"clients_ever_seen"`
	UniqueClients       uint64        `json:"unique_clients"`
	TotalQueries        uint64        `json:"dns_queries_all_types"`
	NODATAReplies       uint64        `json:"reply_NODATA"`
	NXDOMAINReplies     uint64        `json:"reply_NXDOMAIN"`
	CNAMEReplies        uint64        `json:"reply_CNAME"`
	IPReplies           uint64        `json:"reply_IP"`
	PrivacyLevel        uint64        `json:"privacy_level"`
	Status              string        `json:"status"`
	GravityStatus       GravityStatus `json:"gravity_last_updated"`
}

type PiHoleExporter struct {
	status   PiHoleStatus
	Target   *string
	Listener *string
}

func (e *PiHoleExporter) getMetrics() error {
	var httpClient = &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
	req, err := http.NewRequest("GET", "http://"+*e.Target+"/admin/api.php?summaryRaw", nil)
	if err != nil {
		return err
	}
	response, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(body), &e.status)
	return nil
}

func (e *PiHoleExporter) promHandler(w http.ResponseWriter, r *http.Request) {
	err := e.getMetrics()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Fprintf(w, "pihole_blocked_domains %d\n", e.status.BlockedDomains)
	fmt.Fprintf(w, "pihole_dns_queries_today %d\n", e.status.TotalQueriesToday)
	fmt.Fprintf(w, "pihole_ads_blocked_today %d\n", e.status.BlockedQueriesToday)
	fmt.Fprintf(w, "pihole_ads_percentage_today %f\n", e.status.BlockedPcntToday)
	fmt.Fprintf(w, "pihole_unique_domains %d\n", e.status.UniqueDomains)
	fmt.Fprintf(w, "pihole_queries_forwarded %d\n", e.status.ForwardedQueries)
	fmt.Fprintf(w, "pihole_queries_cached %d\n", e.status.CachedQueries)
	fmt.Fprintf(w, "pihole_clients_ever_seen %d\n", e.status.ClientsEverSeen)
	fmt.Fprintf(w, "pihole_unique_clients %d\n", e.status.UniqueClients)
	fmt.Fprintf(w, "pihole_dns_queries_all_types %d\n", e.status.TotalQueries)
	fmt.Fprintf(w, "pihole_reply_nodata %d\n", e.status.NODATAReplies)
	fmt.Fprintf(w, "pihole_reply_nxdomain %d\n", e.status.NXDOMAINReplies)
	fmt.Fprintf(w, "pihole_reply_cname %d\n", e.status.CNAMEReplies)
	fmt.Fprintf(w, "pihole_reply_ip %d\n", e.status.IPReplies)
	fmt.Fprintf(w, "pihole_privacy_level %d\n", e.status.PrivacyLevel)
	fmt.Fprintf(w, "pihole_gravity_last_updated %d\n", e.status.GravityStatus.AbsoluteTS)

}

func GetEnvStr(name, value string) string {
	if os.Getenv(name) != "" {
		return os.Getenv(name)
	}
	return value
}

func main() {
	exporter := &PiHoleExporter{}
	exporter.Target = flag.String("a", GetEnvStr("PIHOLE_ADDR", "127.0.0.1:8080"), "PiHole Address")
	exporter.Listener = flag.String("p", GetEnvStr("PROMETHEUS_PORT", "9333"), "Exporter address to listen to")
	flag.Parse()
	log.Printf("Listening at %s\n", *exporter.Listener)
	http.HandleFunc("/metrics", exporter.promHandler)
	log.Fatal(http.ListenAndServe(":"+*exporter.Listener, nil))
}
