package main

import (
	"html/template"
	"net/http"
	//"github.com/simonz05/godis/redis"
	"github.com/vmihailenco/redis"
	"github.com/gorilla/mux"
	"fmt"
	"code.google.com/p/go.net/websocket"
	"io/ioutil"
	"encoding/json"
	"time"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	t,_ := template.ParseFiles("index.html")
	t.Execute(w, 5)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("zatrazeni path: ", vars["path"])
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
	n := rds.Incr("id_counter").Val()
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
	rds.ZAdd("posts", redis.Z{float64(currentTime), strconv.FormatInt(n, 10)})
	rds.Set("post:" + strconv.FormatInt(n, 10), responseText)
	fmt.Println("returning: ", responseText)
	fmt.Fprint(w, string(responseText))
}

func newCommentHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	n := rds.Incr("id_counter").Val()
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
	rds.LPush("post:" + post_id + ":comments", responseText)
	fmt.Println("returning: ", responseText)
	fmt.Fprint(w, string(responseText))
}

func getPostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	posts := rds.ZRange("posts", 0, 10)
	fmt.Println("u getPostsHandler")
	posts_arr := []Post{}
	for _, post_id := range posts.Val() {
		post_rt := rds.Get("post:" + post_id)
		var post Post
		json.Unmarshal([]byte(string(post_rt.Val())), &post)
		fmt.Println(post)
		
		comments := rds.LRange("post:" + post_id + ":comments", 0, -1)
		fmt.Println(comments.Val())
		for _, comment_str := range comments.Val() {
			var comment Comment
			json.Unmarshal([]byte(comment_str), &comment)
			fmt.Print("Dodajem komentar:   ")
			fmt.Println(comment)
			post.Comments = append(post.Comments, comment)
		}

		posts_arr = append(posts_arr, post)
	}
	marshaled, _ := json.Marshal(posts_arr)
	fmt.Fprint(w, string(marshaled))
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


var rds *redis.Client = redis.NewTCPClient("localhost:6379", "", 0)


func main() {
	defer rds.Close()
	router := mux.NewRouter()
	fmt.Println(rds)
	router.HandleFunc("/", handler)
	router.HandleFunc("/static/{path:.*}", staticHandler)
	router.HandleFunc("/post", newPostHandler).Methods("POST")
	router.HandleFunc("/post", getPostsHandler).Methods("GET")
	router.HandleFunc("/post/{id}/comments", newCommentHandler)
	router.HandleFunc("/post/{id}", postHandler)
	router.Handle("/ws/", websocket.Handler(websocket_handler))
	http.ListenAndServe(":8000", router)
}
