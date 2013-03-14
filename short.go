package main
/*
* My in memory link shortner written in Go 
* This is released under a "you'd be mad to use it" license 
* 
* John Nye
*
*/
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"./utils"
)
const domain = "http://localhost:8080/"
type Data struct {
	Original  string
	Short     string
	FullShort string
	HitCount int
}
var collection []Data
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
		long := getLongURL(r.URL.Path[1:])
		if(len(long)>0){
			http.Redirect(w, r, long, http.StatusFound)
			return
		}
		http.Redirect(w,r,domain,http.StatusNotFound)
		return
	}
	
	err = json.Unmarshal(create, &url)

	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
	newItem := createShortURL(url.URL)
	collection = append(collection, newItem )

	output, err := json.Marshal(newItem)

	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
	//TODO: correct headers 201 
	fmt.Fprintf(w, "%s", output)
}

func createShortURL(url string) Data{
	for _,element := range collection {
		if(element.Original == url){
			return element
		}
	}
	var d Data
	encodedVar := base62.StdEncoding.EncodeToString([]byte(string(len(collection))))
	encodedVar = strings.Trim(encodedVar, "=")
	d.FullShort = strings.Join([]string{domain,encodedVar}, "")
	d.Short = encodedVar
	d.Original = url
	d.HitCount = 0

	log.Print(encodedVar)
	return d
}

func getLongURL(short string)string{
	for _, element := range collection {
		if (element.Short == short){
			return element.Original
		}
	}
	return ""
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
