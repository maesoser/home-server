package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"flag"

        "net/http"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
        "github.com/prometheus/client_golang/prometheus/promhttp"

	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	VERSION = "0.1"
)

type config struct {
	Id    string `json:"id"`
	Token string `json:"token"`
	Pin   int    `json:"pin"`
}

func (r config) Recipient() string {
	return r.Id
}

func (r config) Save(path string) {
	if path == ""{
	  path = "/tmp/config.json"
	}
	configString, err := json.Marshal(r)
	CheckError(err)
	confFile, err := os.Create(path)
	CheckError(err)
	defer confFile.Close()
	confFile.WriteString(string(configString))
}

func (r config) Load(path string) (config, error) {
	content, err := ioutil.ReadFile(path)
	CheckError(err)
	var c config
	err = json.Unmarshal(content, &c)
	CheckError(err)
	return c, err
}

func (r config) isValid(id string) bool {
	if len(r.Id) == 0 {
		return false
	} else if id != r.Id {
		return false
	}
	return true
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}

func cleanEvents(events []time.Time, tw float64) []time.Time {
	nt := time.Now()
	deleted := 0
	for i := range events {
		j := i - deleted
		diff := nt.Sub(events[j]).Seconds()
		if diff > tw {
			events = events[:j+copy(events[j:], events[j+1:])]
			deleted++
		}
	}
	return events
}

func checkEvents(events []time.Time, thr int) bool {
	l := len(events)
	if l >= thr {
		return true
	}
	return false
}

func help() {
	fmt.Printf("Usage:\n\t %s [token]", os.Args[0])
	os.Exit(-1)
}

var (
        doorOpened = promauto.NewCounter(prometheus.CounterOpts{
                Name: "teledoor_door_opened",
                Help: "Number of times the door has been opened",
        })
        sensorTriggered = promauto.NewGauge(prometheus.GaugeOpts{
                Name: "teledoor_sensor_triggered",
                Help: "Number of times the sensor has been triggered",
        })

)

func main() {
	logPathPtr := flag.String("log","/tmp/teledoor.log","Log file")
	configPathPtr := flag.String("config","","Configuration file")
	tokenPtr := flag.String("token","","Token")
        flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	logFile, err := os.OpenFile(*logPathPtr, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	var c config
	if *configPathPtr != "" && *tokenPtr == ""{
		c, _ = c.Load(*configPathPtr)
		if c.Token == "" {
			log.Println("Unable to load config, please restart the server")
			help()
		}
	}else if *configPathPtr == "" && *tokenPtr != "" {
		log.Println("Cold start...")
		c.Id = ""
		c.Pin = rand.Intn(999999)
		c.Token = *tokenPtr
	}else{
		log.Printf("Error. You've to provide a token and a config path or a token")
                help()
	}

	log.Println("Bot token is ", c.Token)
	log.Println("ACTIVATION PIN is ", c.Pin)
	c.Save(*configPathPtr)

	events := make([]time.Time, 3)

	ServerAddr, err := net.ResolveUDPAddr("udp", ":4864")
	CheckError(err)

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	bot, err := tb.NewBot(tb.Settings{
		Token:  c.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	CheckError(err)

	histFile, err := os.Create("history.db")
	CheckError(err)
	defer histFile.Close()

	bot.Handle("/start", func(m *tb.Message) {
		bot.Send(m.Sender, "Hello! I'm your door sensor!")
	})

	bot.Handle("/ifconfig", func(m *tb.Message){
	        bot.Send(m.Sender, "Your IP is 0")
	})
	bot.Handle("/hist", func(m *tb.Message) {
		if c.isValid(m.Sender.Recipient()) == true {
			if len(events) > 0 {
				last_event := events[len(events)-1]
				datestr := last_event.Format("02-01-2006 15:04:05")
				log.Printf("Last event at %s", datestr)
				bot.Send(m.Sender, "Last event at:", datestr)
			}
		} else {
			bot.Send(m.Sender, "Sensor is not activated yet. Please, use \n\"/activate [PIN]\" to register the sensor.")
		}
	})

	bot.Handle("/activate", func(m *tb.Message) {
		if c.isValid(m.Sender.Recipient()) {
			bot.Send(m.Sender, "Sensor is already active")
			return
		}
		pinstr := strconv.Itoa(c.Pin)
		if strings.Contains(m.Payload, pinstr) {
			c.Id = m.Sender.Recipient()
			bot.Send(m.Sender, "Sensor is now active!")
			log.Println("Sensor has been activated to user", m.Sender.Recipient())
			c.Save(*configPathPtr)
		} else {
			bot.Send(m.Sender, "Wrong activation code:")
		}

	})

	serve := func() {
		log.Println("Server Ready")
		buf := make([]byte, 4)
		for {
			_, addr, err := ServerConn.ReadFromUDP(buf)
			CheckError(err)

			now := time.Now()
			events = append(events, now)
			events = cleanEvents(events, 30)
                        sensorTriggered.Inc()
			log.Printf("[%s] Sensor triggered!! %d\n", addr, len(events))

			if checkEvents(events, 3) {
				events = events[len(events):]
                                sensorTriggered.Sub(3)
                                doorOpened.Inc()
				log.Printf(" --- Warning user ---\n")
				bot.Send(c, "Sensor triggered!")
				nowStr := now.Format("02-01-2006 15:04:05")
				_, err := histFile.WriteString(nowStr + "\n")
				CheckError(err)
			}
		}
	}

        http.Handle("/metrics", promhttp.Handler())
        go http.ListenAndServe(":4865", nil)

	go serve()
	bot.Start()

}
