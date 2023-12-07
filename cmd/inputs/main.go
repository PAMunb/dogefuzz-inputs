package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/dogefuzz/inputs/pkg/common"
	"github.com/dogefuzz/inputs/pkg/solc"
	"github.com/dogefuzz/inputs/pkg/vandal"
)

var criticalInstructions = []string{"CALL", "SELFDESTRUCT", "CALLCODE", "DELEGATECALL"}

var CONTRACTS_PATH = "resources/contracts/"

func main() {
	contractSlice := make([]common.ContractInfo, 0)
	data := make([][]string, 0)

	data = common.ReadCsvFile(os.Args[1])

	contractSlice = createContractInfoSlice(data)

	for i := 0; i < len(contractSlice); i++ {
		contractContent, err := readContractContent(contractSlice[i].File)

		if err == nil {

			blocks, branches, criticalInstructions := getNumberOfBlocksAndCriticalInstructions(contractSlice[i].Name, contractContent)

			setNumberOfBlocks(&contractSlice[i], blocks)
			setNumberOfBranches(&contractSlice[i], branches)
			setNumberOfCriticalInstructions(&contractSlice[i], criticalInstructions)
		}
	}

	common.GenerateJsonFileFromContractInfoSlice(contractSlice)

}

func createContractInfoSlice(data [][]string) []common.ContractInfo {
	contractMap := make(map[string]common.ContractInfo)
	for _, row := range data {
		addFileRowToContractInfoMap(contractMap, row)
	}

	return convertContractInfoMapToContractInfoSlice(contractMap)
}

func convertContractInfoMapToContractInfoSlice(contractMap map[string]common.ContractInfo) []common.ContractInfo {
	var contractSlice []common.ContractInfo
	for _, v := range contractMap {
		contractSlice = append(contractSlice, v)
	}

	return contractSlice
}

func addFileRowToContractInfoMap(contractMap map[string]common.ContractInfo, row []string) {
	var contract common.ContractInfo

	contract.File = row[0]
	contract.Name = row[1]
	contract.Link = row[3]
	if v, exists := contractMap[contract.File]; exists {
		contract.Weaknesses = append(v.Weaknesses, row[2])
	} else {
		contract.Weaknesses = append(contract.Weaknesses, row[2])
	}

	contractMap[contract.File] = contract
}

func readContractContent(contractName string) ([]byte, error) {
	data, err := os.ReadFile(CONTRACTS_PATH + contractName)

	return data, err
}

func getNumberOfBlocksAndCriticalInstructions(contractName string, contractContent []byte) (int, int, map[string]int) {
	var blocks, branches int
	criticalInstructionsMap := make(map[string]int)

	compiler := solc.NewSolidityCompiler("/tmp/dogefuzz/")
	name := strings.Split(contractName, ".")
	contract, err := compiler.CompileSource(name[0], string(contractContent))
	if err != nil {
		log.Fatal(err)
	}

	c := vandal.NewVandalClient("http://localhost:5005")
	blockSlice, _, err := c.Decompile(context.Background(), contract.RuntimeBytecode, contractName)
	if err != nil {
		log.Fatal(err)
		return 0, 0, criticalInstructionsMap
	}

	blocks = len(blockSlice)

	for _, block := range blockSlice {
		branches += len(block.Predecessors)
		for _, v := range block.Instructions {
			for i := 0; i < len(criticalInstructions); i++ {
				if v.Op == criticalInstructions[i] {
					criticalInstructionsMap[v.Op]++
					break
				}
			}

		}
	}

	return blocks, branches, criticalInstructionsMap
}

func setNumberOfBlocks(contract *common.ContractInfo, numberOfBlocks int) {
	contract.NumberOfBlocks = numberOfBlocks
}

func setNumberOfBranches(contract *common.ContractInfo, numberOfBranches int) {
	contract.NumberOfBranches = numberOfBranches
}

func setNumberOfCriticalInstructions(contract *common.ContractInfo, numberOfCriticalInstructions map[string]int) {
	contract.NumberOfCriticalInstructions = numberOfCriticalInstructions
}
