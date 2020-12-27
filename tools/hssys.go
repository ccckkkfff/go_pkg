package tools

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

var LastTotalFreed uint64


//打印内存使用信息(堆)
//Alloc=申请 TotalAlloc=总申请 Just Freed=释放 Sys=系统申请 NumGc=GC次数
func PrintMemStats(){
	go func(){
		for{
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			fmt.Printf("Alloc=%vKB TotalAlloc=%vKB Just Freed=%vKB Sys=%vKB NumGc=%v\n",
				m.Alloc/1024,m.TotalAlloc/1024,(m.TotalAlloc-m.Alloc-LastTotalFreed)/1024,
				m.Sys/1024,m.NumGC)

			LastTotalFreed = m.TotalAlloc - m.Alloc
			time.Sleep(time.Second*10)
		}
	}()
}


//获取Heap信息,通过pprof查看
func CatchRoutineMsg(){
	fileName := fmt.Sprintf("heap_%s.pprof",time.Now().Format("2006_01_02_03_04_05"))
	f,err := os.OpenFile(fileName,os.O_RDWR|os.O_CREATE,0777)
	if err != nil{
		fmt.Println("open",fileName,":",err)
		return
	}
	defer f.Close()
	pprof.Lookup("heap").WriteTo(f,1)
}


//获取当前时间
func GetLocalTime()string{
	tm := time.Now()
	return tm.Format("2006-01-02 15:04:05")
}