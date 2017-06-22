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

// ===================================================================================
// Init
// ===================================================================================

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, nil
}

// ===================================================================================
// Invoke
// ===================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 5 {
		return nil errors.New("Incorrect number of arguments. Expecting 5")
	}
	
	var flowID,productID,protocolID,buyer,vendor string
	
	// Initialize the chaincode
	flowID = args[0]
	productID = args[1]
	protocolID = args[2]
	buyer = args[3]
	vendor = args[4]
	
	// ==== Check if transaction already exists ====
	transAsBytes, err := stub.GetState(flowID)
	if err != nil {
		return shim.Error("Failed to get flowID"+err.Error())
	} else if transAsBytes !=nil {
		fmt.Println("This transaction already exists: "+flowID)
		return shim.Error("This transaction already exists: "+flowID)
	}
	
	// ==== Create transInfo object and marshal to JSON ====
	objectType := "transinfo"
	transInfo := &transInfo{objectType,flowID,productID,protocolID,buyer,vendor}
	transInfoJSONasBytes, err := json.Marshal(transInfo)
	if err != nil {
		return shim.Error(err.Error())
	} 
	
	// === Save marble to state ===
	err = stub.PutState(transactionFlowID,transInfoJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
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
// Query callback representing the query of a chaincode
// ===================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var A string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return Avalbytes, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
