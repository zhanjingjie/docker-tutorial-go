package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gomodule/redigo/redis"
)

// PageData is the data needed for the page.
type PageData struct {
	Name     string
	HostName string
	Visits   int
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   10,
		MaxActive: 20, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(os.Getenv("REDIS_CON_TYPE"), os.Getenv("REDIS_CONNECTION"))
			if err != nil {
				panic(err)
			}
			return c, err
		},
	}
}

func main() {
	fmt.Println("Started service.")

	// Start Redis client connection.
	var pool = newPool()
	c := pool.Get()
	defer c.Close()

	t := template.Must(template.ParseFiles("hello.html"))
	// Changed the default URL from "/" to "/view", because when testing in Chrome,
	// the handler for "/" will be called twice.
	// The first time it's called is because Chrome requested /favicon.ico, for the
	// website icon. And it's confusing and will make the counter wrong.
	// The second time the browser requested the website at "/".
	http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
		hostName, _ := os.Hostname()
		name := os.Getenv("NAME")
		counter, err := redis.Int(c.Do("INCR", "counter"))
		if err != nil {
			panic(err)
		}

		data := PageData{
			Name:     name,
			HostName: hostName,
			Visits:   counter,
		}
		t.Execute(w, data)
	})
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
