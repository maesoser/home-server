package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
    
    "github.com/blog/pkg/easycert"
    "github.com/blog/pkg/blog"
	"github.com/gorilla/mux"
)

func Run(addr, sslAddr, cert, key, name string) {
	var wg sync.WaitGroup
	wg.Add(2)
	// Starting HTTP server
	go func() {
		defer wg.Done()
		log.Printf("Starting HTTP service on %s ...", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Println(err)
		}

	}()

	// Starting HTTPS server
	go func() {
		defer wg.Done()
		log.Printf("Starting HTTPS service on %s ...", sslAddr)
		srv := &http.Server{
			Addr:         sslAddr,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
			TLSConfig: &tls.Config{
				PreferServerCipherSuites: true,
				CurvePreferences: []tls.CurveID{
					tls.CurveP256,
					tls.X25519,
				},
			},
		}
		if err := srv.ListenAndServeTLS(cert, key); err != nil {
			log.Println(err)
			log.Println("Generating self signed certificate")
			bundle, err := easycert.GenerateGeneric()
			if err != nil {
				panic(err)
			}
			srv.TLSConfig.Certificates = []tls.Certificate{bundle}
			srv.TLSConfig.ServerName = name
			if err := srv.ListenAndServeTLS("", ""); err != nil {
				log.Println(err)
			}
		}
	}()
	wg.Wait()
}

func GetEnvStr(name, value string) string {
	if os.Getenv(name) != "" {
		return os.Getenv(name)
	}
	return value
}

func main() {

	Port := flag.String("port", GetEnvStr("SERVER_HTTP_PORT", "80"), "Server Port")
	TLSPort := flag.String("tls_port", GetEnvStr("SERVER_HTTPS_PORT", "443"), "Secure Server Port")
	KeyArg := flag.String("key", GetEnvStr("SERVER_KEY", "data/key.pem"), "Server Key")
	CertArg := flag.String("cert", GetEnvStr("SERVER_CERT", "data/certificate.pem"), "Server Certificate")
	Config := flag.String("config", GetEnvStr("CONFIG_FILE", "config.json"), "Configuration file")
	flag.Parse()

	if *Port == "80" || *TLSPort == "443" {
		log.Println("Maybe you would need to set: <setcap 'cap_net_bind_service=+ep' echo>")
	}

	b := blog.Blog{}
	b.Load(*Config)
	b.Compile()
    go b.Watch()

	router := mux.NewRouter()
	router.HandleFunc("/main/{page}", b.ServeMain)
    router.HandleFunc("/tag/{tag}", b.ServeTag)
    router.HandleFunc("/classless.css", b.ServeStyle)
	router.PathPrefix("/posts/").Handler(http.StripPrefix("/posts", http.FileServer(http.Dir(b.Path + "/"))))
	router.NotFoundHandler = http.HandlerFunc(b.ServeMain)
	http.Handle("/", router)
	Run(":"+*Port, ":"+*TLSPort, *CertArg, *KeyArg, b.Domain)

}
