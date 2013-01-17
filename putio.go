package putio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// fix problems with json not handling null as a value
type NString string
type NInt int

func (s *NString) UnmarshalJSON(b []byte) error {
	if string(b) == "null" { // now THIS is dumb.
		s = new(NString)
		return nil
	}
	var tmp string
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*s = NString(tmp)
	return nil
}
func (i *NInt) UnmarshalJSON(b []byte) error {
	if string(b) == "null" { // now THIS is dumb.
		i = new(NInt)
		return nil
	}
	var tmp int
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*i = NInt(tmp)
	return nil
}

const (
	BaseUrl = "https://api.put.io/v2/"
)

var oauthparam = "?oauth_token="
var oathtoken string

type File struct {
	Is_shared          bool    `json: "is_shared"`
	Name               NString `json: "name"`
	Screenshot         NString `json: "screenshot"` // returns url to image
	Created_at         NString `json: "created_at"` // in iso8601 format
	Opensubtitles_hash NString `json: "opensubtitles_hash"`
	Parent_id          NInt    `json: "parent_id"` // parent folder id
	Is_mp4_available   bool    `json: "is_mp4_available"`
	Content_type       NString `json: "content_type"`
	Crc32              NString `json: "crc32"`
	Icon               NString `json: "icon"` // returns url to screenshot image in icon size
	Id                 NInt    `json: "id"`
	Size               int64   `json: "size"`
}

type Files struct {
	Files  []File
	Status string
	Parent File
}

type Transfers struct {
}

type Account struct {
}

type Friends struct {
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
	f := Files{}
	err = json.Unmarshal(bodybytes, &f)
	if err != nil {
		return nil, string(bodybytes), err
	}
	return &f, string(bodybytes), nil
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
