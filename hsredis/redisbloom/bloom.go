//
//a single bloom use redis
//two redis bloom embedded in redis.
//when bloomset == 1,use bloomain.
//when bloomset == 2use bloombak
//
package redisbloom

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"main/pkg/hserror"
	"strings"
)
var (
	ctx = context.Background()
)

//-------------------------------------------------------
// BloomInit
// ------------------------------------------------------
//Return Value: err
//Params:
//		1)rdb * redis.Client:		rdb client
//init redis bloom
func BloomInit(rdb * redis.Client)error{
	if rdb == nil{
		return errors.New("redis.Client failed!")
	}

	_,err := rdb.Do(ctx,"bf.reserve","bloomain","0.01","1000000").Result()
	if err != nil || strings.Index(err.Error(),"exists") == -1 {
		return err
	}

	_,err = rdb.Do(ctx,"bf.reserve","bloombak","0.01","1000000").Result()
	if err != nil || strings.Index(err.Error(),"exists") == -1 {
		return err
	}

	_,err = rdb.Do(ctx,"setnx","bloomset","1").Result()
	if err != nil {
		return err
	}

	return nil
}

//-------------------------------------------------------
// BloomAdd
// ------------------------------------------------------
//Return Value: err
//Params:
//		1)rdb * redis.Client:		rdb client
//		2)val interface{}:          the value
//add value into bloom
func BloomAdd(rdb * redis.Client,val interface{})error{
	if rdb == nil{
		return errors.New("redis.Client failed!")
	}
	err,key := getbloomKey(rdb)
	if err != nil{
		return err
	}

	_,err = rdb.Do(ctx,"bf.add",key,val).Result()
	return err
}

//-------------------------------------------------------
// BloomReBuild
// ------------------------------------------------------
//Return Value: err
//Params:
//		1)rdb * redis.Client:		rdb client
//		2)key string:          		the bloom key
//rebuild bloom
func BloomReBuild(rdb *redis.Client, db *sql.DB,key string)(string,error){
	if rdb == nil || db == nil{
		return "",hserror.New("client err!")
	}

	//1. judge bloom key
	if strings.EqualFold(key,"bloomain"){
		err := SetBloomKey(rdb,2)
		if err != nil{
			return "",hserror.New(err.Error())
		}
	}else if strings.EqualFold(key,"bloombak"){
		err := SetBloomKey(rdb,1)
		if err != nil{
			return "",hserror.New(err.Error())
		}
	}else {
		return "",errors.New("wrong bloom key!")
	}

	//2. delete key
	_,err := rdb.Do(ctx,"del",key).Bool()
	if err != nil{
		return "",hserror.New(err.Error())
	}

	//3. set bloom
	_,err = rdb.Do(ctx,"bf.reserve",key,"0.01","1000000").Result()
	if err != nil{
		return "",hserror.New(err.Error())
	}

	//4. start add bloom key
	ch := make(chan string)
	//4.1 wait for bloom key
	go func() {
		for{
			select {
			case u,ok := <- ch:
				if !ok{
					break
				}
				rdb.Do(ctx,"bf.add",key,u)
			}
		}
	}()

	//4.2 get key from mysql
	var cnt uint32 = 0
	for{
		rows,err := db.Query("select id,user from redis_cache2 where id>? order by id limit 10000",cnt)
		if err != nil{
			return "",hserror.New(err.Error())
		}
		defer rows.Close()

		if !rows.Next(){
			close(ch)
			break
		}

		var u string
		for {
			rows.Scan(&cnt,&u)
			ch <- u

			if !rows.Next(){
				break
			}
		}
	}

	//5. recover bloom key
	if strings.EqualFold(key,"bloomain"){
		SetBloomKey(rdb,1)
	}else if strings.EqualFold(key,"bloombak"){
		SetBloomKey(rdb,2)
	}

	return "BloomReBuild success!",nil
}

//-------------------------------------------------------
// SetBloomKey
// ------------------------------------------------------
//Return Value: err
//Params:
//		1)rdb * redis.Client:		rdb client
//		2)val int:          		the flag of using bloom key
//set bloom key
func SetBloomKey(rdb * redis.Client,flg int)error{
	if rdb == nil{
		return errors.New("redis.Client failed!")
	}
	if flg != 1 && flg != 2{
		return errors.New(fmt.Sprintf("set wrong val:%d!",flg))
	}

	_,err := rdb.Do(ctx,"set","bloomset",flg).Result()
	return err
}

//-------------------------------------------------------
// getbloomKey
// ------------------------------------------------------
//Return Value: err
//Params:
//		1)rdb * redis.Client:		rdb client
//get bloom key
func getbloomKey(rdb * redis.Client)(error,string){
	if rdb == nil{
		return errors.New("redis.Client failed!"),""
	}
	r,err := rdb.Do(ctx,"get","bloomset").Int()
	if err != nil{
		return err,""
	}

	if r == 1{
		return nil,"bloomain"
	}else if r == 2{
		return nil,"bloombak"
	}

	return errors.New("can not found key"),""
}


//-------------------------------------------------------
// BloomIsExists
// ------------------------------------------------------
//Return Value: err
//Params:
//		1)rdb * redis.Client:		rdb client
//		2)val string:				the value to judge
//judge the val
func BloomIsExists(rdb * redis.Client,val string)(error,bool){
	if rdb == nil{
		return errors.New("redis.Client failed!"),false
	}
	if val == ""{
		return errors.New("Bloom key can not nil"),false
	}

	err,key := getbloomKey(rdb)
	if err != nil{
		return  errors.New(err.Error()),false
	}

	r,err := rdb.Do(ctx,"bf.exists", key, val).Bool()
	if err != nil{
		return errors.New("bloom cmd r:"+err.Error()),false
	}

	return nil,r
}