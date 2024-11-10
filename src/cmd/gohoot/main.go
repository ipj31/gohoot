package main

import (
	"fmt"
	"net/http"
)

// https://gohoot.auth.us-east-2.amazoncognito.com/oauth2/authorize?response_type=code&client_id=2ndpsrnnp86dfm4s06kidlq5&redirect_uri=http://localhost/user

func main() {

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(*r)
	})

	http.ListenAndServe("", nil)
}
