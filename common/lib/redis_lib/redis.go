package redis_lib

import (
	"encoding/json"
	"fehu/common/lib"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Cacher 先构建一个Cacher实例，然后将配置参数传入该实例的StartAndGC方法来初始化实例和程序进程退出后的清理工作。
type Cacher struct {
	pool      *redis.Pool
	prefix    string
	marshal   func(v interface{}) ([]byte, error)
	unmarshal func(data []byte, v interface{}) error
}

// Options redis配置参数
type Options struct {
	lib.RedisConf
	Addr string // redis服务的地址，默认为 127.0.0.1:6379
}

var redisClusterMap map[string][]*Cacher

type GlobalRedisUitl int

func InitRedis() {
	redisClusterMap = make(map[string][]*Cacher)
	for confName, cfg := range lib.ConfRedisMap.List {
		cachers := make([]*Cacher, 0)
		for i := 0; i < len(cfg.ProxyList); i++ {
			addr := cfg.ProxyList[i]
			c := &Cacher{
				prefix: cfg.Prefix,
			}
			c.StartAndGC(&Options{
				RedisConf: *cfg,
				Addr:      addr,
			})
			cachers = append(cachers, c)
		}
		redisClusterMap[confName] = cachers
	}
}

//func RedisConnFactory(name string) (redis.Conn, error) {
//	for confName, cfg := range _struct.ConfRedisMap.List {
//		if name == confName {
//			randHost := cfg.ProxyList[rand.Intn(len(cfg.ProxyList))]
//			return redis.Dial(
//				"tcp",
//				randHost,
//				redis.DialConnectTimeout(50*time.Millisecond),
//				redis.DialReadTimeout(100*time.Millisecond),
//				redis.DialWriteTimeout(100*time.Millisecond))
//		}
//	}
//	return nil, errors.New("create redis conn fail")
//}

// StartAndGC 使用 Options 初始化redis，并在程序进程退出时关闭连接池。
func (c *Cacher) StartAndGC(opts *Options) error {

	if opts.Network == "" {
		opts.Network = "tcp"
	}
	if opts.Addr == "" {
		opts.Addr = "127.0.0.1:6379"
	}
	if opts.MaxIdle == 0 {
		opts.MaxIdle = 3
	}
	if opts.IdleTimeout == 0 {
		opts.IdleTimeout = 300
	}
	c.marshal = json.Marshal
	c.unmarshal = json.Unmarshal

	pool := &redis.Pool{
		MaxActive:   opts.MaxActive,
		MaxIdle:     opts.MaxIdle,
		IdleTimeout: time.Duration(opts.IdleTimeout) * time.Second,

		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(opts.Network, opts.Addr)
			if err != nil {
				return nil, err
			}
			if opts.Password != "" {
				if _, err := conn.Do("AUTH", opts.Password); err != nil {
					conn.Close()
					return nil, err
				}
			}
			if _, err := conn.Do("SELECT", opts.Db); err != nil {
				conn.Close()
				return nil, err
			}
			return conn, err
		},

		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
	c.pool = pool
	c.closePool()
	return nil

}

// closePool 程序进程退出时关闭连接池
func (c *Cacher) closePool() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	signal.Notify(ch, syscall.SIGKILL)
	go func() {
		<-ch
		c.pool.Close()
		os.Exit(0)
	}()
}
func randPool(clusterName string) redis.Conn {
	if clusterName == "" {
		clusterName = "default"
	}
	cachers := redisClusterMap["default"]
	// 这里简单做一个轮训
	c := cachers[rand.Intn(len(cachers))]
	return c.pool.Get()
}

// Do 执行redis命令并返回结果。执行时从连接池获取连接并在执行完命令后关闭连接。
func Do(clusterName string, commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := randPool(clusterName)
	defer conn.Close()
	return conn.Do(commandName, args...)
}

// Get 获取键值。一般不直接使用该值，而是配合下面的工具类方法获取具体类型的值，或者直接使用github.com/gomodule/redigo/redis包的工具方法。
func Get(key string, clusterName string) (interface{}, error) {
	return Do(clusterName, "GET", key)
}

// GetString 获取string类型的键值
func GetString(key string, clusterName string) (string, error) {
	return String(Get(key, clusterName))
}

// GetInt 获取int类型的键值
func GetInt(key string, clusterName string) (int, error) {
	return Int(Get(key, clusterName))
}

// GetInt64 获取int64类型的键值
func GetInt64(key string, clusterName string) (int64, error) {
	return Int64(Get(key, clusterName))
}

// GetBool 获取bool类型的键值
func GetBool(key string, clusterName string) (bool, error) {
	return Bool(Get(key, clusterName))
}

// GetObject 获取非基本类型stuct的键值。在实现上，使用json的Marshal和Unmarshal做序列化存取。
func GetObject(key string, val interface{}, clusterName string) error {
	reply, err := Get(key, clusterName)
	return decode(reply, err, val)
}

// Set 存并设置有效时长。时长的单位为秒。
// 基础类型直接保存，其他用json.Marshal后转成string保存。
func Set(key string, val interface{}, expire int64, clusterName string) error {
	value, err := encode(val)
	if err != nil {
		return err
	}
	if expire > 0 {
		_, err := Do(clusterName, "SETEX", key, expire, value)
		return err
	}
	_, err = Do(clusterName, "SET", key, value)
	return err
}

// Exists 检查键是否存在
func Exists(key string, clusterName string) (bool, error) {
	return Bool(Do(clusterName, "EXISTS", key))
}

//Del 删除键
func Del(key string, clusterName string) error {
	_, err := Do(clusterName, "DEL", key)
	return err
}

// Flush 清空当前数据库中的所有 key，慎用！
func Flush(clusterName string) error {
	//TODO 这里应该删除所有
	_, err := Do(clusterName, "db", "FLUSHDB")
	return err
}

// TTL 以秒为单位。当 key 不存在时，返回 -2 。 当 key 存在但没有设置剩余生存时间时，返回 -1
func TTL(key string, clusterName string) (ttl int64, err error) {
	return Int64(Do(clusterName, "TTL", key))
}

// Expire 设置键过期时间，expire的单位为秒
func Expire(key string, expire int64, clusterName string) error {
	_, err := Bool(Do(clusterName, "EXPIRE", key, expire))
	return err
}

// Incr 将 key 中储存的数字值增一
func Incr(key string, clusterName string) (val int64, err error) {
	return Int64(Do(clusterName, "INCR", key))
}

// IncrBy 将 key 所储存的值加上给定的增量值（increment）。
func IncrBy(key string, amount int64, clusterName string) (val int64, err error) {
	return Int64(Do(clusterName, "INCRBY", key, amount))
}

// Decr 将 key 中储存的数字值减一。
func Decr(key string, clusterName string) (val int64, err error) {
	return Int64(Do(clusterName, "DECR", key))
}

// DecrBy key 所储存的值减去给定的减量值（decrement）。
func DecrBy(key string, amount int64, clusterName string) (val int64, err error) {
	return Int64(Do(clusterName, "DECRBY", key, amount))
}

// HMSet 将一个map存到Redis hash，同时设置有效期，单位：秒
// Example:
//
// ```golang
// m := make(map[string]interface{})
// m["name"] = "corel"
// m["age"] = 23
// err := HMSet("user", m, 10)
// ```
func HMSet(key string, val interface{}, expire int, clusterName string) (err error) {
	conn := randPool(clusterName)
	defer conn.Close()
	err = conn.Send("HMSET", redis.Args{}.Add(key).AddFlat(val)...)
	if err != nil {
		return
	}
	if expire > 0 {
		err = conn.Send("EXPIRE", key, int64(expire))
	}
	if err != nil {
		return
	}
	conn.Flush()
	_, err = conn.Receive()
	return
}

/** Redis hash 是一个string类型的field和value的映射表，hash特别适合用于存储对象。 **/

// HSet 将哈希表 key 中的字段 field 的值设为 val
// Example:
//
// ```golang
// _, err := HSet("user", "age", 23)
// ```
func HSet(key, field string, val interface{}, clusterName string) (interface{}, error) {
	value, err := encode(val)
	if err != nil {
		return nil, err
	}
	return Do(clusterName, "HSET", key, field, value)
}

// HGet 获取存储在哈希表中指定字段的值
// Example:
//
// ```golang
// val, err := HGet("user", "age")
// ```
func HGet(key, field string, clusterName string) (reply interface{}, err error) {
	reply, err = Do(clusterName, "HGET", key, field)
	return
}

// HGetString HGet的工具方法，当字段值为字符串类型时使用
func HGetString(key, field string, clusterName string) (reply string, err error) {
	reply, err = String(HGet(clusterName, key, field))
	return
}

// HGetInt HGet的工具方法，当字段值为int类型时使用
func HGetInt(key, field string, clusterName string) (reply int, err error) {
	reply, err = Int(HGet(clusterName, key, field))
	return
}

// HGetInt64 HGet的工具方法，当字段值为int64类型时使用
func HGetInt64(key, field string, clusterName string) (reply int64, err error) {
	reply, err = Int64(HGet(clusterName, key, field))
	return
}

// HGetBool HGet的工具方法，当字段值为bool类型时使用
func HGetBool(key, field string, clusterName string) (reply bool, err error) {
	reply, err = Bool(HGet(clusterName, key, field))
	return
}

// HGetObject HGet的工具方法，当字段值为非基本类型的stuct时使用
func HGetObject(key, field string, val interface{}, clusterName string) error {
	reply, err := HGet(clusterName, key, field)
	return decode(reply, err, val)
}

// HGetAll HGetAll("key", &val)
func HGetAll(key string, val interface{}, clusterName string) error {
	v, err := redis.Values(Do("HGETALL", key))
	if err != nil {
		return err
	}

	if err := redis.ScanStruct(v, val); err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("%+v\n", val)
	return err
}

/**
Redis列表是简单的字符串列表，按照插入顺序排序。你可以添加一个元素到列表的头部（左边）或者尾部（右边）
**/

// BLPop 它是 LPOP 命令的阻塞版本，当给定列表内没有任何元素可供弹出的时候，连接将被 BLPOP 命令阻塞，直到等待超时或发现可弹出元素为止。
// 超时参数 timeout 接受一个以秒为单位的数字作为值。超时参数设为 0 表示阻塞时间可以无限期延长(block indefinitely) 。
func BLPop(clusterName string, key string, timeout int) (interface{}, error) {
	values, err := redis.Values(Do(clusterName, "BLPOP", key, timeout))
	if err != nil {
		return nil, err
	}
	if len(values) != 2 {
		return nil, fmt.Errorf("redisgo: unexpected number of values, got %d", len(values))
	}
	return values[1], err
}

// BLPopInt BLPop的工具方法，元素类型为int时
func BLPopInt(key string, timeout int, clusterName string) (int, error) {
	return Int(BLPop(clusterName, key, timeout))
}

// BLPopInt64 BLPop的工具方法，元素类型为int64时
func BLPopInt64(key string, timeout int, clusterName string) (int64, error) {
	return Int64(BLPop(clusterName, key, timeout))
}

// BLPopString BLPop的工具方法，元素类型为string时
func BLPopString(key string, timeout int, clusterName string) (string, error) {
	return String(BLPop(clusterName, key, timeout))
}

// BLPopBool BLPop的工具方法，元素类型为bool时
func BLPopBool(key string, timeout int, clusterName string) (bool, error) {
	return Bool(BLPop(clusterName, key, timeout))
}

// BLPopObject BLPop的工具方法，元素类型为object时
func BLPopObject(key string, timeout int, val interface{}, clusterName string) error {
	reply, err := BLPop(clusterName, key, timeout)
	return decode(reply, err, val)
}

// BRPop 它是 RPOP 命令的阻塞版本，当给定列表内没有任何元素可供弹出的时候，连接将被 BRPOP 命令阻塞，直到等待超时或发现可弹出元素为止。
// 超时参数 timeout 接受一个以秒为单位的数字作为值。超时参数设为 0 表示阻塞时间可以无限期延长(block indefinitely) 。
func BRPop(clusterName string, key string, timeout int) (interface{}, error) {
	values, err := redis.Values(Do(clusterName, "BRPOP", key, timeout))
	if err != nil {
		return nil, err
	}
	if len(values) != 2 {
		return nil, fmt.Errorf("redisgo: unexpected number of values, got %d", len(values))
	}
	return values[1], err
}

// BRPopInt BRPop的工具方法，元素类型为int时
func BRPopInt(key string, timeout int, clusterName string) (int, error) {
	return Int(BRPop(clusterName, key, timeout))
}

// BRPopInt64 BRPop的工具方法，元素类型为int64时
func BRPopInt64(key string, timeout int, clusterName string) (int64, error) {
	return Int64(BRPop(clusterName, key, timeout))
}

// BRPopString BRPop的工具方法，元素类型为string时
func BRPopString(key string, timeout int, clusterName string) (string, error) {
	return String(BRPop(clusterName, key, timeout))
}

// BRPopBool BRPop的工具方法，元素类型为bool时
func BRPopBool(key string, timeout int, clusterName string) (bool, error) {
	return Bool(BRPop(clusterName, key, timeout))
}

// BRPopObject BRPop的工具方法，元素类型为object时
func BRPopObject(key string, timeout int, val interface{}, clusterName string) error {
	reply, err := BRPop(clusterName, key, timeout)
	return decode(reply, err, val)
}

// LPop 移出并获取列表中的第一个元素（表头，左边）
func LPop(key string, clusterName string) (interface{}, error) {
	return Do(clusterName, "LPOP", key)
}

// LPopInt 移出并获取列表中的第一个元素（表头，左边），元素类型为int
func LPopInt(key string, clusterName string) (int, error) {
	return Int(LPop(clusterName, key))
}

// LPopInt64 移出并获取列表中的第一个元素（表头，左边），元素类型为int64
func LPopInt64(key string, clusterName string) (int64, error) {
	return Int64(LPop(clusterName, key))
}

// LPopString 移出并获取列表中的第一个元素（表头，左边），元素类型为string
func LPopString(key string, clusterName string) (string, error) {
	return String(LPop(clusterName, key))
}

// LPopBool 移出并获取列表中的第一个元素（表头，左边），元素类型为bool
func LPopBool(key string, clusterName string) (bool, error) {
	return Bool(LPop(clusterName, key))
}

// LPopObject 移出并获取列表中的第一个元素（表头，左边），元素类型为非基本类型的struct
func LPopObject(key string, val interface{}, clusterName string) error {
	reply, err := LPop(clusterName, key)
	return decode(reply, err, val)
}

// RPop 移出并获取列表中的最后一个元素（表尾，右边）
func RPop(key string, clusterName string) (interface{}, error) {
	return Do(clusterName, "RPOP", key)
}

// RPopInt 移出并获取列表中的最后一个元素（表尾，右边），元素类型为int
func RPopInt(key string, clusterName string) (int, error) {
	return Int(RPop(clusterName, key))
}

// RPopInt64 移出并获取列表中的最后一个元素（表尾，右边），元素类型为int64
func RPopInt64(key string, clusterName string) (int64, error) {
	return Int64(RPop(clusterName, key))
}

// RPopString 移出并获取列表中的最后一个元素（表尾，右边），元素类型为string
func RPopString(key string, clusterName string) (string, error) {
	return String(RPop(clusterName, key))
}

// RPopBool 移出并获取列表中的最后一个元素（表尾，右边），元素类型为bool
func RPopBool(key string, clusterName string) (bool, error) {
	return Bool(RPop(clusterName, key))
}

// RPopObject 移出并获取列表中的最后一个元素（表尾，右边），元素类型为非基本类型的struct
func RPopObject(key string, val interface{}, clusterName string) error {
	reply, err := RPop(clusterName, key)
	return decode(reply, err, val)
}

// LPush 将一个值插入到列表头部
func LPush(key string, member interface{}, clusterName string) error {
	value, err := encode(member)
	if err != nil {
		return err
	}
	_, err = Do(clusterName, "LPUSH", key, value)
	return err
}

// RPush 将一个值插入到列表尾部
func RPush(key string, member interface{}, clusterName string) error {
	value, err := encode(member)
	if err != nil {
		return err
	}
	_, err = Do(clusterName, "RPUSH", key, value)
	return err
}

// encode 序列化要保存的值
func encode(val interface{}) (interface{}, error) {
	var value interface{}
	switch v := val.(type) {
	case string, int, uint, int8, int16, int32, int64, float32, float64, bool:
		value = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		value = string(b)
	}
	return value, nil
}

// decode 反序列化保存的struct对象
func decode(reply interface{}, err error, val interface{}) error {
	str, err := String(reply, err)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(str), val)
}

//func RedisLogDo(trace *TraceContext, c redis.Conn, commandName string, args ...interface{}) (interface{}, error) {
//	startExecTime := time.Now()
//	reply, err := Do(commandName, args...)
//	endExecTime := time.Now()
//	if err != nil {
//		Log.TagError(trace, "_com_redis_failure", map[string]interface{}{
//			"method":    commandName,
//			"err":       err,
//			"bind":      args,
//			"proc_time": fmt.Sprintf("%fms", endExecTime.Sub(startExecTime).Seconds()),
//		})
//	} else {
//		replyStr, _ := redis.String(reply, nil)
//		Log.TagInfo(trace, "_com_redis_success", map[string]interface{}{
//			"method":    commandName,
//			"bind":      args,
//			"reply":     replyStr,
//			"proc_time": fmt.Sprintf("%fms", endExecTime.Sub(startExecTime).Seconds()),
//		})
//	}
//	return reply, err
//}
//
////通过配置 执行redis
//func RedisConfDo(trace *TraceContext, name string, commandName string, args ...interface{}) (interface{}, error) {
//	c, err := RedisConnFactory(name)
//	if err != nil {
//		Log.TagError(trace, "_com_redis_failure", map[string]interface{}{
//			"method": commandName,
//			"err":    errors.New("RedisConnFactory_error:" + name),
//			"bind":   args,
//		})
//		return nil, err
//	}
//	defer Close()
//
//	startExecTime := time.Now()
//	reply, err := Do(commandName, args...)
//	endExecTime := time.Now()
//	if err != nil {
//		Log.TagError(trace, "_com_redis_failure", map[string]interface{}{
//			"method":    commandName,
//			"err":       err,
//			"bind":      args,
//			"proc_time": fmt.Sprintf("%fms", endExecTime.Sub(startExecTime).Seconds()),
//		})
//	} else {
//		replyStr, _ := redis.String(reply, nil)
//		Log.TagInfo(trace, "_com_redis_success", map[string]interface{}{
//			"method":    commandName,
//			"bind":      args,
//			"reply":     replyStr,
//			"proc_time": fmt.Sprintf("%fms", endExecTime.Sub(startExecTime).Seconds()),
//		})
//	}
//	return reply, err
//}
