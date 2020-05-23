package main

import (
	"fmt"
	"io/ioutil"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
	"path/filepath"
        "math/rand"
	"time"
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

func ListFiles(root string) []string{
    var files []string
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        files = append(files, path)
        return nil
    })
    if err != nil {
        panic(err)
    }
    return files
}

func getFileContent(filename string) string{
    dat, err := ioutil.ReadFile(filename)
    if err != nil{
        log.Println(err)
        return ""
    }
    return string(dat)
}
func HelloServer(w http.ResponseWriter, r *http.Request){

        seed := rand.NewSource(time.Now().UnixNano())
        rg := rand.New(seed)
	files := ListFiles("files")
        file := files[rg.Intn(len(files))]
	fmt.Fprintln(w, "<html><head><title>testpage</title></head>")
	fmt.Fprintln(w, "<body><h1>Static Text</h1><p>")
	fmt.Fprintln(w, getFileContent(files[0]))
	fmt.Fprintln(w, "</p><h1>Dynamic Text</h1><p>")
	fmt.Fprintln(w, getFileContent(file))
	fmt.Fprintln(w, "</p></body></html>")
}

