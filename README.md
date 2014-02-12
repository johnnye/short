Shorter - A Go Link Shortner
===========================
This is a link shortener written in go. 

The shortened hash is calculated by base62 encoding the number of links currently stored.

###Use
Set your application domain at the top of the short.go file. It's a variable called domain

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

- show stats for a URL, if the hitcounter stays 
- correct response codes, 201, 404, 500
- better logging of errors
- tests 
- channels 
- command line options
- make the project structure go-like 
