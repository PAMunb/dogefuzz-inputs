package common

import (
	"encoding/json"
	"log"
	"os"
)

func GenerateJsonFileFromContractInfoSlice(contractSlice []ContractInfo) {
	file, err := json.MarshalIndent(contractSlice, "", "   ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("result.json", file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
