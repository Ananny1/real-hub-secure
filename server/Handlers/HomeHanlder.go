package Handlers

import (
	"fmt"
	"net/http"
)

// w to // talk back to the client
// for sending status code ,
//  or json , or html page
// to let the backend send a repsonse to the fornt

// r to // so the front send a request to the backend
// by checking the the method type ,
// form , button ,click
// any data send by the user
// headers , query parameters,etc

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home!")
}
