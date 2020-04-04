package main

import (
	"io/ioutil"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
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
	Run(":"+*Port, ":"+*TLSPort, map[string]string{
		"cert": *CertArg,
		"key":  *KeyArg,
	})
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["web"]
    	if !ok || len(keys[0]) < 1 {
        	http.Error(w, "can't read query", http.StatusBadRequest)
		return
    	}
	target := keys[0]
	log.Println(target)
	resp, err := http.Get("https://" + target)
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Printf("Error reading body: %v", err)
            http.Error(w, "can't read body", http.StatusBadRequest)
            return
        }
	w.Write(body)
}

