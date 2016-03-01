package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GoogleShortURL struct {
	Kind    string `json:"kind"`
	ID      string `json:"id"`
	LongURL string `json:"longUrl"`
}

func URLShorten(key string, url string) string {
	body := []byte(fmt.Sprintf("{\"longUrl\": \"%s\"}", url))
	api := fmt.Sprintf("https://www.googleapis.com/urlshortener/v1/url?key=%s", key)
	resp, err := http.Post(api, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return url
	}
	defer resp.Body.Close()

	shortURL := GoogleShortURL{}
	if jsonErr := json.NewDecoder(resp.Body).Decode(&shortURL); jsonErr != nil {
		return url
	}

	return shortURL.ID
}

func ArticleShorten(google *GoogleConfigType, article *Article) {
	(*article).Link = URLShorten((*google).Key, (*article).Link)
}
