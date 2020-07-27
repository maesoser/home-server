package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/oschwald/geoip2-golang"
)

func GetEnvStr(name, value string) string {
	if os.Getenv(name) != "" {
		return os.Getenv(name)
	}
	return value
}

type Response struct {
	ASNumber              uint    `json:"as_num"`
	ASOrg                 string  `json:"as_org"`
	City                  string  `json:"city"`
	CountryName           string  `json:"country_name"`
	CountryISO            string  `json:"country_iso"`
	RegisteredCountryName string  `json:"reg_country_name"`
	RegisteredCountryISO  string  `json:"reg_country_iso"`
	Continent             string  `json:"continent"`
	IsEU                  bool    `json:"is_eu"`
	AccuracyRadius        uint16  `json:"accuracy"`
	Latitude              float64 `json:"lat"`
	Longitude             float64 `json:"lon"`
	TimeZone              string  `json:"time_zone"`
	Postal                string  `json:"postal_code"`
}

type DBAPI struct {
	Geo *geoip2.Reader
	AS  *geoip2.Reader
}

func (db *DBAPI) Lookup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["addr"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "no_ip")
		return
	}
	ip := net.ParseIP(vars["addr"])
	if ip == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "bad_ip")
		return
	}
	record, err := db.Geo.City(ip)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "geo_not_found")
		return
	}
	response := Response{
		City:                  record.City.Names["en"],
		CountryName:           record.Country.Names["en"],
		CountryISO:            record.Country.IsoCode,
		IsEU:                  record.Country.IsInEuropeanUnion,
		RegisteredCountryName: record.RegisteredCountry.Names["en"],
		RegisteredCountryISO:  record.RegisteredCountry.IsoCode,
		Continent:             record.Continent.Names["en"],
		AccuracyRadius:        record.Location.AccuracyRadius,
		Latitude:              record.Location.Latitude,
		Longitude:             record.Location.Longitude,
		Postal:                record.Postal.Code,
		TimeZone:              record.Location.TimeZone,
	}
	asrec, err := db.AS.ASN(ip)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "as_not_found")
		return
	}
	response.ASNumber = asrec.AutonomousSystemNumber
	response.ASOrg = asrec.AutonomousSystemOrganization
	JSONResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "marshmalling_err")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(JSONResponse))
}

func (db *DBAPI) Load(geoname, asname string) {
	var err error
	db.Geo, err = geoip2.Open(geoname)
	if err != nil {
		log.Fatal(err)
	}
	db.AS, err = geoip2.Open(asname)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	GeoDBName := flag.String("citydb", GetEnvStr("CITY_DB", "/data/GeoLite2-City.mmdb"), "City Database")
	ASDBName := flag.String("asdb", GetEnvStr("AS_DB", "/data/GeoLite2-ASN.mmdb"), "AS Database")
	Port := flag.String("port", GetEnvStr("HTTTP_PORT", "8080"), "HTTP Port")
	flag.Parse()

	db := DBAPI{}
	db.Load(*GeoDBName, *ASDBName)

	router := mux.NewRouter()
	router.HandleFunc("/addr/{addr}", db.Lookup)
	http.Handle("/", router)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("SIGTERM Received, closing DBs")
		db.AS.Close()
		db.Geo.Close()
		os.Exit(1)
	}()

	log.Printf("Starting HTTP service on %s ...", *Port)
	if err := http.ListenAndServe(":"+*Port, nil); err != nil {
		log.Println(err)
	}
}
