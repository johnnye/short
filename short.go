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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"github.com/garyburd/redigo/redis"
)

const domain = "http://localhost:8080/"

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
	newItem := createShortURL(url.URL)
	//collection = append(collection, newItem)

	output, err := json.Marshal(newItem)

	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
	//TODO: correct headers 201 
	fmt.Fprintf(w, "%s", output)
}

func createShortURL(url string) Data {
	conn, err := redis.Dial("tcp", ":6379")
	//TOOD: Here check the redis instance. 
	
	n, err := conn.Do("HGETALL", url)
	if err == nil {
		log.Print(n)

		other, err := conn.Do("HINCRBY",url, "HitCount", "1")
		log.Print(other)
		if err != nil{
			log.Print("try something here.")
			log.Print(err)
		}
		var x Data
		return x
	}	
		
		
	noKeys, err := conn.Do("DBSIZE")

	var d Data
	encodedVar := base62.EncodeInt(int64(noKeys.(int)))
	d.FullShort = strings.Join([]string{domain, encodedVar}, "")
	d.Short = encodedVar
	d.Original = url
	d.HitCount = 0
	
	newurlindb, err := conn.Do("HMSET", d.Short, d)
	log.Print(err)	
	log.Print(d)
	log.Print(newurlindb)
	return d
}

func getLongURL(short string) Data {
	conn, err := redis.Dial("tcp", ":6379")
	//TOOD: Here check the redis instance. 
	
	n, err := conn.Do("HGETALL", short)
	if err == nil {
		log.Print(n)
		other, err := conn.Do("HINCRBY", short,"HitCount", "1")
		log.Print(other)
		if err !=nil {
			log.Print(err)
		}
		var x Data
		return x
	}

	var d Data
	return d
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
