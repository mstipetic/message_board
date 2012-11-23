package main

import (
	"html/template"
	"net/http"
	"github.com/simonz05/godis/redis"
	"github.com/gorilla/mux"
	"fmt"
	"code.google.com/p/go.net/websocket"
	"io/ioutil"
	"encoding/json"
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

func postHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "1")
}

type Comment struct {
	Author string `json:"author"`
	Text string `json:"text"`
}

type Post struct {
	Id int64 `json:"id"`
	Author string `json:"author"`
	Text string `json:"text"`
	Comments []Comment `json:"comments"`
}

func newPostHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	n, _ := rds.Incr("id_counter")
	fmt.Println("Body: " + string(b))
	var post Post
	json.Unmarshal([]byte(string(b)), &post)
	fmt.Println("Unmarshaled: ", post)
	post.Author = "mislav"
	post.Id = n
	responseText, _ := json.Marshal(post)
	fmt.Fprint(w, string(responseText))
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


var rds *redis.Client = redis.New("", 0, "")


func main() {
	router := mux.NewRouter()
	fmt.Println(rds)
	router.HandleFunc("/", handler)
	router.HandleFunc("/static/{path}", staticHandler)
	router.HandleFunc("/post", newPostHandler)
	router.HandleFunc("/post/{id}", postHandler)
	router.Handle("/ws/", websocket.Handler(websocket_handler))
	http.ListenAndServe(":8000", router)
}
