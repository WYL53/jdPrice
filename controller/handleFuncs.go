package controller

import (
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"jdPrice/model"
	"jdPrice/redisDAO"
	"jdPrice/view"
	"jdPrice/mediator"
	"jdPrice/log"
)

func IndexServer(w http.ResponseWriter, req *http.Request) {
	t := template.New("template")       //创建一个模板
	t, _ = t.Parse(view.TPL_INDEX_PAGE) //解析模板文件
	i := 0
	models := mediator.CopyModelsStandardPrice()
	prices := make([]*model.GoodPrices, len(models))
	for _, v := range models {
		prices[i] = v
		i++
	}
	data := struct {
		Selects map[string][]string
		Prices  []*model.GoodPrices
	}{
		Selects: mediator.CopyBrandModel(),
		Prices:  prices,
	}
	t.Execute(w, data) //执行模板的merger操作
}


func AddModelServer(w http.ResponseWriter, req *http.Request) {
	brand := trimSpace(req.FormValue("brand"))
	modelName := trimSpace(req.FormValue("modelName"))
	standardPrice := trimSpace(req.FormValue("standardPrice"))
	minPrice := trimSpace(req.FormValue("minPrice"))
	maxPrice := trimSpace(req.FormValue("maxPrice"))
	if prices := mediator.GetModelsStandardPrice(modelName); prices == nil {
		err:= mediator.SetBrandModelItem(brand,modelName,standardPrice, minPrice, maxPrice)
		if err != nil {
			log.Println(err)
		} else {
			if err == nil {
				w.Write([]byte("添加成功"))
			}
			return
		}
	}
	w.Write([]byte("添加失败"))
}

func DelModelServer(w http.ResponseWriter, req *http.Request) {
	brand := trimSpace(req.FormValue("brand"))
	modelName := trimSpace(req.FormValue("modelName"))
	if  prices := mediator.GetModelsStandardPrice(modelName); prices != nil {
		mediator.DelBrand2ModelItem(brand,modelName)
	}
	w.Write([]byte("删除成功"))
}

func UpdatePriceServer(w http.ResponseWriter, req *http.Request) {
	//	oldbrand := trimSpace(req.FormValue("oldbrand"))
	//	brand := trimSpace(req.FormValue("brand"))
	modelName := trimSpace(req.FormValue("modelName"))
	standardPrice := trimSpace(req.FormValue("standardPrice"))
	minPrice := trimSpace(req.FormValue("minPrice"))
	maxPrice := trimSpace(req.FormValue("maxPrice"))
	if price:= mediator.GetModelsStandardPrice(modelName); price != nil {
		sp, err1 := strconv.Atoi(standardPrice)
		minp, err2 := strconv.Atoi(standardPrice)
		maxp, err3 := strconv.Atoi(standardPrice)
		if err1 == nil && err2 == nil && err3 == nil {
			//if price.StandardPrice != sp || price.MinPrice != minp || price.MaxPrice != maxp {
				err := mediator.SetModelsStandardPrice(modelName, standardPrice, minPrice, maxPrice)
				if err == nil {
					price.StandardPrice = sp
					price.MinPrice = minp
					price.MaxPrice = maxp
					w.Write([]byte("价格更新成功"))
					return
				}
				log.Println("价格更新失败", err)
				w.Write([]byte("价格更新失败"))
				return
			//}

		}
		w.Write([]byte("价格参数有问题"))
		return
	}
	w.Write([]byte("型号不存在"))
	return
}

//价格显示
func modelPriceShow(w http.ResponseWriter, req *http.Request) {
	modelName := trimSpace(req.URL.Query().Get("model"))
	t := template.New("template")      //创建一个模板
	t, _ = t.Parse(view.TPL_SHOW_PAGE) //解析模板文件
	data := struct {
		Model         string
		Goods         map[string]*model.JdGood
		StandardPrice int
	}{
		Model:         modelName,
		Goods:         mediator.CopyModelCurrentData(modelName),
		StandardPrice: mediator.GetModelStandardPrice(modelName),
	}
	t.Execute(w, data) //执行模板的merger操作
}

//价格走势
func priceChange(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	thisMonthPrefix := req.URL.Query().Get("month")
	if len(thisMonthPrefix) < 7 {
		thisMonthPrefix = time.Now().Format("2006-01")
	}
	lastPrice := 0
	var lastTime string
	prices := make([]int, 0)
	times := make([]string, 0)
	contents := redisDAO.ReadPrice(id)
	//	log.Println(contents)
	for _, line := range contents {

		if !strings.HasPrefix(line, thisMonthPrefix) {
			continue
		}
		ss := strings.Split(line, "|")
		if len(ss) != 2 {
			continue
		}
		p, err := strconv.Atoi(ss[1])
		if err == nil {
			if p != lastPrice {
				times = append(times, ss[0])
				prices = append(prices, p)
				lastTime = ss[0]
				lastPrice = p
			}

		}
	}
	times = append(times, lastTime)
	prices = append(prices, lastPrice)
	//	log.Println(times, prices)
	t := template.New("template")            //创建一个模板
	t, _ = t.Parse(view.TPL_PRICE_LINE_PAGE) //解析模板文件
	data := struct {
		Times  []string
		Prices []int
	}{
		Times:  times,
		Prices: prices,
	}
	t.Execute(w, data) //执行模板的merger操作
}

func trimSpace(s string) string {
	return strings.TrimSpace(s)
}
