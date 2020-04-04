package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
        "crypto/tls"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Interval string   `json:"interval"`
	Stocks   []string `json:"stocks"`
}

type StockData struct {
	Date  time.Time `json:"date"`
	Open  float64   `json:"open"`
	High  float64   `json:"high"`
	Low   float64   `json:"low"`
	Close float64   `json:"close"`
	Vol   float64   `json:"vol"`
}

func pf(s string) float64 {
	multipl := 1.0
	s = strings.Replace(s, ",", ".", -1)
	if strings.Contains(s, "K") {
		multipl = 1000.0
		s = strings.Replace(s, "K", "", -1)
	}
	if strings.Contains(s, "M") {
		multipl = 1000000.0
		s = strings.Replace(s, "M", "", -1)
	}
	v, e := strconv.ParseFloat(s, 64)
	if e != nil {
		log.Printf("Error Parsing: %s", e)
	}
	return v * multipl
}

func parseEntries(str string) (StockData, error) {
	var dateRegex = regexp.MustCompile(`(?m)\d{2}.\d{2}.\d{4}`)
	DateStr := dateRegex.FindAllString(str, -1)[0]
	var valuesRegex = regexp.MustCompile(`(?m)\d{1,4},\d{2,4}[KM]*`)
	matches := valuesRegex.FindAllString(str, -1)

	var data StockData
	Date, err := time.Parse("02.01.2006", DateStr)
	if err != nil {
		return data, err
	}
	//log.Printf("%s, %v\n", str, matches)
	if len(matches) < 5 {
		return data, errors.New("Error parsing: " + str)
	}

	data.Date = Date
	data.Close = pf(matches[0])
	data.Open = pf(matches[1])
	data.High = pf(matches[2])
	data.Low = pf(matches[3])
	data.Vol = pf(matches[4])
	return data, nil
}

func insertToDb(dbname string, stockname string, stockdata []StockData) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Println(err)
		return
	}

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS stocks (date TEXT, name TEXT, close REAL, open REAL, high REAL, low REAL, vol REAL, PRIMARY KEY (date, name))")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println(err)
		return
	}

	for _, s := range stockdata {
		stmt, err := db.Prepare("INSERT OR REPLACE INTO stocks(date, name, close, open, high, low, vol) values(?,?,?,?,?,?,?)")
		if err != nil {
			log.Println(err)
		}
		_, err = stmt.Exec(s.Date, stockname, s.Close, s.Open, s.High, s.Low, s.Vol)
		if err != nil {
			log.Println(err)
		}
	}
	db.Close()

}

func parsePage(stockValue string) []StockData {
	tr := &http.Transport{
        	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    	}
    	client := &http.Client{Transport: tr}
	// log.Printf("[%s] Scrapping\n", stockValue)
	req, err := http.NewRequest("GET", "https://m.es.investing.com/equities/"+stockValue+"-historical-data", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1)")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	s := string(body)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	var re = regexp.MustCompile(`(?m)((<td((\w*\s*)+="(\w*\s*)+")*>)\d{2}.\d{2}.\d{4}</td>\s*)((<td((\w*\s*)+="(\w*\s*)+")*>)\d{1,4},\d{1,3}[a-zA-Z]*%*</td>\s*)*`)
	var dataList []StockData
        good_entries := 0
        bad_entries := 0
	for _, match := range re.FindAllString(s, -1) {
		data, err := parseEntries(match)
		if err == nil {
                        good_entries += 1
			dataList = append(dataList, data)
			// log.Printf("[%s] %s\tO/C: %f/%f\tH/L: %f/%f\n", stockValue, data.Date, data.Open, data.Close, data.High, data.Low)
		} else {
                        bad_entries += 1
			// log.Println(err)
		}
	}
        log.Printf("[%s] Loaded %d entries, %d entries failed.\n", stockValue, good_entries, bad_entries)
	return dataList
}

func main() {
	var config Config

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if len(os.Args) != 2 {
		log.Printf("Usage:\n\t%s [config.json]", os.Args[0])
		os.Exit(-1)
	}

	/*
	f, err := os.OpenFile("/tmp/stocks.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	*/

	log.Printf("Loading config file <%s>\n", os.Args[1])
	jsonFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)

	intv, err := time.ParseDuration(config.Interval)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for {
		for _, stock := range config.Stocks {
			data := parsePage(stock)
			insertToDb("stocks.db", stock, data)
			sleeptime := r.Intn(90) + 10
			log.Println("Starting sleep for", sleeptime, "seconds")
			time.Sleep(time.Duration(sleeptime) * time.Second)
		}
		log.Println("Done")
		time.Sleep(intv)
	}
}
