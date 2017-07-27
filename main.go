package main

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
	"time"
	"os"
	"os/signal"
	"syscall"

	"jdPrice/controller"
	"jdPrice/model"
	"jdPrice/mediator"
	"jdPrice/log"
	"fmt"
)



func main() {
	defer log.Clear()
	go loop(int(conf.FrequencyOfDay))
	//	log.Println("start http server")
	go controller.StartHttpServer(int(conf.Port))
	exitC := make(chan os.Signal,1)
	signal.Notify(exitC,syscall.SIGINT,syscall.SIGTERM)
	for  {
		select {
		case <-exitC:
			log.Println("app exit!")
			os.Exit(0)
		}
	}
}

func loop(frequencyOfDay int) {
	//ips,_ := mediator.GetIpPool()
	for model, prices := range mediator.CopyModelsStandardPrice() {
		startMission(model, prices.StandardPrice, prices.MinPrice, prices.MaxPrice,"")
	}
	ticker := time.NewTicker(time.Hour * time.Duration(24) / time.Duration(frequencyOfDay))
	for  range ticker.C {
		for model, prices := range mediator.CopyModelsStandardPrice() {
			startMission(model, prices.StandardPrice, prices.MinPrice, prices.MaxPrice,"")
		}
	}
}

//获取一次数据
func startMission(modelName string, standardPirce int, priceMin int, priceMax int,ip string) {
	log.Printf("%s开始获取数据。\n", modelName)
	var cmdResp *bytes.Buffer
	if ip != ""{
		proxy := fmt.Sprintf("--proxy=%s",ip)
		cmdResp = execCmd("phantomjs", "jd_spider.js", proxy,modelName)
	}else {
		cmdResp = execCmd("phantomjs", "jd_spider.js", modelName)
	}

	if cmdResp == nil{
		log.Printf("%s获取数据超时。\n", modelName)
		return
	}
	data := formatData(cmdResp, standardPirce, priceMin, priceMax)
	//	log.Println(data)
	if data != nil {
		mediator.UpdateCurrentData(modelName,data)
	}
	log.Printf("%s获取数据完毕。\n", modelName)
}

// 店名:商品名:价格评论:其他:店铺href：商品href
//商品href 作为返回的map的key
func formatData(data *bytes.Buffer, standardPrice, priceMin int, priceMax int) map[string]*model.JdGood {
	m := make(map[string]*model.JdGood)
	reader := bufio.NewReader(data)
	for {
		line, _, err := reader.ReadLine()
		if line == nil || err != nil {
			break
		}
		ss := strings.Split(string(line), ":")
		if len(ss) < 6 {
			log.Println("line error :",string(line))
			continue
		}

		priceAndSales := strings.Split(ss[2], ".00")
		priceString := priceAndSales[0]
		priceStartIndex := strings.IndexFunc(priceString, func(r rune) bool {
			if r >= '0' && r <= '9' {
				return true
			}
			return false
		})
		p, _ := strconv.Atoi(priceString[priceStartIndex:len(priceString)])
		if !isRange(p, priceMin, priceMax) {
			continue
		}
		sales := priceAndSales[1]

		shopName := ss[0]
		shopHref := ss[4]
		if shopHref != "undefined" {
			cacheShopName := mediator.GetShopName(shopHref)
			if cacheShopName == ""{
				mediator.UpdateShopName(shopHref,shopName)
			}else {
				if cacheShopName != shopName {
					if cacheShopName == "京东自营" {
						mediator.UpdateShopName(shopHref, shopName)
					} else {
						shopName = cacheShopName
					}
				}
			}
		}
		if shopName == "京东自营" {
			mediator.GetShopName()
		}
		newGood := model.NewJdGood(shopName, ss[1], p, sales, ss[3], ss[5])
		newGood.SetPriceDiff(standardPrice, p)
		m[ss[5]] = newGood
	}
	return m
}


func isRange(n, min, max int) bool {
	if n < min {
		return false
	}
	if n > max {
		return false
	}
	return true
}
