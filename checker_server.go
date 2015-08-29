package proxychecker

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const PROXY_CHECKER_SERVICE = "http://127.0.0.1/check"

func CheckAnonyGet(proxyAddr string) (bool, float64) { //验证Get方法
	start := time.Now()
	proxy, err1 := url.Parse(proxyAddr)
	if err1 != nil {
		// log.Println(err1)
		return false, 0
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(10 * time.Second)
				c, err2 := net.DialTimeout(netw, addr, 10*time.Second)
				if err2 != nil {
					return nil, err2
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: 10 * time.Second,
			DisableCompression:    false,
			Proxy:                 http.ProxyURL(proxy),
		},
	}
	respGet, err3 := client.Get(PROXY_CHECKER_SERVICE + "?username=credit") //验证Get方法,cugb
	if err3 != nil {
		// log.Println(err3)
		return false, 0
	}
	bodyGet, _ := ioutil.ReadAll(respGet.Body)
	defer respGet.Body.Close()
	return respGet.StatusCode == http.StatusOK && string(bodyGet) == "ok", time.Now().Sub(start).Seconds()
}

func CheckAnonyPost(proxyAddr string) (bool, float64) { //验证Post方法
	start := time.Now()
	proxy, err1 := url.Parse(proxyAddr)
	if err1 != nil {
		// log.Println(err1)
		return false, 0
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(10 * time.Second)
				c, err2 := net.DialTimeout(netw, addr, 10*time.Second)
				if err2 != nil {
					return nil, err2
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: 10 * time.Second,
			DisableCompression:    false,
			Proxy:                 http.ProxyURL(proxy),
		},
	}
	respPost, err4 := client.Post(PROXY_CHECKER_SERVICE, "application/x-www-form-urlencoded", strings.NewReader("username=credit")) //验证Post方法
	if err4 != nil {
		// log.Println(err4)
		return false, 0
	}
	bodyPost, _ := ioutil.ReadAll(respPost.Body)
	defer respPost.Body.Close()
	return respPost.StatusCode == http.StatusOK && string(bodyPost) == "ok", time.Now().Sub(start).Seconds()
}
