package main

import(
	"fmt"
    "log"
    "net/http"
)

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

func main(){
	http.HandleFunc("/", homePage)
	fmt.Println("Server Started at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}