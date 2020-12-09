package hshttp

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"
	"context"
)

//-------------------------------------------------------
// HttpRunSafe
// ------------------------------------------------------
//
//Return Value: None
//Params:
//		1)addr string:			http sever bind addr
//		2)Handler http.Handler: http mux
//Safe start/Close http server
//------------------------------------------------------
func HttpRunSafe(addr string, Handler http.Handler){
	server := &http.Server{
		Addr:            addr,
		Handler:         Handler,
	}
	go func() {
		if err := server.ListenAndServe();err!=nil&&err!=http.ErrServerClosed{
			log.Fatal("listen:%s\n",err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<-quit
	fmt.Println("shutting down server...")

	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx);err!=nil{
		fmt.Println("server forced to shutdown:",err)
	}
	fmt.Println("server exiting")
}