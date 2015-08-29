package proxychecker

import (
	//	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type JsonProxy struct {
	Link            string
	Region          string
	GetMethod       bool
	PostMethod      bool
	AveResponseTime float64
	Stability       int
}

type jproxy []JsonProxy

//用于结构体排序的三个方法
func (u jproxy) Len() int {
	return len(u)
}

func (u jproxy) Less(i, j int) bool {
	return u[i].AveResponseTime < u[j].AveResponseTime
}

func (u jproxy) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func HandleProxyChecker(w http.ResponseWriter, r *http.Request) { //判断是否是高匿名代理
	r.ParseForm()
	username, err := r.Form["username"]
	if err == false {
		fmt.Fprint(w, "fail")
		return
	}
	// if r.Method == "Get" {
	// 	log.Println(r.Method, string(username[1])) //从这里开始
	// }
	if len(r.Header.Get("X-Forwarded-For")) == 0 && len(r.Header.Get("Via")) == 0 && strings.EqualFold(username[0], "credit") {
		fmt.Fprint(w, "ok")
	} else {
		fmt.Fprint(w, "fail")
	}
}

func HandleIntroduce(w http.ResponseWriter, r *http.Request) { //http服务，显示introduce.html页面
	t, err := template.ParseFiles("introduce.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}

func HandleGet(w http.ResponseWriter, r *http.Request) { //"/get"路径的http服务方法
	copyCP := make(map[string]*Proxy)
	PC.lock.Lock()
	copyCP = PC.checkedProxys
	PC.lock.Unlock()
	flag := true
	for link, p := range copyCP {
		if time.Now().Sub(p.LastCheckOKTime).Minutes() < 5 && p.AveResponseTime <= 1.0 && p.Stability == 1 {
			fmt.Fprint(w, link)
			flag = false
			break
		}
	}

	if flag {
		fmt.Fprint(w, "No high Anonymous Proxy!")
	}
}

func HandleAPI(w http.ResponseWriter, r *http.Request) { //"/api"路径的http服务方法
	params := r.URL.Query()
	number := params.Get("number")
	num := 0
	if len(number) == 0 {
		num = 20 //默认提取20个
	} else {
		num, _ = strconv.Atoi(number)
	}
	port := params.Get("port")
	removePort := params.Get("removePort")
	region := params.Get("region")
	removeRegion := params.Get("removeRegion")
	method := params.Get("method")
	stability := params.Get("stability")
	speed := params.Get("speed")
	ipstart := params.Get("ipStart")
	todayip := params.Get("todayIP")
	isort := params.Get("sort")
	// ?number=&port=&removePort=&region=&removeRegion=&method=&stability=&speed=&ipStart=&todayIP=&sort=

	copyCP := make(map[string]*Proxy)
	PC.lock.Lock()
	copyCP = PC.checkedProxys
	PC.lock.Unlock()
	if len(copyCP) == 0 {
		fmt.Fprint(w, "No high Anonymous Proxy!")
	} else {
		jsp := make([]JsonProxy, 0)
		k := 0
		for link, p := range copyCP {
			if UrlCondition(p, port, removePort, region, removeRegion, method, stability, speed, ipstart, todayip, isort) {
				continue
			}
			PC.lock.Lock()
			PC.checkedProxys[link].TodayIP = false
			PC.lock.Unlock()
			jp := JsonProxy{
				Link:            link,
				Region:          p.Region,
				GetMethod:       p.GetMethod,
				PostMethod:      p.PostMethod,
				AveResponseTime: p.AveResponseTime,
				Stability:       p.Stability,
			}
			jsp = append(jsp, jp)
			k += 1
			if k == num {
				break
			}
		}
		if strings.EqualFold(isort, "1") {
			sort.Sort(jproxy(jsp))
		}
		if len(jsp) == 0 {
			fmt.Fprint(w, "No high Anonymous Proxy!")
		} else {
			str := ""
			for i := 0; i < len(jsp); i++ {
				str = str + jsp[i].Link + "\n"
			}
			fmt.Fprint(w, str)
		}
		// bJsonProxy, JsonError := json.Marshal(jsp)
		// if JsonError != nil {
		// 	fmt.Fprint(w, "")
		// } else {
		// 	fmt.Fprint(w, string(bJsonProxy))
		// }
	}
}

func UrlCondition(p *Proxy, port, removePort, region, removeRegion, method, stability, speed, ipstart, todayip, isort string) bool {
	if p.LastCheckOK == false {
		return true
	}

	if len(port) != 0 {
		iport := strings.Split(port, ",")
		i := 0
		for ; i < len(iport); i++ {
			if strings.EqualFold(p.Port, iport[i]) {
				break
			}
		}
		if i == len(iport) {
			return true
		}
	}

	if len(removePort) != 0 {
		iremovePort := strings.Split(removePort, ",")
		i := 0
		for ; i < len(iremovePort); i++ {
			if strings.EqualFold(p.Port, iremovePort[i]) {
				return true
			}
		}
	}

	if len(region) != 0 {
		iregion := strings.Split(region, ",")
		i := 0
		for ; i < len(iregion); i++ {
			if strings.Contains(p.Region, iregion[i]) {
				break
			}
		}
		if i == len(iregion) {
			return true
		}
	}

	if len(removeRegion) != 0 {
		iremoveRegion := strings.Split(removeRegion, ",")
		for i := 0; i < len(iremoveRegion); i++ {
			if strings.Contains(p.Region, iremoveRegion[i]) {
				return true
			}
		}
	}

	if len(method) != 0 {
		iMethod := strings.Split(method, ",")
		iMethod[0] = strings.ToLower(iMethod[0])
		if len(iMethod) == 1 {
			if iMethod[0] == "get" {
				if p.GetMethod == false {
					return true
				}
			} else if iMethod[0] == "post" {
				if p.PostMethod == false {
					return true
				}
			}
		} else if len(iMethod) == 2 {
			iMethod[1] = strings.ToLower(iMethod[1])
			if iMethod[0] == "post" {
				iMethod[0], iMethod[1] = iMethod[1], iMethod[0]
			}
			if iMethod[0] == "get" && iMethod[1] == "post" {
				if !(p.GetMethod == true && p.PostMethod == true) {
					return true
				}
			}
		}
	}

	if len(stability) != 0 {
		istability := strings.Split(stability, ",")
		verySta, err := strconv.Atoi(istability[0])
		if err == nil {
			if len(istability) == 1 {
				if verySta == 1 {
					if p.Stability == 2 {
						return true
					}
				} else if verySta == 2 {
					if p.Stability == 1 {
						return true
					}
				}
			}
		}
	}

	if len(speed) != 0 {
		if len(speed) == 1 {
			if strings.Contains(speed, "1") {
				if p.AveResponseTime > 1.0 {
					return true
				}
			} else if strings.Contains(speed, "2") {
				if p.AveResponseTime < 1.0 || p.AveResponseTime > 3.0 {
					return true
				}
			} else if strings.Contains(speed, "3") {
				if p.AveResponseTime <= 3.0 {
					return true
				}
			}
		} else if len(speed) == 3 {
			if !strings.Contains(speed, "1") {
				if p.AveResponseTime < 1.0 {
					return true
				}
			} else if !strings.Contains(speed, "2") {
				if p.AveResponseTime >= 1.0 && p.AveResponseTime <= 3.0 {
					return true
				}
			} else if !strings.Contains(speed, "3") {
				if p.AveResponseTime > 3.0 {
					return true
				}
			}
		}
	}

	if len(ipstart) != 0 && !strings.HasPrefix(p.Link, "http://"+ipstart) {
		return true
	}
	if strings.EqualFold(todayip, "1") && p.TodayIP == false {
		return true
	}
	return false
}
