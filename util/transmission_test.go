package util

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func TestAdd(t *testing.T) {
	tr := Client{"http://192.168.1.110" + ":" + strconv.Itoa(9091), "", ""}
	TrAdd(tr, TransmissionAdd{"",
		"https://totheglory.im/",
		"/volume2/pt/ttg", true},
		func(body []byte, res *http.Response, result TransmissionAddResult, err error) {
			fmt.Println(string(body))
			if result.Result.Arguments != nil {
				torrent := result.Result.Arguments.(map[string]interface{})
				fmt.Printf("hashString: %s tid: %.f name: %s\n", torrent["hashString"], torrent["id"], torrent["name"]);
			}
		})
}

func TestConnect(t *testing.T) {
	tr := Client{"http://192.168.1.110" + ":" + strconv.Itoa(9091), "zzyandzzy", "zzyianzhi1128"}
	TrConnect(tr, func(body []byte, res *http.Response, result TransmissionResult, err error) {
		fmt.Println(string(body))
		fmt.Println(result.Arguments.(map[string]interface{})["activeTorrentCount"])
	})
}
