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
	//_________________________________________________
	fmt.Println("\nFilesList Test..")
	files, jsonstr, err := p.FilesList()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("json len : %v\n", len(jsonstr))
	//fmt.Println(jsonstr)
	fmt.Println(files)
	fmt.Println(files.Files)

	//_________________________________________________
	fmt.Println("\nFilesSearch Test..")
	files, jsonstr, err = p.FilesSearch("Who", "-1")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("json len : %v\n", len(jsonstr))
	//fmt.Println(jsonstr)
	fmt.Println(files)
	fmt.Println(files.Next)

	//_________________________________________________
	fmt.Println("\nFilesId Test..")
	files, jsonstr, err = p.FilesId("55113191")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("json len : %v\n", len(jsonstr))
	//fmt.Println(jsonstr)
	fmt.Println(files.File)

	//_________________________________________________
	fmt.Println("\nFilesMP4 Test..")
	files, jsonstr, err = p.FilesMP4("55113191")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("json len : %v\n", len(jsonstr))
	fmt.Println(jsonstr)
	fmt.Println(files.Mp4.Status)

	//_________________________________________________
	fmt.Println("\nFilesDownload Test..")
	urlstr, err := p.FilesDownload("40704262")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(urlstr)
	//_________________________________________________
	fmt.Println("\nTransfersList Test..")
	transfers, jsonstr, err := p.TransfersList()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("json len : %v\n", len(jsonstr))
	//fmt.Println(jsonstr)
	fmt.Println(transfers)
	fmt.Println(transfers.Transfers)

	//_________________________________________________
	fmt.Println("\nTransfersId Test..")
	transfers, jsonstr, err = p.TransfersId("4485593")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("json len : %v\n", len(jsonstr))
	//fmt.Println(jsonstr)
	fmt.Println(transfers.Transfer)

	//_________________________________________________
	fmt.Println("\nAccount Info Test..")
	account, jsonstr, err := p.AccountInfo()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("json len : %v\n", len(jsonstr))
	//fmt.Println(jsonstr)
	fmt.Println(account)
	fmt.Println(account.Info.Disk)

	//_________________________________________________
	fmt.Println("\nAccount Settings Test..")
	account, jsonstr, err = p.AccountSettings()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("json len : %v\n", len(jsonstr))
	fmt.Println(jsonstr)
	fmt.Println(account.Settings)

	//_________________________________________________
	fmt.Println("\nAccount Settings Test..")
	friends, jsonstr, err := p.FriendsList()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("json len : %v\n", len(jsonstr))
	fmt.Println(friends)

	//_________________________________________________
	fmt.Println("\nFiles Create Folder Test..")
	files, jsonstr, err = p.FilesCreateFolder("apitest", 0)
	if err != nil {
		t.Error(err.Error())
	}
	//fmt.Printf("json len : %v\n", len(jsonstr))
	fmt.Println(jsonstr)
	fmt.Println(files)

	folderid := int(files.File.Id)
	// rename it
	files, jsonstr, err = p.FilesRename(folderid, "ApiRenamed")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(jsonstr)
	fmt.Println(files)

	// now delte it
	files, jsonstr, err = p.FilesDelete(folderid)
	if err != nil {
		t.Error(err.Error())
	}
	//fmt.Printf("json len : %v\n", len(jsonstr))
	fmt.Println(jsonstr)
	fmt.Println(files)

}
