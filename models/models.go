package models 

import (
    "gopkg.in/mgo.v2/bson"
    "time"
)

type User struct {
    Id bson.ObjectId `bson:"_id"`
    Name string 
    Created time.Time
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

type StatCache struct {
    PreviousStat             *ServerStatus
    OpCommands                []float64
    OpInserts                []float64
    First                     bool
}

type ServerStatus struct {
    Host               string                 `bson:"host"`
    Opcounters         *OpcountStats          `bson:"opcounters"`
    OpcountersRepl     *OpcountStats          `bson:"opcountersRepl"`
    Repl               *ReplStatus            `bson:"repl"`
}

// ReplStatus stores data related to replica sets.
type ReplStatus struct {
    SetName      string      `bson:"setName"`
    IsMaster     interface{} `bson:"ismaster"`
    Secondary    interface{} `bson:"secondary"`
    IsReplicaSet interface{} `bson:"isreplicaset"`
    ArbiterOnly  interface{} `bson:"arbiterOnly"`
    Hosts        []string    `bson:"hosts"`
    Primary      string      `bson:"primary"`
    Me           string      `bson:"me"`
}

type OpcountStats struct {
    Insert  int64 `bson:"insert"`
    Query   int64 `bson:"query"`
    Update  int64 `bson:"update"`
    Delete  int64 `bson:"delete"`
    GetMore int64 `bson:"getmore"`
    Command int64 `bson:"command"`
}
