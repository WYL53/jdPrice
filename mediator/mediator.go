package mediator

import (
	"sync"

	"jdPrice/model"
	"jdPrice/redisDAO"
)

//店铺id -> 店铺名
var shopId2Name map[string]string
var shopId2NameLock sync.RWMutex

//品牌 -> 型号
var brand2Model map[string][]string
var brand2ModelLock sync.RWMutex

//型号的参考价
var modelsStandardPrice map[string]*model.GoodPrices
var modelsStandardPriceLock sync.RWMutex

//最新的价格数据
var currentData map[string]map[string]*model.JdGood = make(map[string]map[string]*model.JdGood)
var currentDataLock sync.RWMutex

//品牌
func UpdateBrand2Model(m map[string][]string)  {
	brand2ModelLock.Lock()
	defer  brand2ModelLock.Unlock()
	brand2Model = m
}

func CopyBrandModel()  map[string][]string {
	brand2ModelLock.Lock()
	defer brand2ModelLock.Unlock()

	m := make(map[string][]string)
	for k,v := range brand2Model{
		models := make([]string,len(v))
		copy(models,v)
		m[k] = models
	}
	return m
}

func SetBrandModelItem(brand string,modelName, standardPrice, minPrice, maxPrice string) error {
	brand2ModelLock.Lock()
	defer brand2ModelLock.Unlock()

	if _, ok := brand2Model[brand]; !ok {
		models := make([]string, 1)
		models[0] = modelName
		brand2Model[brand] = models
		err := redisDAO.WriteModel(brand,modelName)
		if err != nil{
			return err
		}
		return SetModelsStandardPrice(modelName, standardPrice, minPrice, maxPrice)
	}
	for _, v := range brand2Model[brand] {
		if v == modelName {
			return nil
		}
	}
	brand2Model[brand] = append(brand2Model[brand], modelName)
	return SetModelsStandardPrice(modelName, standardPrice, minPrice, maxPrice)
}

func DelBrand2ModelItem(brand, modelName string) error {
	brand2ModelLock.Lock()
	defer brand2ModelLock.Unlock()

	models, ok := brand2Model[brand]
	if !ok {
		return nil
	}
	for i, m := range models {
		if m == modelName {
			newModels := make([]string, len(models)-1)
			copy(newModels, models[:i])
			copy(newModels[i:], models[i+1:])
			brand2Model[brand] = newModels
			delModelsStandardPrice(modelName)
			redisDAO.RemoveModel(brand, modelName)
		}
	}
	return nil
}
//
//func GetBrand2Model()  {
//	brand2ModelLock.RLock()
//	defer brand2ModelLock.RUnlock()
//}

func UpdateModelsStandardPrice(m map[string]*model.GoodPrices)  {
	modelsStandardPrice = m
}

func CopyModelsStandardPrice() map[string]*model.GoodPrices {
	modelsStandardPriceLock.Lock()
	defer modelsStandardPriceLock.Unlock()

	m := make(map[string]*model.GoodPrices)
	for k,v := range modelsStandardPrice{
		m[k] = model.CopyGoodPrices(v)
	}
	return m
}

func GetModelStandardPrice(modelName string) int {
	modelsStandardPriceLock.RLock()
	defer modelsStandardPriceLock.RUnlock()
	if prices,ok := modelsStandardPrice[modelName];ok{
		return int(prices.StandardPrice)
	}
	return -1
}


func SetModelsStandardPrice(modelName, standardPrice, minPrice, maxPrice string) error {
	modelsStandardPriceLock.Lock()
	modelsStandardPrice[modelName] = model.NewGoodPrices2(modelName, standardPrice, minPrice, maxPrice)
	modelsStandardPriceLock.Unlock()
	return redisDAO.WriteStandardPrice(modelName, standardPrice, minPrice, maxPrice)
}

func GetModelsStandardPrice(modelName string) *model.GoodPrices {
	modelsStandardPriceLock.RLock()
	defer modelsStandardPriceLock.RUnlock()
	if prices,ok := modelsStandardPrice[modelName];ok{
		return prices
	}
	return nil
}

func delModelsStandardPrice(modelName string) error {
	modelsStandardPriceLock.Lock()
	defer modelsStandardPriceLock.Unlock()
	delete(modelsStandardPrice,modelName)
	return nil
}

func UpdatShopId2Name(m map[string]string)  {
	shopId2NameLock.Lock()
	defer shopId2NameLock.Unlock()
	shopId2Name = m
}

func GetShopName(id string) string {
	shopId2NameLock.RLock()
	defer shopId2NameLock.RUnlock()

	if name,ok := shopId2Name[id];ok{
		return name
	}
	return ""
}


func UpdateShopName(id,name string)  {
	shopId2NameLock.Lock()
	shopId2Name[id] = name
	shopId2NameLock.Unlock()

	redisDAO.WiretShopId(id, name)
}

//当前价格
func UpdateCurrentData(modelName string,prices map[string]*model.JdGood)  {
	currentDataLock.Lock()
	currentData[modelName] = prices
	currentDataLock.Unlock()

	redisDAO.WritePrice(prices)
}

func CopyModelCurrentData(modelName string) map[string]*model.JdGood {
	currentDataLock.RLock()
	currentDataLock.RUnlock()
	if data,ok := currentData[modelName];ok{
		m := make(map[string]*model.JdGood)
		for k,v := range data{
			m[k] = model.CopyJdGood(v)
		}
		return m
	}
	return nil
}