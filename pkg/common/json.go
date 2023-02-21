package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func GenerateJsonFileFromContractInfoSlice(contractSlice []ContractInfo) {
	file, err := json.MarshalIndent(contractSlice, "", "   ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("contracts/result.json", file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
