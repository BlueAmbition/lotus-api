package redis

import (
	"github.com/astaxie/beego"
	redigo "github.com/gomodule/redigo/redis"
	"strconv"
)

var redisPool *redigo.Pool

func init() {
	var (
		con redigo.Conn
		err error
	)
	host := beego.AppConfig.String("redis::host")
	port, _ := beego.AppConfig.Int("redis::port")
	password := beego.AppConfig.String("redis::password")
	maxIdle, _ := beego.AppConfig.Int("redis::max_idle")
	maxActive, _ := beego.AppConfig.Int("redis::max_active")
	db := 0
	redisPool = &redigo.Pool{
		MaxIdle:   maxIdle,   //最大空闲数
		MaxActive: maxActive, // 最大连接数
		Wait:      true,
		Dial: func() (redigo.Conn, error) {
			if password != "" {
				con, err = redigo.Dial("tcp", host+":"+strconv.Itoa(port),
					redigo.DialPassword(password),
					redigo.DialDatabase(db))
			} else {
				con, err = redigo.Dial("tcp", host+":"+strconv.Itoa(port),
					redigo.DialDatabase(db))
			}
			//redis.DialConnectTimeout(timeout*time.Second),
			//redis.DialReadTimeout(timeout*time.Second),
			//redis.DialWriteTimeout(timeout*time.Second))
			if err != nil {
				return nil, err
			}
			return con, err
		},
	}
}

//获取连接
func GetCon(db uint) redigo.Conn {
	con := redisPool.Get()
	con.Do("SELECT", db)
	return con
}

// Key是否存在
func KeyExists(db uint, key string) bool {
	var (
		err      error
		flag     int
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	flag, err = redigo.Int(redisCon.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return flag > 0
}

//设置Key
func SetString(db uint, key string, value interface{}, expireSeconds int) bool {
	var (
		err      error
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	if expireSeconds > 0 {
		_, err = redisCon.Do("SET", key, value, "EX", expireSeconds)
	} else {
		_, err = redisCon.Do("SET", key, value)
	}
	return err == nil
}

//获取Key
func GetString(db uint, key string, vType string) (interface{}, error) {
	var (
		value    interface{}
		err      error
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()

	switch vType {
	case "string":
		value, err = redigo.String(redisCon.Do("GET", key))
		break
	case "int64":
		value, err = redigo.Int64(redisCon.Do("GET", key))
		break
	case "int":
		value, err = redigo.Int(redisCon.Do("GET", key))
		break
	case "float64":
		value, err = redigo.Float64(redisCon.Do("GET", key))
		break
	}
	return value, err
}

//设置过期时间
func ExpireKey(db uint, key string, expireSeconds int) bool {
	var (
		err      error
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	_, err = redigo.Int(redisCon.Do("expire ", key, expireSeconds))
	if err != nil {
		return false
	}
	return true
}

//获取Key剩余有效秒数
func GetTTLKey(db uint, key string) int {
	var (
		value    int
		err      error
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	value, err = redigo.Int(redisCon.Do("TTL", key))
	if err != nil || value == -2 {
		return 0
	}
	return value
}

//删除Key
func DelKey(db uint, key string) bool {
	var (
		err  error
		flag int
	)
	redisCon := GetCon(db)
	defer redisCon.Close()
	flag, err = redigo.Int(redisCon.Do("DEL", key))
	if err != nil {
		return false
	}
	return flag > 0
}

//Hash Key是否存在
func HashExists(db uint, key string) bool {
	var (
		err      error
		flag     int
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	flag, err = redigo.Int(redisCon.Do("HEXISTS", key))
	if err != nil {
		return false
	}
	return flag > 0
}

//删除Key
func DelHash(db uint, key string, fields ...string) bool {
	var (
		err      error
		flag     int
		redisCon redigo.Conn
		args     []interface{}
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	args = []interface{}{key}
	for _, v := range fields {
		args = append(args, v)
	}
	flag, err = redigo.Int(redisCon.Do("HDEL", args...))
	if err != nil {
		return false
	}
	return flag > 0
}

//设置Hash
func SetHash(db uint, key string, kvs ...string) bool {
	var (
		err error
		//flag     interface{}
		redisCon redigo.Conn
		args     []interface{}
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	args = []interface{}{key}
	for _, v := range kvs {
		args = append(args, v)
	}
	_, err = redisCon.Do("HMSET", args...)
	return err == nil
}

//获取Hash
func GetHash(db uint, key string, fields ...string) []interface{} {
	var (
		err      error
		data     []interface{}
		redisCon redigo.Conn
		args     []interface{}
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	args = []interface{}{key}
	for _, v := range fields {
		args = append(args, v)
	}
	data, err = redigo.Values(redisCon.Do("HMGET", args...))
	if err != nil {
		return nil
	}
	return data
}

//设置List
func SetList(db uint, key string, value string, pre bool) bool {
	var (
		err      error
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	if pre {
		_, err = redisCon.Do("LPUSH", key, value)
	} else {
		_, err = redisCon.Do("RPUSH", key, value)
	}
	return err == nil
}

//获取List
func GetList(db uint, key string, begin int, end int) []interface{} {
	var (
		err      error
		redisCon redigo.Conn
		data     []interface{}
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	data, err = redigo.Values(redisCon.Do("LRANGE", key, begin, end))
	if err != nil {
		return nil
	}
	return data
}

//删除list所有value值的项
func DelList(db uint, key string, value string) bool {
	var (
		err      error
		redisCon redigo.Conn
		flag     int
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	flag, err = redigo.Int(redisCon.Do("LREM", key, 0, value))
	if err != nil {
		return false
	}
	return flag > 0
}

//设置有序集合
func SetSortSet(db uint, key string, score int, value string) bool {
	var (
		err      error
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	_, err = redisCon.Do("ZADD", key, score, value)
	return err == nil
}

//获取有序集合
func GetSortSet(db uint, key string, begin int, end int, withScore bool) []interface{} {
	var (
		err      error
		redisCon redigo.Conn
		data     []interface{}
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	if withScore {
		data, err = redigo.Values(redisCon.Do("ZRANGE", key, begin, end, "WITHSCORES"))
	} else {
		data, err = redigo.Values(redisCon.Do("ZRANGE", key, begin, end))
	}
	if err != nil {
		return nil
	}
	return data
}

//删除SortSet所有value值的项
func DelSortSet(db uint, key string, value string) bool {
	var (
		err      error
		redisCon redigo.Conn
		flag     int
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	flag, err = redigo.Int(redisCon.Do("ZREM", key, value))
	if err != nil {
		return false
	}
	return flag > 0
}
