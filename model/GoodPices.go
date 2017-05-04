package model

import (
	"strconv"
)

type GoodPrices struct {
	Name          string
	StandardPrice int
	MinPrice      int
	MaxPrice      int
}

func NewGoodPrices(n string, sp int, minP int, maxP int) *GoodPrices {
	return &GoodPrices{
		Name:          n,
		StandardPrice: sp,
		MinPrice:      minP,
		MaxPrice:      maxP,
	}
}

func NewGoodPrices2(n, sp, minP, maxP string) *GoodPrices {
	standardPrice, _ := strconv.Atoi(sp)
	minPrice, _ := strconv.Atoi(minP)
	maxPrice, _ := strconv.Atoi(maxP)
	return &GoodPrices{
		Name:          n,
		StandardPrice: standardPrice,
		MinPrice:      minPrice,
		MaxPrice:      maxPrice,
	}
}
