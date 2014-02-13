Shorter - A Go Link Shortner
===========================
This is a link shortener written in go. 

The shortened hash is calculated by base62 encoding the number of links currently stored.

###Use
There are 5 commandline options for short

- **-h**: Host address to listen on, defaults to localhost 
- **-p**: Port to listen on, defaults to 8080
- **-b**: Base Address: this is base address of the shortener, defaults to http://localhost:8080/
- **-c**: Max Connections: defaults to 512. Not Currently Used
- **-r**: Redis Address: defaults to localhost:6379 

Example Serving directly
````bash
short -h shortdoma.in -p 80 -b http://shortdoma.in/
````

Behind a proxy with a remote Redis box:
````bash 
short -h localhost -p 12345 -b http://shrtdm.in -c 5 -r http://redis.myserver:999393
````


####Create
POST request to the root URL with the following JSON <code>{"url":"http://example.com"}</code> the folowing response is sent to you

```JSON
{
	"Original":"http://example.com",
	"Short":"A",
	"FullShort":"http://localhost:8080/A",
	"HitCount":0
}
```

###TODO

- show stats for a URL
- correct response codes, 201, 404, 500
- tests 
- channels 
