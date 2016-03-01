package common

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"net/url"
)

func UpdateStatus(config *TwitterConfigType, article *Article) {
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
	status := fmt.Sprintf("%s %s", article.Title, article.Link)
	api.PostTweet(status, url.Values{})
}
