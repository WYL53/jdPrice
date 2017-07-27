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

func TestHash(t *testing.T)  {
	err := hashSetValue("test","k1","v1","k2","v2")
	if err != nil{
		t.Error(err)
	}
	v1,err := hashGetValue("test","k1")
	if err != nil{
		t.Error(err)
	}
	if v1 != "v1"{
		t.Error("v1 != v1")
	}
	vs,err := hashGetValues("test")
	if err != nil{
		t.Error(err)
	}
	if vs["k2"].(string) != "v2"{
		t.Error("v2 != v2")
	}
}