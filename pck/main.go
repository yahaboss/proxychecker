package main

import (
	"git.bdp.cc/termite/proxychecker"
	"net/http"
)

func main() {
	proxychecker.PC = proxychecker.NewProxyChecker()
	go proxychecker.Run()

	http.HandleFunc("/check", proxychecker.HandleProxyChecker)
	http.HandleFunc("/get", proxychecker.HandleGet) //响应速度在1.5s以内，且存在时间在30分钟以内定为非常稳定的代理IP
	http.HandleFunc("/api", proxychecker.HandleAPI)
	http.HandleFunc("/introduce", proxychecker.HandleIntroduce)
	// ?number=&port=&removePort=&region=&removeRegion=&method=&stability=&speed=&ipStart=&todayIP=&sort=
	http.ListenAndServe(":29840", nil)
}
