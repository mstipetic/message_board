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
	"time"
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
	Id int64 `json:"id"`
	Author string `json:"author"`
	Text string `json:"text"`
	Time int64 `json:"timestamp"`
}

type Post struct {
	Id int64 `json:"id"`
	Author string `json:"author"`
	Text string `json:"text"`
	Url string `json:"url"`
	Time int64 `json:"timestamp"`
	Comments []Comment `json:"comments"`
}

func newPostHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	n, _ := rds.Incr("id_counter")
	fmt.Println("returned identifier: ", n)
	fmt.Println("Body: " + string(b))
	var post Post
	json.Unmarshal([]byte(string(b)), &post)
	fmt.Println("Unmarshaled: ", post)
	post.Author = "mislav"
	post.Id = n
	currentTime := time.Now().Unix()
	post.Time = currentTime
	responseTextByte, _ := json.Marshal(post)
	responseText := string(responseTextByte)
	rds.Zadd("posts", float64(currentTime), string(n))
	rds.Set("post:" + string(n), responseText)
	fmt.Println("returning: ", responseText)
	fmt.Fprint(w, string(responseText))
}

func newCommentHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	n, _ := rds.Incr("id_counter")
	fields := mux.Vars(r)
	post_id := fields["id"]
	fmt.Println("returned identifier: ", n)
	fmt.Println("Body: " + string(b))
	var comment Comment 
	json.Unmarshal([]byte(string(b)), &comment)
	fmt.Println("Unmarshaled: ", comment)
	comment.Author = "mislav"
	currentTime := time.Now().Unix()
	comment.Time = currentTime
	comment.Id = n
	responseTextByte, _ := json.Marshal(comment)
	responseText := string(responseTextByte)
	rds.Lpush("post:" + post_id + ":comments", responseText)
	fmt.Println("returning: ", responseText)
	fmt.Fprint(w, string(responseText))
}

func getPostsHandler(w http.ResponseWriter, r *http.Request) {
	
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
	router.HandleFunc("/post", newPostHandler).Methods("POST")
	router.HandleFunc("/post", getPostsHandler).Methods("GET")
	router.HandleFunc("/post/{id}/comments", newCommentHandler)
	router.HandleFunc("/post/{id}", postHandler)
	router.Handle("/ws/", websocket.Handler(websocket_handler))
	http.ListenAndServe(":8000", router)
}
