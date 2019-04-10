package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type QBAddUrlData struct {
	Urls                               string
	AutoTMM                            bool
	SavePath, Cookie, Rename, Category string
	Paused, Skip_checking, Root_folder bool
}

const (
	qbFooter        = "/api/v2/auth/login"
	qbAddApi        = "/api/v2/torrents/add"
	cookieSetHeader = "Set-Cookie"
	cookieHeader    = "Cookie"
	fromHeader      = "application/x-www-form-urlencoded"
	multipartHeader = "multipart/form-data"
)

//200
//Fails.
//Ok.
func QbGetCookie(qb Client) (string, error) {
	qbUrl := qb.Local + qbFooter
	client := &http.Client{}
	v := url.Values{}
	v.Add("username", qb.User)
	v.Add("password", qb.Pwd)
	u := ioutil.NopCloser(strings.NewReader(v.Encode()))
	req, _ := http.NewRequest("POST", qbUrl, u)
	req.Header.Set(contentTypeHeader, fromHeader)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	} else {
		defer res.Body.Close()
		code := res.StatusCode
		// 用户的IP因许多失败的登录尝试而被禁止
		if code == 403 {
			return "403", nil
		} else {
			body, _ := ioutil.ReadAll(res.Body)
			// 登录成功
			if string(body) == "Ok." {
				return strings.Split(res.Header.Get(cookieSetHeader), ";")[0], nil
			}
			return "", nil
		}
	}
}

//urls url地址
//autoTMM: 管理模式
//savepath: 保存路径
//cookie:
//rename: 重命名
//category: 分类
//paused: 是否暂停
//skip_checking: 是否跳过检查
//root_folder: 是否创建子文件夹
//dlLimit: 限制下载速度 B/s
//upLimit: 限制上传速度 b/S
func QbAddFromUrl(baseUrl string, cookie string, data QBAddUrlData) {
	qbUrl := baseUrl + qbAddApi
	client := &http.Client{}
	var bufReader bytes.Buffer
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("urls", data.Urls)
	if data.AutoTMM {
		mpWriter.WriteField("autoTMM", "true")
	} else {
		mpWriter.WriteField("autoTMM", "false")
	}
	mpWriter.WriteField("savepath", data.SavePath)
	mpWriter.WriteField("category", data.Category)
	if data.Paused {
		mpWriter.WriteField("paused", "true")
	} else {
		mpWriter.WriteField("paused", "false")
	}
	if data.Skip_checking {
		mpWriter.WriteField("skip_checking", "true")
	} else {
		mpWriter.WriteField("skip_checking", "false")
	}
	if data.Root_folder {
		mpWriter.WriteField("root_folder", "true")
	} else {
		mpWriter.WriteField("root_folder", "false")
	}
	mpWriter.Close()
	fmt.Println(bufReader.String())

	req, _ := http.NewRequest("POST", qbUrl, strings.NewReader(bufReader.String()))
	req.Header.Set(contentTypeHeader, multipartHeader)
	req.Header.Set(cookieHeader, cookie)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	} else {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(res.Header.Get(cookieSetHeader))
		fmt.Println(res.StatusCode)
		fmt.Println(string(body))
	}
}