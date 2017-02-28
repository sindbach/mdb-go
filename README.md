# mgo-intro
Introduction to mgo (MongoDB Go Driver) connecting to replica set


Starts replica set as: 

```
mongod --dbpath /data/gomongo/r1 --replSet gopher --wiredTigerCacheSizeGB 1 --port 30001
mongod --dbpath /data/gomongo/r2 --replSet gopher --wiredTigerCacheSizeGB 1 --port 30002
mongod --dbpath /data/gomongo/r3 --replSet gopher --wiredTigerCacheSizeGB 1 --port 30003
```

Initiate replica set
```
rs.initiate({_id:"gopher" members:[{_id:0, host:"localhost:30001"}]})
```

Run go server: 
```
go run server.go
```

Run go client insert:
```
while true; do echo 'inserting'; go run insert.go; sleep 2; done
```
