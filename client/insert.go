package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sindbach/mdb-go/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {

	session, err := mgo.Dial("localhost:27000,localhost:27001,localhost:27002")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer session.Close()

	// init pseudorandom generator
	rand.Seed(time.Now().Unix())

	// Define list of users
	userLists := []string{
		"John Smith",
		"Mr. Gopher",
		"Stephen Steneker",
		"Wan Bachtiar",
	}
	c := session.DB("gopher").C("users")

	i := bson.NewObjectId()
	err = c.Insert(&models.User{Id: i, Name: userLists[rand.Intn(len(userLists))], Created: time.Now()})
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted ObjectId:", i)
}
