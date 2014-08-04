# Faye + Go


Experimental

## Usage

```go

package main

import (
	"github.com/roncohen/faye"
	"github.com/roncohen/faye/adapters"
	"log"
	"net/http"
)

func OurLoggingHandler(pattern string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v: %+v", pattern, *r.URL)
		h.ServeHTTP(w, r)
	})
}

func main() {
	fayeServer := faye.NewServer(faye.NewEngine())
	http.Handle("/faye", adapters.FayeHandler(fayeServer))

	// Also serve up some static files and show off 
	// the wonderful go http handler chain
	http.Handle("/", OurLoggingHandler("/", http.FileServer(http.Dir("src/github.com/roncohen/faye-test/static"))))

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
```