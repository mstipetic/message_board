package main

import (
	"html/template"
	"net/http"
	"github.com/simonz05/godis/redis"
	"github.com/gorilla/mux"
	"fmt"
	"code.google.com/p/go.net/websocket"
)

func handler(w http.ResponseWriter, r *http.Request) {
	t,_ := template.ParseFiles("index.html")
	t.Execute(w, 5)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := "static/" + string(vars["path"])
	fmt.Println(path)
	http.ServeFile(w, r, path)
}

func websocket_handler(ws *websocket.Conn) {
	for {
		var buf string;
		e := websocket.Message.Receive(ws, &buf)
		if (e != nil) {
			fmt.Println(e)
			break
		}

		e = websocket.Message.Send(ws, buf)
		if (e != nil) {
			fmt.Println(e)
			break
		}
	}
}

func main() {
	r := redis.New("", 0, "")
	router := mux.NewRouter()
	fmt.Println(r)
	router.HandleFunc("/", handler)
	router.HandleFunc("/static/{path}", staticHandler)
	router.Handle("/ws/", websocket.Handler(websocket_handler))
	http.ListenAndServe(":8000", router)
}
