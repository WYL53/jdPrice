package redisDAO

import (
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"

	"jdPrice/model"
)

const (
	REDIS_HOST                = ":6379"
	RedisTargetModelSet       = "targetModel_%s"
	RedisTargetModelSetPre    = "targetModel_"
	RedisShopIdTable          = "shopIdTable"
	RedisGoodPriceTableFormat = "goodPriceTable_%s"
	RedisPriceZSetFormat      = "prices_%s"
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
			fmt.Println("ZADD err:", err.Error())
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
		fmt.Println(err)
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
			fmt.Println(err)
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
func WiretShopId(id, shopName string) {
	conn := RedisClient.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", RedisShopIdTable, id, shopName)
	if err != nil {
		fmt.Println(err)
	}
}

//id -> 店名
func ReadShopIds() map[string]string {
	conn := RedisClient.Get()
	defer conn.Close()

	reply, err := conn.Do("HGETALL", RedisShopIdTable)
	if err == nil {
		objs, ok := reply.([]interface{})
		if ok {
			ret := make(map[string]string)
			for i := 0; i < len(objs)-2; i += 2 {
				key := string(objs[i].([]uint8))
				value := string(objs[i+1].([]uint8))
				ret[key] = value
			}
			return ret
		}
	} else {
		fmt.Println(err)
	}
	return nil
}

//参考价
func WriteStandardPrice(model, standardPrice, minPrice, maxPrice string) error {
	conn := RedisClient.Get()
	defer conn.Close()
	key := fmt.Sprintf(RedisGoodPriceTableFormat, model)
	_, err := conn.Do("HMSET", key, "standardPrice", standardPrice, "minPrice", minPrice, "maxPrice", maxPrice)
	return err
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

func getKeys(parrent string) []string {
	conn := RedisClient.Get()
	defer conn.Close()

	reply, err := conn.Do("KEYS", parrent)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	lines, err := convStringArray(reply)
	return lines
}

func setAddValue(setName, value string) error {
	conn := RedisClient.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", setName, value)
	return err
}

func convStringArray(reply interface{}) ([]string, error) {
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
