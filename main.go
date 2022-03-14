package main

import "net/http"

func main() {
	http.HandleFunc("/v1/weather", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	http.ListenAndServe(":8080", nil)
}
