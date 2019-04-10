package util

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type TransmissionData struct {
	Method    string      `json:"method"`
	Arguments interface{} `json:"arguments"`
	Tas       string      `json:"tas"`
}

type TransmissionAdd struct {
	MetaInfo    string `json:"metainfo,omitempty"`
	Filename    string `json:"filename"`
	DownloadDir string `json:"download-dir,omitempty"`
	Paused      bool   `json:"paused"`
}

type TransmissionResult struct {
	Arguments interface{} `json:"arguments"`
	Result    string      `json:"result"`
}

type TransmissionAddResult struct {
	Result TransmissionResult
	Flag   string
}

const (
	footer     string = "/transmission/rpc"
	xSessionId string = "X-Transmission-Session-Id"
)

var (
	url, sessionId string
	isSessionId    bool
)

type TransmissionCallback func(body []byte, res *http.Response, result TransmissionResult, err error)
type TransmissionAddCallback func(body []byte, res *http.Response, result TransmissionAddResult, err error)

func TrExec(tr Client, data TransmissionData, callback TransmissionCallback) {
	url = tr.Local + footer
	jsonData, _ := json.Marshal(data)
	payload := strings.NewReader(string(jsonData))
	req, _ := http.NewRequest("POST", url, payload)
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(tr.User+":"+tr.Pwd))
	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/json")
	if isSessionId {
		req.Header.Add(xSessionId, sessionId)
	}
	res, err := http.DefaultClient.Do(req)
	var result TransmissionResult
	if err != nil {
		callback(nil, res, result, err)
	} else {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		code := res.StatusCode
		if code == 409 {
			isSessionId = true
			sessionId = res.Header.Get(xSessionId)
			TrExec(tr, data, callback)
		} else {
			if res.StatusCode == 200 {
				json.Unmarshal(body, &result)
			}
			callback(body, res, result, nil)
		}
	}
}

func TrAdd(tr Client, data TransmissionAdd, callback TransmissionAddCallback) {
	TrExec(tr, TransmissionData{"torrent-add", data, ""},
		func(body []byte, res *http.Response, result TransmissionResult, err error) {
			var addResult TransmissionAddResult
			if result.Result == "success" {
				args := result.Arguments.(map[string]interface{})
				if args["torrent-added"] != nil {
					result.Arguments = args["torrent-added"]
					addResult.Flag = "add"
				} else if args["torrent-duplicate"] != nil {
					result.Arguments = args["torrent-duplicate"]
					addResult.Flag = "duplicate"
				}
			}
			addResult.Result = result
			callback(body, res, addResult, err)
		})
}

func TrConnect(tr Client, callback TransmissionCallback) {
	TrExec(tr, TransmissionData{"session-stats", "", ""}, callback)
}
