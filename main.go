package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    b, err := os.ReadFile("./public/index.html")
    if err != nil {
        return
    }
    fmt.Fprintf(w, string(b))
}

func handler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/" {
        homeHandler(w, r)
        return
    }
    fmt.Fprintf(w, "I love %v", r.URL.Path)
}

func main() {
    fmt.Println("Hello, World!")
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":80", nil))
}
