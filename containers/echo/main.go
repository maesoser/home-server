package main

import (
    "bytes"
    "crypto/rand"
    "crypto/rsa"
    "crypto/tls"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/json"
    "encoding/pem"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "math/big"
    "net"
    "net/http"
    "os"
    "strconv"
    "sync"
    "time"

    "github.com/gorilla/mux"
)

type requestInfo struct {
    Src        string            `json:"src"`
    Dst        string            `json:"dst"`
    Method     string            `json:"method"`
    Proto      string            `json:"protocol"`
    Path       string            `json:"path"`
    Size       int               `json:"size"`
    TLSVersion string            `json:"tls_ver"`
    TLSCipher  string            `json:"tls_cipher"`
    TLSServer  string            `json:"tls_servername"`
    Headers    map[string]string `json:"headers"`
}

func GetEnvStr(name, value string) string {
    if os.Getenv(name) != "" {
        return os.Getenv(name)
    }
    return value
}

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
        if err := srv.ListenAndServeTLS(cert ,key); err != nil {
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

func main() {

    Port := flag.String("port", GetEnvStr("SERVER_HTTP_PORT", "80"), "Server Port")
    TLSPort := flag.String("tls_port", GetEnvStr("SERVER_HTTPS_PORT", "443"), "Secure Server Port")
    KeyArg := flag.String("key", GetEnvStr("SERVER_KEY", "data/key.pem"), "Server Key")
    CertArg := flag.String("cert", GetEnvStr("SERVER_CERT", "data/certificate.pem"), "Server Certificate")
    DomainName := flag.String("domain", GetEnvStr("SNI_NAME", "localhost"), "Server Name")
    //Verbose := flag.Bool("save", false, "Save NSG data to json file.")
    flag.Parse()

    if *Port == "80" || *TLSPort == "443" {
        log.Println("Maybe you would need to set: <setcap 'cap_net_bind_service=+ep' echo>")
    }

    router := mux.NewRouter()
    router.HandleFunc("/ping", MainWebpage)
    router.HandleFunc("/pong", PongResponse)
    router.HandleFunc("/wait/{seconds}", Wait)
    router.HandleFunc("/code/{code}", ReturnCode)
    router.NotFoundHandler = http.HandlerFunc(ReturnHeaders)
    http.Handle("/", router)
    Run(":"+*Port, ":"+*TLSPort, *CertArg, *KeyArg, *DomainName)
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
    info := requestInfo{
        Src:    r.RemoteAddr,
        Dst:    r.Host,
        Size:   GetBodySize(r),
        Method: r.Method,
        Proto:  r.Proto,
        Path:   r.URL.Path,
    }
    info.Headers = make(map[string]string)
    for name, value := range r.Header {
        if len(value) == 1 {
            info.Headers[name] = fmt.Sprintf("%s", value[0])
        } else {
            info.Headers[name] = fmt.Sprintf("%v", value)
        }
    }
    if r.TLS != nil {
        info.TLSVersion = TLSVersionName(r.TLS.Version)
        info.TLSCipher = tls.CipherSuiteName(r.TLS.CipherSuite)
        info.TLSServer = r.TLS.ServerName
    }
    e, err := json.MarshalIndent(info, "", "  ")
    if err != nil {
        fmt.Println(err)
    }
    return string(e)
}

func MainWebpage(w http.ResponseWriter, r *http.Request) {
    log.Println("Trying to serve main.html")
    data, err := ioutil.ReadFile("/app/main.html")
    if err != nil {
        log.Println(err)
    }
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, string(data))
}

func ReturnHeaders(w http.ResponseWriter, r *http.Request) {
    output := GetInfo(r)
    log.Println(output)
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, output)
}

func ReturnCode(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    code, err := strconv.Atoi(vars["code"])
    if err != nil {
        log.Println(err)
    } else {
        w.WriteHeader(code)
    }
    log.Printf("Sending code %s to %s (%v)\n", vars["code"], r.RemoteAddr, r.Header.Get("User-Agent"))
    fmt.Fprintf(w, vars["code"])
}

func Wait(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    seconds, err := strconv.Atoi(vars["seconds"])
    log.Printf("Waiting %s seconds before answering to %s (%v)\n", vars["seconds"], r.RemoteAddr, r.Header.Get("User-Agent"))
    if err != nil {
        log.Println(err)
    } else {
        time.Sleep(time.Duration(seconds) * time.Second)
    }
    fmt.Fprintf(w, vars["seconds"])
}

func PongResponse(w http.ResponseWriter, r *http.Request) {
    log.Printf("Ping from %s (%v)\n", r.RemoteAddr, r.Header.Get("User-Agent"))
    fmt.Fprintf(w, "pong")
}

