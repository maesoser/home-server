package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

func rndSerial() (*big.Int, error) {
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return n, err
	}
	return n, nil
}

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
			bundle, err := certsetup()
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

func certsetup() (tls.Certificate, error) {
	data := pkix.Name{
		Organization:  []string{"Company, Inc."},
		Country:       []string{"US"},
		Province:      []string{""},
		Locality:      []string{"San Francisco"},
		StreetAddress: []string{"Golden Gate Bridge"},
		PostalCode:    []string{"94016"},
	}
	// set up our CA certificate
	serial, _ := rndSerial()
	ca := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               data,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, 7),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return tls.Certificate{}, err
	}
	serial, _ = rndSerial()
	cert := &x509.Certificate{
		SerialNumber: serial,
		Subject:      data,
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, 7),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return tls.Certificate{}, err
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return tls.Certificate{}, err
	}
	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey)})

	bundle, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	return bundle, nil
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

	blog := Blog{}
	blog.Load(*Config)
	blog.Compile()

	router := mux.NewRouter()
	router.HandleFunc("/main/{page}", blog.ServeMain)
	fs := http.FileServer(http.Dir(blog.Path + "/"))
	router.PathPrefix("/posts/").Handler(http.StripPrefix("/posts", fs))
	router.NotFoundHandler = http.HandlerFunc(blog.ServeMain)
	http.Handle("/", router)
	Run(":"+*Port, ":"+*TLSPort, *CertArg, *KeyArg, blog.Domain)

}
