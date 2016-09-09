/*
Copyright IBM Corp 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		 http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// BienChaincode is  a Chaincode for bien application implementation
type BienChaincode struct {
}
var orderIndexStr ="_orderindex"

type Bien struct{
		id int64 `json:"orderId"`
		name string `json:"name"`
		state string `json:"state"`
		price int `json:"price"`
		postage int `json:"postage"`
		owner string `json:"owner"`
}
var logger = shim.NewLogger("SimpleChaincode")
func main() {
    logger.SetLevel(shim.LogInfo) 
	err := shim.Start(new(BienChaincode))
	if err != nil {
		fmt.Printf("Error starting BienChaincode chaincode: %s", err)
	}
}

// Init resets all the things
func (t *BienChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Printf("hello init chaincode, it is for testing")
	var Aval int
	var err error
    logger.Warning("init logger should be 1 string") 
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval)))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	logger.Infof("init logger arg0=%v", args[0])
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(orderIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *BienChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "set_owner" {
		return t.set_owner(stub, args)
	} else if function == "change_state" {
		return t.change_state(stub, args)
	} else if function == "add_goods" {
		return t.add_goods(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *BienChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query")
}

// write - invoke function to write key/value pair
func (t *BienChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] 
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *BienChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	
	valAsbytes, err := stub.GetState(key)
	logger.Infof("query.read logger valAsbytes=%v", valAsbytes)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// read - query function to read key/value pair
func (t *BienChaincode) set_owner(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error
	
	if len(args)<2 {
	 return nil,errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	fmt.Println("- start set owner-")
	fmt.Println(args[0] + " - " + args[1])
	bienAsBytes, err := stub.GetState(args[0])
	if err != nil {
			return nil, errors.New("Failed to get item")
		}
		res := Bien{}
		json.Unmarshal(bienAsBytes, &res)										//un stringify it aka JSON.parse()
		res.owner = args[1]
		
		jsonAsBytes, _ := json.Marshal(res)
		err = stub.PutState(args[0], jsonAsBytes)								//rewrite the marble with id as key
		if err != nil {
			return nil, err
		}
		
		fmt.Println("- end set owner-")
		return t.read(stub,args)
		//return bienBytes, nil
}

// read - query function to read key/value pair, then change the data structure's state field
func (t *BienChaincode) change_state(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error
	
	if len(args)<2 {
	 return nil,errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	fmt.Println("- start change state -")
	fmt.Println("id:" + args[0] + " - state" + args[1])
	bienAsBytes, err := stub.GetState(args[0])
	logger.Infof("change_state getState: logger bienAsBytes=%v", bienAsBytes)
	if err != nil {
			return nil, errors.New("Failed to get thing")
		}
		res := Bien{}
		json.Unmarshal(bienAsBytes, &res)	
		logger.Infof("change_state before set res: logger res=%v", res)									//un stringify it aka JSON.parse()
		res.state = args[1]
		res.name = "hi"
		logger.Infof("change_state res: logger res=%v", res)
		jsonAsBytes, _ := json.Marshal(res)
		err = stub.PutState(args[0], jsonAsBytes)								//rewrite the goods with name as key
		logger.Infof("change_state putstate: logger jsonAsBytes=%v", jsonAsBytes)
		if err != nil {
			return nil, err
		}
		bien2AsBytes, err := stub.GetState(args[0])
		logger.Infof("query.read after bien2AsBytes=%v", bien2AsBytes)
		fmt.Println("- end change state-")

		return nil, nil
}

func (t *BienChaincode) add_goods(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
var err error
fmt.Println("hello add goods")
	//   0       1       2          3       4
	// "name", "owner", "state", "price"  "postage"
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	fmt.Println("- start add goods")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return nil, errors.New("5th argument must be a non-empty string")
	}
	
	timestamp := time.Now().Unix()
	str := `{"id":"`+strconv.FormatInt(timestamp , 10)+`","name": "` + args[0] + `", "owner": "` + args[1] + `", "state": "` + args[2]+ `", "price": ` + args[3] + `, "postage": ` + args[4] +`}`
	
	err = stub.PutState(strconv.FormatInt(timestamp , 10), []byte(str))								//store marble with id as key
	if err != nil {
		return nil, err
	}
	
	//get the  index
	bienAsBytes, err := stub.GetState(orderIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get bien index")
	}
	var orderIndex []string
	json.Unmarshal(bienAsBytes, &orderIndex)							//un stringify it aka JSON.parse()
	
	//append
	orderIndex = append(orderIndex,strconv.FormatInt(timestamp , 10))								//add bien id to index list
	fmt.Println("! order(bien) index: ", orderIndex)
	jsonAsBytes, _ := json.Marshal(orderIndex)
	err = stub.PutState(orderIndexStr, jsonAsBytes)						//store id of bien

	fmt.Println("- end add goods")
	return nil, nil
}