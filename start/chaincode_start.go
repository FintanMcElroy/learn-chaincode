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

// multiple imports using import ()
import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	// err is assigned the result of shim.Start and will either be nil (no error) or an error object
	err := shim.Start(new(SimpleChaincode))
	// Note you don't have () around the if test, however the statements that follow must be in {} even if only 1 statement
	if err != nil {
		// call the Printf method on fmt reference which we imported
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
// the (t *SimpleChaincode) designates the method receiver - it will receive a pointer to a struct of type SimpleChaincode and store that as variable t
// func name is initial capital to show it is a method on a method receiver (I think also to make it a public method that can be invoked from outside)
// method receives 3 parameters - stub of type shim.ChaincodeStubInterface, function of type string, and args of type []string (string array)
// method will return 2 paramters, the first of type []byte (byte array) and the second of type error
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	// write to the blockchain by invoking stub.PutState()
	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}
	// Return null values for each of the declared return parameters
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
// the method signature is the same as Init - see there for comments on the different parts
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "write" {
		// call the write() method on the t pointer
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function) //error

	// return nil for first return value, for second create an error object using the errors.New() method
	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is called to read blocks from the blockchain
// the method signature is the same as Init - see there for comments on the different parts
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		// call the read() method on the t pointer
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query")
}

// Perform the write of a new block to the blockchain
// the method signature is the same as Init - see there for comments on the different parts
// I think the name is initial lowercase to make this a private method
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// declare variables of type string
	var name, value string
	// declare variable of type error
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}
	// assign new values to name, value
	name = args[0]
	value = args[1]
	//write the variable into the chaincode state
	err = stub.PutState(name, []byte(value))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Perform the read of a block on the blockchain
// the method signature is the same as Init - see there for comments on the different parts
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	// This is an interesting one. I believe from Go tutorial that you either do var x, x = 1 or x :=1 to initialise a variable
	// Once initiliased you can't then do := again. However in the below we have 2 target variables, one not declared before and one already declared (err)
	// so maybe this is a special case where it is OK to use := as one of the targets has not been declared and initialised previously
	valAsbytes, err := stub.GetState(name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
