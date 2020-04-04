package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
        "os"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)


type RelayClient struct {
	Interval     time.Duration
	Target       *url.URL
        Hostname     string
	HMACKey      string
	HMACHash     string
	Exporters    exporters
	httpClient   http.Client
	Compress     bool
        Oneshot      bool
	Data         []byte
}

// CleanData  removes the comments and the go_ related metrics to just expose the relevant metrics
func (c *RelayClient) ShrinkPayload() error {
	cleanedData := ""
	slicedData := strings.Split(string(c.Data), "\n")
	for _, line := range slicedData {
		if len(line) > 3 {
			if line[0:3] != "go_" && line[0] != '#' {
				cleanedData = cleanedData + c.Hostname + "_" + line + "\n"
			}
		}
	}
	if c.Compress {
		var gzipData bytes.Buffer
		g := gzip.NewWriter(&gzipData)
		if _, err := g.Write([]byte(cleanedData)); err != nil {
			return err
		}
		if err := g.Close(); err != nil {
			return err
		}
		c.Data = gzipData.Bytes()
	} else {
		c.Data = []byte(cleanedData)
	}
	if c.HMACKey != "" {
		h := hmac.New(sha256.New, []byte(c.HMACKey))
		h.Write(c.Data)
		c.HMACHash = base64.StdEncoding.EncodeToString(h.Sum(nil))
	}
	return nil
}

// ScrapeExporter scrapes the prometheus exporter
func (c *RelayClient) ScrapeExporter(ExporterAddr string) ([]byte, error) {
	req, err := http.NewRequest("GET", ExporterAddr, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Got %d bytes from exporter %s", len(bodyBytes), ExporterAddr)
	return bodyBytes, nil
}

// ForwardData sends the data to the remote prometheus server
func (c *RelayClient) ForwardData() error {
	qParams := url.Values{}
	if c.HMACHash != "" {
		qParams.Add("mac",string(c.HMACHash))
		c.Target.RawQuery = qParams.Encode()
	}
	req, err := http.NewRequest("POST", c.Target.String(), bytes.NewBuffer([]byte(c.Data)))
	if err != nil {
		return err
	}
	if c.Compress {
		req.Header.Set("Content-Encoding", "gzip")
	}
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(c.Data)))
	req.Header.Set("Accept", "*/*")
        req.Header.Set("Content-Type", "text/html; charset=utf-8")
	req.Header.Set("User-Agent", "PromRelay Ver 1.0")
	req.Header.Set("Host", c.Target.Hostname())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server did not answered OK, (%d)", resp.StatusCode)
	}
	log.Printf("%d %d bytes sent to %s", resp.StatusCode, len(c.Data), c.Target.Hostname())
	return nil
}

// AddTarget parses the url string into a URL object and add it to the client Structure
func (c *RelayClient) AddTarget(target string) {
	var err error
	c.Target, err = url.Parse(target)
	if err != nil {
		log.Printf("Error parsing target endpoint: %s\n", err)
	}
}

// Run runs the client
func (c *RelayClient) Run() {
	c.httpClient = http.Client{Timeout: time.Second * 5}
        c.HMACHash = ""
	if c.Hostname == "" {
	name, err := os.Hostname()
    	if err != nil {
        	log.Printf("Error parsing target endpoint: %s\n", err)
		c.Hostname = "unknown"
    	}
    	c.Hostname = name
	}
	for {
		c.Data = nil
                for _, exporter := range(c.Exporters){
			payload, err := c.ScrapeExporter(exporter)
			if err != nil || len(payload) == 0 {
				log.Printf("%s\n", err)
			}else{
				c.Data = append(c.Data, payload...)
			}
		}
		if len(c.Data) != 0{
			c.ShrinkPayload()
			err := c.ForwardData()
			if err != nil {
				log.Printf("%s\n", err)
			}
		}
		if c.Oneshot {
			return
		}
		time.Sleep(c.Interval)
	}
}

