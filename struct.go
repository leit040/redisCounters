package redisCounters

import (
	"context"
	"fmt"
	RedisPull "github.com/leit040/RedisPull"
	"github.com/redis/rueidis"
	"log"
	"strconv"
	"time"
)

type Counter struct {
	Key string
}
type MapCounters map[string]CountersGroup

type CountersGroup struct {
	prefix  string
	keys    []string
	connect *RedisPull.Connect
}

func GetTimestamps(interval string) (int64, int64) {
	now := time.Now()
	switch interval {
	case "day":
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix() * 1000, 3600000
	case "hour":
		return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location()).Unix() * 1000, 86400000

	default:
		return int64(0), 0
	}

}

func (c *CountersGroup) GetCounterValue(key string, interval string) int {
	if !inArray(key, c.keys) {
		log.Printf("Key %s not found in array\n", key)
		return -1
	}
	counterKey := fmt.Sprintf("%s_%s", c.prefix, key)
	startFrom, duration := GetTimestamps(interval)
	res, err := c.connect.GoRedis.Do(context.Background(), "ts.range", counterKey, startFrom, "+", "AGGREGATION", "SUM", strconv.Itoa(int(duration))).Result()
	if err != nil {
		log.Printf("Something wrong, %s\n", err)
		return -1
	}
	result := res.([]interface{})
	if len(result) == 0 {
		return 0
	}
	lastItem := result[len(result)-1]
	value := lastItem.([]interface{})[1].(float64)
	return int(value)
}

func (c *CountersGroup) IncreaseCounter(key string) {
	counterKey := fmt.Sprintf("%s_%s", c.prefix, key)
	cmd := c.connect.Ruedis.B().TsAdd().Key(counterKey).Timestamp(strconv.FormatInt(time.Now().Unix()*1000, 10)).Value(1).OnDuplicateSum().Labels().Build()
	err := c.connect.Ruedis.Do(context.Background(), cmd).Error()
	if err != nil {
		log.Fatal(err)
	}
}

func (c *CountersGroup) CreateCounters() {
	for _, key := range c.keys {
		counterKey := fmt.Sprintf("%s_%s", c.prefix, key)
		createFeedTimeSeries(counterKey, c.connect.Ruedis)
	}
}

type CountersMap struct {
	cm MapCounters
}

func NewCountersMap() CountersMap {
	var cm CountersMap
	cm.cm = make(map[string]CountersGroup)
	return cm
}

func (cm *CountersMap) AddCountersGroup(prefix string, keys []string, connect *RedisPull.Connect) {
	var cg CountersGroup
	cg.prefix = prefix
	cg.connect = connect
	cg.keys = keys
	cg.CreateCounters()
	cm.cm[prefix] = cg
}

func (cm *CountersMap) GetCountersGroup(prefix string) CountersGroup {
	return cm.cm[prefix]
}

func createFeedTimeSeries(key string, r rueidis.Client) {
	cl := r.B().TsCreate().Key(key).DuplicatePolicyMax()
	err := r.Do(context.Background(), cl.Build()).Error()
	if err != nil {
		fmt.Printf("Time series for %s already exists\n", key)
	}
}
