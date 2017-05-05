package controller

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"jdPrice/model"
	"jdPrice/redisDAO"
	"jdPrice/view"
)

var BrandModelMap map[string][]string
var TargetModels map[string]*model.GoodPrices
var CurrentData map[string]map[string]*model.JdGood = make(map[string]map[string]*model.JdGood)

func setBrandModelMap(brand, model string) {
	if _, ok := BrandModelMap[brand]; !ok {
		models := make([]string, 1)
		models[0] = model
		BrandModelMap[brand] = models
		return
	}
	for _, v := range BrandModelMap[brand] {
		if v == model {
			return
		}
	}
	BrandModelMap[brand] = append(BrandModelMap[brand], model)
}

func delBrandModelMap(brand, model string) {
	models, ok := BrandModelMap[brand]
	if !ok {
		return
	}
	for i, m := range models {
		if m == model {
			newModels := make([]string, len(models)-1)
			copy(newModels, models[:i])
			copy(newModels[i:], models[i+1:])
			BrandModelMap[brand] = newModels
		}
	}
}

func StartHttpServer(prot int) {
	addr := fmt.Sprintf("0.0.0.0:%d", prot)
	fmt.Println("http server listen address:", addr)
	http.HandleFunc("/", IndexServer)
	http.HandleFunc("/addModel", AddModelServer)
	http.HandleFunc("/delModel", DelModelServer)
	http.HandleFunc("/updatePrice", UpdatePriceServer)
	http.HandleFunc("/jd", HomeServer)
	http.HandleFunc("/price", PriceServer)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("start http server error:%s", err.Error())
		os.Exit(1)
	}
}

func IndexServer(w http.ResponseWriter, req *http.Request) {
	t := template.New("template")       //创建一个模板
	t, _ = t.Parse(view.TPL_INDEX_PAGE) //解析模板文件
	i := 0
	prices := make([]*model.GoodPrices, len(TargetModels))
	for _, v := range TargetModels {
		prices[i] = v
		i++
	}
	data := struct {
		Selects map[string][]string
		Prices  []*model.GoodPrices
	}{
		Selects: BrandModelMap,
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
	if _, ok := TargetModels[modelName]; !ok {
		setBrandModelMap(brand, modelName)
		TargetModels[modelName] = model.NewGoodPrices2(modelName, standardPrice, minPrice, maxPrice)
		err := redisDAO.WriteGoodPrice(modelName, standardPrice, minPrice, maxPrice)
		if err != nil {
			fmt.Println(err)
		} else {
			err = redisDAO.WriteModel(brand, modelName)
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
	if _, ok := TargetModels[modelName]; ok {
		delete(TargetModels, modelName)
		delBrandModelMap(brand, modelName)
		redisDAO.RemoveModel(brand, modelName)
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
	if price, ok := TargetModels[modelName]; ok {
		sp, err1 := strconv.Atoi(standardPrice)
		minp, err2 := strconv.Atoi(standardPrice)
		maxp, err3 := strconv.Atoi(standardPrice)
		if err1 == nil && err2 == nil && err3 == nil {
			if price.StandardPrice != sp || price.MinPrice != minp || price.MaxPrice != maxp {
				err := redisDAO.WriteGoodPrice(modelName, standardPrice, minPrice, maxPrice)
				if err == nil {
					price.StandardPrice = sp
					price.MinPrice = minp
					price.MaxPrice = maxp
					w.Write([]byte("价格更新成功"))
					return
				}
				fmt.Println("价格更新失败", err)
				w.Write([]byte("价格更新失败"))
				return
			}

		}
		w.Write([]byte("价格参数有问题"))
		return
	}
	w.Write([]byte("型号不存在"))
	return
}

//价格显示
func HomeServer(w http.ResponseWriter, req *http.Request) {
	modelName := trimSpace(req.URL.Query().Get("model"))
	t := template.New("template")      //创建一个模板
	t, _ = t.Parse(view.TPL_SHOW_PAGE) //解析模板文件
	data := struct {
		Model         string
		Goods         map[string]*model.JdGood
		StandardPrice int
	}{
		Model:         modelName,
		Goods:         CurrentData[modelName],
		StandardPrice: int(TargetModels[modelName].StandardPrice),
	}
	t.Execute(w, data) //执行模板的merger操作
}

//价格走势
func PriceServer(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	//	model := req.URL.Query().Get("model")
	//	fmt.Println(model)

	thisMonthPrefix := req.URL.Query().Get("month")
	if len(thisMonthPrefix) < 7 {
		thisMonthPrefix = time.Now().Format("2006-01")
	}
	lastPrice := 0
	var lastTime string
	prices := make([]int, 0)
	times := make([]string, 0)
	contents := redisDAO.ReadPrice(id)
	//	fmt.Println(contents)
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
	//	fmt.Println(times, prices)
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
