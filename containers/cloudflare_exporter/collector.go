package main

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/prometheus/client_golang/prometheus"
)

// Based on https://github.com/prometheus/haproxy_exporter/blob/master/haproxy_exporter.go

// CfCollector is the structure that stores all the information related to the collector
type CfCollector struct {
	apiKey    string
	apiEmail  string
	dataset   []string
	accountID string
	zoneName  string

	api   *cloudflare.API
	zones []cloudflare.Zone
	account cloudflare.Account

	startDate string
	endDate   string

	fetchFailed prometheus.Counter
	fetchDone   prometheus.Counter

	requestBytes         *prometheus.Desc
	requestResponseCodes *prometheus.Desc

	requestCountry *prometheus.Desc
	bytesCountry   *prometheus.Desc
	threatsCountry *prometheus.Desc

	totalBytes     *prometheus.Desc
	cachedBytes    *prometheus.Desc
	encryptedBytes *prometheus.Desc

	totalRequests     *prometheus.Desc
	cachedRequests    *prometheus.Desc
	encryptedRequests *prometheus.Desc

	requestSSLType     *prometheus.Desc
	requestHTTPVersion *prometheus.Desc

	WAFEvents *prometheus.Desc

	requestContentType *prometheus.Desc
	bytesContentType   *prometheus.Desc

	workerCPUTime    *prometheus.Desc
	workerErrors    *prometheus.Desc
	workerRequests    *prometheus.Desc
	workerSubRequests    *prometheus.Desc

	networkBits    *prometheus.Desc
	networkPackets *prometheus.Desc

	mutex sync.Mutex
}

// NewCfCollector returns an initialized Collector.
func NewCfCollector(apiKey, apiMail, AccountID, zoneName, dataset string) *CfCollector {

	c := CfCollector{
		apiKey:    apiKey,
		apiEmail:  apiMail,
		accountID: AccountID,
		zoneName:  zoneName,
		dataset:   strings.Split(dataset, ","),
	}

	err := c.Validate()
	if err != nil {
		log.Fatal(err)
	}
	
    log.Printf("Datasets: %v\n", c.dataset)
    if c.zoneName != ""{
    	log.Printf("Zone: %s\n", c.zoneName)
    }
    
	c.fetchFailed = prometheus.NewCounter(prometheus.CounterOpts{Name: "cloudflare_failed_fetches", Help: "The total number of failed fetches"})
	c.fetchDone = prometheus.NewCounter(prometheus.CounterOpts{Name: "cloudflare_done_fetches", Help: "The total number of done fetches"})

	c.requestBytes = prometheus.NewDesc("cloudflare_bytes_per_cache_status", "The total number of processed bytes, labelled per cache status", []string{"cacheStatus", "method", "contentType", "country", "zoneName"}, nil)
	c.requestResponseCodes = prometheus.NewDesc("cloudflare_requests_per_response_code", "The total number of request, labelled per HTTP response codes", []string{"responseCode", "zoneName"}, nil)

	c.requestCountry = prometheus.NewDesc("cloudflare_requests_per_country", "The total number of request, labeled per Country", []string{"country", "zoneName"}, nil)
	c.bytesCountry = prometheus.NewDesc("cloudflare_bytes_per_country", "The total number of request, labeled per Country", []string{"country", "zoneName"}, nil)
	c.threatsCountry = prometheus.NewDesc("cloudflare_threats_per_country", "The total number of threats, labeled per Country", []string{"country", "zoneName"}, nil)

	c.requestContentType = prometheus.NewDesc("cloudflare_requests_per_content_type", "The total number of request, labeled per content type", []string{"contentType", "zoneName"}, nil)
	c.bytesContentType = prometheus.NewDesc("cloudflare_bytes_per_content_type", "The total number of bytes, labeled per content type", []string{"contentType", "zoneName"}, nil)

	c.requestSSLType = prometheus.NewDesc("cloudflare_requests_per_ssl_version", "The total number of requests, labeled per SSL type", []string{"version", "zoneName"}, nil)
	c.requestHTTPVersion = prometheus.NewDesc("cloudflare_requests_per_http_version", "The total number of requests, labeled per HTTP version", []string{"version", "zoneName"}, nil)

	c.totalBytes = prometheus.NewDesc("cloudflare_total_bytes", "The total number of bytes sent", []string{"zoneName"}, nil)
	c.cachedBytes = prometheus.NewDesc("cloudflare_cached_bytes", "The total number of bytes cached", []string{"zoneName"}, nil)
	c.encryptedBytes = prometheus.NewDesc("cloudflare_encrypted_bytes", "The total number of bytes encrypted", []string{"zoneName"}, nil)

	c.totalRequests = prometheus.NewDesc("cloudflare_total_requests", "The total number of requests served", []string{"zoneName"}, nil)
	c.cachedRequests = prometheus.NewDesc("cloudflare_cached_requests", "The total number of requests cached", []string{"zoneName"}, nil)
	c.encryptedRequests = prometheus.NewDesc("cloudflare_encrypted_requests", "The total number of requests encrypted", []string{"zoneName"}, nil)

	c.WAFEvents = prometheus.NewDesc("cloudflare_waf_events", "Cloudflare WAF Hits", []string{"as", "country", "action", "ruleID", "zoneName"}, nil)

	c.workerCPUTime = prometheus.NewDesc("cloudflare_worker_cputime", "CPU time consumed by worker", []string{"workerName", "accountName", "percentile"}, nil)
	c.workerErrors = prometheus.NewDesc("cloudflare_worker_errors", "Errors trigered by worker", []string{"workerName", "accountName"}, nil)
	c.workerRequests = prometheus.NewDesc("cloudflare_worker_requests", "Requests received by worker", []string{"workerName", "accountName"}, nil)
	c.workerSubRequests = prometheus.NewDesc("cloudflare_worker_subrequests", "Subrequests performed by worker", []string{"workerName", "accountName"}, nil)

	c.networkBits = prometheus.NewDesc("cloudflare_network_bits", "Number of bits, labelled per AttackID", []string{"attackID", "accountName", "attackProtocol", "mitigationType", "country", "destinationPort", "attackType"}, nil)
	c.networkPackets = prometheus.NewDesc("cloudflare_network_packets", "Number of packets, labelled per AttackID", []string{"attackID", "accountName", "attackProtocol", "mitigationType", "country", "destinationPort", "attackType"}, nil)

	return &c
}

// Describe describes all the metrics ever exported by the Cloudflare exporter. It
// implements prometheus.Collector.
func (collector *CfCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.requestBytes
	ch <- collector.requestResponseCodes
	ch <- collector.requestContentType
	ch <- collector.bytesContentType
	ch <- collector.requestCountry
	ch <- collector.bytesCountry
	ch <- collector.threatsCountry

	ch <- collector.requestSSLType
	ch <- collector.requestHTTPVersion

	ch <- collector.totalBytes
	ch <- collector.cachedBytes
	ch <- collector.encryptedBytes

	ch <- collector.totalRequests
	ch <- collector.cachedRequests
	ch <- collector.encryptedRequests

	ch <- collector.WAFEvents

	ch <- collector.networkBits
	ch <- collector.networkPackets

	ch <- collector.workerCPUTime
	ch <- collector.workerErrors
	ch <- collector.workerRequests
	ch <- collector.workerSubRequests

	ch <- collector.fetchDone.Desc()
	ch <- collector.fetchFailed.Desc()
}

// Validate checks the configuration parameters given to the Collector
func (collector *CfCollector) Validate() error {

	if collector.apiKey == "" || collector.apiEmail == "" {
		return errors.New("Must provide both api-key and api-email")
	}
	if len(collector.dataset) == 0 {
		collector.dataset = append(collector.dataset, "http")
	}
	if contains(collector.dataset, "net") && collector.accountID == "" {
		return errors.New("You must provide an accountID when exporting network analytics")
	}
	if contains(collector.dataset, "workers") && collector.accountID == "" {
		return errors.New("You must provide an accountID when exporting worker analytics")
	}
	return nil
}

// Collect fetches the stats from Cloudflare zones and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (collector *CfCollector) Collect(ch chan<- prometheus.Metric) {
	collector.mutex.Lock()
	defer collector.mutex.Unlock()

	var err error
	err = collector.login()
	if err != nil {
		log.Println(err)
		return
	}

	if contains(collector.dataset, "net") {
		collector.fetchDone.Inc()
		err = collector.collectNetwork(ch)
		if err != nil {
			log.Println(err)
			collector.fetchFailed.Inc()
		}
	}

	if contains(collector.dataset, "http") {
		collector.fetchDone.Inc()
		err = collector.collectHTTP(ch)
		if err != nil {
			log.Println(err)
			collector.fetchFailed.Inc()
		}
	}

	if contains(collector.dataset, "waf") {
		collector.fetchDone.Inc()
		err = collector.collectWAF(ch)
		if err != nil {
			log.Println(err)
			collector.fetchFailed.Inc()
		}
	}

	if contains(collector.dataset, "workers") {
		collector.fetchDone.Inc()
		err = collector.collectWorkers(ch)
		if err != nil {
			log.Println(err)
			collector.fetchFailed.Inc()
		}
	}
}

func (collector *CfCollector) login() error {
	collector.startDate = time.Now().Add(time.Duration(-20) * time.Minute).Format(time.RFC3339)
	collector.endDate = time.Now().Add(time.Duration(-5) * time.Minute).Format(time.RFC3339)

	var err error
	collector.api, err = cloudflare.New(collector.apiKey, collector.apiEmail)
	if err != nil {
		return err
	}
	collector.zones, err = collector.api.ListZones()
	if err != nil {
		return err
	}
	if collector.accountID != ""{
		collector.account, _, err = collector.api.Account(collector.accountID)
		if err != nil {
			return err
		}
	}
	return err
}

func (collector *CfCollector) collectHTTP(ch chan<- prometheus.Metric) error {
	for _, zone := range collector.zones {
		if zone.Plan.ZonePlanCommon.Name != "Enterprise Website" {
			continue
		}
		if zone.Name != "" && zone.Name != collector.zoneName {
			continue
		}
		log.Printf("Getting HTTP metrics for %s from %s to %s \n", zone.Name, collector.startDate, collector.endDate)
		resp, err := getCloudflareHTTPMetrics(
			buildHttpGraphQLQuery(collector.startDate, collector.endDate, zone.ID),
			collector.apiEmail,
			collector.apiKey,
		)
		if err == nil {
			for _, node := range resp.Viewer.Zones[0].Caching {
				ch <- prometheus.MustNewConstMetric(collector.requestBytes, prometheus.GaugeValue, float64(node.SumEdgeResponseBytes.EdgeResponseBytes), node.Dimensions.CacheStatus, node.Dimensions.HTTPMethod, node.Dimensions.ContentTypeName, node.Dimensions.CountryName, zone.Name)
			}

			RequestsData := resp.Viewer.Zones[0].Requests[0].RequestsData

			ch <- prometheus.MustNewConstMetric(collector.totalBytes, prometheus.GaugeValue, float64(RequestsData.Bytes), zone.Name)
			ch <- prometheus.MustNewConstMetric(collector.cachedBytes, prometheus.GaugeValue, float64(RequestsData.CachedBytes), zone.Name)
			ch <- prometheus.MustNewConstMetric(collector.encryptedBytes, prometheus.GaugeValue, float64(RequestsData.EncryptedBytes), zone.Name)
			ch <- prometheus.MustNewConstMetric(collector.totalRequests, prometheus.GaugeValue, float64(RequestsData.Requests), zone.Name)
			ch <- prometheus.MustNewConstMetric(collector.cachedRequests, prometheus.GaugeValue, float64(RequestsData.CachedRequests), zone.Name)
			ch <- prometheus.MustNewConstMetric(collector.encryptedRequests, prometheus.GaugeValue, float64(RequestsData.EncryptedRequests), zone.Name)

			for _, node := range RequestsData.ResponseStatusMap {
				ch <- prometheus.MustNewConstMetric(collector.requestResponseCodes, prometheus.GaugeValue, float64(node.Requests), strconv.Itoa(node.EdgeResponseStatus), zone.Name)
			}
			for _, node := range RequestsData.CountryMap {
				ch <- prometheus.MustNewConstMetric(collector.requestCountry, prometheus.GaugeValue, float64(node.Requests), node.CountryName, zone.Name)
				ch <- prometheus.MustNewConstMetric(collector.bytesCountry, prometheus.GaugeValue, float64(node.Bytes), node.CountryName, zone.Name)
				ch <- prometheus.MustNewConstMetric(collector.threatsCountry, prometheus.GaugeValue, float64(node.Threats), node.CountryName, zone.Name)
			}
			for _, node := range RequestsData.ContentTypeMap {
				ch <- prometheus.MustNewConstMetric(collector.requestContentType, prometheus.GaugeValue, float64(node.Requests), node.ContentTypeName, zone.Name)
				ch <- prometheus.MustNewConstMetric(collector.bytesContentType, prometheus.GaugeValue, float64(node.Bytes), node.ContentTypeName, zone.Name)
			}
			for _, node := range RequestsData.ClientSSLMap {
				ch <- prometheus.MustNewConstMetric(collector.requestSSLType, prometheus.GaugeValue, float64(node.Requests), node.ClientSSLProtocol, zone.Name)
			}
			for _, node := range RequestsData.ClientHTTPVersionMap {
				ch <- prometheus.MustNewConstMetric(collector.requestHTTPVersion, prometheus.GaugeValue, float64(node.Requests), node.ClientHTTPProtocol, zone.Name)
			}
		} else {
			log.Println("Fetch failed :", err)
		}
	}
	return nil
}

func (collector *CfCollector) collectWAF(ch chan<- prometheus.Metric) error {
	for _, zone := range collector.zones {
		if zone.Plan.ZonePlanCommon.Name != "Enterprise Website" {
			continue
		}
		if zone.Name != "" && zone.Name != collector.zoneName {
			continue
		}
		log.Printf("Getting WAF metrics for %s from %s to %s \n", zone.Name, collector.startDate, collector.endDate)
		resp, err := getCloudflareWAFMetrics(
			buildWAFGraphQLQuery(collector.startDate, collector.endDate, zone.ID),
			collector.apiEmail,
			collector.apiKey,
		)
		if err == nil {
			for _, node := range resp.Viewer.Zones[0].FwEvents {
				ch <- prometheus.MustNewConstMetric(collector.WAFEvents, prometheus.GaugeValue, float64(node.Count), node.Dimensions.ASName, node.Dimensions.Country, node.Dimensions.Action, node.Dimensions.RuleID, zone.Name)
			}
		} else {
			log.Println("Fetch failed :", err)
		}
	}
	return nil
}

func (collector *CfCollector) collectWorkers(ch chan<- prometheus.Metric) error {

	log.Printf("Getting Worker metrics for %s from %s to %s \n", collector.accountID, collector.startDate, collector.endDate)
	resp, err := getCloudflareWorkerMetrics(
		buildWorkersGraphQLQuery(collector.startDate, collector.endDate, collector.accountID),
		collector.apiEmail,
		collector.apiKey,
	)
	if err != nil {
	  log.Println("Fetch Failed:", err)
	  return err
	}
	for _, node := range resp.WorkersViewer.Accounts[0].Workers {
			ch <- prometheus.MustNewConstMetric(collector.workerCPUTime, prometheus.GaugeValue, float64(node.Quantiles.CpuTimeP50), node.Info.Name, collector.account.Name, "50")
			ch <- prometheus.MustNewConstMetric(collector.workerCPUTime, prometheus.GaugeValue, float64(node.Quantiles.CpuTimeP75), node.Info.Name, collector.account.Name, "75")
			ch <- prometheus.MustNewConstMetric(collector.workerCPUTime, prometheus.GaugeValue, float64(node.Quantiles.CpuTimeP99), node.Info.Name, collector.account.Name, "99")
			ch <- prometheus.MustNewConstMetric(collector.workerCPUTime, prometheus.GaugeValue, float64(node.Quantiles.CpuTimeP999), node.Info.Name, collector.account.Name, "99.9")

			ch <- prometheus.MustNewConstMetric(collector.workerErrors, prometheus.GaugeValue, float64(node.Sum.Errors), node.Info.Name, collector.account.Name)
			ch <- prometheus.MustNewConstMetric(collector.workerRequests, prometheus.GaugeValue, float64(node.Sum.Requests), node.Info.Name, collector.account.Name)
			ch <- prometheus.MustNewConstMetric(collector.workerSubRequests, prometheus.GaugeValue, float64(node.Sum.SubRequests), node.Info.Name, collector.account.Name)
	}
	return nil
}

func (collector *CfCollector) collectNetwork(ch chan<- prometheus.Metric) error {

	resp, err := getCloudflareNetworkMetrics(buildNetworkGraphQLQuery(collector.startDate, collector.endDate, collector.accountID), collector.apiEmail, collector.apiKey)
	if err == nil {
		for _, node := range resp.NetworkViewer.Accounts[0].AttackHistory {

			ch <- prometheus.MustNewConstMetric(
				collector.networkBits,
				prometheus.CounterValue,
				float64(node.Sum.Bits),
				node.NetworkDimensions.AttackID,
				collector.account.Name,
				node.NetworkDimensions.AttackProtocol,
				node.NetworkDimensions.AttackMitigationType,
				node.NetworkDimensions.ColoCountry,
				strconv.Itoa(node.NetworkDimensions.DestinationPort),
				node.NetworkDimensions.AttackType,
			)

			ch <- prometheus.MustNewConstMetric(
				collector.networkPackets,
				prometheus.CounterValue,
				float64(node.Sum.Packets),
				node.NetworkDimensions.AttackID,
				collector.account.Name,
				node.NetworkDimensions.AttackProtocol,
				node.NetworkDimensions.AttackMitigationType,
				node.NetworkDimensions.ColoCountry,
				strconv.Itoa(node.NetworkDimensions.DestinationPort),
				node.NetworkDimensions.AttackType,
			)
		}
	} else {
		log.Println("Fetch failed :", err)
	}
	return nil
}

func contains(elements []string, element string) bool {
	for _, e := range elements {
		if element == e {
			return true
		}
	}
	return false
}
