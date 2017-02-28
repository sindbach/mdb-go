package main

import (
    "time"
    "fmt"
    "math/rand"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)


type User struct {
    Id bson.ObjectId `bson:"_id"`
    Name string 
    Assigned int32
}

func main() {
    session, err := mgo.Dial("localhost:30001,localhost:30002,localhost:30003")
    if err != nil {
        panic(err)
    }
    defer session.Close()

    // init pseudorandom generator
    rand.Seed(time.Now().Unix())

    // Define list of users
    userLists := []string {
        "Markus Thielsch",
        "Rick Sanchez",
        "Stephen Steneker",
        "Kevin Adistambha",
        "Wan Bachtiar",
    }
    c := session.DB("gopher").C("users")
    
    i := bson.NewObjectId()
    err = c.Insert(&User{Id: i,Name: userLists[rand.Intn(len(userLists))], Assigned: rand.Int31()})
    if err != nil {
        panic(err)
    }
    fmt.Println("Inserted ", i)
}
