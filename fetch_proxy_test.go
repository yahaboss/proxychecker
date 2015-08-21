package proxychecker

import (
	"fmt"
	"testing"
)

func TestKuaiDaili(t *testing.T) {
	body, goqueryError := fetchProxyPage("http://www.kuaidaili.com/proxylist/1")
	if goqueryError != nil {
		t.Error(goqueryError)
		return
	}

	proxys, fetchError := parseKuaiDaili(body)
	if fetchError != nil {
		t.Error(fetchError)
		return
	}
	fmt.Println(proxys)
	t.Log(proxys)
}

func TestNianshao(t *testing.T) {
	body, goqueryError := fetchProxyPage("http://www.nianshao.me/?page=1")
	if goqueryError != nil {
		t.Error(goqueryError)
		return
	}

	proxys, fetchError := nianshao(body)
	if fetchError != nil {
		t.Error(fetchError)
		return
	}
	fmt.Println(proxys)
	t.Log(proxys)
}

func TestHaodailiip(t *testing.T) {
	body, goqueryError := fetchProxyPage("http://www.haodailiip.com/other/1/1")
	if goqueryError != nil {
		t.Error(goqueryError)
		return
	}

	proxys, fetchError := haodailiip(body)
	if fetchError != nil {
		t.Error(fetchError)
		return
	}
	fmt.Println(proxys)
	t.Log(proxys)
}

func TestKuaidailiInha(t *testing.T) {
	body, goqueryError := fetchProxyPage("http://www.kuaidaili.com/free/inha/1")
	if goqueryError != nil {
		t.Error(goqueryError)
		return
	}

	proxys, fetchError := parseKuaiDailiInha(body)
	if fetchError != nil {
		t.Error(fetchError)
		return
	}
	fmt.Println(proxys)
	t.Log(proxys)
}

func TestKuaidailiOutha(t *testing.T) {
	body, goqueryError := fetchProxyPage("http://www.kuaidaili.com/free/outha/1")
	if goqueryError != nil {
		t.Error(goqueryError)
		return
	}

	proxys, fetchError := parseKuaiDailiOutha(body)
	if fetchError != nil {
		t.Error(fetchError)
		return
	}
	fmt.Println(proxys)
	t.Log(proxys)
}

func TestXicinn(t *testing.T) {
	body, goqueryError := fetchProxyPage("http://www.xici.net.co/nn/")
	if goqueryError != nil {
		t.Error(goqueryError)
		return
	}

	proxys, fetchError := xicinn(body)
	if fetchError != nil {
		t.Error(fetchError)
		return
	}
	fmt.Println(proxys)
	t.Log(proxys)
}

func TestXiciwn(t *testing.T) {
	body, goqueryError := fetchProxyPage("http://www.xici.net.co/wn/")
	if goqueryError != nil {
		t.Error(goqueryError)
		return
	}

	proxys, fetchError := xiciwn(body)
	if fetchError != nil {
		t.Error(fetchError)
		return
	}
	fmt.Println(proxys)
	t.Log(proxys)
}

func TestBaizhongsou(t *testing.T) {
	body, goqueryError := fetchProxyPage("http://ip.baizhongsou.com/")
	if goqueryError != nil {
		t.Error(goqueryError)
		return
	}

	proxys, fetchError := baizhongsou(body)
	if fetchError != nil {
		t.Error(fetchError)
		return
	}
	fmt.Println(proxys)
	t.Log(proxys)
}
