package util

import (
	"fmt"
	"testing"
)

func TestQbExec(t *testing.T) {
	client := Client{"http://192.168.1.110:8080", "zzyandzzy", "zzyianzhi1128"}
	cookie, err := QbGetCookie(client)
	if err != nil {
		fmt.Printf("无法连接: %s\n", err.Error())
	} else if cookie == "403" || cookie == "" {
		fmt.Println("账号或密码错误 或者 用户的IP因许多失败的登录尝试而被禁止")
	} else {
		fmt.Printf("登录成功: %s\n", cookie)
	}
}

func TestQbAddFromUrl(t *testing.T) {
	urls := "https://totheglory.im/rssdd.php?par=dnZ2Mzc1NjE1fHx8M2Q0MjliODBlZGJiNGUzZGVjZGExYzRiYWY5Yjc4ODh6eg==&ssl=yes"
	QbAddFromUrl("http://192.168.1.110:8080", "SID=/Xvzad1Erwnkb1vnI01FUTs5cP6vAvh2",
		QBAddUrlData{urls, false, "/pt/ttg/",
			"", "", "ttg", false, true, true,})
}