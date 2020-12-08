// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario

 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
 */

package main

/* Imports  
* 4 utility libraries for handling bytes, reading and writing JSON, 
formatting, and string manipulation  
* 2 specific Hyperledger Fabric specific libraries for Smart Contracts  
*/ 
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"github.com/iota"
)

// Define the Smart Contract structure
type SmartContract struct {
}

/* Define Container structure, with 4 properties.  
Structure tags are used by encoding/json library
*/
type Container struct {
	Description string `json:"description"`
	Timestamp string `json:"timestamp"`
	Location  string `json:"location"`
	Holder  string `json:"holder"`
}

type Sample struct {
	Force string `json:"force"`
	Stretching string `json:"stretching"`
	Holder  string `json:"holder"`
	Timestamp string `json:"timestamp"`
}

type IotaWallet struct {
	Seed        string `json:"seed"`
	Address     string `json:"address"`
	KeyIndex    uint64 `json:"keyIndex"`
}

type Participant struct {
	Role string `json:"role"`
	Description string `json:"description"`
	IotaWallet
}

type IotaPayload struct {
	Seed        string `json:"seed"`
	MamState    string `json:"mamState"`
	Root        string `json:"root"`
	Mode       	string `json:"mode"`
	SideKey     string `json:"sideKey"`
}

/*
 * The Init method *
 called when the Smart Contract "container-chaincode" is instantiated by the network
 * Best practice is to have any Ledger initialization in separate function 
 -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "container-chaincode"
 The app also specifies the specific smart contract function to call with args
 */

 // https://hyperledger-fabric.readthedocs.io/en/release-1.4/chaincode4ade.html
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "querySample" {
		return s.querySample(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "recordSample" {
		return s.recordSample(APIstub, args)
	} else if function == "queryAllContainers" {
		return s.queryAllContainers(APIstub)
	} else if function == "deleteSample" {
		return s.deleteSample(APIstub, args)
	} else if function == "getHistoryForSample" {
		return s.getHistoryForSample(APIstub, args)
	} 

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The initLedger method *
Will add test data (10 containers) to our network
 */
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	timestamp := strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)

	// Insert samples
	samples := []Sample{
		Sample{Force: "58.22", Stretching: "1", Holder: "A", Timestamp: timestamp},
		Sample{Force: "104, 25", Stretching: "2", Holder: "B", Timestamp: timestamp},
		Sample{Force: "143.29", Stretching: "3", Holder: "C", Timestamp: timestamp},
	}

	i := 0

	for i < len(samples) {
		sampleAsBytes, _ := json.Marshal(samples[i])
		APIstub.PutState(strconv.Itoa(i+1), sampleAsBytes)

		// Define own values for IOTA MAM message mode and MAM message encryption key
		// If not set, default values from iota/config.go file will be used
		mode := iota.MamMode
		sideKey := iota.PadSideKey(iota.MamSideKey) // iota.PadSideKey(iota.GenerateRandomSeedString(50))
		
		mamState, root, seed := iota.PublishAndReturnState(string(sampleAsBytes), false, "", "", mode, sideKey)
		iotaPayload := IotaPayload{Seed: seed, MamState: mamState, Root: root, Mode: mode, SideKey: sideKey}
		iotaPayloadAsBytes, _ := json.Marshal(iotaPayload)
		APIstub.PutState("IOTA_" + strconv.Itoa(i+1), iotaPayloadAsBytes)

		fmt.Println("New Asset", strconv.Itoa(i+1), samples[i], root, mode, sideKey)
		
		i = i + 1
	}

	participants := []Participant{
		Participant{Role: "A", Description: "Participant A"},
		Participant{Role: "B", Description: "Participant B"},
		Participant{Role: "C", Description: "Participant C"},
	}

	for i := range participants {
		walletAddress, walletSeed := iota.CreateWallet()
		participants[i].Seed = walletSeed
		participants[i].Address = walletAddress
		participants[i].KeyIndex = 0
		participantAsBytes, _ := json.Marshal(participants[i])
		APIstub.PutState(participants[i].Role, participantAsBytes)
	}

	iotaWallet := IotaWallet{Seed: iota.DefaultWalletSeed, KeyIndex: iota.DefaultWalletKeyIndex, Address: ""}
	iotaWalletAsBytes, _ := json.Marshal(iotaWallet)
	APIstub.PutState("IOTA_WALLET", iotaWalletAsBytes)

	return shim.Success(nil)
}

/*
 * The recordSample method *
Container owners like Sarah would use to record each of her containers. 
This method takes in five arguments (attributes to be saved in the ledger). 
 */
func (s *SmartContract) recordSample(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	timestamp := strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)
	sample := Sample{ Force: args[1], Stretching: args[2], Holder: args[3], Timestamp: timestamp }

	sampleAsBytes, _ := json.Marshal(sample)
	err := APIstub.PutState(args[0], sampleAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record sample: %s", args[0]))
	}

	// Define own values for IOTA MAM message mode and MAM message encryption key
	// If not set, default values from iota/config.go file will be used
	mode := iota.MamMode
	sideKey := iota.PadSideKey(iota.MamSideKey) // iota.PadSideKey(iota.GenerateRandomSeedString(50))
	
	mamState, root, seed := iota.PublishAndReturnState(string(sampleAsBytes), false, "", "", mode, sideKey)	
	iotaPayload := IotaPayload{Seed: seed, MamState: mamState, Root: root, Mode: mode, SideKey: sideKey}
	iotaPayloadAsBytes, _ := json.Marshal(iotaPayload)
	APIstub.PutState("IOTA_" + args[0], iotaPayloadAsBytes)

	fmt.Println("New Asset", args[0], sample, root, mode, sideKey)

	return shim.Success(nil)
}

/*
 * The querySample method *
Used to view the records of one particular sample
It takes one argument -- the key for the sample in question
 */
 func (s *SmartContract) querySample(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	sampleAsBytes, _ := APIstub.GetState(args[0])
	if sampleAsBytes == nil {
		return shim.Error("Could not locate sample")
	}
	sample := Sample{}
	json.Unmarshal(sampleAsBytes, &sample)

	iotaPayloadAsBytes, _ := APIstub.GetState("IOTA_" + args[0])
	if iotaPayloadAsBytes == nil {
		return shim.Error("Could not locate IOTA state object")
	}
	iotaPayload := IotaPayload{}
	json.Unmarshal(iotaPayloadAsBytes, &iotaPayload)

	mamstate := map[string]interface{}{}
	mamstate["root"] = iotaPayload.Root
	mamstate["sideKey"] = iotaPayload.SideKey

	// IOTA MAM stream values
	messages := iota.Fetch(iotaPayload.Root, iotaPayload.Mode, iotaPayload.SideKey)

	participantAsBytes, _ := APIstub.GetState(sample.Holder)
	if participantAsBytes == nil {
		return shim.Error("Could not locate participant")
	}
	participant := Participant{}
	json.Unmarshal(participantAsBytes, &participant)

	out := map[string]interface{}{}
	out["sample"] = sample
	out["mamstate"] = mamstate
	out["messages"] = strings.Join(messages, ", ")
	out["wallet"] = participant.Address
	
	result, _ := json.Marshal(out)

	return shim.Success(result)
}

/*
 * The queryAllContainers method *
allows for assessing all the records added to the ledger(all containers)
This method does not take any arguments. Returns JSON string containing results. 
 */
func (s *SmartContract) queryAllContainers(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := "0"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		
		// iotaPayloadAsBytes, _ := APIstub.GetState("IOTA_" + queryResponse.Key)
		// if iotaPayloadAsBytes == nil {
		// 	return shim.Error("Could not locate IOTA state object")
		// }
		// iotaPayload := IotaPayload{}
		// json.Unmarshal(iotaPayloadAsBytes, &iotaPayload)

		// buffer.WriteString(", \"Root\":")
		// buffer.WriteString("\"" + string(iotaPayload.Root) + "\"")
		// buffer.WriteString(", \"SideKey\":")
		// buffer.WriteString("\"" + string(iotaPayload.SideKey) + "\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllContainers:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The changeSample method *
The data in the world state can be updated with who has possession. 
This function takes in 2 arguments, container ID and new holder name. 
 */
func (s *SmartContract) changeSample(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	updateData := map[string]interface{}{}
	json.Unmarshal([]byte(args[1]), &updateData)
	sampleAsBytes, _ := APIstub.GetState(args[0])

	if sampleAsBytes == nil {
		return shim.Error("Could not locate sample")
	}

	sample := Sample{}

	json.Unmarshal(sampleAsBytes, &sample)
	// Normally check that the specified argument is a valid holder of a sample
	// we are skipping this check for this example

	if updateData["force"] != nil {
		sample.Force = updateData["force"].(string)
	}
	
	if updateData["stretching"] != nil {
		sample.Stretching = updateData["stretching"].(string)
	} 
	
	if updateData["holder"] != nil {
		sample.Holder = updateData["holder"].(string)
	}

	timestamp := strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)
	sample.Timestamp = timestamp

	sampleAsBytes, _ = json.Marshal(sample)
	err := APIstub.PutState(args[0], sampleAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change sample holder: %s", args[0]))
	}

	iotaPayloadAsBytes, _ := APIstub.GetState("IOTA_" + args[0])
	if iotaPayloadAsBytes == nil {
		return shim.Error("Could not locate IOTA state object")
	}
	iotaPayload := IotaPayload{}
	json.Unmarshal(iotaPayloadAsBytes, &iotaPayload)

	mamState, _, _ := iota.PublishAndReturnState(string(sampleAsBytes), true, iotaPayload.Seed, iotaPayload.MamState, iotaPayload.Mode, iotaPayload.SideKey)
	iotaPayloadNew := IotaPayload{Seed: iotaPayload.Seed, MamState: mamState, Root: iotaPayload.Root, Mode: iotaPayload.Mode, SideKey: iotaPayload.SideKey}
	iotaPayloadNewAsBytes, _ := json.Marshal(iotaPayloadNew)
	APIstub.PutState("IOTA_" + args[0], iotaPayloadNewAsBytes)

	return shim.Success([]byte("changeSample success"))
}

// Deletes an entity from state
func (t *SmartContract) deleteSample(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := args[0]

	sampleAsBytes, _ := APIstub.GetState(args[0])

	if sampleAsBytes == nil {
		return shim.Error("Could not locate sample")
	}

	// Delete the key from the state in ledger
	err := APIstub.DelState(key)
	if err != nil {
		return shim.Error("Failed to delete sample")
	}

	return shim.Success([]byte("deleteSample success"))
}

func (t *SmartContract) getHistoryForSample(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	sampleKey := args[0]

	fmt.Printf("- start getHistoryForSample: %s\n", sampleKey)

	resultsIterator, err := APIstub.GetHistoryForKey(sampleKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForSample returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * main function *
calls the Start function 
The main function starts the chaincode in the container during instantiation.
 */
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
