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
	Tag       string      `json:"tag"`
}

// 处理tr add结构体
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
	trFooter              string = "/transmission/rpc"
	trSessionIdHeader     string = "X-Transmission-Session-Id"
	authorizationHeader          = "Authorization"
	contentTypeHeader            = "Content-Type"
	jsonHeader        = "application/json"
)

var (
	trSessionId string
)

type TransmissionCallback func(body []byte, res *http.Response, result TransmissionResult, err error)
type TransmissionAddCallback func(body []byte, res *http.Response, result TransmissionAddResult, err error)

func TrExec(tr Client, data TransmissionData, callback TransmissionCallback) {
	trUrl := tr.Local + trFooter
	jsonData, _ := json.Marshal(data)
	payload := strings.NewReader(string(jsonData))
	req, _ := http.NewRequest("POST", trUrl, payload)
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(tr.User+":"+tr.Pwd))
	req.Header.Add(authorizationHeader, auth)
	req.Header.Add(contentTypeHeader, jsonHeader)
	if trSessionId != "" {
		req.Header.Add(trSessionIdHeader, trSessionId)
	}
	res, err := http.DefaultClient.Do(req)
	var result TransmissionResult
	if err != nil {
		callback(nil, res, result, err)
	} else {
		defer res.Body.Close()
		code := res.StatusCode
		if code == 409 {
			trSessionId = res.Header.Get(trSessionIdHeader)
			TrExec(tr, data, callback)
		} else {
			body, _ := ioutil.ReadAll(res.Body)
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
