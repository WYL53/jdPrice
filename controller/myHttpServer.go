package controller

import (
	"fmt"
	"os"
	"net/http"
)

type recoverHttpServer struct {
	next http.Handler
}

func newRecoverHttpServer(next http.Handler) *recoverHttpServer {
	return &recoverHttpServer{
		next:next,
	}
}

func (this *recoverHttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	defer func() {
		if err := recover();err != nil{
			fmt.Printf("handle request error:%v",err)
		}
	}()
	this.next.ServeHTTP(w,r)
}

func StartHttpServer(prot int) {
	addr := fmt.Sprintf("0.0.0.0:%d", prot)
	fmt.Println("http server listen address:", addr)
	http.HandleFunc("/", IndexServer)
	http.HandleFunc("/addModel", AddModelServer)
	http.HandleFunc("/delModel", DelModelServer)
	http.HandleFunc("/updatePrice", UpdatePriceServer)
	http.HandleFunc("/jd", HomeServer)
	http.HandleFunc("/price", PriceServer)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	server := newRecoverHttpServer(http.DefaultServeMux)
	err := http.ListenAndServe(addr, server)
	if err != nil {
		fmt.Printf("start http server error:%s", err.Error())
		os.Exit(1)
	}
}

