/*
	简易打log模块
	使用前声明一个Hslogger变量，通过该变量来输出日志
	ex:
	var LOG hslog.Hslogger  //定义LOG变量
	LOG.HslogInit()         //初始化
	defer LOG.HslogClose()  //析构
	LOG.Printf()
	LOG.Println()

*/
package hslog

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

//对外提供日志结构体
type Hslogger struct {
	out      io.Writer //抽象写模块
	filename string    //日志文件名(包含路径)
}

//初始化对外结构体
func (l *Hslogger) HslogInit(args ...interface{}) {
	//处理参数，判断文件名等操作
	filename := fmt.Sprint(args)
	filename = strings.Replace(filename, " ", "", -1)
	filename = strings.Replace(filename, "\n", "", -1)
	filename = filename[1 : len(filename)-1]
	l.filename = filename
	writer := rotateNew(l.filename)
	l.out = writer
}

//关闭log，并退出
func (l *Hslogger) HslogClose() {
	var p *hsRotate = (l.out).(*hsRotate)

	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.fileHandler.Sync()
	p.fileHandler.Close()
}

//输出格式日志
func (l *Hslogger) Printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args) //在这里转进来是args变成里[]interface{},要加上...把[]interface{}中的元素拆分出来
	t, m, s := time.Now().Clock()
	logmsg := fmt.Sprintf("[%02d:%02d:%02d] %s\r\n", t, m, s, msg)

	l.out.Write([]byte(logmsg))
}

//输出完整log
func (l *Hslogger) Println(args ...interface{}) error {
	msg := fmt.Sprint(args...)
	msg = msg[1 : len(msg)-1]

	t, m, s := time.Now().Clock()
	logmsg := fmt.Sprintf("[%02d:%02d:%02d] %s\r\n", t, m, s, msg)

	_, err := l.out.Write([]byte(logmsg))
	return err
}

//核心结构体，抽象文件读写操作和筛选文件操作
type hsRotate struct {
	mutex       sync.RWMutex
	logFile     string
	fileHandler *os.File
}

//简单初始化核心结构体
func rotateNew(filenName string) *hsRotate {
	var rl hsRotate
	rl.logFile = filenName
	return &rl
}

//自定义读取文件
func (r *hsRotate) getFile() (*os.File, error) {
	now := time.Now()
	if r.logFile != "" { //已经设置文件名
		if r.fileHandler == nil { //打开文件
			fh, err := os.OpenFile(r.logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0774)
			if err != nil {
				return nil, err
			}
			r.fileHandler = fh
		}
		return r.fileHandler, nil
	} else { //使用默认的文件名
		filenName := fmt.Sprintf("%04d%02d%02d.log", now.Year(), now.Month(), now.Day())
		if r.fileHandler == nil { //第一次启动
			fh, err := os.OpenFile(filenName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0774)
			if err != nil {
				return nil, err
			}
			r.fileHandler = fh
			return r.fileHandler, nil
		} else {
			if filenName == r.fileHandler.Name() { //不需要更换log文件
				return r.fileHandler, nil
			} else {
				fh, err := os.OpenFile(filenName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0774)
				if err != nil {
					return nil, err
				}

				r.fileHandler.Close()
				r.fileHandler = fh
				return r.fileHandler, nil
			}
		}
	}
}

//写日志操作
func (r *hsRotate) Write(p []byte) (n int, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	out, err := r.getFile()
	if err != nil {
		return -1, err
	}

	n, err = out.Write(p)
	if err != nil {
		fmt.Printf("write log err : %s \n", err)
		r.fileHandler.Close()
		r.fileHandler = nil
	}

	return n, err
}
