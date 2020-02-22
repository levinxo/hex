package main

import (
    "fmt"
    "log"
    "net/http"
)

var BlogMakerInstance = NewBlogMaker()

func main(){
    fmt.Println(BlogMakerInstance.Config.DevSiteUrl)
    dir := BlogMakerInstance.OutputDir
    log.Fatal(http.ListenAndServe(":8081", http.FileServer(http.Dir(dir))))
}
