package tools

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//mysql相关测评对象
type mysqlInfo struct {
	/*基本信息*/
	Questions       int32     //语句执行总数
	Uptime    		int32     //总运行时间

	/*测试重点*/
	TCH       		float32   //TCH值(连接缓存比)
	QPS       		float32   //QPS值(查询处理量)
	TPS       		float32   //TPS值(事务处理量)

	/*连接池相关*/
	ThreadConnected int32     //当前线程连接数
	ThreadCreated   int32     //总创建线程数
	Connections     int32     //总连接数
	ConnectAborted  int32     //连接失败数
	MaxConnections  int32     //最大连接数
	MaxConnectUse   int32     //已经使用的最大连接数

	/*事务相关，针对innodb*/
	Commit    		int32     //事务提交总数
	Rollback  		int32     //事务回滚总数

	/*读:select(命中和非命中) 写:insert、update、delete等*/
	ReadNum   		int32     //总查询数(命中和非命中)
	WriteNum  		int32     //总写数

	/*特殊的查询*/
	SlowQuery 		int32     //慢查询数
	FullJoin  		int32     //full_join查询数
	SelectScan      int32     //全表扫描

	/*innodb buffer Pool缓存所有索引、数据等记录*/
	InnodbPoolSize  int32     //innodb buffer Pool 缓冲池的大小
	InnodbReadDisk  int32     //缓存池无法满足的请求数
	InnodbReadPool  int32     //缓存池满足的请求数

	/*OpenTableCache应该大于OpenTables且小于OpenedTables*/
	OpenTables      int32      //当前打开的表数量
	OpenedTables    int32      //总共打开的表数量
	OpenTableCache  int32      //表缓存大小

	/*临时表的数量*/
	TmpTable        int32     //内存临时表
	TmpDiskTable    int32     //磁盘临时表

	/*innodb 行锁相关*/
	InnodbRows      int16     //当前等待锁的行数
	InnodbRoWaits   int32     //等待行锁必须的时间
	InnodbRowAvg    int32     //等待行锁的平均时间

	/*Qcache结果集缓存(关闭)*/
	QcacheHits     int16      //Qcache缓存命中数
	QcacheInserts  int16      //插入Qcache缓存数

	/*key buffer相关，针对MyISAM表*/
	KeyBufRReq    int32       //读索引请求数
	KeyBufRDisk   int32       //从磁盘IO读索引请求数
	KeyBufWReq    int32       //写索引请求数
	KeyBufWDisk   int32       //从磁盘IO写索引请求数

	/*发送接受数据大小*/
	BytesReceive  int64       //客户端发送的数据大小
	BytesSend     int64       //服务器端发送给客户端的数据大小
}

//全局记录
var MysqlInfo mysqlInfo = mysqlInfo{}

/*得到mysql TCH(线程新创建占比)
TCH<90% 左右建议修改Thread_cache*/
func (q * mysqlInfo)getTCH(db * sql.DB)float32{
	var agment string
	var AbortedClientNum,AbortedConntedNum int32 = 0,0

	err := db.QueryRow("show global status like \"Threads_created\";").Scan(&agment,&q.ThreadCreated)
	if err != nil{
		q.ThreadCreated = -1
	}

	err = db.QueryRow("show global status like \"Threads_connected\";").Scan(&agment,&q.ThreadConnected)
	if err != nil{
		q.ThreadConnected = -1
	}

	err = db.QueryRow("show global variables like \"max_connections\";").Scan(&agment,&q.MaxConnections)
	if err != nil{
		q.MaxConnections = -1
	}

	err = db.QueryRow("show global status like \"Max_used_connections\";").Scan(&agment,&q.MaxConnectUse)
	if err != nil{
		q.MaxConnectUse = -1
	}

	err = db.QueryRow("show global status like \"Connections\";").Scan(&agment,&q.Connections)
	if err != nil{
		q.Connections = -1
	}

	err = db.QueryRow("show global status like \"Aborted_clients\";").Scan(&agment,&AbortedClientNum)
	if err != nil{
		AbortedClientNum = -1
	}

	err = db.QueryRow("show global status like \"Aborted_connects\";").Scan(&agment,&AbortedConntedNum)
	if err != nil{
		AbortedConntedNum = -1
	}

	q.ConnectAborted = AbortedConntedNum + AbortedClientNum

	if q.Connections > 0{
		return 1-float32(q.ThreadCreated)/float32(q.Connections)
	}else{
		return -1
	}
}

/*获取QPS(每秒查询(所有包含select等)处理(成功)率)*/
func (q * mysqlInfo)getQPS(db *sql.DB)float32{
	var agment string
	var QuestionsNum,UptimeNum int32 = 0,0

	err := db.QueryRow("show global status like \"Questions\";").Scan(&agment,&QuestionsNum)
	if err != nil{
		QuestionsNum = -1
	}

	err = db.QueryRow("show global status like \"Uptime\";").Scan(&agment,&UptimeNum)
	if err != nil{
		UptimeNum = -1
	}

	if UptimeNum > 0{
		if q.Questions == 0{
			q.Questions = QuestionsNum
			q.Uptime = UptimeNum
			return float32(QuestionsNum)/float32(UptimeNum)
		}else{
			var result float32 = 0
			//fmt.Println(QuestionsNum,q.Questions,UptimeNum,q.Uptime)
			if UptimeNum == q.Uptime{
				result = float32(QuestionsNum-q.Questions)/float32(q.Uptime)
			}else{
				result = float32(QuestionsNum-q.Questions)/float32(UptimeNum-q.Uptime)
			}
			q.Questions = QuestionsNum
			q.Uptime = UptimeNum
			return result
		}
	}else{
		return -1
	}
}

/*获取总读写数*/
func (q * mysqlInfo)getRWNum(db *sql.DB)(int32,int32){
	var agment string
	var QcacheHitsNum,SelectNum,InsertNum int32 = 0,0,0
	var UpdateNum,DeleteNum,ReplaceNum int32 = 0,0,0

	/*查询命中数*/
	err := db.QueryRow("show global status like \"Qcache_hits\";").Scan(&agment,&QcacheHitsNum)
	if err != nil{
		QcacheHitsNum = -1
	}

	/*查询未命中数*/
	err = db.QueryRow("show global status like \"Com_select\";").Scan(&agment,&SelectNum)
	if err != nil{
		SelectNum = -1
	}

	err = db.QueryRow("show global status like \"Com_insert\";").Scan(&agment,&InsertNum)
	if err != nil{
		InsertNum = -1
	}

	err = db.QueryRow("show global status like \"Com_update\";").Scan(&agment,&UpdateNum)
	if err != nil{
		UpdateNum = -1
	}

	err = db.QueryRow("show global status like \"Com_delete\";").Scan(&agment,&DeleteNum)
	if err != nil{
		DeleteNum = -1
	}

	err = db.QueryRow("show global status like \"Com_replace\";").Scan(&agment,&ReplaceNum)
	if err != nil{
		ReplaceNum = -1
	}

	return QcacheHitsNum+SelectNum,UpdateNum+DeleteNum+ReplaceNum+InsertNum
}

/*获取TPS(每秒处理事务数)*/
func (q * mysqlInfo)getTPS(db *sql.DB)float32{
	var agment string
	var commitNum,rollbackNum,UptimeNum int32 = 0,0,0

	err := db.QueryRow("show global status like \"Com_commit\";").Scan(&agment,&commitNum)
	if err != nil{
		commitNum = -1
	}

	err = db.QueryRow("show global status like \"Com_rollback\";").Scan(&agment,&rollbackNum)
	if err != nil{
		rollbackNum = -1
	}

	err = db.QueryRow("show global status like \"Uptime\";").Scan(&agment,&UptimeNum)
	if err != nil{
		UptimeNum = 0
	}

	if UptimeNum != 0{
		if q.Commit == 0{
			q.Commit = commitNum
			q.Rollback = rollbackNum
			if UptimeNum == q.Uptime{
				return (float32(commitNum)+float32(rollbackNum))/float32(q.Uptime)
			}else{
				return (float32(commitNum)+float32(rollbackNum))/float32(UptimeNum-q.Uptime)
			}
		}else{
			var result float32 = 0
			if UptimeNum == q.Uptime{
				result = (float32(commitNum-q.Commit)+float32(rollbackNum-q.Rollback))/float32(q.Uptime)
			}else{
				result = (float32(commitNum-q.Commit)+float32(rollbackNum-q.Rollback))/float32(UptimeNum-q.Uptime)
				q.Uptime = UptimeNum
			}
			q.Commit = commitNum
			q.Rollback = rollbackNum
			return result
		}
	}else{
		return -1
	}
}

/*记录慢查询、full_join和SelectScan的语句*/
func (q * mysqlInfo)getSlow(db * sql.DB){
	var agment string

	err := db.QueryRow("show global status like \"Slow_queries\";").Scan(&agment,&q.SlowQuery)
	if err != nil{
		q.SlowQuery = -1
	}

	err = db.QueryRow("show global status like \"Select_full_join\";").Scan(&agment,&q.FullJoin)
	if err != nil{
		q.FullJoin = -1
	}

	err = db.QueryRow("show global status like \"Select_scan\";").Scan(&agment,&q.SelectScan)
	if err != nil{
		q.SelectScan = -1
	}
}

/*获取Innodb缓存池的相关信息*/
func (q * mysqlInfo)getInnodbBuffer(db * sql.DB){
	var agment string
	/*缓存区总大小*/
	err := db.QueryRow("show global variables like \"innodb_buffer_pool_size\";").Scan(&agment,&q.InnodbPoolSize)
	if err != nil{
		q.InnodbPoolSize = -1
	}

	/*缓存区无法满足的请求数，此请求直接在磁盘操作*/
	err = db.QueryRow("show global status like \"innodb_buffer_pool_reads\";").Scan(&agment,&q.InnodbReadDisk)
	if err != nil{
		q.InnodbReadDisk = -1
	}

	/*缓存区能满足的请求数*/
	err = db.QueryRow("show global status like \"innodb_buffer_pool_read_requests\";").Scan(&agment,&q.InnodbReadPool)
	if err != nil{
		q.InnodbReadPool = -1
	}
}

/*获取表缓存相关信息*/
func (q * mysqlInfo)getTableOpen(db * sql.DB){
	var agment string
	/*打开表缓存大小*/
	err := db.QueryRow("show global variables like \"table_open_cache\";").Scan(&agment,&q.OpenTableCache)
	if err != nil{
		q.OpenTableCache = -1
	}

	err = db.QueryRow("show global status like \"open_tables\";").Scan(&agment,&q.OpenTables)
	if err != nil{
		q.OpenTables = -1
	}

	err = db.QueryRow("show global status like \"opened_tables\";").Scan(&agment,&q.OpenedTables)
	if err != nil{
		q.OpenedTables = -1
	}

	err = db.QueryRow("show global status like \"Created_tmp_tables\";").Scan(&agment,&q.TmpTable)
	if err != nil{
		q.TmpTable = -1
	}

	err = db.QueryRow("show global status like \"Created_tmp_disk_tables\";").Scan(&agment,&q.TmpDiskTable)
	if err != nil{
		q.TmpDiskTable = -1
	}

}

/*获取Innodb行锁相关信息*/
func (q * mysqlInfo)getRowLock(db * sql.DB){
	var agment string

	/*需要等待行锁的总行数*/
	err := db.QueryRow("show global status like \"Innodb_row_lock_current_waits\";").Scan(&agment,&q.InnodbRows)
	if err != nil{
		q.InnodbRows = -1
	}

	/*当前行锁等待的必须时间*/
	err = db.QueryRow("show global status like \"Innodb_row_lock_waits\";").Scan(&agment,&q.InnodbRoWaits)
	if err != nil{
		q.InnodbRoWaits = -1
	}

	/*当前行锁等待的评价时间*/
	err = db.QueryRow("show global status like \"Innodb_row_lock_time_avg\";").Scan(&agment,&q.InnodbRowAvg)
	if err != nil{
		q.InnodbRowAvg = -1
	}
}

/*获取Qcache缓存相关*/
func (q * mysqlInfo)getQcache(db * sql.DB){
	var agment string

	err := db.QueryRow("show global status like \"Qcache_hits\";").Scan(&agment,&q.QcacheHits)
	if err != nil{
		q.QcacheHits = -1
	}

	err = db.QueryRow("show global status like \"Qcache_inserts\";").Scan(&agment,&q.QcacheInserts)
	if err != nil{
		q.QcacheInserts = -1
	}
}

/*获取Key Buffer相关*/
func (q * mysqlInfo)getKeyBuffer(db * sql.DB){
	var agment string

	err := db.QueryRow("show global status like \"Key_read_requests\";").Scan(&agment,&q.KeyBufRReq)
	if err != nil{
		q.KeyBufRReq = -1
	}

	err = db.QueryRow("show global status like \"Key_reads\";").Scan(&agment,&q.KeyBufRDisk)
	if err != nil{
		q.KeyBufRDisk = -1
	}

	err = db.QueryRow("show global status like \"Key_write_requests\";").Scan(&agment,&q.KeyBufWReq)
	if err != nil{
		q.KeyBufWReq = -1
	}

	err = db.QueryRow("show global status like \"Key_writes\";").Scan(&agment,&q.KeyBufWDisk)
	if err != nil{
		q.KeyBufWDisk = -1
	}
}

/*收发数据大小*/
func (q * mysqlInfo)getDataSize(db * sql.DB)(int64,int64){
	var agment string
	var ReceivedSize,SentSize int64 = 0,0

	err := db.QueryRow("show global status like \"Bytes_received\";").Scan(&agment,&ReceivedSize)
	if err != nil{
		ReceivedSize = 0
	}

	err = db.QueryRow("show global status like \"Bytes_sent\";").Scan(&agment,&SentSize)
	if err != nil{
		SentSize = 0
	}

	if ReceivedSize > 0{
		Tmp := ReceivedSize
		ReceivedSize = ReceivedSize - q.BytesReceive
		q.BytesReceive = Tmp
	}

	if SentSize > 0{
		Tmp := SentSize
		SentSize = SentSize - q.BytesSend
		q.BytesSend = Tmp
	}

	return SentSize,ReceivedSize
}


/*显示相关参数*/
func ShowMysqlInfo(db * sql.DB){
	MysqlInfo.TCH = MysqlInfo.getTCH(db)
	MysqlInfo.QPS = MysqlInfo.getQPS(db)
	MysqlInfo.TPS = MysqlInfo.getTPS(db)
	MysqlInfo.ReadNum,MysqlInfo.WriteNum = MysqlInfo.getRWNum(db)
	MysqlInfo.getSlow(db)
	MysqlInfo.getInnodbBuffer(db)
	MysqlInfo.getTableOpen(db)
	MysqlInfo.getRowLock(db)
	Received,Sent := MysqlInfo.getDataSize(db)

	fmt.Println("-----------------------------------------------------")
	fmt.Printf("runtime(h):%d\n",MysqlInfo.Uptime/60/60)
	fmt.Printf("ThreadConnected:%d Failed:%d maxConnections:%d used:%d TCH:%f\n",MysqlInfo.ThreadConnected,MysqlInfo.ConnectAborted,MysqlInfo.MaxConnections,MysqlInfo.MaxConnectUse,MysqlInfo.TCH)
	fmt.Printf("Questions:%d QPS(c/s):%f\n",MysqlInfo.Questions,MysqlInfo.QPS)
	fmt.Printf("Commit:%d Rollback:%d TPS(c/s):%f\n",MysqlInfo.Commit,MysqlInfo.Rollback,MysqlInfo.TPS)
	fmt.Printf("R:%d W:%d\n",MysqlInfo.ReadNum,MysqlInfo.WriteNum)
	fmt.Printf("Qcache hits:%d inserts:%d\n",MysqlInfo.QcacheHits,MysqlInfo.QcacheInserts)
	fmt.Printf("slow queries:%d full join select:%d Select All Scan:%d\n",MysqlInfo.SlowQuery,MysqlInfo.FullJoin,MysqlInfo.SelectScan)
	fmt.Printf("Innodb buffer pool size(M):%d 命中率:%f\n",MysqlInfo.InnodbPoolSize/1024/1024,float32(MysqlInfo.InnodbReadPool)/float32(MysqlInfo.InnodbReadDisk+MysqlInfo.InnodbReadPool))
	fmt.Printf("Key Buffer RReq:%d RDisek:%d WReq:%d WDisek:%d\n", MysqlInfo.KeyBufRReq,MysqlInfo.KeyBufRDisk,MysqlInfo.KeyBufWReq,MysqlInfo.KeyBufWDisk)
	fmt.Printf("OpenTables:%d OpenedTable:%d OpenTableCache:%d\n",MysqlInfo.OpenTables,MysqlInfo.OpenedTables,MysqlInfo.OpenTableCache)
	fmt.Printf("TmpTalble:%d TmpDiskTable:%d\n",MysqlInfo.TmpTable,MysqlInfo.TmpDiskTable)
	fmt.Printf("CurrentLock:%d CurrenTime(ms):%d AvgTime(ms):%d\n",MysqlInfo.InnodbRows,MysqlInfo.InnodbRoWaits,MysqlInfo.InnodbRowAvg)
	fmt.Printf("Receive(KB):%d Send(KB):%d\n",Received/1024,Sent/1024)
}

func MysqlInof(){
	//1.连接数据库
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

	for{
		ShowMysqlInfo(db)
		time.Sleep(time.Second*5)
	}
}
