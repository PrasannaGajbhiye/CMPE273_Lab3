# CMPE273_Lab3
# Consistent Hashing


## Usage

### Install

```
go get github.com/PrasannaGajbhiye/CMPE273_Lab3/
```


### Start the  server:
Start three instances of the servers by executing the below command in three terminal windows
```
go run server.go
```

### Start the client 


For putting a key
```
go run client.go "PUT" 1 "a"
```


For getting a key
```
go run client.go "GET" 1
```

For getting all keys
```
go run client.go
```
