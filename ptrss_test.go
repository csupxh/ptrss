package main

import (
	"fmt"
	"github.com/zzyandzzy/ptrss/util"
	"net/http"
	"testing"
)

func TestCmd(t *testing.T) {
	util.GetBody("https://totheglory.im/",
		func(channel util.RSSChannel) {
			fmt.Printf("title: %s language: %s\n", channel.Title, channel.Language)
			for i := 0; i < len(channel.Items); i++ {
				fmt.Println(channel.Items[i])
			}
		})
}

func AddClient() {
	local := "http://192.168.1.110:9091"
	user := ""
	pwd := ""
	tr := util.Client{local, user, pwd}
	util.TrConnect(tr, func(body []byte, res *http.Response, result util.TransmissionResult, err error) {
		if res.StatusCode == 200 {
			if result.Result == "success" {
				if InsertClient("tr", local, user, pwd) {
					fmt.Printf("add client: %s success\n", "tr")
				}
			}
		} else {
			fmt.Printf("add client: %s fail\n", "tr")
		}
	})
}
