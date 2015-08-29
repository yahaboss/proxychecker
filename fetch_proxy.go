package proxychecker

import (
	"github.com/PuerkitoBio/goquery"
	iconv "github.com/djimenez/iconv-go" //编码转换
	"net"
	"net/http"
	"strings"
	"time"
)

func GetProxyRegion(url string) string { //得到IP的地区，从另一个网站爬的
	c := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(time.Second * 10)
				c, err := net.DialTimeout(netw, addr, time.Second*10)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			ResponseHeaderTimeout: time.Second * 10,
		},
	}
	resp, getError := c.Get(url)
	body, goqueryError := goquery.NewDocumentFromResponse(resp)
	if getError != nil || goqueryError != nil {
		return ""
	}
	preRegion := body.Find(".ul1").Find("li").Eq(0).Text() //爬取的网站正好是gb2312格式的
	address, err11 := iconv.ConvertString(preRegion, "gb2312", "utf-8")
	if err11 != nil {
		return ""
	}
	region := strings.Split(address, "：")
	if strings.Contains(region[1], "省") || strings.Contains(region[1], "市") {
		region[1] = "中国" + region[1]
	}
	return region[1]
}

func fetchProxyPage(link string) (*goquery.Document, error) {
	c := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(time.Second * 10)
				c, err := net.DialTimeout(netw, addr, time.Second*10)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			ResponseHeaderTimeout: time.Second * 10,
		},
	}
	resp, getError := c.Get(link)
	if getError != nil {
		return nil, getError
	}
	body, goqueryError := goquery.NewDocumentFromResponse(resp)
	return body, goqueryError
}

func parseKuaiDaili(body *goquery.Document) (map[string]struct{}, error) {
	ret := make(map[string]struct{})
	body.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
		td := tr.Find("td")
		ip := td.Eq(0).Text()
		port := td.Eq(1).Text()
		anony := td.Eq(2).Text()
		if anony == "高匿名" {
			proxy := "http://" + ip + ":" + port
			ret[proxy] = struct{}{}
		}
	})
	return ret, nil
}

func nianshao(body *goquery.Document) (map[string]struct{}, error) {
	ret := make(map[string]struct{})
	body.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
		td := tr.Find("td")
		ip := td.Eq(0).Text()
		port := td.Eq(1).Text()
		proxy := "http://" + ip + ":" + port
		ret[proxy] = struct{}{}
	})
	return ret, nil
}

func haodailiip(body *goquery.Document) (map[string]struct{}, error) {
	ret := make(map[string]struct{})
	body.Find(".proxy_table tbody tr").Each(func(i int, tr *goquery.Selection) { //通过测试
		td := tr.Find("td")
		ip := strings.TrimSpace(td.Eq(0).Text())
		port := strings.TrimSpace(td.Eq(1).Text())
		anony := td.Eq(4).Text()
		if anony == "高匿" {
			proxy := "http://" + ip + ":" + port
			ret[proxy] = struct{}{}
		}
	})
	return ret, nil
}

func parseKuaiDailiInha(body *goquery.Document) (map[string]struct{}, error) {
	ret := make(map[string]struct{})
	body.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
		td := tr.Find("td")
		ip := td.Eq(0).Text()
		port := td.Eq(1).Text()
		proxy := "http://" + ip + ":" + port
		ret[proxy] = struct{}{}
	})
	return ret, nil
}

func parseKuaiDailiOutha(body *goquery.Document) (map[string]struct{}, error) {
	ret := make(map[string]struct{})
	body.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
		td := tr.Find("td")
		ip := td.Eq(0).Text()
		port := td.Eq(1).Text()
		proxy := "http://" + ip + ":" + port
		ret[proxy] = struct{}{}
	})
	return ret, nil
}

func xicinn(body *goquery.Document) (map[string]struct{}, error) {
	ret := make(map[string]struct{})
	body.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
		if i != 0 {
			td := tr.Find("td")
			ip := td.Eq(2).Text()
			port := td.Eq(3).Text()
			proxy := "http://" + ip + ":" + port
			ret[proxy] = struct{}{}
		}
	})
	return ret, nil
}

func xiciwn(body *goquery.Document) (map[string]struct{}, error) {
	ret := make(map[string]struct{})
	body.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
		if i != 0 {
			td := tr.Find("td")
			ip := td.Eq(2).Text()
			port := td.Eq(3).Text()
			proxy := "http://" + ip + ":" + port
			ret[proxy] = struct{}{}
		}
	})
	return ret, nil
}

func baizhongsou(body *goquery.Document) (map[string]struct{}, error) {
	ret := make(map[string]struct{})
	body.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
		if i%2 == 1 {
			td := tr.Find("td")
			proxy := "http://" + td.Eq(0).Text()
			ret[proxy] = struct{}{}
		}
	})
	return ret, nil
}
