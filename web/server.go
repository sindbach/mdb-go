package main

import (
    "path"
    "fmt"
	"net/http"
    "gopkg.in/mgo.v2"
    "html/template"
    "github.com/sindbach/gomongo/models"
)

func ReadCollection(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Handling request ...")

    session, err := mgo.Dial("localhost:27000,localhost:27001,localhost:27002")
    if err != nil {
        panic(err)
    }
    session.SetMode(mgo.Eventual, true)
    session.Copy()
    defer session.Close()

    c := session.DB("gopher").C("users")
    results := make([]models.User, 0, 10)
    err = c.Find(nil).Limit(10).Sort("-_id").All(&results)
    ip := path.Join("web/templates", "index.html")
    tmpl, err := template.ParseFiles(ip)
    if err != nil {
        panic(err)
    }

    status := models.ReplSetGetStatus{}
    if err := session.DB("admin").Run("replSetGetStatus", &status); err != nil {
        panic(err)
    } 

    webdata := models.WebData{results, status}

    if err:= tmpl.Execute(w, webdata); err != nil{
        panic(err)
    }
}

func handlerIcon(w http.ResponseWriter, r *http.Request) {} 
func main() {
    http.HandleFunc("/favicon.ico", handlerIcon)
	http.HandleFunc("/", ReadCollection)
    port := ":8000"
    fmt.Println(fmt.Sprintf("Waiting request http://localhost%s", port))
	http.ListenAndServe(port, nil)
}
