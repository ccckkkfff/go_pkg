package hsmysql

import (
	"main/pkg/tools"
	"testing"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"fmt"
)

func TestInsertCols(t *testing.T) {
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

	type redis_cache struct {
		User		  string
		Gmt_create    time.Time
		Gmt_modify	  time.Time
	}

	type redis_cache2 struct {
		User		  string
		Gmt_create    string
		Gmt_modify	  string
	}

	/*测试组*/
	type test struct {
		db			*sql.DB
		sql			string
		data		interface{}
	}
	tests := map[string]test{
		"nil db":{nil,"",[]redis_cache{}},
		"nil data":{db,"",nil},
		"data":{db,"insert into Redis_Cache2 (user,gmt_create,gmt_modify) values ",[]redis_cache{
			redis_cache{tools.RandomMillonToString(),time.Now(),time.Now()},
			redis_cache{tools.RandomMillonToString(),time.Now(),time.Now()},
		}},
		"data pointer":{db,"insert into Redis_Cache2 (user,gmt_create,gmt_modify) values ",&[]redis_cache{
			redis_cache{tools.RandomMillonToString(),time.Now(),time.Now()},
			redis_cache{tools.RandomMillonToString(),time.Now(),time.Now()},
		}},
		"data string":{db,"insert into Redis_Cache2 (user,gmt_create,gmt_modify) values ",[]redis_cache2{
			redis_cache2{tools.RandomMillonToString(),tools.GetLocalTime(),tools.GetLocalTime()},
			redis_cache2{tools.RandomMillonToString(),tools.GetLocalTime(),tools.GetLocalTime()},
		}},
	}

	for name,tc := range tests{
		t.Run(name, func(t *testing.T) {
			err,r := InsertCols(tc.db,tc.sql,tc.data)
			t.Log(err,r)
		})
	}
}

func BenchmarkInsertCols(b *testing.B) {
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

	type redis_cache struct {
		User		  string
		Gmt_create    time.Time
		Gmt_modify	  time.Time
	}
	data := []redis_cache{
	//	redis_cache{tools.RandomMillonToString(), time.Now(), time.Now()},
	//	redis_cache{tools.RandomMillonToString(), time.Now(), time.Now()},
	//	redis_cache{tools.RandomMillonToString(), time.Now(), time.Now()},
	//	redis_cache{tools.RandomMillonToString(), time.Now(), time.Now()},
	//	redis_cache{tools.RandomMillonToString(), time.Now(), time.Now()},
	//	redis_cache{tools.RandomMillonToString(), time.Now(), time.Now()},
	//	redis_cache{tools.RandomMillonToString(), time.Now(), time.Now()},
	//	redis_cache{tools.RandomMillonToString(), time.Now(), time.Now()},
	//	redis_cache{tools.RandomMillonToString(), time.Now(), time.Now()},
	//	redis_cache{tools.RandomMillonToString(), time.Now(), time.Now()},
	//	redis_cache{tools.RandomMillonToString(), time.Now(), time.Now()},
	//	redis_cache{tools.RandomMillonToString(), time.Now(), time.Now()},
	//	redis_cache{tools.RandomMillonToString(), time.Now(), time.Now()},
	//	redis_cache{tools.RandomMillonToString(), time.Now(), time.Now()},
	}

	for i:=0;i<100;i++{
		data = append(data,redis_cache{tools.RandomMillonToString(), time.Now(), time.Now()})
	}
	//设置CPU数
	b.SetParallelism(1)
	b.ResetTimer()

	//并行方式测试
	//b.RunParallel()

	for i:=0;i<b.N;i++{
		InsertCols(db,"insert into Redis_Cache2 (user,gmt_create,gmt_modify) values ",data)
		//err,r := InsertCols(db,"insert into Redis_Cache2 (user,gmt_create,gmt_modify) values ",data)
		//b.Log(err,r)
	}
}