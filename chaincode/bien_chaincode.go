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
const (
	millisPerSecond     = int64(time.Second / time.Millisecond)
	nanosPerMillisecond = int64(time.Millisecond / time.Nanosecond)
)
func generateCUSIPSuffix(issueDate string, days int) (string, error) {

	t, err := msToTime(issueDate)
	if err != nil {
		return "", err
	}

	maturityDate := t.AddDate(0, 0, days)
	month := int(maturityDate.Month())
	day := maturityDate.Day()

	suffix := seventhDigit[month] + eigthDigit[day]
	return suffix, nil

}

func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(msInt/millisPerSecond,
		(msInt%millisPerSecond)*nanosPerMillisecond), nil
}
var orderIndexStr ="_orderindex"
//var openTradesStr = "_opentrades"				//name for the key/value that will store all open trades
var goodsPrefix = "goods:"
//var accountPrefix = "acct:"
type Owner struct {
	Company string    `json:"company"`
}

type Goods struct{
		GDSID string `json:"goodsId"`
		name string `json:"name"`	
		price float64 `json:"price"`
		postage float64 `json:"postage"`
		Owners    []Owner `json:"owner"`
	    Issuer    string  `json:"issuer"`
	    state string `json:"state"`
}

type Transaction struct {
	GDSID       string   `json:"gdsid"`
	FromCompany string   `json:"fromCompany"`
	ToCompany   string   `json:"toCompany"`
	postage    float64  `json:"discount"`
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
	} else if function == "add_goods" {
		//return t.add_goods(stub, args)
		return t.issueCommercialGoods(stub, args)
	}
	//else if function == "set_owner" {
	//	return t.set_owner(stub, args)
	//} else if function == "change_state" {
	//	return t.change_state(stub, args)
	//} 
	
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

func (t *BienChaincode) issueCommercialGoods(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	/*		0
	
	    GDSID int64 `json:"goodsId"`
		name string `json:"name"`	
		price int `json:"price"`
		postage int `json:"postage"`
		Owners    []Owner `json:"owner"`
	    Issuer    string  `json:"issuer"`
	    state string `json:"state"`
		
		json
	  	{
			"name":  "string",
			"price": 0.00,
			"postage": 7.5,
			"owners": [ // This one is not required
				{
					"company": "company1",
					"quantity": 5
				},
				{
					"company": "company3",
					"quantity": 3
				},
				{
					"company": "company4",
					"quantity": 2
				}
			],				
			"issuer":"company2",
			"state":"new"  

		}
	*/
	//need one arg
	if len(args) != 1 {
		fmt.Println("error invalid arguments")
		return nil, errors.New("Incorrect number of arguments. Expecting commercial paper record")
	}

	var goods Goods
	var err error
	//var account Account
    timestamp := time.Now().Unix()
	fmt.Println("Unmarshalling goods")
	err = json.Unmarshal([]byte(args[0]), &goods)
	if err != nil {
		fmt.Println("error invalid goods issue")
		return nil, errors.New("Invalid commercial goods issue")
	}

	
	// Set the issuer to be the owner of all quantity
	var owner Owner
	owner.Company = goods.Issuer
	
	goods.Owners = append(goods.Owners, owner)

	suffix, err := generateCUSIPSuffix(strconv.FormatInt(timestamp, 10), 15)
	if err != nil {
		fmt.Println("Error generating gdsid")
		return nil, errors.New("Error generating GDSID")
	}

	fmt.Println("Marshalling goods bytes")
	goods.GDSID = goods.Issuer + suffix
	
	fmt.Println("Getting State on goods " + goods.GDSID)
	cpRxBytes, err := stub.GetState(goodsPrefix+goods.GDSID)
	if cpRxBytes == nil {
		fmt.Println("GDSID does not exist, creating it")
		goodsBytes, err := json.Marshal(&goods)
		if err != nil {
			fmt.Println("Error marshalling cp")
			return nil, errors.New("Error issuing commercial goods")
		}
		err = stub.PutState(goodsPrefix+goods.GDSID, goodsBytes)
		if err != nil {
			fmt.Println("Error issuing paper")
			return nil, errors.New("Error issuing commercial paper")
		}

	fmt.Println("Getting goods Keys")
	GoodsAsBytes, err := stub.GetState(orderIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Goods index")
	}
	var orderIndex []string
	json.Unmarshal(GoodsAsBytes, &orderIndex)							//un stringify it aka JSON.parse()
	fmt.Println("get order(Goods) index: ", orderIndex)
	//append
	fmt.Println("Appending the new goods GDSID to order Keys")
							//add Goods id to index list
		
		foundKey := false
		for _, index := range orderIndex {
			if index == goodsPrefix+goods.GDSID {
				foundKey = true
			}
		}
		if foundKey == false {
			orderIndex = append(orderIndex,goods.GDSID)		
			keysBytesToWrite, err := json.Marshal(&orderIndex)
			if err != nil {
				fmt.Println("Error marshalling orderIndex")
				return nil, errors.New("Error marshalling the orderIndex")
			}
			fmt.Println("Put state on orderIndex")
			err = stub.PutState(orderIndexStr, keysBytesToWrite)
			if err != nil {
				fmt.Println("Error writting orderIndexStr back")
				return nil, errors.New("Error writing the orderIndexStr back")
			}
		}
		
		fmt.Println("Issue commercial paper %+v\n", goods)
		return nil, nil
	}else{
	fmt.Println("GDSID exists")
	}
	return nil, nil
}


// read - query function to read key/value pair
/*func (t *BienChaincode) set_owner(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error
	
	if len(args)<2 {
	 return nil,errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	fmt.Println("- start set owner-")
	fmt.Println(args[0] + " - " + args[1])
	GoodsAsBytes, err := stub.GetState(args[0])
	if err != nil {
			return nil, errors.New("Failed to get item")
		}
		res := Goods{}
		json.Unmarshal(GoodsAsBytes, &res)										//un stringify it aka JSON.parse()
		res.owner = args[1]
		
		jsonAsBytes, _ := json.Marshal(res)
		err = stub.PutState(args[0], jsonAsBytes)								//rewrite the marble with id as key
		if err != nil {
			return nil, err
		}
		
		fmt.Println("- end set owner-")
		
		return nil, nil
}*/

// read - query function to read key/value pair, then change the data structure's state field
/*func (t *BienChaincode) change_state(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
//   0       1       2          3       4     5
	//id  "name", "owner", "state", "price"  "postage"
	var err error
	
	if len(args)<2 {
	 return nil,errors.New("Incorrect number of arguments. Expecting 2")
	}

	GoodsAsBytes, err := stub.GetState(args[0])
	logger.Infof("change_state getState: logger GoodsAsBytes=%v", GoodsAsBytes)
	if err != nil {
			return nil, errors.New("Failed to get thing")
		}
	
    var res Goods

		json.Unmarshal(GoodsAsBytes, &res)	
		
		res.state = args[1]

		logger.Infof("change_state res: logger res=%v", res)
		fmt.Println(res.id, ":",res.name, ":", res.owner, ":", res.state, ":", res.price, ":", res.postage)
		jsonAsBytes, _ := json.Marshal(res)
		err = stub.PutState(args[0], jsonAsBytes)								//rewrite the goods with name as key

		if err != nil {
			return nil, err
		}
		
		fmt.Println("- end change state-")

	return nil, nil
		
}*

/*func (t *BienChaincode) add_goods(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
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
	GoodsAsBytes, err := stub.GetState(orderIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Goods index")
	}
	var orderIndex []string
	json.Unmarshal(GoodsAsBytes, &orderIndex)							//un stringify it aka JSON.parse()
	fmt.Println("get order(Goods) index: ", orderIndex)
	//append
	orderIndex = append(orderIndex,strconv.FormatInt(timestamp , 10))								//add Goods id to index list
	fmt.Println("append:! order(Goods) index: ", orderIndex)
	jsonAsBytes, _ := json.Marshal(orderIndex)
	err = stub.PutState(orderIndexStr, jsonAsBytes)						//store id of Goods

	fmt.Println("- end add goods")
	return nil, nil
}*/
var seventhDigit = map[int]string{
	1:  "A",
	2:  "B",
	3:  "C",
	4:  "D",
	5:  "E",
	6:  "F",
	7:  "G",
	8:  "H",
	9:  "J",
	10: "K",
	11: "L",
	12: "M",
	13: "N",
	14: "P",
	15: "Q",
	16: "R",
	17: "S",
	18: "T",
	19: "U",
	20: "V",
	21: "W",
	22: "X",
	23: "Y",
	24: "Z",
}

var eigthDigit = map[int]string{
	1:  "1",
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "A",
	11: "B",
	12: "C",
	13: "D",
	14: "E",
	15: "F",
	16: "G",
	17: "H",
	18: "J",
	19: "K",
	20: "L",
	21: "M",
	22: "N",
	23: "P",
	24: "Q",
	25: "R",
	26: "S",
	27: "T",
	28: "U",
	29: "V",
	30: "W",
	31: "X",
}
