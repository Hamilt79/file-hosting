package main

import (
    "sync"
    "strings"
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
    if strings.Contains(r.URL.Path[10:], "/") {
        fmt.Println("There is a / in the path")
        return
    }
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
    // Generating a new unique id for the file dir
    uid := uuid.New().String()
    // Printing that id to console
    fmt.Println(uid)
    // Creating the directory for the file
    os.MkdirAll("./files/" + uid, 0777)
    // Parsing the form
    r.ParseForm()

    // Getting the file the user sent from the form
    file, header, err := r.FormFile("file")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()
    
    // Printing the file name to the console
    fmt.Println(header.Filename)
    // Making a buffer to store the file
    fileBuffer := bytes.NewBuffer(nil)
    // Writing the file to the buffer
    _, err = io.Copy(fileBuffer, file)
    if err != nil {
        fmt.Println(err)
        return 
    }

    // Making sure the filename doesn't contain a / to prevent directory traversal
    // There are still prolly other ways to do directory traversal but this is a start
    if strings.Contains(header.Filename, "/") {
        fmt.Println("There is a / in the path")
        return
    }
    // Writing the file to the directory
    os.WriteFile("./files/" + uid + "/" + header.Filename, fileBuffer.Bytes(), 0777)

    hostName := r.Host

    // Writing the link the user can use to download the file to the response page
    fmt.Fprintf(w, "<p>%v</p>", hostName + "/download/" + uid)
}

func main() {
    fmt.Println("Server starting up.")
    fs := http.FileServer(http.Dir("./public"))
    http.Handle("/", fs)
    http.HandleFunc("/post-file", fileHandler)
    http.HandleFunc("/download/*", downloadFile)

    var waitG sync.WaitGroup

    waitG.Add(1)
    go func() {
        log.Fatal(http.ListenAndServe(":8081", nil))
        waitG.Done()
    }()
    waitG.Add(1)
    go func() {
        log.Fatal(http.ListenAndServeTLS(":8080", "../certs/certificate.crt", "../certs/private.key" , nil))
        waitG.Done()
    }()
    fmt.Println("Server started.")
    waitG.Wait()
}
