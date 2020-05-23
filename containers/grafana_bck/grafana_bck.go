package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	VERSION = "0.1"
)

func CheckError(err error, fatal bool) {
	if err != nil {
		log.Fatal(err)
		if fatal {
			os.Exit(-1)
		}
	}
}

type weatherData struct {
	New         bool
	Timestamp   time.Time
	Pressure    float64
	Humidity    float64
	Temperature float64
	Light       float64
	DewPoint    float64
}

type weatherHandler struct {
	Data    weatherData
	conn    *net.UDPConn
	logfile *os.File
}

func (h *weatherHandler) Init(listenPort, logfile string) {
	ServerAddr, err := net.ResolveUDPAddr("udp", ":"+listenPort)
	CheckError(err, true)

	h.conn, err = net.ListenUDP("udp", ServerAddr)
	CheckError(err, true)

	h.logfile, err = os.Create(logfile)
	CheckError(err, true)
}

func (h *weatherHandler) updateData(buffer []byte) {
	h.Data.Humidity = float64(binary.LittleEndian.Uint32(buffer[0:4])) / 100.0
	h.Data.Temperature = float64(binary.LittleEndian.Uint32(buffer[4:8])) / 100.0
	h.Data.Pressure = float64(binary.LittleEndian.Uint32(buffer[8:12])) / 100.0
	h.Data.Light = (float64(binary.LittleEndian.Uint32(buffer[12:16])) / 1024.0) * 100.0
	h.Data.Timestamp = time.Now()
	h.Data.DewPoint = h.Data.Temperature - ((100.0 - h.Data.Humidity) / 5.0)
	h.Data.New = true
}

func (h *weatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Data.New == true {
		fmt.Fprintf(w, "weather_pressure %f\n", h.Data.Pressure)
		fmt.Fprintf(w, "weather_temperature %f\n", h.Data.Temperature)
		fmt.Fprintf(w, "weather_humidity %f\n", h.Data.Humidity)
		fmt.Fprintf(w, "weather_light %f\n", h.Data.Light)
		fmt.Fprintf(w, "weather_dp %f\n", h.Data.DewPoint)
		fmt.Fprintf(w, "weather_refresh_ofset_ns %d\n", time.Now().Sub(h.Data.Timestamp))
		h.Data.New = false
	}else{
                log.Println("No new data to deliver to prometheus.")
	}
}

func (h *weatherHandler) ListenForUpdates() {
	buffer := make([]byte, 16)
	for {
		n, _, err := h.conn.ReadFromUDP(buffer)
		CheckError(err, true)
		if n != 0 {
			// log.Printf("Received %d bytes from %v\n", n, peer)
			h.updateData(buffer)

			line := fmt.Sprintf("%v\n", h.Data)
			log.Println(line)
			_, err = h.logfile.WriteString(line)
			CheckError(err, false)
		}
	}
}

func main() {
	PromPort := flag.String("port", "9200", "Prometheus Port")
	ListenPort := flag.String("listen", "4865", "UDP Port")
	LogFile := flag.String("log", "/tmp/history.db", "Log file")
	flag.Parse()

	handler := &weatherHandler{}
	handler.Init(*ListenPort, *LogFile)
	go handler.ListenForUpdates()
	log.Printf("Serving exporter at port %s, listening for weather station at port %s\n", *PromPort, *ListenPort)
	err := http.ListenAndServe(":"+*PromPort, handler)
	log.Panic(err)

}

