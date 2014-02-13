package main

/*
* Link Shortener, with a Redis backend.
*
* Released under and MIT License, please see the LICENSE.md file.
*
* John Nye
*
 */
import (
	"./utils"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var host = flag.String("h", "localhost", "Bind address to listen on")
var base = flag.String("b", "http://localhost/", "Base URL for the shortener")
var port = flag.String("p", "8080", "Port you want to listen on, defaults to 8080")
var maxConnections = flag.Int("c", 512, "The maximum number of active connections") //Currently Not Used
var redisConn = flag.String("r", "localhost:6379", "Redis Address, defaults to localhost:6379")

type Data struct {
	Original  string
	Short     string
	FullShort string
	HitCount  int
}

var redisPool = &redis.Pool{
	MaxIdle:   3,
	MaxActive: 50, // max number of connections
	Dial: func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", *redisConn)
		if err != nil {
			panic(err.Error())
		}
		return c, err
	},
}

func handler(w http.ResponseWriter, r *http.Request) {

	log.Println(r.UserAgent())

	type NewURL struct {
		URL string
	}
	var url NewURL
	var domain Data

	conn := redisPool.Get()

	if r.Method == "GET" {
		domain = getLongURL(r.URL.Path[1:], conn)
		if len(domain.Original) > 0 {
			http.Redirect(w, r, domain.Original, http.StatusFound)
			return
		}
		http.ServeFile(w, r, "./index.html")		
		log.Println("Served Homepage")
		return
	}

	create, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(create, &url)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	search := strings.Join([]string{"*||", url.URL}, "")

	keys, err := redis.Strings(conn.Do("KEYS", search))

	if len(keys) < 1 {
		domain = createShortURL(url.URL, conn)
	} else {
		domain = getInfoForKey(keys[0], conn)
	}

	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(domain)

	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "%s", output)
	conn.Close()
}

func getInfoForKey(key string, conn redis.Conn) Data {
	var d Data
	parts := strings.Split(key, "||")
	d.Short = parts[0]
	d.Original = parts[1]
	d.FullShort = strings.Join([]string{*base, parts[0]}, "")
	newCount, err := redis.Int(conn.Do("HGET", key, "count"))
	if err != nil {
		log.Print(err)
	}
	d.HitCount = newCount
	return d
}

func createShortURL(url string, conn redis.Conn) Data {
	var d Data
	count, err := redis.Int(conn.Do("INCR", "global:size"))
	if err != nil {
		log.Print(err)
		return d
	}
	log.Print("Total: ",count)
	encodedVar := base62.EncodeInt(int64(count))
	key := strings.Join([]string{encodedVar, url}, "||")
	conn.Send("MULTI")
	conn.Send("HSET", key, "count", 0)
	_, err2 := conn.Do("EXEC")

	if err2 != nil {
		log.Print(err)
		return d
	}

	d.Original = url
	d.HitCount = 0
	d.Short = encodedVar
	d.FullShort = strings.Join([]string{*base, encodedVar}, "")

	return d
}

func getLongURL(short string, conn redis.Conn) Data {
	var d Data

	search := strings.Join([]string{short, "||*"}, "")
	fmt.Println(search)
	n, err := redis.Strings(conn.Do("KEYS", search))

	if err != nil {
		log.Print(err)
		return d
	}

	if len(n) < 1 {
		log.Print("Nothing Found")
	} else {
		parts := strings.Split(n[0], "||")

		d.Short = parts[0]
		d.Original = parts[1]
		d.FullShort = strings.Join([]string{*base, parts[0]}, "")
		newCount, err := redis.Int(conn.Do("HINCRBY", n[0], "count", 1))
		if err != nil {
			log.Println(err)
		}
		d.HitCount = newCount
	}
	log.Println("Served: ",d.Original)
	return d
}

func main() {
	flag.Parse()

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(*host+":"+*port, nil)
	if err != nil {
		fmt.Println(err)
	}
}
