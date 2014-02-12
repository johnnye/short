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
	"flag"
	"./utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"net/http"
	"github.com/garyburd/redigo/redis"
)

var host = flag.String("h", "localhost", "Bind address to listen on")
var base = flag.String("b", "http://localhost/", "Base URL for the shortener")
var port = flag.String("p", "8080", "Port you want to listen on, defaults to 8080")
var maxConnections = flag.Int("c", 512, "The maximum number of active connections")
var redisConn = flag.String("r", "localhost:6379", "Redis Address, defaults to localhost:6379")

type Data struct {
	Original  string
	Short     string
	FullShort string
	HitCount  int
}

//var collection []Data

func handler(w http.ResponseWriter, r *http.Request) {

	type NewURL struct {
		URL string
	}
	var url NewURL

	create, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if len(create) == 0 {
		var domain Data
		domain = getLongURL(r.URL.Path[1:])
		fmt.Println(domain)
		if len(domain.Original) > 0 {
			http.Redirect(w, r, domain.Original, http.StatusFound)
			return
		}
		http.Redirect(w, r, domain.Original, http.StatusNotFound)
		return
	}

	err = json.Unmarshal(create, &url)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	conn, err := redis.Dial("tcp", *redisConn)

	search := strings.Join([]string{"*||", url.URL},"")

	n, err := redis.Strings(conn.Do("KEYS", search));
	
	var newItem Data
	if len(n) < 1{
		newItem = createShortURL(url.URL)	
	}else{
		parts :=  strings.Split(n[0], "||")
		
		newItem.Short = parts[0]
		newItem.Original = parts[1]
		newItem.FullShort = strings.Join([]string{*base, parts[0]}, "")
		newCount, err := redis.Int(conn.Do("HGET", n[0], "count"))
		if err == nil {
			//TODO..
		}
		newItem.HitCount = newCount
	}

	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
	
	//collection = append(collection, newItem)

	output, err := json.Marshal(newItem)

	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "%s", output)
}

func createShortURL(url string) Data {

	conn, err := redis.Dial("tcp", *redisConn)
	var d Data
	count, err := redis.Int(conn.Do("INCR", "global:size"))
	if err != nil {
		log.Print(err)
		return d
	}
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

func getLongURL(short string) Data {
	var d Data
	
	conn, err := redis.Dial("tcp", *redisConn)

	search := strings.Join([]string{short, "||*"},"")
	fmt.Println(search)
	n, err := redis.Strings(conn.Do("KEYS", search));
	
	if err != nil {
		fmt.Println("Errors")
		log.Print(err)
		return d
	}

	if len(n) < 1{
		//Return an error
		
	}else{
		parts :=  strings.Split(n[0], "||")
		
		d.Short = parts[0]
		d.Original = parts[1]
		d.FullShort = strings.Join([]string{*base, parts[0]}, "")
		newCount, err := redis.Int(conn.Do("HINCRBY", n[0], "count",1))
		if err == nil {
			//TODO..
		}
		d.HitCount = newCount
		return d
	}

	return d
}

func main() {
	flag.Parse()
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(*host+":"+*port, nil)
	if err != nil{
		fmt.Println(err)
	}
}
