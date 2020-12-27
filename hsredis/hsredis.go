package hsredis

import (
	"bytes"
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"main/pkg/hserror"
	"main/pkg/hsmysql"
	"main/pkg/hsredis/redisbloom"
	"main/pkg/tools"
)

var (
	ctx = context.Background()
)


func HsredisAddCache(rdb *redis.Client,k,v string)error{
	if rdb == nil{
		return hserror.New("redis client err!")
	}

	rdb.Do(ctx,"set",k,v)

	return nil
}

func HsredisGetCache(rdb *redis.Client,db *sql.DB, k string)(error,string){
	if rdb == nil{
		return hserror.New("redis client err!"),""
	}

	//1. judget key from bloom
	err,r := redisbloom.BloomIsExists(rdb,k)
	if err != nil || !r{
		return hserror.New("cannot judge from bloom!"),""
	}

	//2. get from cache
	b := new(bytes.Buffer)
	b.WriteString("select user from redis_cache2 where user="+"\""+k+"\"")

	hash := md5.New()
	sqlkey := hex.EncodeToString(hash.Sum(b.Bytes()))

	rst,err := rdb.Do(ctx,"get",sqlkey).Result()
	if rst == nil{
		//3. set mutex
		mutex := fmt.Sprintf("%smutex",k)
		r,err := rdb.Do(ctx,"setnx",mutex,1).Bool()
		if r && err == nil{
			//4. get from db
			cache := hsmysql.GetDBVal(db,b.String())
			if cache != ""{
				rdb.Do(ctx,"setex",sqlkey,tools.RandomNum1000()%300,string(cache)).Bool()
			}

			//5. recover mutex
			rdb.Do(ctx,"del",mutex).Bool()

			return errors.New("get from db"),string(cache)
		}else{
			return errors.New("can not get mutex!"),""
		}
	}

	return errors.New("get from redis"),rst.(string)
}
