package main

import (
	"net/http"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "html/template"
    "path"
    "fmt"
)

type User struct {
    Id bson.ObjectId `bson:"_id"`
    Name string 
    Assigned int32
}
type ReplSetGetStatus struct {
    Members []struct {
        Name     string `bson:"name"`
        StateStr string `bson:"stateStr"`
    } `bson:"members"`
}
type WebData struct {
    Users []User
    Status ReplSetGetStatus
}


func ReadCollection(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Handling request ...")
    session, err := mgo.Dial("localhost:30001,localhost:30002,localhost:30003")
    session.SetMode(mgo.Secondary, true)
    if err != nil {
        panic(err)
    }
    session.Copy()
    defer session.Close()

    c := session.DB("gopher").C("users")
    results := make([]User, 0, 10)
    err = c.Find(nil).Limit(10).Sort("-_id").All(&results)
    ip := path.Join("templates", "index.html")
    tmpl, err := template.ParseFiles(ip)
    if err != nil {
        panic(err)
    }

    status := ReplSetGetStatus{}
    if err := session.DB("admin").Run("replSetGetStatus", &status); err != nil {
        panic(err)
    } 

    webdata := WebData{results, status}

    if err:= tmpl.Execute(w, webdata); err != nil{
        panic(err)
    }
}

func handlerIcon(w http.ResponseWriter, r *http.Request) {} 
func main() {
    http.HandleFunc("/favicon.ico", handlerIcon)
	http.HandleFunc("/", ReadCollection)
	http.ListenAndServe(":8000", nil)
}
