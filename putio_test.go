package putio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

type config struct {
	Appid          string
	Appsecret      string
	Appcallbackurl string
	Usercode       string
}

func TestPutio(t *testing.T) {
	config := config{}
	f, err := ioutil.ReadFile("config.json")
	if err != nil {
		t.Errorf("Config file not read : %s", err.Error())
	}
	if err = json.Unmarshal(f, &config); err != nil {
		t.Errorf("Error reading json from config : " + err.Error())
	}

	p, _ := NewPutio(config.Appid, config.Appsecret, config.Appcallbackurl, config.Usercode)
	exstr := "ABV9KDHN"
	if p.OauthToken != exstr {
		t.Errorf("OAuth token appears invalid.  Expected: %s got :%s", exstr, p.OauthToken)
	}
}
