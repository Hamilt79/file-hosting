package main

import (
    "bytes"
    "io"
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/google/uuid"
)

func exists(path string) bool {
    _, err := os.Stat(path)
    if err == nil { return true }
    if os.IsNotExist(err) { return false }
    return false
}

func downloadFile(w http.ResponseWriter, r *http.Request) {

    fmt.Println("Here")
    dirPath := "./files/" + r.URL.Path[10:]
    fmt.Println(r.URL.Path[10:])
    
    dirExists := exists(dirPath)
    if (!dirExists) {
        return;
    }
    entry, err := os.ReadDir(dirPath)
    if err != nil {
        fmt.Println(err)
        return
    }
    if len(entry) != 1 {
        fmt.Println("There is more than one file in the directory")
        return
    }
    // Set the content type
    w.Header().Set("Content-Type", "application/octet-stream")
    
    // Set the content disposition to force download
    w.Header().Set("Content-Disposition", "attachment; filename=" + entry[0].Name())

    file, err := os.OpenFile(dirPath + "/" + entry[0].Name(), os.O_CREATE, 0777)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()

    _, err = io.Copy(w, file)
    if err != nil {
        fmt.Println(err)
        return
    }
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.Header)
    fmt.Println(r.Body)
    uid := uuid.New().String()
    fmt.Println(uid)
    os.MkdirAll("./files/" + uid, 0777)

    r.ParseForm()

    file, header, err := r.FormFile("file")
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(header.Filename)
    fileBuffer := bytes.NewBuffer(nil)
    _, err = io.Copy(fileBuffer, file)
    if err != nil {
        fmt.Println(err)
        return 
    }
    os.WriteFile("./files/" + uid + "/" + header.Filename, fileBuffer.Bytes(), 0777)

    hostName := r.Host

    fmt.Fprintf(w, "<p>%v</p>", hostName + "/download/" + uid)

}

func main() {
    fmt.Println("Hello, World!")
    fs := http.FileServer(http.Dir("./public"))
    http.Handle("/", fs)
    http.HandleFunc("/post-file", fileHandler)
    http.HandleFunc("/download/*", downloadFile)
    log.Fatal(http.ListenAndServe(":80", nil))
}
