package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly"
	"log"
	"net/url"
	"strings"
	"time"
)

func findAppointments(dg *discordgo.Session) {
	for {
		c := colly.NewCollector()
		// Find and visit all links
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			rawURL := e.Attr("href")
			siteURL, _ := url.Parse(rawURL)
			path := siteURL.Path
			pathSubstrings := strings.Split(path, "/")
			if(siteURL.Host == "www.signupgenius.com" && pathSubstrings[1] == "go"){
				e.Request.Visit(e.Attr("href"))

			}

		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})

		c.OnHTML("td.SUGtable", func(table *colly.HTMLElement) {
			table.ForEachWithBreak("span.SUGbigbold", func(_ int, elem *colly.HTMLElement) bool {
				if strings.Contains(elem.Text, "slots filled") {
					messageText := "<@&" + VACCINEROLE + "> I've observed an available vaccination appointment at: \n" + table.Request.URL.Scheme + "://" + table.Request.URL.Host + table.Request.URL.Path
					message := discordgo.MessageSend{
						Content:   messageText,
					}
					_, err := dg.ChannelMessageSendComplex(VAXCHANNEL, &message)
					if err != nil {
						log.Print(err)
					}
					//we return false here to prevent the for loop from continuing
					//to search for available slots on this page
					return false
				}
				return false
			})
		})

		c.Visit(VAXOKCURL)
		time.Sleep(time.Minute)
	}
}
