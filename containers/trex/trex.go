package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	tb "gopkg.in/tucnak/telebot.v2"
)

type tconfig struct {
	Token     *string `json:"token"`
	ChannelID *string `json:"id"`
}

func (r tconfig) Recipient() string {
	return *r.ChannelID
}

func generateURL(tweet *twitter.Tweet) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%d", tweet.User.ScreenName, tweet.ID)
}

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		//fmt.Println("Error: " , err)
		log.Fatal(err)
		os.Exit(-1)
	}
}

func GetEnvStr(name, value string) string {
	if os.Getenv(name) != "" {
		return os.Getenv(name)
	}
	return value
}

func main() {
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)

	consumerKey := flags.String("consumer-key", GetEnvStr("TWITTER_CONSUMER_KEY", ""), "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", GetEnvStr("TWITTER_CONSUMER_SECRET", ""), "Twitter Consumer Secret")
	accessToken := flags.String("access-token", GetEnvStr("TWITTER_ACCESS_KEY", ""), "Twitter Access Token")
	accessSecret := flags.String("access-secret", GetEnvStr("TWITTER_ACCESS_SECRET", ""), "Twitter Access Secret")

	telegramConfig := new(tconfig)

	telegramConfig.Token = flags.String("telegram-token", GetEnvStr("TELEGRAM_TOKEN", ""), "Telegram Token")
	telegramConfig.ChannelID = flags.String("telegram-channel", GetEnvStr("TELEGRAM_CHANNEL", ""), "Telegram Channel")

	flags.Parse(os.Args[1:])

	/*
		logFile, err := os.OpenFile("/tmp/"+os.Args[0]+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	*/

	log.Printf("Accessing telegram with token %s\n", *telegramConfig.Token)
	log.Printf("Dumping to Channel %s\n", *telegramConfig.ChannelID)

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	bot, err := tb.NewBot(tb.Settings{
		Token:  *telegramConfig.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	CheckError(err)

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		if tweet.RetweetedStatus == nil {
			url := generateURL(tweet)
			log.Printf("%s\n", url)
			_, err = bot.Send(telegramConfig, url)
			if err != nil {
				log.Printf("ERROR: %v\n", err)
			}
		}
	}
	demux.DM = func(dm *twitter.DirectMessage) {
		log.Printf("Received DM from %d\n", dm.SenderID)
	}
	demux.Event = func(event *twitter.Event) {
		log.Printf("Twitter event: %#v\n", event)
	}

	log.Println("Starting Stream...")

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Language:      []string{"es"},
		Track:         []string{"T-Rex", "t-rex", "T-rex", "t-Rex"},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)
	go bot.Start()

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	log.Println("Stopping Stream...")
	stream.Stop()
}
