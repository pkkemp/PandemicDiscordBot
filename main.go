package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const DOG_API_URL = "https://api.thedogapi.com/"
const CAT_API_URL = "https://api.thecatapi.com/"
const NASA_API_URL = "https://api.nasa.gov/planetary/"
const XKCD_URL = "https://xkcd.com/"

//findAppointment constants
const VAXCHANNEL = "819118034903236628"
const VACCINEROLE = "819282075164737577"
const VAXOKCURL = "https://www.vaxokc.com/"

// Variables used for command line parameters
var (
	Token      string
	DogToken   string
	CatToken   string
	NASAAPIKey string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}


func main() {
	//get our API keys from the OS env vars
	Token = os.Getenv("DISCORDTOKEN")
	DogToken = os.Getenv("DOGAPITOKEN")
	CatToken = os.Getenv("CATTOKEN")
	NASAAPIKey = os.Getenv("NASAAPIKEY")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
			fmt.Println("error creating Discord session,", err)
			return
	}

	//start the findAppointments function in a goroutine
	go findAppointments(dg)
	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)


	// Set the Discord gateway to notify us when messages are sent
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// replace weird iOS single quotes/apostrophe
	messageContent := strings.Replace(m.Content, "’", "'", -1)
	// split the words into an array
	words := strings.Split(messageContent, " ")
	// switch on the first word of the message
	switch strings.ToLower(words[0]) {
	case "ping":
		s.ChannelMessageSend(m.ChannelID, "Pong!")
		s.UpdateGameStatus(0, "the Stock Market")
	case "$subscribe":
		allRoles, _ := s.GuildRoles(m.GuildID)
		name := strings.Join(words[1:], " ")
		for _, role := range allRoles {
			if strings.ToLower(role.Name) == strings.ToLower(name) {
				//TODO fix the permissions here and add the rest of discord bitwise flags
				//right now the || role.Permissions > 0 will always make this statement false if the role has associated permissions
				if role.Permissions > 0 ||
					(role.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator) ||
					(role.Permissions&discordgo.PermissionManageRoles == discordgo.PermissionManageRoles) {
					s.ChannelMessageSendReply(m.ChannelID, "I'm sorry, but I can't grant you roles with permissions", m.Reference())
				} else {
					s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, role.ID)
					s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
				}
				//we found a match, so we can stop the loop
				break
			}

		}
	case "$unsubscribe":
		allRoles, _ := s.GuildRoles(m.GuildID)
		name := strings.Join(words[1:], " ")
		for _, role := range allRoles {
			if strings.ToLower(role.Name) == strings.ToLower(name) {
				//TODO fix the permissions here and add the rest of discord bitwise flags
				//right now the || role.Permissions > 0 will always make this statement false if the role has associated permissions
				if role.Permissions > 0 ||
					(role.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator) ||
					(role.Permissions&discordgo.PermissionManageRoles == discordgo.PermissionManageRoles) {
					s.ChannelMessageSendReply(m.ChannelID, "I'm sorry, but I can't remove roles with permissions", m.Reference())
				} else {
					s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, role.ID)
					s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
				}
				//we found a match, so we can stop the loop
				break
			}

		}
	case "$deleterole":
		allRoles, _ := s.GuildRoles(m.GuildID)
		name := strings.Join(words[1:], " ")
		for _, role := range allRoles {
			if role.Name == name {
				if role.Permissions > 0 {
					s.ChannelMessageSendReply(m.ChannelID, "I'm sorry, but I can't delete roles with permissions", m.Reference())
				} else {
					s.GuildRoleDelete(m.GuildID, role.ID)
					s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
				}
				//we found a match, so we can stop the loop
				break
			}

		}
	case "$createrole":
		role, err := s.GuildRoleCreate(m.GuildID)
		if err != nil {
			log.Println(err)
		}
		rand.Seed(time.Now().UnixNano())
		min := 0
		max := 0xFFFFFF
		randomColor := rand.Intn(max-min+1) + min

		name := strings.Join(words[1:], " ")
		s.GuildRoleEdit(m.GuildID, role.ID, name, randomColor, false, 0, true)
		s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, role.ID)
		s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	case "fetch!":
		s.ChannelMessageSend(m.ChannelID, "Ruff!")
		s.UpdateGameStatus(0, "Fetch")
	case "pong":
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	case "woof":
		Dogs := loadDogImage()
		if Dogs == nil {
			//API failed
			emoji := discordgo.Emoji{
				ID:            "807498175886786580",
				Name:          "broken_wifi",
				User:          nil,
				RequireColons: true,
				Managed:       false,
				Animated:      false,
				Available:     true,
			}
			s.MessageReactionAdd(m.ChannelID, m.ID, emoji.APIName())
			break
		}
		dog := Dogs[0]
		// Create the file
		tempFile, err  := ioutil.TempFile("", dog.ID + ".*.jpg")
		if err != nil {
			log.Print(err)
			break
		}
		defer os.Remove(tempFile.Name()) // clean up

		resp, err := http.Get(dog.URL)
		// Write the body to file
		_, err = io.Copy(tempFile, resp.Body)
		f, err := os.Open(tempFile.Name())
		if err != nil {
			//something bad happened, exit this case
			log.Print(err)
			break
		}
		defer f.Close()
		var r io.Reader
		r = f
		file := discordgo.File{
			Name:        dog.Breeds[0].Name + ".jpg",
			ContentType: "image/jpg",
			Reader:      r,
		}
		message := discordgo.MessageSend{
			Content:   "***" + dog.Breeds[0].Name + "***\r*" + dog.Breeds[0].Temperament + "* ",
			File:      &file,
			Reference: m.Reference(),
		}
		s.ChannelMessageSendComplex(m.ChannelID, &message)
		// delete the image from disk now that we've sent it
		os.Remove("./" + dog.ID + ".jpg")
	case "nasa":
		NASA := loadNASAImage()
		//create the file
		tempFile, err  := ioutil.TempFile("", NASA.Title + ".*.jpg")

		if err != nil {
			//something bad happened, exit this case
			log.Print(err)
			break
		}
		defer os.Remove(tempFile.Name()) // clean up

		resp, err := http.Get(NASA.URL)
		// Write the body to file
		_, err = io.Copy(tempFile, resp.Body)
		f, err := os.Open(tempFile.Name())
		if err != nil {
			//something bad happened, exit this case
			log.Print(err)
			break
		}
		defer f.Close()
		var r io.Reader
		r = f
		file := discordgo.File{
			Name:        NASA.Title + ".jpg",
			ContentType: "image/jpg",
			Reader:      r,
		}
		message := discordgo.MessageSend{
			Content:   "***NASA Astronomy Picture of the Day***\r***" + NASA.Title + "***\r*" + NASA.Explanation + "*" + "\r_© " + NASA.Copyright + "_",
			File:      &file,
			Reference: m.Reference(),
		}
		s.ChannelMessageSendComplex(m.ChannelID, &message)
		os.Remove("./" + NASA.Title + ".jpg")
	case "xkcd":
		var comicNum uint64 = 0
		if len(words) > 1 && words[1] != "" {
			var err error
			comicNum, err = strconv.ParseUint(words[1], 10, 32)
			if err != nil {
				//something bad happened, exit this case
				log.Print(err)
				break
			}
		}
		comic, APIerr := xkcdImage(comicNum)
		if APIerr != nil {
			//API failed
			break
		}
		// Create the file
		//out, err := os.Create("./" + comic.Month + comic.Day + comic.Year + ".png")
		tempFile, err  := ioutil.TempFile("", comic.Month + comic.Day + comic.Year +  ".*.png")
		if err != nil {
			//something bad happened, exit this case
			log.Print(err)
			break
		}
		defer os.Remove(tempFile.Name()) // clean up

		resp, err := http.Get(comic.Img)
		// Write the body to file
		_, err = io.Copy(tempFile, resp.Body)
		f, err := os.Open(tempFile.Name())
		defer f.Close()
		var r io.Reader
		r = f
		file := discordgo.File{
			Name:        comic.Month + comic.Day + comic.Year + ".png",
			ContentType: "image/png",
			Reader:      r,
		}
		message := discordgo.MessageSend{
			Content:   "***" + comic.Title + "***\r*" + comic.Month + "-" + comic.Day + "-" + comic.Year + "* ",
			File:      &file,
			Reference: m.Reference(),
		}
		//s.ChannelMessageSend(m.ChannelID, "***"+dog.Breeds[0].Name + "*** \r *"+dog.Breeds[0].Temperament+"* " + dog.URL)
		s.ChannelMessageSendComplex(m.ChannelID, &message)
		os.Remove("./" + comic.Month + comic.Day + comic.Year + ".png")
	case "moviequote":
		//read the contents of the quotes file into memory
		quotesFile, err := ioutil.ReadFile("/data/json-tv-quotes/quotes.json")
		if err != nil {
			//something bad happened, exit this case
			log.Print(err)
			break
		}
		var quotes MovieQuotes
		err2 := json.Unmarshal(quotesFile, &quotes)
		if err2 != nil {
			//something bad happened, exit this case
			log.Print(err2)
			break
		}
		rand.Seed(time.Now().UnixNano())
		min := 0
		max := len(quotes)
		num := rand.Intn(max-min+1) + min
		//log.Println(quotes[num])
		messageText := "_" + quotes[num].Quote + "_" + "\r" + "***—" + quotes[num].Author + "***"
		if quotes[num].Source != "" {
			messageText += " _(" + quotes[num].Source + ")_"
		}
		message := discordgo.MessageSend{
			Content:   messageText,
			Reference: m.Reference(),
		}
		s.ChannelMessageSendComplex(m.ChannelID, &message)
	case "quote":
		go getQuote(s, m)
	case "it's":
		if strings.ToLower(words[1]) == "thursday"  {
			weekday := time.Now().Weekday()
			var message discordgo.MessageSend
			if int(weekday) != 4 {
				message = discordgo.MessageSend{
					Content: "*no it's not*",
				}
			} else {
				f, _ := os.Open("/data/" + "thursday.gif")
				defer f.Close()
				var r io.Reader
				r = f
				file := discordgo.File{
					Name:        "thursday.gif",
					ContentType: "image/gif",
					Reader:      r,
				}
				message = discordgo.MessageSend{
					Content: "***what a concept.***",
					File:    &file,
				}
			}
			s.ChannelMessageSendComplex(m.ChannelID, &message)
		}
	case "weather":
		s.ChannelMessageSend(m.ChannelID, "There are currently " + strconv.Itoa(getActiveAlertCount()) + " active watches and warnings.")
		alerts := getActiveAlerts()
		ZONE_CODES_MONITORED := []string{"OKZ015", "AKZ125", "ARC035"}
		filteredAlerts := searchAlertsByZoneCode(ZONE_CODES_MONITORED, alerts)
		messages := constructAlertDiscordMessages(filteredAlerts)
		bulkSendDiscordMessages(messages, s,  m.ChannelID)

	default:
		// don't run the last case
		break
		words := strings.Split(messageContent, " ")
		re, err := regexp.Compile(`[^\w]`)
		if err != nil {
			log.Fatal(err)
		}
		for _, word := range words {
			tempWord := re.ReplaceAllString(word, "")
			if tempWord == "why" {
				message := discordgo.MessageSend{
					Content:   "_You see things; and you say 'Why?' But I dream things that never were; and I say 'Why not?'☄️✨_" + "\r" + "***—George Bernard Shaw***",
					Reference: m.Reference(),
				}
				s.ChannelMessageSendComplex(m.ChannelID, &message)
				//make sure we only send it once per message
				break
			}

		}
	}
}

