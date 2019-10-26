package main

import (
	_ "my_go_learn/memory"  //这里修改成你存放menory.go相应的目录
	"my_go_learn/session" //这里修改成你存放session.go相应的目录
	"fmt"
	"log"
	"net/http"
)

var globalSessions *session.Manager

func init() {
	var err error
	globalSessions, err = session.NewSessionManager("memory", "goSessionid", 3600)
	if err != nil {
		fmt.Println(err)
		return
	}
	go globalSessions.GC()
}

func sayHelloHandler(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("name")
	if err == nil {
		fmt.Println(cookie.Value)
		fmt.Println(cookie.Domain)
		fmt.Println(cookie.Expires)
	}
	//fmt.Fprintf(w, "Hello world!\n") //这个写入到w的是输出到客户端的
}
func login(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	val := sess.Get("username")
	if val != nil {
		fmt.Println(val)
	} else {
		sess.Set("username", "jerry")
		fmt.Println("set session")
	}
}
func loginOut(w http.ResponseWriter, r *http.Request) {
	//销毁
	globalSessions.SessionDestroy(w, r)
	fmt.Println("session destroy")
}

func main() {
	http.HandleFunc("/", sayHelloHandler) //	设置访问路由
	http.HandleFunc("/login", login)
	http.HandleFunc("/loginout", loginOut) //销毁
	log.Fatal(http.ListenAndServe(":9090", nil))
}
