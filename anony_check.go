package proxychecker

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var fetchUrl []string
var page [6]int

type Proxy struct {
	Link              string       //代理ip:port
	Port              string       //端口
	Region            string       //IP所在国家城市地址
	GetMethod         bool         //是否支持Get方法
	PostMethod        bool         //是否支持Post方法
	LastCheckOK       bool         //最后一次验证的结果
	FirstCheckOKTime  time.Time    //第一次检测的时间
	LastCheckOKTime   time.Time    //最后一次验证的时间
	TotalTimeout      float64      //总的已经不是高匿名IP的持续时间
	LastResponseSpeed float64      //最后一次的响应速度
	TotalSeconds      float64      //总的验证时间
	TotalTimes        int          //总的验证次数
	AveResponseTime   float64      //平均响应时间
	Stability         int          //稳定性，1为非常稳定，2为稳定;稳定性由（最后一次检测时间-第一次检测时间）及平均响应时间 确定
	TodayIP           bool         //今天是否已经取过
	Ticker            *time.Ticker //自己检测自己的时间间隔
}

type ProxyChecker struct {
	checkedProxys map[string]*Proxy
	lock          *sync.RWMutex
}

var PC *ProxyChecker

func NewProxyChecker() *ProxyChecker {
	ret := &ProxyChecker{
		checkedProxys: make(map[string]*Proxy),
		lock:          &sync.RWMutex{},
	}
	return ret
}

func check(link string) {
	getOK, getSecs := CheckAnonyGet(link) //判断是否是代理ip
	postOK, postSecs := CheckAnonyPost(link)
	if getOK || postOK {
		log.Println("check ok: ", link)
		ipPort := strings.Split(link, ":")
		var sec float64
		if getOK {
			sec = getSecs
		} else {
			sec = postSecs
		}
		url := "http://www.ip138.com/ips1388.asp?ip=" + ipPort[1][2:] + "&action=2"
		ipRegion := GetProxyRegion(url)
		p := &Proxy{
			Link:              link,
			Port:              ipPort[2],
			Region:            ipRegion,
			GetMethod:         getOK,
			PostMethod:        postOK,
			LastCheckOK:       true,
			FirstCheckOKTime:  time.Now(),
			LastCheckOKTime:   time.Now(),
			TotalTimeout:      0.5,
			LastResponseSpeed: sec,
			TotalSeconds:      sec,
			TotalTimes:        1,
			AveResponseTime:   sec,
			Stability:         2,
			TodayIP:           true,
			Ticker:            time.NewTicker(time.Minute * 3), //每5分钟重新检测此IP
		}
		PC.lock.Lock()
		PC.checkedProxys[link] = p
		PC.lock.Unlock() //解锁

		go func() { //每个是高匿名IP之后，自己新建goroutine检测自己
			for t := range p.Ticker.C {
				log.Println(t, link)
				aliveGetOK, aliveGet := CheckAnonyGet(link)
				alivePostOK, alivePost := CheckAnonyPost(link)
				var isecs float64
				if aliveGetOK {
					isecs = aliveGet
				} else {
					isecs = alivePost
				}

				if aliveGetOK || alivePostOK {
					PC.lock.Lock()
					p.GetMethod = aliveGetOK
					p.PostMethod = alivePostOK
					if p.LastCheckOK == false {
						p.FirstCheckOKTime = time.Now()
						p.TotalSeconds = 0.0
						p.TotalTimes = 0
					}
					p.LastCheckOK = true
					p.LastCheckOKTime = time.Now()
					p.LastResponseSpeed = isecs
					p.TotalSeconds = p.TotalSeconds + isecs
					p.TotalTimes = p.TotalTimes + 1
					p.AveResponseTime = p.TotalSeconds / float64(p.TotalTimes)
					if p.LastCheckOKTime.Sub(p.FirstCheckOKTime).Minutes() >= 25 && p.AveResponseTime <= 2.0 {
						p.Stability = 1
					} else {
						p.Stability = 2
					}
					if p.TodayIP == false && p.LastCheckOKTime.Sub(p.FirstCheckOKTime).Hours() > 24 {
						p.TodayIP = true
					}
					PC.lock.Unlock()
				} else {
					PC.lock.Lock()
					p.LastCheckOK = false
					p.TotalTimeout = p.TotalTimeout + 10 //每次检测不是高匿名的话加10s
					if p.TotalTimeout >= 250 {           //如果超过2小时仍然不是高匿名IP，则删除
						p.Ticker.Stop()
						delete(PC.checkedProxys, link)
						PC.lock.Unlock()
						break
					}
					PC.lock.Unlock()
				}
			}
		}()
	}
}

func crawler_CheckAll(i, j int) {
	proxys := make(map[string]struct{})
	var body *goquery.Document
	var goqueryError error
	if j == 0 {
		body, goqueryError = fetchProxyPage(fetchUrl[i])
	} else {
		body, goqueryError = fetchProxyPage(fetchUrl[i] + strconv.Itoa(j))
	}
	if goqueryError != nil {
		return
	}
	ps := make(map[string]struct{}) //模拟set，因为struct{}作为value几乎不占内存
	var fetchError error
	fetchError = nil
	switch i {
	case 0:
		ps, fetchError = parseKuaiDaili(body) //把这里变了就行
	case 1:
		ps, fetchError = nianshao(body)
	case 2:
		ps, fetchError = haodailiip(body)
	case 3:
		ps, fetchError = parseKuaiDailiInha(body)
	case 4:
		ps, fetchError = parseKuaiDailiOutha(body)
	case 5:
		ps, fetchError = xicinn(body)
	case 6:
		ps, fetchError = xiciwn(body)
	default:
		ps, fetchError = baizhongsou(body)
	}

	if fetchError != nil {
		return
	}
	copyCP := make(map[string]*Proxy)
	PC.lock.Lock()
	copyCP = PC.checkedProxys
	PC.lock.Unlock()
	for p, _ := range ps {
		if _, ok := copyCP[p]; ok == false {
			proxys[p] = struct{}{}
		}
	}

	for link, _ := range proxys {
		go check(link) //go
	}
}

func fetchAllProxys() {
	for i := 0; i <= 7; i++ {
		// log.Printf("第 %d 个网址开始运行……", i)
		if i <= 4 {
			pageNumber := page[i]
			for j := 1; j <= pageNumber; j++ {
				go crawler_CheckAll(i, j) //爬取每个网页独立协程
				// log.Printf("第 %d 页爬取完毕！", j)
			}
		} else {
			go crawler_CheckAll(i, 0)
		}
	}
}

func Run() {
	fetchUrl = make([]string, 10)
	fetchUrl[0] = "http://www.kuaidaili.com/proxylist/"  //扫描10页
	fetchUrl[1] = "http://www.nianshao.me/?page="        //扫描15页
	fetchUrl[2] = "http://www.haodailiip.com/other/1/"   //扫描20页
	fetchUrl[3] = "http://www.kuaidaili.com/free/inha/"  //扫描2页
	fetchUrl[4] = "http://www.kuaidaili.com/free/outha/" //扫描2页
	fetchUrl[5] = "http://www.xici.net.co/nn/"           //扫描当前页
	fetchUrl[6] = "http://www.xici.net.co/wn/"           //扫描当前页
	fetchUrl[7] = "http://ip.baizhongsou.com/"           //扫描当前页
	page[0] = 10                                         //首次运行链接可扫描的页数可以多点
	page[1] = 40
	page[2] = 30
	page[3] = 20
	page[4] = 20
	fetchAllProxys() //第一次运行
	page[1] = 15     //重复刷新的链接可扫描的页数少点
	page[2] = 20
	page[3] = 2
	page[4] = 2

	timer := time.NewTicker(5 * time.Minute)
	for t := range timer.C {
		log.Println(t)
		go fetchAllProxys() //go
	}
}
