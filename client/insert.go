package main

import (
    "time"
    "fmt"
    "math/rand"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/sindbach/gomongo/models"
    //"crypto/tls"
    //"net"
)

func main() {

    /*tlsConfig := &tls.Config{}
    dialInfo := &mgo.DialInfo{
        Addrs: []string{"dagobah-shard-00-02-nesbp.mongodb.net:27017"},
        Database: "admin",
        Username: "skywalker",
        Password: "1amsecure",
        Timeout:  5 * time.Second, 
    }

    dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
        conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
        return conn, err
    }
    session, err := mgo.DialWithInfo(dialInfo)
    if err != nil {
        panic(err)
    }
    */

    session, err := mgo.Dial("localhost:27000,localhost:27001,localhost:27002")
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
    defer session.Close()

    // init pseudorandom generator
    rand.Seed(time.Now().Unix())

    // Define list of users
    userLists := []string {
        "Markus Thielsch",
        "Stephen Steneker",
        "Wan Bachtiar",
    }
    c := session.DB("gopher").C("users")
    
    i := bson.NewObjectId()
    err = c.Insert(&models.User{Id: i,Name: userLists[rand.Intn(len(userLists))], Assigned: rand.Int31()})
    if err != nil {
        panic(err)
    }
    fmt.Println("Inserted ", i)
}
