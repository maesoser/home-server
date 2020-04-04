package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	Port := flag.String("port", "80", "Port to serve")
	File := flag.String("file", "index.html", "File to serve")
	flag.Parse()

	log.Printf("Serving file %s at port %s\n", *File, *Port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request from %s\n", r.RemoteAddr)
		data, err := ioutil.ReadFile(*File)
		if err != nil {
			fmt.Fprintf(w,err.Error())
		}else{
			fmt.Fprintf(w, string(data))
		}
	})

	err := http.ListenAndServe(":"+*Port, nil)
	if err != nil {
		log.Panic(err)
	}
}
