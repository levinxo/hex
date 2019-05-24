package main

import (
    "fmt"
    "log"
    "net/http"
)

var BlogMakerInstance = NewBlogMaker()

func main(){
    fmt.Println("http://localhost:8081")
    log.Fatal(http.ListenAndServe(":8081", http.FileServer(http.Dir(BlogMakerInstance.OutputDir))))
}
