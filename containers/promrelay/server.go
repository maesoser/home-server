package main

import (
	"encoding/base64"
	"net/url"
	"crypto/hmac"
	"crypto/sha256"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type PromData struct {
	Payload   string
	Timestamp time.Time
	Used      bool
	Reuse     bool
}

func (d *PromData) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//log.Printf("Received request from %s\n", r.RemoteAddr)
	if len(d.Payload) == 0 {
		//log.Println("Payload is empty.")
		return
	}
	if d.Used == true && d.Reuse == false {
		//log.Println("Data has been already sent to prometheus.")
		return
	}
	fmt.Fprintf(w, "%s\n", d.Payload)
	d.Used = true
}

func (d *PromData) Save(filename string) error{
    err := ioutil.WriteFile(filename, []byte(d.Payload), 0644)
    return err
}

// RelayServer defines the structure used to Serve
type RelayServer struct {
	Reuse       bool
	ExposedPort int
	ListenPort  int
	Verbose     bool
        Output      string
	HMACKey     string
	Data        *PromData
}

func (s *RelayServer) Run() error {
	if s.ListenPort == 0 {
		s.ListenPort = 8080
	}
	s.Data = &PromData{}
	go s.ListenForUpdates()
	log.Printf("Serving exporter at port %d, listening for updates at port %d\n", s.ExposedPort, s.ListenPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.ExposedPort), s.Data)
	return err
}

func (s *RelayServer) ListenForUpdates() {
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.ListenPort), s)
	if err != nil {
		log.Fatalf("Unable to listen at port %d: %s\n", s.ListenPort, err)
	}
}

func (s *RelayServer) HMACVerified(URL *url.URL, body []byte) error{
    keys, ok := URL.Query()["mac"]
    if !ok || len(keys[0]) < 1 {
        return fmt.Errorf("No HMAC received.")
    }
    messageMAC, err := base64.StdEncoding.DecodeString(keys[0])
    if err != nil{
	return err
    }
    mac := hmac.New(sha256.New, []byte(s.HMACKey))
    mac.Write(body)
    expectedMAC := mac.Sum(nil)

    if hmac.Equal(messageMAC, expectedMAC) == false{
        return fmt.Errorf("Expected HMAC does not coincide with received HMAC")
    }
    return nil
}

func (s *RelayServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error parsing body: %s\n", err)
		return
	}
	if s.HMACKey != ""{
		err := s.HMACVerified(r.URL, bodyBytes)
		if err != nil{
			log.Printf("Error verifying packet: %s\n", err)
			return
		}
	}
	if r.Header.Get("Content-Encoding") == "gzip" && len(bodyBytes) != 0{
		gr, err := gzip.NewReader(bytes.NewBuffer(bodyBytes))
		defer gr.Close()
		data, err := ioutil.ReadAll(gr)
		if err != nil {
			log.Printf("Error uncompressing update: %s\n", err)
			return
		} else {
			compressPcnt := 100.0 - 100.0 * float64(len(bodyBytes))/float64(len(data))
			log.Printf("Received %d bytes (%.1f%% compression) update from %s\n", len(bodyBytes), compressPcnt, r.RemoteAddr)
			s.Data.Payload = string(data)
		}
	}else{
                log.Printf("Received %d bytes update from %s\n", len(bodyBytes), r.RemoteAddr)
		s.Data.Payload = string(bodyBytes)
	}
	s.Data.Used = false
	s.Data.Timestamp = time.Now()
	s.Data.Reuse = s.Reuse
        if s.Output != ""{
		err := s.Data.Save(s.Output)
		if err != nil {
			log.Println(err)
		}
	}
}

