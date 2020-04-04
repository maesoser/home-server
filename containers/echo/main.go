package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
        "io/ioutil"
)

func GetEnvStr(name, value string) string {
	if os.Getenv(name) != "" {
		return os.Getenv(name)
	}
	return value
}

func Run(addr string, sslAddr string, ssl map[string]string) {

	var wg sync.WaitGroup
	wg.Add(2)
	// Starting HTTP server
	go func() {
		defer wg.Done()
		log.Printf("Staring HTTP service on %s ...", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Println(err)
		}

	}()

	// Starting HTTPS server
	go func() {
		defer wg.Done()
		log.Printf("Staring HTTPS service on %s ...", sslAddr)
		if err := http.ListenAndServeTLS(sslAddr, ssl["cert"], ssl["key"], nil); err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()

}

func main() {

	Port := flag.String("port", GetEnvStr("SERVER_HTTP_PORT", "8080"), "Server Port")
	TLSPort := flag.String("tls_port", GetEnvStr("SERVER__HTTPS_PORT", "8443"), "Secure Server Port")
	KeyArg := flag.String("key", GetEnvStr("SERVER_KEY", "data/key.pem"), "Server Key")
	CertArg := flag.String("cert", GetEnvStr("SERVER_CERT", "data/certificate.pem"), "Server Certificate")
	//Verbose := flag.Bool("save", false, "Save NSG data to json file.")
	flag.Parse()

	log.Printf("Using:\n\tKey: %s\n\tCert: %s\n", *KeyArg, *CertArg)
	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/wait", Wait)
	http.HandleFunc("/code", ReturnCode)

	Run(":"+*Port, ":"+*TLSPort, map[string]string{
		"cert": *CertArg,
		"key":  *KeyArg,
	})
}

func TLSVersionName(version uint16) string {
	switch version {
	case 0x0300:
		return "SSLv3.0"
	case 0x0301:
		return "TLSv1.0"
	case 0x0302:
		return "TLSv1.1"
	case 0x0303:
		return "TLSv1.2"
	case 0x0304:
		return "TLSv1.3"
	}
	return fmt.Sprintf("%d", version)
}

func GetBodySize(r *http.Request) int {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return -1
    }
    return len(body)
}

func GetInfo(r *http.Request) string {
	output := ""
	output += fmt.Sprintf("[%s] %s %s\n", r.Method, r.Proto, r.URL.Path)
	output += fmt.Sprintf("\tConnection: %s  to  %s\n", r.RemoteAddr, r.Host)
	output += fmt.Sprintf("\tBody size: %d\n", GetBodySize(r))
        output += fmt.Sprintf("\tHeaders:\n")
	for k, v := range r.Header {
		if k == "Cookie" {
			for _, c := range v {
				output += fmt.Sprintf("\t\tCookie: %s\n", c)
			}
		} else {
			output += fmt.Sprintf("\t\t%s: %s\n", k, strings.Join(v, ","))
		}
	}
	if r.TLS != nil {
		output += fmt.Sprintf("\tTLS:\n\t\tVersion: %s\n", TLSVersionName(r.TLS.Version))
		output += fmt.Sprintf("\t\tCipherSuite: %s\n", tls.CipherSuiteName(r.TLS.CipherSuite))
	}
	return output
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	output := GetInfo(r)
	log.Println(output)
	fmt.Fprintf(w, output)
}

func ReturnCode(w http.ResponseWriter, r *http.Request) {
	output := GetInfo(r)
	log.Println(output)
	fmt.Fprintf(w, output)
}

func Wait(w http.ResponseWriter, r *http.Request) {
	output := GetInfo(r)
	log.Println(output)
	for {
		time.Sleep(1 * time.Second)
	}
}
