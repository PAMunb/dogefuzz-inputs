package main

import (
	"fmt"
	"os"

	"github.com/dogefuzz/inputs/pkg/common"
)

func main() {
	fmt.Println("Hello World")

	data := common.ReadCsvFile(os.Args[1])

	contractSlice := createContractSlice(data)

	fmt.Println(contractSlice)
}

func createContractSlice(data [][]string) []common.ContractInfo {
	contractMap := make(map[string]common.ContractInfo)
	for _, line := range data {
		addLineToContractMap(contractMap, line)
	}

	return convertMapToSlice(contractMap)
}

func convertMapToSlice(contractMap map[string]common.ContractInfo) []common.ContractInfo {
	var contractSlice []common.ContractInfo
	for _, v := range contractMap {
		contractSlice = append(contractSlice, v)
	}

	return contractSlice
}

func addLineToContractMap(contractMap map[string]common.ContractInfo, line []string) {
	var contract common.ContractInfo

	contract.Name = line[0]
	contract.Link = line[2]
	if v, exists := contractMap[contract.Name]; exists {
		contract.Weaknesses = append(v.Weaknesses, line[1])
	} else {
		contract.Weaknesses = append(contract.Weaknesses, line[1])
	}

	contractMap[contract.Name] = contract
}
