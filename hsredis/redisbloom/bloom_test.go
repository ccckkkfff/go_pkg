package redisbloom

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

var (
	db 		*sql.DB
	rdb 	*redis.Client
)

func bloomInit() {
	var err error
	//1.连接数据库
	strConnCmd := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s","root","ckf123456","118.190.53.54",3306,"gr_test")
	db, err = sql.Open("mysql", strConnCmd)
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
	//defer db.Close()

	opt := &redis.Options{
		Addr:     "192.168.0.105:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}
	rdb = redis.NewClient(opt)
	//defer rdb.Close()
}

func TestBloomReBuild(t *testing.T) {
	bloomInit()
	r,err := BloomReBuild(rdb,db,"bloombak")
	fmt.Println(r,err)

	r,err = BloomReBuild(nil,db,"bloombak")
	fmt.Println(r,err)

	r,err = BloomReBuild(nil,nil,"bloombak")
	fmt.Println(r,err)
}
