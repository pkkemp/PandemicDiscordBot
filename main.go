package main

import (
        "encoding/json"
        "errors"
        "flag"
        "fmt"
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

        "github.com/bwmarrin/discordgo"
)
const DOG_API_URL   = "https://api.thedogapi.com/"
const CAT_API_URL   = "https://api.thecatapi.com/"
const NASA_API_URL = "https://api.nasa.gov/planetary/"
const XKCD_URL = "https://xkcd.com/"
// Variables used for command line parameters
var (
        Token string
        DogToken string
        CatToken string
        NASAAPIKey string
)

func init() {

        flag.StringVar(&Token, "t", "", "Bot Token")
        flag.Parse()
}

func main() {

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

        // Register the messageCreate func as a callback for MessageCreate events.
        dg.AddHandler(messageCreate)

        // In this example, we only care about receiving message events.
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
        // This isn't required in this specific example but it's a good practice.
        if m.Author.ID == s.State.User.ID {
                return
        }
        //messageContent := strings.ToLower(m.Content)
        //replace weird iOS single quotes/apostrophe
        messageContent := strings.Replace(m.Content, "’", "'", -1)
        //split the words into an array
        words := strings.Split(messageContent, " ")
        switch strings.ToLower(words[0]) {
        case "ping":
                s.ChannelMessageSend(m.ChannelID, "Pong!")
                s.UpdateGameStatus(0, "the Stock Market")
        case "$subscribe":
                allRoles, _ := s.GuildRoles(m.GuildID)
                name := strings.Join(words[1:], " ")
                for _, role := range allRoles {
                        if(role.Name == name) {
                                //TODO fix the permissions here and add the rest of discord bitwise flags
                                //right now the || role.Permissions > 0 will always make this statement false if the role has associated permissions
                                if(role.Permissions > 0 ||
                                        (role.Permissions & discordgo.PermissionAdministrator == discordgo.PermissionAdministrator) ||
                                        (role.Permissions & discordgo.PermissionManageRoles == discordgo.PermissionManageRoles)) {
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
                        if(role.Name == name) {
                                //TODO fix the permissions here and add the rest of discord bitwise flags
                                //right now the || role.Permissions > 0 will always make this statement false if the role has associated permissions
                                if(role.Permissions > 0 ||
                                    (role.Permissions & discordgo.PermissionAdministrator == discordgo.PermissionAdministrator) ||
                                    (role.Permissions & discordgo.PermissionManageRoles == discordgo.PermissionManageRoles)) {
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
                        if(role.Name == name) {
                                if(role.Permissions > 0) {
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
                if (err != nil) {
                        log.Println(err)
                }
                rand.Seed(time.Now().UnixNano())
                min := 0
                max := 0xFFFFFF
                randomColor := rand.Intn(max - min + 1) + min

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
                Dogs := loadImage()
                dog := Dogs[0]
                // Create the file
                out, err := os.Create("./"+dog.ID+".jpg")
                if err != nil {

                }
                defer out.Close()

                resp, err := http.Get(dog.URL)
                // Write the body to file
                _, err = io.Copy(out, resp.Body)
                f, err := os.Open("./"+dog.ID+".jpg")
                defer f.Close()
                var r io.Reader
                r = f
                file := discordgo.File{
                        Name:        dog.Breeds[0].Name+".jpg",
                        ContentType: "image/jpg",
                        Reader:      r,
                }
                message := discordgo.MessageSend{
                        Content:         "***"+dog.Breeds[0].Name + "***\r*"+dog.Breeds[0].Temperament+"* ",
                        File:           &file,
                        Reference: m.Reference(),
                }
                //s.ChannelMessageSend(m.ChannelID, "***"+dog.Breeds[0].Name + "*** \r *"+dog.Breeds[0].Temperament+"* " + dog.URL)
                s.ChannelMessageSendComplex(m.ChannelID, &message)
                os.Remove("./"+dog.ID+".jpg")
        case "nasa":
                NASA := loadNASAImage()
                out, err := os.Create("./"+NASA.Title+".jpg")
                if err != nil {

                }
                defer out.Close()

                resp, err := http.Get(NASA.URL)
                // Write the body to file
                _, err = io.Copy(out, resp.Body)
                f, err := os.Open("./"+NASA.Title+".jpg")
                defer f.Close()
                var r io.Reader
                r = f
                file := discordgo.File{
                        Name:        NASA.Title+".jpg",
                        ContentType: "image/jpg",
                        Reader:      r,
                }
                message := discordgo.MessageSend{
                        Content:         "***NASA Astronomy Picture of the Day***\r***"+NASA.Title + "***\r*"+NASA.Explanation+"*" + "\r_© " + NASA.Copyright + "_",
                        File:           &file,
                        Reference: m.Reference(),
                }
                //s.ChannelMessageSend(m.ChannelID, "***"+dog.Breeds[0].Name + "*** \r *"+dog.Breeds[0].Temperament+"* " + dog.URL)
                s.ChannelMessageSendComplex(m.ChannelID, &message)
                os.Remove("./"+NASA.Title+".jpg")
        case "xkcd":
                var comicNum uint64 = 0
                if(len(words) > 1 && words[1] != "") {
                        var err error
                        comicNum, err = strconv.ParseUint(words[1], 10, 32)
                        if(err != nil) {
                                //something bad happened, exit this case
                                break
                        }
                }
                comic, APIerr := xkcdImage(comicNum)
                if(APIerr != nil) {
                        //API failed
                        break
                }
                // Create the file
                out, err := os.Create("./"+comic.Month+comic.Day+comic.Year+".png")
                if err != nil {

                }
                defer out.Close()

                resp, err := http.Get(comic.Img)
                // Write the body to file
                _, err = io.Copy(out, resp.Body)
                f, err := os.Open("./"+comic.Month+comic.Day+comic.Year+".png")
                defer f.Close()
                var r io.Reader
                r = f
                file := discordgo.File{
                        Name:        comic.Month+comic.Day+comic.Year+".png",
                        ContentType: "image/png",
                        Reader:      r,
                }
                message := discordgo.MessageSend{
                        Content:         "***"+comic.Title + "***\r*"+comic.Month + "-" + comic.Day + "-" + comic.Year + "* ",
                        File:           &file,
                        Reference: m.Reference(),
                }
                //s.ChannelMessageSend(m.ChannelID, "***"+dog.Breeds[0].Name + "*** \r *"+dog.Breeds[0].Temperament+"* " + dog.URL)
                s.ChannelMessageSendComplex(m.ChannelID, &message)
                os.Remove("./"+comic.Month+comic.Day+comic.Year+".png")
        case "moviequote":
                //read the contents of the quotes file into memory
                quotesFile, err := ioutil.ReadFile("./json-tv-quotes/quotes.json")
                if(err != nil) {
                        break
                }
                var quotes MovieQuotes
                err2 := json.Unmarshal(quotesFile, &quotes)
                if(err2 != nil) {
                        log.Println(err2)
                }
                rand.Seed(time.Now().UnixNano())
                min := 0
                max := len(quotes)
                num := rand.Intn(max - min + 1) + min
                //log.Println(quotes[num])
                messageText := "_" + quotes[num].Quote + "_" +"\r"+"***—"+quotes[num].Author+"***"
                if(quotes[num].Source != "") {
                        messageText += " _(" + quotes[num].Source + ")_"
                }
                message := discordgo.MessageSend{
                        Content:         messageText,
                        Reference: m.Reference(),
                }
                s.ChannelMessageSendComplex(m.ChannelID, &message)
        case "quote":
                //read the contents of the quotes file into memory
                quotesFile, err := ioutil.ReadFile("./quotes/quotes.json")
                if(err != nil) {
                        break
                }
                var quotes QuoteData
                err2 := json.Unmarshal(quotesFile, &quotes)
                if(err2 != nil) {
                        log.Println(err2)
                }
                rand.Seed(time.Now().UnixNano())
                min := 0
                max := len(quotes)
                num := rand.Intn(max - min + 1) + min
                //log.Println(quotes[num])
                message := discordgo.MessageSend{
                        Content:         "_" + quotes[num].Text + "_" +"\r"+"***—"+quotes[num].Author+"***",
                        Reference: m.Reference(),
                }
                s.ChannelMessageSendComplex(m.ChannelID, &message)

        case "it's":
                if(strings.ToLower(words[1]) == "thursday") {
                        weekday := time.Now().Weekday()
                        var message discordgo.MessageSend
                        if (int(weekday) != 4) {
                                message = discordgo.MessageSend{
                                        Content: "*no it's not*",
                                }
                        } else {
                                f, _ := os.Open("./" + "thursday.gif")
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
                        //s.ChannelMessageSend(m.ChannelID, "***"+dog.Breeds[0].Name + "*** \r *"+dog.Breeds[0].Temperament+"* " + dog.URL)
                        s.ChannelMessageSendComplex(m.ChannelID, &message)
                }
        default:
                words := strings.Split(messageContent, " ")
                re, err := regexp.Compile(`[^\w]`)
                if err != nil {
                        log.Fatal(err)
                }
                for _, word := range words {
                        tempWord := re.ReplaceAllString(word, "")
                        if(tempWord == "why") {
                                message := discordgo.MessageSend{
                                        Content:         "_You see things; and you say 'Why?' But I dream things that never were; and I say 'Why not?'☄️✨_"+"\r"+"***—George Bernard Shaw***",
                                        Reference: m.Reference(),
                                }
                                s.ChannelMessageSendComplex(m.ChannelID, &message)
                                //make sure we only send it once per message
                                break
                        }

                }
        }
}

type QuoteData []struct {
        Author string `json:"author"`
        Text   string `json:"text"`
}

type MovieQuotes []struct {
        Type     string `json:"type"`
        Language string `json:"language"`
        Quote    string `json:"quote"`
        Author   string `json:"author,omitempty"`
        Source   string `json:"source,omitempty"`
}

type Dog []struct {
        Breeds []struct {
                Weight struct {
                        Imperial string `json:"imperial"`
                        Metric   string `json:"metric"`
                } `json:"weight"`
                Height struct {
                        Imperial string `json:"imperial"`
                        Metric   string `json:"metric"`
                } `json:"height"`
                ID               int    `json:"id"`
                Name             string `json:"name"`
                BredFor          string `json:"bred_for"`
                BreedGroup       string `json:"breed_group"`
                LifeSpan         string `json:"life_span"`
                Temperament      string `json:"temperament"`
                ReferenceImageID string `json:"reference_image_id"`
        } `json:"breeds"`
        ID     string `json:"id"`
        URL    string `json:"url"`
        Width  int    `json:"width"`
        Height int    `json:"height"`
}

type NASAImage struct {
        Copyright      string `json:"copyright"`
        Date           string `json:"date"`
        Explanation    string `json:"explanation"`
        Hdurl          string `json:"hdurl"`
        MediaType      string `json:"media_type"`
        ServiceVersion string `json:"service_version"`
        Title          string `json:"title"`
        URL            string `json:"url"`
}

type XKCDComic struct {
        Month      string `json:"month"`
        Num        int    `json:"num"`
        Link       string `json:"link"`
        Year       string `json:"year"`
        News       string `json:"news"`
        SafeTitle  string `json:"safe_title"`
        Transcript string `json:"transcript"`
        Alt        string `json:"alt"`
        Img        string `json:"img"`
        Title      string `json:"title"`
        Day        string `json:"day"`
}
func loadImage() Dog {
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
        if(comic > 0) {
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
        if(resp.StatusCode > 200) {
                APIError = errors.New("NX")
        }
        defer resp.Body.Close()
        json.NewDecoder(resp.Body).Decode(&theImage)
        return theImage, APIError
}

