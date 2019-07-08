package main

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

type DataMap struct {
	key    string
	leaf   string
	parent *DataMap
	count  int
	data   map[string]*DataMap
}

func (dm *DataMap) Dump(n int) {
	s := ""
	for i := 0; i < n; i++ {
		//s += "  "
	}
	if len(dm.leaf) > 0 && len(dm.data) > 0 {
		leaf := fmt.Sprintf("%s:*", dm.leaf)
		fmt.Printf("%s%-35s => %d(%d)\n", s,leaf, dm.Count(), len(dm.data))
	}
	if len(dm.data) > viper.GetInt("max_split") {
		return
	}
	for _, d := range dm.data {
		d.Dump(n + 1)
	}
}

func (dm *DataMap) Count() int {
	if dm.count > 0 {
		return dm.count
	}
	dm.count = len(dm.data)
	for _, d := range dm.data {
		dm.count += d.Count()
	}
	return dm.count
}

func (dm *DataMap) pre() {
	if dm.data == nil {
		dm.data = make(map[string]*DataMap)
	}
}

func (dm *DataMap) Key(key string) *DataMap {
	dm.pre()
	if m, ok := dm.data[key]; ok {
		return m
	} else {
		m = &DataMap{}
		m.key = key
		if len(dm.key) > 0 {
			m.leaf = dm.leaf + ":" + key
		} else {
			m.leaf = key
		}
		m.parent = dm
		dm.data[key] = m
		return m
	}
}

func main() {
	flag.String("server", "127.0.0.1:6379", "redis地址")
	flag.String("password", "", "授权")
	flag.Int("db", 0, "数据库")
	flag.Int("max_split", 100, "输出最大分割")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	opt := &redis.Options{
		Addr:     viper.GetString("server"),
		Password: viper.GetString("password"), // no password set
		DB:       viper.GetInt("db"),          // use default DB
		PoolSize: 1,
	}
	fmt.Println(opt)
	r := redis.NewClient(opt)

	root := &DataMap{}
	//TODO 这里就是直接keys，标准一点应该用scan来做

	keysCount := 0
	cursor := uint64(0)
	var keys []string
	var err error
	for {
		keys, cursor, err = r.Scan(cursor, "*", 10000).Result()
		if err != nil {
			fmt.Println(err)
			return
		}
		keysCount += len(keys)
		for _, v := range keys {
			ss := strings.Split(v, ":")

			r := root
			for _, s0 := range ss {
				sss := strings.Split(s0, ",")
				for _, s1 := range sss {
					r = r.Key(s1)
				}
			}
		}
		if cursor == 0 {
			break
		}
	}

	/*
		keys, err := r.Keys("*").Result()
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, v := range keys {
			ss := strings.Split(v, ":")

			r := root
			for _, s := range ss {
				r = r.Key(s)
			}
		}
	*/

	fmt.Printf("KeysCount:%d\n", keysCount)
	root.Dump(0)
}
