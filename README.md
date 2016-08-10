# Bien-Chaincode
This project is a demo for using IBM blockchain.
It is a chaincode that could be deployed into a network of Hyperledger fabric peer nodes that enables interaction with that network's shared ledger.
#Dependencies

The import statement lists a few dependencies that you will need for your chaincode to build successfully.

fmt - contains Println for debugging/logging.
errors - standard go error format.
github.com/hyperledger/fabric/core/chaincode/shim - the code that interfaces your golang code with a peer.
#Init()

Init is called when you first deploy your chaincode. As the name implies, this function should be used to do any initialization your chaincode needs. In our example, we use Init to configure the initial state of one variable on the ledger.

In your chaincode_start.go file, change the Init function so that it stores the first element in the args argument to the key "hello_world".

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }

    err := stub.PutState("hello_world", []byte(args[0]))
    if err != nil {
        return nil, err
    }

    return nil, nil
}
This is done by using the shim function stub.PutState. The first argument is the key as a string, and the second argument is the value as an array of bytes. This function may return an error which our code inspects and returns if present.

#Invoke()

Invoke is called when you want to call chaincode functions to do real work. Invocation transactions will be captured as blocks on the chain. The structure of Invoke is simple. It receives a function argument and based on this argument calls Go functions in the chaincode.

In your chaincode_start.go file, change the Invoke function so that it calls a generic write function.

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

    // Handle different functions
    if function == "init" {
        return t.Init(stub, "init", args)
    } else if function == "write" {
        return t.write(stub, args)
    }
    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation")
}
Now that it’s looking for write let’s make that function somewhere in your chaincode_start.go file.

func (t *SimpleChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var key, value string
    var err error
    fmt.Println("running write()")

    if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
    }

    key = args[0]                            //rename for fun
    value = args[1]
    err = stub.PutState(key, []byte(value))  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
    return nil, nil
}
This write function should look similar to the Init change you just did. One major difference is that you can now set the key and value for PutState. This function allows you to store any key/value pair you want into the blockchain ledger.

#Query()

As the name implies, Query is called whenever you query your chaincode state. Queries do not result in blocks being added to the chain. You will use Query to read the value of your chaincode state's key/value pairs.

In your chaincode_start.go file, change the Query function so that it calls a generic read function.

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)

    // Handle different functions
    if function == "read" {                            //read a variable
        return t.read(stub, args)
    }
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query")
}
Now that it’s looking for read, make that function somewhere in your chaincode_start.go file.

func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var key, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil
}
This read function is using the complement to PutState called GetState. This shim function just takes one string argument. The argument is the name of the key to retrieve. Next, this function returns the value as an array of bytes back to Query, who in turn sends it back to the REST handler.

#Main()

Finally, you need to create a short main function that will execute when each peer deploys their instance of the chaincode. It just starts the chaincode and registers it with the peer. You don’t need to add any code for this function. Both chaincode_start.go and chaincode_finished.go have a main function that lives at the top of the file. The function looks like this:

func main() {
    err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}

