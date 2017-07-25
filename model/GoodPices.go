package model

import (
	"strconv"
)

//型号的参考价
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

func CopyGoodPrices(src *GoodPrices) *GoodPrices {
	return &GoodPrices{
		Name:src.Name,
		StandardPrice:src.StandardPrice,
		MinPrice:src.MinPrice,
		MaxPrice:src.MaxPrice,
	}
}