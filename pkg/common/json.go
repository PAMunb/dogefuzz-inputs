package common

import (
	"encoding/json"
	"log"
	"os"
)

const RESULT_JSON_PATH = "result.json"

func GenerateJsonFileFromContractInfoSlice(contractSlice []ContractInfo) {
	file, err := json.MarshalIndent(contractSlice, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(RESULT_JSON_PATH, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
