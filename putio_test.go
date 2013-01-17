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
	ExpectedToken  string
}

func TestPutio(t *testing.T) {
	// Load Config from config.json
	config := config{}
	f, err := ioutil.ReadFile("config.json")
	if err != nil {
		t.Errorf("Config file not read : %s", err.Error())
	}
	if err = json.Unmarshal(f, &config); err != nil {
		t.Errorf("Error reading json from config : " + err.Error())
	}

	// create new putio object
	p, _ := NewPutio(config.Appid, config.Appsecret, config.Appcallbackurl, config.Usercode)

	if p.OauthToken == "" {
		t.Error("OAuth token is empty")
	}

	expstr := config.ExpectedToken
	if p.OauthToken != expstr {
		t.Errorf("OAuth token appears invalid.  Expected: %s got :%s", expstr, p.OauthToken)
	}

	// now test the most basic function, listing your files
	files, _, err := p.ListFiles()
	if err != nil {
		t.Error(err.Error())
	}
	//fmt.Println(jsonstr)
	fmt.Println(files)
	fmt.Println(files.Files)
}
