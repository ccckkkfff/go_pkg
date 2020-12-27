package hsredis

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestHsredisGetCache(t *testing.T) {
	strConnCmd := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s","root","ckf123456","118.190.53.54",3306,"gr_test")
	db, err := sql.Open("mysql", strConnCmd)
	if err != nil {
		fmt.Println(err)
		return
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()


	opt := &redis.Options{
		Addr:     "192.168.0.105:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}
	rdb := redis.NewClient(opt)

	tests := map[string]string{
		"15af61c024":"15af61c024",
		"786149f00a":"786149f00a",
		"284c36050b":"284c36050b",
		"890b31c10a":"890b31c10a",
		"b703bb3f97":"b703bb3f97",
		"cac73efa10":"cac73efa10",
		"cc6a25d620":"cc6a25d620",
	}

	for k,v := range tests{
		t.Run(k, func(t *testing.T) {
			err,s := HsredisGetCache(rdb,db,v)
			t.Log(err,s)
		})
	}
}
