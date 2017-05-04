package main

import (
	"crypto/md5"
	"encoding/hex"

	"jdPriceShowWeb/model"
	"jdPriceShowWeb/redisDAO"
)

var writeLoop = true

var priceChan = make(chan string, 1024)
var shopInfoChan = make(chan string, 64)

// 格式 id:name
func loadShopId() map[string]string {
	return redisDAO.ReadShopIds()
}

// 格式 id:name
func loadTargetModel() map[string]*model.GoodPrices {
	ms := redisDAO.ReadModels()
	return redisDAO.ReadGoodPrices(ms)
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
