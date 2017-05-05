package redisDAO

import (
	"testing"
)

func TestReadBrand(t *testing.T) {
	ss := ReadBrands()
	for _, v := range ss {
		t.Log(v)
	}
}
