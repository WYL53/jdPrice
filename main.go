package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"jdPrice/controller"
	"jdPrice/model"
	"jdPrice/redisDAO"
)

const (
	HISTORY_PATH      = "history"
	SHOP_INFO_PATH    = HISTORY_PATH + "/shopInfo.txt"
	TARGET_SHOPS_PATH = HISTORY_PATH + "/targetShops.txt"
)

//var lastData map[string]*JdGood

var shopIdData map[string]string

func main() {
	controller.BrandModelMap = loadBrands()
	controller.TargetModels = loadTargetModel(controller.BrandModelMap)
	shopIdData = loadShopId()
	fmt.Println("TargetModels:", controller.TargetModels)
	//	fmt.Println("shopIdData:", shopIdData)
	go loop(int(conf.FrequencyOfDay))
	//	fmt.Println("start http server")
	controller.StartHttpServer(int(conf.Port))
}

func loop(frequencyOfDay int) {
	for model, prices := range controller.TargetModels {
		//		fmt.Println(model, prices.StandardPrice, prices.MinPrice, prices.MaxPrice)
		startMession(model, prices.StandardPrice, prices.MinPrice, prices.MaxPrice)
	}
	ticlker := time.NewTicker(time.Hour * time.Duration(24) / time.Duration(frequencyOfDay))
	for _ = range ticlker.C {
		for model, prices := range controller.TargetModels {

			startMession(model, prices.StandardPrice, prices.MinPrice, prices.MaxPrice)
		}
	}
}

func startMession(model string, standardPirce int, priceMin int, priceMax int) {
	fmt.Printf("%s %s获取数据。\n", time.Now().Format("2006-01-02 15:04"), model)
	cmdResp := execCmd("phantomjs", "jd_spider.js", model)
	data := formatData(cmdResp, standardPirce, priceMin, priceMax)
	//	fmt.Println(data)
	if data != nil {
		controller.CurrentData[model] = data
		redisDAO.WritePrice(data)
	}
	fmt.Printf("%s %s获取数据完毕。\n", time.Now().Format("2006-01-02 15:04"), model)
}

// 店名:商品名:价格评论:其他:href
func formatData(data *bytes.Buffer, standardPirce, priceMin int, priceMax int) map[string]*model.JdGood {
	m := make(map[string]*model.JdGood)
	reader := bufio.NewReader(data)
	for {
		line, _, err := reader.ReadLine()
		if line == nil || err != nil {
			break
		}
		ss := strings.Split(string(line), ":")
		if len(ss) < 6 {
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
			shopIdStartIndex := strings.Index(shopHref, "-")
			shopIdLastIndex := strings.LastIndex(shopHref, ".")
			if shopIdStartIndex < shopIdLastIndex {
				shopId := shopHref[shopIdStartIndex+1 : shopIdLastIndex]
				if shopName != "京东自营" {
					if _, ok := shopIdData[shopId]; !ok {
						shopIdData[shopId] = shopName
						redisDAO.WiretShopId(shopId, shopName)
					}
				} else {
					shopName = shopIdData[shopId]
				}
			}

		}

		key := getMd5(shopName + ss[1])
		//		fmt.Println(shopName, ss[1], p, sales, ss[3], ss[5])
		newGood := model.NewJdGood(shopName, ss[1], p, sales, ss[3], ss[5])
		newGood.SetPriceDiff(standardPirce, p)
		m[key] = newGood
	}
	return m
}
