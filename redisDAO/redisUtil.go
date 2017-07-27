package redisDAO

import (
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"

	"jdPrice/model"
	"jdPrice/log"
)

const (
	REDIS_HOST                = ":6379"
	RedisTargetModelSet       = "targetModel_%s"
	RedisTargetModelSetPre    = "targetModel_"
	RedisShopIdTable          = "shopIdTable"
	RedisGoodPriceTableFormat = "goodPriceTable_%s"
	RedisPriceZSetFormat      = "prices_%s"

	//good id -> shop info,hashtable
	RedisShopInfoFormat = "shopInfo_%s"

	RedisIpPoolName = "ipPool"
)

var RedisClient *redis.Pool

func init() {
	RedisClient = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     1,
		MaxActive:   5,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", REDIS_HOST)
			if err != nil {
				return nil, err
			}
			// 选择db
			//			c.Do("SELECT", REDIS_DB)
			return c, nil
		},
	}
}

func WritePrice(m map[string]*model.JdGood) {
	conn := RedisClient.Get()
	defer conn.Close()
	timeStamp := time.Now().Format("2006-01-02 15:04")
	for key, good := range m {
		s := fmt.Sprintf("%s|%d", timeStamp, good.Price)
		key = fmt.Sprintf(RedisPriceZSetFormat, key)
		_, err := conn.Do("ZADD", key, 1, s)
		if err != nil {
			log.Println("ZADD err:", err.Error())
			continue
		}
	}
}

func ReadPrice(goodId string) []string {
	conn := RedisClient.Get()
	defer conn.Close()
	var ret []string
	key := fmt.Sprintf(RedisPriceZSetFormat, goodId)
	reply, err := conn.Do("ZRANGE", key, 1, -1)
	if err == nil {
		objs, ok := reply.([]interface{})
		if ok {
			ret = make([]string, len(objs))
			for i := range objs {
				ret[i] = string(objs[i].([]uint8))
			}
		}
	} else {
		log.Println(err)
	}
	return ret
}

//获取品牌，如美的、小天鹅
func ReadBrands() []string {
	parrent := fmt.Sprintf(RedisTargetModelSet, "*")
	ss := getKeys(parrent)
	brands := make([]string, len(ss))
	for i, v := range ss {
		brand := v[len(RedisTargetModelSetPre):]
		brands[i] = brand
	}
	return brands
}

//添加品牌、型号
func WriteModel(brand, model string) error {
	conn := RedisClient.Get()
	defer conn.Close()
	key := fmt.Sprintf(RedisTargetModelSet, brand)
	err := setAddValue(key, model)
	return err
}

//获取某个品牌下的所有型号
func ReadModels(brands []string) map[string][]string {
	conn := RedisClient.Get()
	defer conn.Close()
	m := make(map[string][]string)
	for _, brand := range brands {
		key := fmt.Sprintf(RedisTargetModelSet, brand)
		reply, err := conn.Do("SMEMBERS", key)
		if err != nil {
			log.Println(err)
			continue
		}
		objs, ok := reply.([]interface{})
		if ok {
			models := make([]string, len(objs))
			for i, v := range objs {
				models[i] = string(v.([]uint8))
			}
			m[brand] = models
		}
	}
	return m
}

//删除型号
func RemoveModel(brand, model string)  {
	conn := RedisClient.Get()
	defer conn.Close()
	key := fmt.Sprintf(RedisTargetModelSet, brand)
	conn.Do("SREM", key, model)

	//del good price
	key = fmt.Sprintf(RedisGoodPriceTableFormat, model)
	conn.Do("DEL", key)
}

//id -> 店名
func WiretShopId(id, shopName string) error {
	return hashSetValue(RedisShopIdTable, id, shopName)
}

//id -> 店名
func ReadAllShopName() (map[string]string,error) {
	return hashGetValues(RedisShopIdTable)
}

//参考价
func WriteStandardPrice(model, standardPrice, minPrice, maxPrice string) error {
	key := fmt.Sprintf(RedisGoodPriceTableFormat, model)
	return hashSetValue( key, "standardPrice", standardPrice, "minPrice", minPrice, "maxPrice", maxPrice)
}

func ReadStandardPrice(modelNames []string) map[string]*model.GoodPrices {
	conn := RedisClient.Get()
	defer conn.Close()
	ret := make(map[string]*model.GoodPrices)
	for i := range modelNames {
		key := fmt.Sprintf(RedisGoodPriceTableFormat, modelNames[i])
		reply, err := conn.Do("HMGET", key, "standardPrice", "minPrice", "maxPrice")
		if err == nil {
			objs, ok := reply.([]interface{})
			if ok && len(objs) == 3 {
				spb, ok1 := objs[0].([]uint8)
				if ok1 {
					sp := string(spb)
					mp := string(objs[1].([]uint8))
					maxp := string(objs[2].([]uint8))

					standardPrice, _ := strconv.Atoi(sp)
					minPrice, _ := strconv.Atoi(mp)
					maxPrice, _ := strconv.Atoi(maxp)
					gp := model.NewGoodPrices(modelNames[i], standardPrice, minPrice, maxPrice)
					ret[modelNames[i]] = gp
				}

			}
		}
	}
	return ret
}

//good id -> shop info
func SetShopInfo(goodId string, shopName string, shopHref string) error {
	key := fmt.Sprintf(RedisShopInfoFormat,goodId)
	return hashSetValue(key,"shopName",shopName,"shopHref",shopHref)
}

func GetShopInfos(goodId string) (shopName string, shopHref string,err error) {
	key := fmt.Sprintf(RedisShopInfoFormat,goodId)
	var m map[string]string
	m,err = hashGetValues(key)
	if err != nil{
		return
	}
	shopName = m["shopName"]
	shopHref = m["shopHref"]
	return
}

func GetIps() ([]string,error) {
	return setGetValues(RedisIpPoolName)
}

func hashSetValue(key string, pairs ...string) error {
	conn := RedisClient.Get()
	defer conn.Close()

	if len(pairs)%2 != 0{
		return fmt.Errorf("len(pairs)%2 != 0")
	}
	args := make([]interface{},len(pairs)+1)
	args[0] = key
	for i := range pairs {
		args[i+1] = pairs[i]
	}
	_, err := conn.Do("HMSET", args...)
	return err
}

func hashGetValue(key string, args string)(string, error) {
	conn := RedisClient.Get()
	defer conn.Close()

	reply, err := conn.Do("HGET", key,args)
	if err != nil{
		return "",err
	}
	return replyConvString(reply)
}

func hashGetValues(key string) (map[string]string,error) {
	conn := RedisClient.Get()
	defer conn.Close()

	reply, err := conn.Do("HGETALL", key)
	if err != nil{
		return nil,err
	}
	return replyConvMap(reply)
}

func getKeys(parrent string) []string {
	conn := RedisClient.Get()
	defer conn.Close()

	reply, err := conn.Do("KEYS", parrent)
	if err != nil {
		log.Println(err)
		return nil
	}
	lines, err := replyConvStringArray(reply)
	return lines
}

func setAddValue(setName, value string) error {
	conn := RedisClient.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", setName, value)
	return err
}

func setGetValues(setName string) ([]string,error) {
	conn := RedisClient.Get()
	defer conn.Close()
	reply, err := conn.Do("SMEMBERS ", setName)
	if err != nil{
		return nil,err
	}
	return replyConvStringArray(reply)
}

func replyConvStringArray(reply interface{}) ([]string, error) {
	objs, ok := reply.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not array")
	}
	lines := make([]string, len(objs))
	for i, v := range objs {
		lines[i] = string(v.([]uint8))
	}
	return lines, nil
}

func replyConvString(reply interface{}) (string, error) {
	objs, ok := reply.([]uint8)
	if !ok {
		return "", fmt.Errorf("not []uint8")
	}
	return string(objs), nil
}

func replyConvMap(reply interface{}) (map[string]string, error) {
	objs, ok := reply.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not array")
	}
	if len(objs)%2 != 0{
		return nil,fmt.Errorf("len(objs)%2 != 0")
	}
	m := make(map[string]string)
	for i:=0;i<len(objs);i+=2{
		k := objs[i]
		v := objs[i+1]
		key := string(k.([]uint8))
		value := string(v.([]uint8))
		m[key] = value
	}
	return m, nil
}
