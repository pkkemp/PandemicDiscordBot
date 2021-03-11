package main

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

