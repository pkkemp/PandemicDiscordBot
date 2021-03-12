package main

import (
	"encoding/xml"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
	"os"
)

const NOAAAlertsURL = "https://alerts.weather.gov/cap/us.php?x=0"


func getActiveAlerts() NOAAAlertFeed {
	alerts := NOAAAlertFeed{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", NOAAAlertsURL, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	xml.NewDecoder(resp.Body).Decode(&alerts)
	return alerts
}

func getActiveAlertCount() int {
	return len(getActiveAlerts().Entry)
}

func constructAlertDiscordMessages(feed NOAAAlertFeed) []discordgo.MessageSend {
	var messages []discordgo.MessageSend
	for _, alert := range feed.Entry {
		messages = append(messages,
			discordgo.MessageSend{
				Content:         alert.Event + "\n" + alert.Title + "\n" + "***" + alert.Summary + "***\n" + "_Effective: " + alert.Effective + "_\n_Expires: " + alert.Expires + "_\n",
			})

	}
	return messages
}

func searchAlertsByZoneCode(zones []string, feed NOAAAlertFeed) NOAAAlertFeed {
	var filteredAlerts NOAAAlertFeed
	for _, alert := range feed.Entry {
		_, applicableZone := find(zones, alert.Geocode.Value[1])
		if applicableZone {
			filteredAlerts.Entry = append(filteredAlerts.Entry, alert)
		}
	}
	return filteredAlerts
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
