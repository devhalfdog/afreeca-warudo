package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/devhalfdog/afreecachat"
	"github.com/hypebeast/go-osc/osc"
	"github.com/joho/godotenv"
	"github.com/tidwall/gjson"
)

const (
	stationUrl = "https://bjapi.afreecatv.com/api/%s/station"
	infoUrl    = "https://api.m.afreecatv.com/broad/a/watch"
	dataUrl    = "https://live.afreecatv.com/afreeca/player_live_api.php?bjid=%s"

	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0"
)

var (
	chatClient *afreecachat.Client
	oscClient  *osc.Client
	client     *http.Client = &http.Client{Timeout: 5 * time.Second}
	stream     bool         = false

	bj string
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	bj = os.Getenv("BJ_ID")

	oscClient = osc.NewClient("localhost", 19190)

	checkStream()
}

func checkStream() {
	for {
		station, err := getStation()
		if err != nil {
			log.Println(err)
		}

		isLive := station.BroadNo != 0

		if !stream && isLive {
			err := watchChat(station)
			if err != nil {
				log.Println(err)
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func getStation() (Station, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(stationUrl, bj), nil)
	if err != nil {
		return Station{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return Station{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Station{}, fmt.Errorf("an unknown error occurred during the request: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Station{}, err
	}

	var station Station
	result := gjson.GetBytes(body, "broad")

	station.BroadNo = int(result.Get("broad_no").Int())

	return station, nil
}

func watchChat(station Station) error {
	token, err := setToken(station)
	if err != nil {
		return err
	}

	return setupChat(token)
}

func setupChat(token afreecachat.Token) error {
	var err error
	chatClient, err = afreecachat.NewClient(token)
	if err != nil {
		return err
	}

	chatClient.OnConnect(func(connect bool) {
		if connect {
			log.Println("Chatting connect")
		}
	})

	chatClient.OnBalloon(func(balloon afreecachat.Balloon) {
		log.Printf("nick : %s, count : %d\n", balloon.User.Name, balloon.Count)
		msg := osc.NewMessage("/osc/ballon")
		msg.Append(balloon.User.Name)
		msg.Append(int32(balloon.Count))
		err := oscClient.Send(msg)
		if err != nil {
			log.Printf("OnBalloon error : %s\n", err.Error())
		}
	})

	return chatClient.Connect()
}

func setToken(station Station) (afreecachat.Token, error) {
	token := afreecachat.Token{
		BJID: bj,
		Flag: "524304",
	}

	return token, nil
}
