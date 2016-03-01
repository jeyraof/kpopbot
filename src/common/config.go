package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type DBConfigType struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	SSLMode  string `json:"sslmode"`
}

type IRCConfigType struct {
	Nick            string `json:"nick"`
	NickAlternative string `json:"nickAlternative"`
	Server          string `json:"server"`
	Channel         string `json:"channel"`

	Nickserv         string `json:"nickserv"`
	NickservNick     string `json:"nickservNick"`
	NickservPassword string `json:"nickservPassword"`
}

type GoogleConfigType struct {
	Key string `json:"key"`
}

type TwitterConfigType struct {
	ConsumerKey       string `json:"consumerKey"`
	ConsumerSecret    string `json:"consumerSecret"`
	AccessToken       string `json:"accessToken"`
	AccessTokenSecret string `json:"accessTokenSecret"`
}

type ConfigType struct {
	DB  DBConfigType  `json:"db"`
	IRC IRCConfigType `json:"irc"`

	Google  GoogleConfigType  `json:"google"`
	Twitter TwitterConfigType `json:"twitter"`

	Period     uint   `json:"period"`
	Repository string `json:"repository"`
}

var Config ConfigType

func InitConfig() error {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println(err)
		return err
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Println(err)
	}
	return err
}

func BuildURIFrom(db *DBConfigType) string {
	config := *db
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		config.User, config.Password,
		config.Host, config.Port,
		config.Name,
		config.SSLMode)
}
