package common

import (
	"github.com/ChimeraCoder/anaconda"
	"net/url"
)

func UpdateStatus(config *TwitterConfigType, message string) {
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
	api.PostTweet(message, url.Values{})
}
