package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"putio"
)

func main() {
	fmt.Println("-- Starting Example --")
	config := config{}
	f, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Config file not read : %s", err.Error())
	}
	if err = json.Unmarshal(f, &config); err != nil {
		fmt.Println("Error reading json from config : " + err.Error())
	}

	// create new putio object
	p, _ := NewPutio(config.Appid, config.Appsecret, config.Appcallbackurl, config.Usercode)

	if p.OauthToken == "" {
		fmt.Println("OAuth token is empty")
	} else {
		fmt.Println(p.OauthToken)
	}
}
