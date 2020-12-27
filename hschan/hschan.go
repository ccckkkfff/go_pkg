//用于处理批量插入的redis和mysql
//将(interface{}->interface{})当做参数传入无法用reflect调用
package hschan

import (
	"database/sql"
	"fmt"
	"main/pkg"
	"main/pkg/hsmysql"
	"time"
)

//用于保存批量插入记录的结构体
type HsChanElem struct {
	StateMent	string
	//Data 		interface{}		//具体实现时修改该interface{}为对应的机构体
	Data		[]pkg.RedisCache
	Priority	uint8
	Cap			uint16
}

//对应的通道结构体
type HsChan struct {
	list chan 	HsChanElem
	close chan 	bool
	db         	*sql.DB
}

//创建对于的Chan
func NewHsChan(db *sql.DB,size int)* HsChan{
	return &HsChan{
		list:	make(chan HsChanElem,size),
		close:  make(chan bool),
		db:		db,
	}
}

//启动goroutine缓存处理插入任务
func (h *HsChan)Start(n int){
	for i:=0;i<n;i++{
		go func() {
			var data = make([]pkg.RedisCache,0,200)
			var maxPriority uint8 = 1

			for{
				select {
				case e := <-h.list:
					if e.Cap > 100{
						err,_ := hsmysql.InsertCols(h.db,"insert into Redis_Cache2 (user,gmt_create,gmt_modify) values ",e.Data)
						if err != nil{
							if e.Priority < 10{
								e.Priority++
								h.list <- e
							}
							fmt.Println("InsertCols 100:",err)
						}else{
							fmt.Println("InsertCols 100 success")
						}
					}else{
						if int(e.Cap) + len(data) > 100{
							values := append(data, (e.Data)...)
							data = data[0:0]
							err,_ := hsmysql.InsertCols(h.db,"insert into Redis_Cache2 (user,gmt_create,gmt_modify) values ",values)
							if err != nil{
								if maxPriority < 10{
									e := HsChanElem{
										StateMent: "",
										Data:      values,
										Priority:  maxPriority,
										Cap:       uint16(len(values)),
									}
									h.list <- e
									fmt.Println("InsertCols 100 success but reinsert")
								}else{
									fmt.Println("InsertCols 100:",err)
								}
							}
							maxPriority = 0
						}else{
							data = append(data, (e.Data)...)
							if e.Priority < maxPriority{
								maxPriority = e.Priority
							}
						}
					}
				case <-time.After(time.Second*10):
					if len(data) > 0{
						err,_ := hsmysql.InsertCols(h.db,"insert into Redis_Cache2 (user,gmt_create,gmt_modify) values ",data)
						if err != nil{
							if maxPriority < 10{
								e := HsChanElem{
									StateMent: "",
									Data:      data,
									Priority:  maxPriority,
									Cap:       uint16(len(data)),
								}
								h.list <- e
								fmt.Println("InsertCols timeout but reinsert")
							}else{
								fmt.Println("InsertCols timeout success!")
							}
						}
						data = data[0:0]
						maxPriority = 0
					}
				case <-h.close:
					break
				}
			}
		}()
	}
}

//添加到chan中
func (h *HsChan)Insert(data pkg.RedisCache, stateMent string){
	e := HsChanElem{
		StateMent:	stateMent,
		Data:		[]pkg.RedisCache{data},
		Priority:	1,
		Cap:		1,
	}

	h.list <- e
}

//退出
func (h *HsChan)Close(){
	h.close <- true
}