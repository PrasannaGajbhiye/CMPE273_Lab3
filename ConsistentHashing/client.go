package main

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type KeyValue struct {
	Key   int64  `json:"key"`
	Value string `json:"value"`
}

type KeysStructure struct {
	Keys []KeyValue
}

var crc32q *crc32.Table
var circle map[uint32]string
var servers []string

func main() {

	servers = []string{
		"http://localhost:3000/",
		"http://localhost:3001/",
		"http://localhost:3002/"}
	crc32q = crc32.MakeTable(0xD5828281)

	circle = make(map[uint32]string)

	for i := 0; i < len(servers); i++ {
		val := crc32.Checksum([]byte(servers[i]), crc32q)
		circle[val] = servers[i]
	}

	cmdArgs := os.Args[1:]
	//	Start Client

	StartClient(cmdArgs)

}

func get(key int64, circle map[uint32]string) string {
	keyStr := strconv.Itoa(int(key))
	if len(circle) == 0 {
		return ""
	} else {
		hash := crc32.Checksum([]byte(keyStr), crc32q)
		var keys []int
		for k := range circle {
			keys = append(keys, int(k))
		}
		sort.Ints(keys)
		var serverHash uint32
		foundHash := false
		for i := 0; i < len(keys); i++ {
			if keys[i] >= int(hash) {
				serverHash = uint32(keys[i])
				foundHash = true
				break
			}
		}
		if foundHash == true {
			return circle[serverHash]
		} else {
			return circle[uint32(keys[0])]
		}

	}
}

func StartClient(indCmdArgs []string) {
	if len(indCmdArgs) == 3 {
		httpMethod := indCmdArgs[0]
		key, err := strconv.ParseInt(indCmdArgs[1], 10, 64)
		if err != nil {
			fmt.Println("Invalid key.")
			return
		}

		serverHostName := get(key, circle)
		if serverHostName != "" {
			keyStr := strconv.Itoa(int(key))
			value := indCmdArgs[2]

			if httpMethod == "PUT" {

				client := &http.Client{}
				request, _ := http.NewRequest("PUT", serverHostName+"keys/"+keyStr+"/"+value, nil)
				resp, err := client.Do(request)
				if err != nil {
					fmt.Println("Error in requesting method.")
				}
				defer resp.Body.Close()

				fmt.Println(resp.StatusCode)
			} else {
				fmt.Println("Invalid HTTP Method - Expected Method: PUT")
			}

		} else {
			fmt.Println("No node to add this object.")
		}

	} else if len(indCmdArgs) == 2 {
		httpMethod := indCmdArgs[0]
		key, _ := strconv.ParseInt(indCmdArgs[1], 10, 64)
		keyStr := strconv.Itoa(int(key))
		serverHostName := get(key, circle)
		if serverHostName != "" {
			if httpMethod == "GET" {
				client := &http.Client{}
				request, _ := http.NewRequest("GET", serverHostName+"keys/"+keyStr, nil)
				resp, err := client.Do(request)
				if err != nil {
					fmt.Println("Error")
				}
				defer resp.Body.Close()

				body, _ := ioutil.ReadAll(resp.Body)
				var msgRes interface{}
				_ = json.Unmarshal(body, &msgRes)

				var keyValuePair KeyValue
				k := msgRes.(map[string]interface{})["key"].(float64)
				v := msgRes.(map[string]interface{})["value"].(string)

				keyValuePair.Key = int64(k)
				keyValuePair.Value = v
				jsonOutput, _ := json.Marshal(keyValuePair)

				fmt.Println("\n" + string(jsonOutput) + "\n")
			} else {
				fmt.Println("Invalid HTTP Method. Expected Method: GET")
			}
		} else {
			fmt.Println("No node to add this object.")
		}

	} else if len(indCmdArgs) == 0 {
		client := &http.Client{}
		var AllKey KeysStructure
		for i := 0; i < len(servers); i++ {
			serverHostName := servers[i]
			request, _ := http.NewRequest("GET", serverHostName+"keys", nil)
			resp, err := client.Do(request)
			if err != nil {
				fmt.Println("Error")
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			var msgRes interface{}
			_ = json.Unmarshal(body, &msgRes)

			if msgRes != nil {
				keyValArr := msgRes.([]interface{})

				for i := 0; i < len(keyValArr); i++ {
					var keyValuePair KeyValue
					k := keyValArr[i].(map[string]interface{})["key"].(float64)
					v := keyValArr[i].(map[string]interface{})["value"].(string)

					keyValuePair.Key = int64(k)
					keyValuePair.Value = v
					AllKey.Keys = append(AllKey.Keys, keyValuePair)
				}
			}

		}

		jsonOutput, _ := json.Marshal(AllKey.Keys)
		fmt.Println("\n" + string(jsonOutput) + "\n")

	}

}
