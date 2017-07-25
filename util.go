package main

import (
	"crypto/md5"
	"encoding/hex"

	"jdPrice/model"
	"jdPrice/redisDAO"
)

var writeLoop = true

var priceChan = make(chan string, 1024)
var shopInfoChan = make(chan string, 64)

// 格式 id:name
func loadShopId() map[string]string {
	return redisDAO.ReadShopIds()
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

func getMd5(s string) string {
	md5Byte16 := md5.Sum([]byte(s))
	return hex.EncodeToString(md5Byte16[:])
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
