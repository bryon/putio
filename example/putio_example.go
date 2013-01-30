package main

import (
	"encoding/json"
	"fmt"
	"github.com/bryon/putio"
	"io/ioutil"
)

type config struct {
	Appid          string
	Appsecret      string
	Appcallbackurl string
	Usercode       string
}

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
	p, _ := putio.NewPutio(config.Appid, config.Appsecret, config.Appcallbackurl, config.Usercode)

	if p.OauthToken == "" {
		fmt.Println("OAuth token is empty")
	} else {
		fmt.Println(p.OauthToken)
	}

	// now that we've got our account all set up and verified lets run a sample of items
	// create a folder for us to play around in 
	files, s, err := p.FilesCreateFolder("apitest", 0)
	fmt.Println("created.. " + files.Status)

	folderid := files.File.Id

	// now lets put a file into it 
	_, s, err = p.TransfersAdd("magnet:?xt=urn:btih:e1e90d4166168f6f2790fc3a0a61772ed27ab8cc&dn=The+Avengers+-+Clip&tr=http://tracker.publicbt.com:80/announce", folderid, true)
	fmt.Println(s)

	// rename the folder
	id := files.File.Id
	files, s, err = p.FilesRename(id, "apitest_renamed")
	fmt.Println("renamed.. " + files.Status)
	_, s, err = p.FilesId(id)
	fmt.Println(s)

	// delete the folder
	files, s, err = p.FilesDelete(id)
	fmt.Println("deleted.. " + files.Status)
	_, s, err = p.FilesId(id)
	fmt.Println(s)

	// list all your files
	files, jsonstr, err := p.FilesList()
	fmt.Printf("json len : %v\n", len(jsonstr))
	//fmt.Println(jsonstr)
	fmt.Println(files)
	fmt.Println(files.Files)
}
