package putio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// fix problems with json not handling null as a value
type NString string
type NInt int
type NInt64 int64

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
func (i *NInt64) UnmarshalJSON(b []byte) error {
	if string(b) == "null" { // now THIS is dumb.
		i = new(NInt64)
		return nil
	}
	var tmp int64
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*i = NInt64(tmp)
	return nil
}

const (
	BaseUrl = "https://api.put.io/v2/"
)

var oauthparam = "?oauth_token="
var oathtoken string

type MP4 struct {
	Status       NString
	Stream_url   NString
	Download_url NString
	Size         NInt64
	Percent_done NInt
}

type File struct {
	Is_shared          bool    `json: "is_shared"`
	Name               NString `json: "name"`
	Screenshot         NString `json: "screenshot"` // returns url to image
	Created_at         NString `json: "created_at"` // in iso8601 format
	Opensubtitles_hash NString `json: "opensubtitles_hash"`
	Parent_id          NInt64  `json: "parent_id"` // parent folder id
	Is_mp4_available   bool    `json: "is_mp4_available"`
	Content_type       NString `json: "content_type"`
	Crc32              NString `json: "crc32"`
	Icon               NString `json: "icon"` // returns url to screenshot image in icon size
	Id                 NInt64  `json: "id"`
	Size               NInt64  `json: "size"`
}

type Files struct {
	Files  []File // for multi file results
	File   File   // for single file result like files/id
	Mp4    MP4    // for mp4 streaming results
	Status string
	Parent File
	Next   NString
}

type Transfer struct {
	Uploaded        NInt64  `json: "uploaded"`
	EstimatedTime   NInt    `json: "estimated_time"`
	PeersGetting    NInt    `json: "peers_getting_from_us"`
	Extract         bool    `json: "extract"`
	CurrentRatio    float64 `json: "current_ratio"`
	Size            NInt64  `json: "size"`
	UpSpeed         NInt64  `json: "up_speed"`
	Id              NInt64  `json: "id"`
	Source          NString `json: "source"`
	Subscription_id NInt64  `json: "subscription_id"`
	StatusMessage   NString `json: "status_message"`
	Status          NString `json: "status"`
	DownSpeed       NString `json: "down_speed"`
	PeersConnected  NInt    `json: "peers_connected"`
	Downloaded      NInt64  `json: "downloaded"`
	FileId          NInt64  `json: "file_id"`
	PeersSending    NInt    `json: "peers_sending_to_us"`
	PercentDone     NInt    `json: "percent_done"`
	IsPrivate       bool    `json: "is_private"`
	TrackerMessage  NString `json: "tracker_message"`
	Name            NString `json: "name"`
	CreatedAt       NString `json: "created_at"`
	ErrorMessage    NString `json: "error_message"`
	SaveParentId    NInt64  `json: "save_parent_id"`
	CallbackUrl     NString `json: "callback_url"`
}

type Transfers struct {
	Status    string
	Transfers []Transfer
	Transfer  Transfer
}

type Disk struct {
	Available NInt64
	Used      NInt64
	Size      NInt64
}

type UserInfo struct {
	Username string
	Mail     string
	Disk     Disk
}

type Settings struct {
	Routing               string `json: "routing"`
	HideItemsShared       string `json: "hide_items_shared"`
	DefaultDownloadFolder int    `json: "default_download_folder"`
	SSLEnabled            bool   `json: "ssl_enabled"`
	IsInvisible           bool   `json: "is_invisible"`
	ExtractionDefault     string `json: "extraction_default"`
}

type Account struct {
	Status   string
	Info     UserInfo
	Settings Settings
}

type Friend struct {
	Name string
}

type Friends struct {
	Status  string
	Friends []Friend
	Friend  Friend
}

type Putio struct {
	OauthToken string
}

func (p *Putio) GetReqBody(path string) (bodybytes []byte, err error) {
	url := BaseUrl + path + oauthparam + p.OauthToken
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// read in the body of the response
	defer resp.Body.Close()
	bodybytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodybytes, nil
}

func (p *Putio) PostFilesReq(path string, data url.Values) (files *Files, jsonstr string, err error) {
	posturl := BaseUrl + path + oauthparam + p.OauthToken
	resp, err := http.PostForm(posturl, data)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	bodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	if err = json.Unmarshal(bodybytes, &files); err != nil {
		return nil, string(bodybytes), err
	}
	return files, string(bodybytes), nil
}

func (p *Putio) GetFilesReq(path string) (files *Files, jsonstr string, err error) {
	bodybytes, err := p.GetReqBody(path)
	if err != nil {
		return nil, string(bodybytes), err
	}
	if err = json.Unmarshal(bodybytes, &files); err != nil {
		return nil, string(bodybytes), err
	}
	return files, string(bodybytes), nil
}

// https://api.put.io/v2/docs/#files-list
func (p *Putio) FilesList() (files *Files, jsonstr string, err error) {
	return p.GetFilesReq("files/list")
}

// https://api.put.io/v2/docs/#files-search
func (p *Putio) FilesSearch(query string, pageno string) (files *Files, jsonstr string, err error) {
	return p.GetFilesReq("files/search/" + query + "/page/" + string(pageno))
}

// https://api.put.io/v2/docs/#files-create-folder
func (p *Putio) FilesCreateFolder(name string, parent_id NInt64) (files *Files, jsonstr string, err error) {
	data := make(url.Values)
	data.Set("name", name)
	data.Set("parent_id", string(parent_id))
	return p.PostFilesReq("files/create-folder", data)
}

// https://api.put.io/v2/docs/#files-id
func (p *Putio) FilesId(id NInt64) (files *Files, jsonstr string, err error) {
	return p.GetFilesReq("files/" + strconv.FormatInt(int64(id), 10))
}

// https://api.put.io/v2/docs/#files-delete
func (p *Putio) FilesDelete(file_id NInt64) (files *Files, jsonstr string, err error) {
	return p.PostFilesReq("files/delete", url.Values{"file_ids": {strconv.FormatInt(int64(file_id), 10)}})
}

// https://api.put.io/v2/docs/#files-rename
func (p *Putio) FilesRename(file_id NInt64, name string) (files *Files, jsonstr string, err error) {
	return p.PostFilesReq("files/rename", url.Values{"file_id": {strconv.FormatInt(int64(file_id), 10)}, "name": {name}})
}

// https://api.put.io/v2/docs/#files-move
func (p *Putio) FilesMove(file_id NInt64, parent_id NInt64) (files *Files, jsonstr string, err error) {
	return p.PostFilesReq("files/move", url.Values{"file_id": {strconv.FormatInt(int64(file_id), 10)}, "parent_id": {strconv.FormatInt(int64(parent_id), 10)}})
}

// https://api.put.io/v2/docs/#files-mp4-post
func (p *Putio) FilesMP4(id NInt64) (files *Files, jsonstr string, err error) {
	return p.PostFilesReq("files/"+strconv.FormatInt(int64(id), 10)+"/mp4", url.Values{"id": {strconv.FormatInt(int64(id), 10)}})
}

// https://api.put.io/v2/docs/#files-mp4-post
func (p *Putio) FilesMP4Status(id NInt64) (files *Files, jsonstr string, err error) {
	return p.GetFilesReq("files/" + strconv.FormatInt(int64(id), 10) + "/mp4")
}

// https://api.put.io/v2/docs/#files-id-download
// in this case we will just return the url to download from and leave it up to 
// the client to actually download it. It's a redirect so can't use the usual request method
func (p *Putio) FilesDownload(id NInt64) (urlstr string, err error) {
	path := "download"

	url := BaseUrl + "files/" + strconv.FormatInt(int64(id), 10) + "/" + path + oauthparam + p.OauthToken
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	finalURL := resp.Request.URL.String()
	return finalURL, nil
}

func (p *Putio) PostTransfersReq(path string, data url.Values) (transfers *Transfers, jsonstr string, err error) {
	posturl := BaseUrl + path + oauthparam + p.OauthToken
	resp, err := http.PostForm(posturl, data)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	bodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	if err = json.Unmarshal(bodybytes, &transfers); err != nil {
		return nil, string(bodybytes), err
	}
	return transfers, string(bodybytes), nil
}

func (p *Putio) GetTransfersReq(path string) (transfers *Transfers, jsonstr string, err error) {
	bodybytes, err := p.GetReqBody(path)
	if err != nil {
		return nil, string(bodybytes), err
	}
	if err = json.Unmarshal(bodybytes, &transfers); err != nil {
		return nil, string(bodybytes), err
	}
	return transfers, string(bodybytes), nil
}

// https://api.put.io/v2/docs/#transfers-list
func (p *Putio) TransfersList() (transfers *Transfers, jsonstr string, err error) {
	return p.GetTransfersReq("transfers/list")
}

// https://api.put.io/v2/docs/#transfers-add
func (p *Putio) TransfersAdd(transfer_url string, save_parent_id NInt64, extract bool) (transfers *Transfers, jsonstr string, err error) {
	return p.PostTransfersReq("transfers/add", url.Values{"url": {transfer_url}, "save_parent_id": {strconv.FormatInt(int64(save_parent_id), 10)}, "extract": {strconv.FormatBool(extract)}})
}

// https://api.put.io/v2/docs/#transfers-add
func (p *Putio) TransfersCancel(transfer_id NInt64) (transfers *Transfers, jsonstr string, err error) {
	return p.PostTransfersReq("transfers/cancel", url.Values{"transfer_ids": {strconv.FormatInt(int64(transfer_id), 10)}})
}

// https://api.put.io/v2/docs/#transfers-id
func (p *Putio) TransfersId(id NInt64) (transfers *Transfers, jsonstr string, err error) {
	return p.GetTransfersReq("transfers/" + strconv.FormatInt(int64(id), 10))
}

func (p *Putio) GetAccountReq(path string) (account *Account, jsonstr string, err error) {
	bodybytes, err := p.GetReqBody(path)
	if err != nil {
		return nil, string(bodybytes), err
	}
	if err = json.Unmarshal(bodybytes, &account); err != nil {
		return nil, string(bodybytes), err
	}
	return account, string(bodybytes), nil
}

// https://api.put.io/v2/docs/#account-info
func (p *Putio) AccountInfo() (account *Account, jsonstr string, err error) {
	return p.GetAccountReq("account/info")
}

// https://api.put.io/v2/docs/#account-settings
func (p *Putio) AccountSettings() (account *Account, jsonstr string, err error) {
	return p.GetAccountReq("account/settings")
}

func (p *Putio) PostFriendsReq(path string, data url.Values) (friends *Friends, jsonstr string, err error) {
	posturl := BaseUrl + path + oauthparam + p.OauthToken
	resp, err := http.PostForm(posturl, data)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	bodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	if err = json.Unmarshal(bodybytes, &friends); err != nil {
		return nil, string(bodybytes), err
	}
	return friends, string(bodybytes), nil
}

func (p *Putio) GetFriendReq(path string) (friends *Friends, jsonstr string, err error) {
	bodybytes, err := p.GetReqBody(path)
	if err != nil {
		return nil, string(bodybytes), err
	}
	if err = json.Unmarshal(bodybytes, &friends); err != nil {
		return nil, string(bodybytes), err
	}
	return friends, string(bodybytes), nil
}

// https://api.put.io/v2/docs/#friends-list
func (p *Putio) FriendsList() (friends *Friends, jsonstr string, err error) {
	return p.GetFriendReq("friends/list")
}

// https://api.put.io/v2/docs/#friends-username-request
func (p *Putio) FriendsRequest(username string) (friends *Friends, jsonstr string, err error) {
	return p.PostFriendsReq("friends/"+username+"/request", nil)
}

// https://api.put.io/v2/docs/#friends-username-deny
func (p *Putio) FriendsDeny(username string) (friends *Friends, jsonstr string, err error) {
	return p.PostFriendsReq("friends/"+username+"/deny", nil)
}

// https://api.put.io/v2/docs/#friends-waiting-requests
func (p *Putio) FriendsWaiting() (friends *Friends, jsonstr string, err error) {
	return p.GetFriendReq("friends/waiting-requests")
}

// NewPutio takes in the apps oauth information and gets the token that will be used for all other calls
// This function doesn't have to be used if you provied the OauthToken when creating a Putio struct.
func NewPutio(appid, appsecret, appredirect, usercode string) (*Putio, error) {
	// get the user token using the calling apps credentials
	url := "https://api.put.io/v2/oauth2/access_token?client_id=" + appid + "&client_secret=" + appsecret + "&grant_type=authorization_code&redirect_uri=" + appredirect + "&code=" + usercode
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
