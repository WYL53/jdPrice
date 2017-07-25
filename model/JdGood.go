package model

import (
	"fmt"
)

type JdGood struct {
	ShopName  string
	Name      string
	GoodHref  string
	Price     int
	PriceDiff int
	Sales     string
	Etc       string //优惠券、自营、满减等等
	UpOrDown  int
}

func NewJdGood(shopName string, name string, price int, sales string, etc string, goodHref string) *JdGood {
	return &JdGood{
		Name:     name,
		ShopName: shopName,
		Price:    price,
		Sales:    sales,
		Etc:      etc,
		GoodHref: fmt.Sprintf("http:%s", goodHref),
	}
}

func CopyJdGood(src *JdGood) *JdGood {
	return &JdGood{
		ShopName:src.ShopName,
		Name:src.Name,
		GoodHref:src.GoodHref,
		Price:src.Price,
		PriceDiff:src.PriceDiff,
		Sales:src.Sales,
		Etc:src.Etc,
		UpOrDown:src.UpOrDown,
	}

}

func (good JdGood) ThanOther(other *JdGood) (diff int) {
	return (good.Price - other.Price) * 100 / other.Price
}

func (good *JdGood) SetPriceDiff(standPrice, price int) {
	good.PriceDiff = price - standPrice
}
