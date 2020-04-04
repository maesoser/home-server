package main

import (
	"flag"
	"github.com/maesoser/glan/pkg/nmap"
	"log"
	"os"
	"strings"
)

var Scans []nmap.NMAPScan

func GetEnvStr(name, value string) string {
	if os.Getenv(name) != "" {
		return os.Getenv(name)
	}
	return value
}

func main() {

	PortVar := flag.String("port", GetEnvStr("GLAN_PORT", "8080"), "Web Interface Port")
	TargetsVar := flag.String("targets", GetEnvStr("GLAN_TARGETS", ""), "List of targets to scan from the beginning.")
	flag.Parse()
	log.Printf("Listening on %s\n", *PortVar)

	if *TargetsVar != "" {
		Scan := nmap.NewScan(strings.Split(*TargetsVar, ";"))
		Scan.Run("")
		log.Println("Parsing")
		log.Println(Scan.Parse())
	}

}
