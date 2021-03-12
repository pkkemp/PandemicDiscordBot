package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func fillChannelIDForMessages(messages []discordgo.Message, channelID string) []discordgo.Message {
	for _, message := range messages {
		message.ChannelID = channelID
	}
	return messages
}

func bulkSendDiscordMessages(messages []discordgo.MessageSend, dg *discordgo.Session, channelID string) {
	for _, message := range messages {
		dg.ChannelMessageSendComplex(channelID, &message)
	}
}

func loadDogImage() Dog {
	theDog := Dog{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", DOG_API_URL+"v1/images/search", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	req.Header.Set("X-API-KEY", DogToken)
	q := req.URL.Query()
	q.Add("has_breeds", "true")
	q.Add("mime_types", "jpg,png")
	q.Add("size", "small")
	q.Add("limit", "1")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&theDog)
	return theDog
}

func getQuote(s *discordgo.Session, m *discordgo.MessageCreate) {
	//read the contents of the quotes file into memory
	quotesFile, err := ioutil.ReadFile("./quotes/quotes.json")
	if err != nil {
		return
	}
	var quotes QuoteData
	err2 := json.Unmarshal(quotesFile, &quotes)
	if err2 != nil {
		log.Println(err2)
	}
	//generate a random number
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := len(quotes)
	num := rand.Intn(max-min+1) + min
	//build the message and select the random quote from our random number
	message := discordgo.MessageSend{
		Content:   "_" + quotes[num].Text + "_" + "\r" + "***—" + quotes[num].Author + "***",
		Reference: m.Reference(),
	}
	s.ChannelMessageSendComplex(m.ChannelID, &message)
}

func loadNASAImage() NASAImage {
	theImage := NASAImage{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", NASA_API_URL+"apod", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	q := req.URL.Query()
	q.Add("api_key", NASAAPIKey)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&theImage)
	return theImage
}

func xkcdImage(comic uint64) (XKCDComic, error) {
	comicNum := ""
	theImage := XKCDComic{}
	client := &http.Client{}
	if comic > 0 {
		comicNum = fmt.Sprint(comic)
	}
	req, err := http.NewRequest("GET", XKCD_URL+comicNum+"/info.0.json", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)

	var APIError error
	if resp.StatusCode > 200 {
		APIError = errors.New("NX")
	}

	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&theImage)
	return theImage, APIError
}
