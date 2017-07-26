package mediator

import (
	//"crypto/md5"
	//"encoding/hex"

	"jdPrice/model"
	"jdPrice/redisDAO"
)

func init()  {
	brands := loadBrands()
	UpdateBrand2Model(brands)
	UpdateModelsStandardPrice(loadTargetModel(brands))
	UpdatShopId2Name(loadShopId())
}

// 根据id（href） -> 店名
func loadShopId() map[string]string {
	id2shopName,err:= redisDAO.ReadAllShopName()
	if err != nil{
		panic(err)
	}
	return id2shopName
}

func loadBrands() map[string][]string {
	brands := redisDAO.ReadBrands()
	return redisDAO.ReadModels(brands)
}

// 格式 id:name
func loadTargetModel(brands map[string][]string) map[string]*model.GoodPrices {
	ms := make([]string, 0)
	for _, v := range brands {
		ms = append(ms, v...)
	}
	return redisDAO.ReadStandardPrice(ms)
}

//func getMd5(s string) string {
//	md5Byte16 := md5.Sum([]byte(s))
//	return hex.EncodeToString(md5Byte16[:])
//}
