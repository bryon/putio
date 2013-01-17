package putio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	BaseUrl = "https://api.put.io/v2/"
)

var oauthparam = "?oauth_token="
var oathtoken string

type File struct {
	Is_shared          *bool
	Name               *string
	Screenshot         *string // returns url to image
	Created_at         *string // in iso8601 format
	Opensubtitles_hash *string
	Parent_id          *int // parent folder id
	Is_mp4_available   *bool
	Content_type       *string
	Crc32              *string
	Icon               *string // returns url to screenshot image in icon size
	Id                 *int
	Size               *int64
}

type Files struct {
	Status string `json: "status"`
	List   []File `json:"Files"`
	Parent File   `json:"parent"`
}

type Transfers struct {
}

type Account struct {
}

type Friends struct {
	token int
}

type Putio struct {
	OauthToken string
}

func (p *Putio) ListFiles() (files *Files, jsonstr string, err error) {
	url := BaseUrl + "/files/list" + oauthparam + p.OauthToken
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	// read in the body of the response
	defer resp.Body.Close()
	bodybytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodybytes, &files)
	if err != nil {
		return nil, string(bodybytes), err
	}

	return files, string(bodybytes), nil
}

// NewPutio takes in the apps oauth information and gets the token that will be used for all other calls
// This function doesn't have to be used if you provied the OauthToken when creating a Putio struct.
func NewPutio(appid, appsecret, appredirect, usercode string) (*Putio, error) {
	// get the user token using the calling apps credentials
	url := "https://api.put.io/v2/oauth2/access_token?client_id=" + appid + "&client_secret=" + appsecret + "&grant_type=authorization_code&redirect_uri=" + appredirect + "&code=" + usercode
	//fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// read in the body of the response
	defer resp.Body.Close()
	bodybytes, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// token returns as json result like { "access_token": "ABV9KDHN" }
	type oauthtoken struct {
		Access_token string
	}
	token := oauthtoken{}
	if err = json.Unmarshal(bodybytes, &token); err != nil {
		return nil, err
	}

	return &Putio{OauthToken: token.Access_token}, nil
}
