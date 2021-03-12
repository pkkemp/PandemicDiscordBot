package main

import "encoding/xml"

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

type NOAAAlertFeed struct {
	XMLName   xml.Name `xml:"feed"`
	Text      string   `xml:",chardata"`
	Xmlns     string   `xml:"xmlns,attr"`
	Cap       string   `xml:"cap,attr"`
	Ha        string   `xml:"ha,attr"`
	ID        string   `xml:"id"`
	Logo      string   `xml:"logo"`
	Generator string   `xml:"generator"`
	Updated   string   `xml:"updated"`
	Author    struct {
		Text string `xml:",chardata"`
		Name string `xml:"name"`
	} `xml:"author"`
	Title string `xml:"title"`
	Link  struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Entry []struct {
		Text      string `xml:",chardata"`
		ID        string `xml:"id"`
		Updated   string `xml:"updated"`
		Published string `xml:"published"`
		Author    struct {
			Text string `xml:",chardata"`
			Name string `xml:"name"`
		} `xml:"author"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
		} `xml:"link"`
		Summary   string `xml:"summary"`
		Event     string `xml:"event"`
		Effective string `xml:"effective"`
		Expires   string `xml:"expires"`
		Status    string `xml:"status"`
		MsgType   string `xml:"msgType"`
		Category  string `xml:"category"`
		Urgency   string `xml:"urgency"`
		Severity  string `xml:"severity"`
		Certainty string `xml:"certainty"`
		AreaDesc  string `xml:"areaDesc"`
		Polygon   string `xml:"polygon"`
		Geocode   struct {
			Text      string   `xml:",chardata"`
			ValueName []string `xml:"valueName"`
			Value     []string `xml:"value"`
		} `xml:"geocode"`
		Parameter struct {
			Text      string `xml:",chardata"`
			ValueName string `xml:"valueName"`
			Value     string `xml:"value"`
		} `xml:"parameter"`
	} `xml:"entry"`
}

