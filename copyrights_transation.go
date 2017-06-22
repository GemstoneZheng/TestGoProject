// copyrights_transation
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// ===================================================================================
// Invoke
// ===================================================================================

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "initTradeInfo" {
		return t.initTradeInfo(stub, args)
	} else if function == "readTransInfo" {
		return t.readTransInfo(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *SimpleChaincode) initTradeInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init marble")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	}

	var flowID, productID, protocolID, buyer, vendor string

	// Initialize the chaincode
	flowID = args[0]
	productID = args[1]
	protocolID = args[2]
	buyer = args[3]
	vendor = args[4]

	// ==== Check if transaction already exists ====
	transAsBytes, err := stub.GetState(flowID)
	if err != nil {
		return shim.Error("Failed to get flowID" + err.Error())
	} else if transAsBytes != nil {
		fmt.Println("This transaction already exists: " + flowID)
		return shim.Error("This transaction already exists: " + flowID)
	}

	// ==== Create transInfo object and marshal to JSON ====
	objectType := "transinfo"
	transInfo := &transInfo{objectType, flowID, productID, protocolID, buyer, vendor}
	transInfoJSONasBytes, err := json.Marshal(transInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save marble to state ===
	err = stub.PutState(transactionFlowID, transInfoJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== trandeInfo saved and indexed. Return success ====
	fmt.Println("- end init marble")
	return shim.Success(nil)
}

// ==================================================
// delete
// ==================================================
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	transInfo := args[0]

	valAsbytes, err := stub.GetState(transInfo)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + transInfo + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"TransInfo does not exist: " + transInfo + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(transInfo)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	return shim.Success(nil)
}

// ===================================================================================
// Deletes an entity from state
// ===================================================================================
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// ===================================================================================
// Query
// ===================================================================================
func (t *SimpleChaincode) readTransInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the marble to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
