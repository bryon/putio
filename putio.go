package putio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	BaseUrl = "https://api.put.io/v2/"
)

var Paths = map[string]string{
	// values in {{ }} are variables.  Marked thusly for template use
	"FilesList":         "/files/list",
	"FilesSearch":       "/files/search/{{query}}/page/{{pageno}}",
	"FilesUpload":       "/files/upload",
	"FilesCreateFolder": "/files/create-folder",
	"FilesId":           "/files/{{id}}",
	"FilesDelete":       "/files/delete",
	"FilesRename":       "/files/rename",
	"FilesMove":         "/files/move",
	"FilesMP4":          "/files/{{id}}/mp4",
	"FilesDowload":      "/files/{{id}}/download",
	"FilesZip":          "/files/zip",
	"TransfersList":     "/transfers/list",
	"TransfersAdd":      "/transfers/add",
	"TransfersId":       "/transfers/{{id}}",
	"TransfersCancel":   "/transfers/cancel",
	"AccountInfo":       "/account/info",
	"AccountSettings":   "/account/settings",
	"FriendsList":       "/friends/list",
	"FriendsWaiting":    "/friends/waiting-requests",
	"FriendsRequest":    "/friends/{{username}}/request",
	"FriendsDeny":       "/friends/{{username}}/deny",
}

type Files struct {
}

type Transfers struct {
}

type Account struct {
}

type Friends struct {
}

type Putio struct {
	OauthToken string
	Files
	Transfers
	Account
	Friends
}

func NewPutio(appid, appsecret, appredirect, usercode string) (*Putio, error) {
	// get the user token using the calling apps credentials
	//token := "ABV9KDHN"
	url := "https://api.put.io/v2/oauth2/access_token?client_id=" + appid + "&client_secret=" + appsecret + "&grant_type=authorization_code&redirect_uri=" + appredirect + "&code=" + usercode

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error calling oauth service : " + err.Error())
		return nil, err
	}
	// read in the body of the response
	defer resp.Body.Close()
	bodybytes, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading oauth response : " + err.Error())
		return nil, err
	}
	//token := string(bodybytes)

	// token returns as json like { "access_token": "ABV9KDHN" }
	type oauthtoken struct {
		Access_token string
	}
	token := oauthtoken{}
	err = json.Unmarshal(bodybytes, &token)
	if err != nil {
		fmt.Println("Error reading json from oauth : " + err.Error())
		return nil, err
	}

	return &Putio{OauthToken: token.Access_token}, nil
}
