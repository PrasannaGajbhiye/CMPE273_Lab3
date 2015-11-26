// ConsistentHashing project main.go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type KeyValue struct {
	Key   int64  `json:"key"`
	Value string `json:"value"`
}

type KeysStructure struct {
	Keys []KeyValue
}

var AllKeys KeysStructure

func main() {
	mux := httprouter.New()
	mux.PUT("/keys/:key_id/:value", putKey)
	mux.GET("/keys/:key_id", getKey)
	mux.GET("/keys", getAllKeys)
	server1 := http.Server{
		Addr:    "0.0.0.0:3000",
		Handler: mux,
	}
	server2 := http.Server{
		Addr:    "0.0.0.0:3001",
		Handler: mux,
	}
	server3 := http.Server{
		Addr:    "0.0.0.0:3002",
		Handler: mux,
	}
	server1.ListenAndServe()
	server2.ListenAndServe()
	server3.ListenAndServe()
}

// Put a key value pair - PUT Request
func putKey(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	fmt.Println("Putting object on host:" + req.Host)
	key_idStr := p.ByName("key_id")
	key_id, _ := strconv.ParseInt(key_idStr, 10, 64)
	value := p.ByName("value")

	var keyValuePair KeyValue
	keyValuePair.Key = key_id
	keyValuePair.Value = value

	AllKeys.Keys = append(AllKeys.Keys, keyValuePair)

	jsonOutput, _ := json.Marshal(&keyValuePair)
	fmt.Fprintf(rw, string(jsonOutput))

	jsonOutput, _ = json.Marshal(&AllKeys.Keys)
	fmt.Println(string(jsonOutput))
	fmt.Println()

}

// Get a key value pair - GET Request
func getKey(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	fmt.Println("Getting object on host:" + req.Host)
	key_idStr := p.ByName("key_id")
	key_id, _ := strconv.ParseInt(key_idStr, 10, 64)

	var searchKeyValue KeyValue
	found := false
	for i := 0; i < len(AllKeys.Keys); i++ {
		if AllKeys.Keys[i].Key == key_id {
			searchKeyValue.Key = AllKeys.Keys[i].Key
			searchKeyValue.Value = AllKeys.Keys[i].Value
			found = true
			break
		}
	}
	if found == true {
		jsonOutput, _ := json.Marshal(&searchKeyValue)
		fmt.Fprintf(rw, string(jsonOutput))
	} else {
		fmt.Fprintf(rw, "Key not found.")
	}
}

// Get all key value pairs - GET Request
func getAllKeys(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	fmt.Println()
	fmt.Println("Getting keys for host:" + req.Host)

	if len(AllKeys.Keys) > 0 {
		jsonOutput, _ := json.Marshal(&AllKeys.Keys)
		fmt.Fprintf(rw, "\n"+string(jsonOutput)+"\n")
	}
}

/*
You will be implementing a simple RESTful key-value data store with the following features:
PUT http://localhost:3000/keys/{key_id}/{value}
E.g. http://localhost:3000/keys/1/foobar
Response: 200
GET http://localhost:3000/keys/{key_id}
E.g. http://localhost:3000/keys/1
Response: {
                     “key” : 1,
                     “value” : “foobar”
                   }

GET http://localhost:3000/keys
E.g. http://localhost:3000/keys
Response: [
          {
                     “key” : 1,
                     “value” : “foobar”
           },
                       {
                                 “key” : 2,
                                “value” : “b”
                       }
            ]

*/
