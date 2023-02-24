package main

import (
	"context"
	"os"
	"strings"

	"github.com/dogefuzz/inputs/pkg/common"
	"github.com/dogefuzz/inputs/pkg/solc"
	"github.com/dogefuzz/inputs/pkg/vandal"
)

var criticalInstructions = []string{"CALL", "SELFDESTRUCT", "CALLCODE", "DELEGATECALL"}

func main() {
	contractSlice := make([]common.ContractInfo, 0)
	data := make([][]string, 0)

	data = common.ReadCsvFile(os.Args[1])

	contractSlice = createContractInfoSlice(data)

	for i := 0; i < len(contractSlice); i++ {
		contractContent, err := readContractContent(contractSlice[i].Name)

		if err == nil {
			contractName := strings.Split(contractSlice[i].Name, ".")

			blocks, branches, criticalInstructions := getNumberOfBlocksAndCriticalInstructions(contractName[0], contractContent)

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

	contract.Name = row[0]
	contract.Link = row[2]
	if v, exists := contractMap[contract.Name]; exists {
		contract.Weaknesses = append(v.Weaknesses, row[1])
	} else {
		contract.Weaknesses = append(contract.Weaknesses, row[1])
	}

	contractMap[contract.Name] = contract
}

func readContractContent(contractName string) (string, error) {
	data, err := os.ReadFile(contractName)

	return string(data), err
}

func getNumberOfBlocksAndCriticalInstructions(contractName string, contractContent string) (int, int, map[string]int) {
	var blocks, branches int
	criticalInstructionsMap := make(map[string]int)

	compiler := solc.NewSolidityCompiler("/tmp/dogefuzz/")
	name := strings.Split(contractName, ".")
	contract, _ := compiler.CompileSource(name[0], contractContent)

	c := vandal.NewVandalClient("http://localhost:5005")
	blockSlice, _, _ := c.Decompile(context.Background(), contract.CompiledCode)

	blocks = len(blockSlice)

	for _, block := range blockSlice {
		branches += len(block.Predecessors)
		for _, v := range block.Instructions {
			for i := 0; i < len(criticalInstructions); i++ {
				if v.Op == criticalInstructions[i] {
					criticalInstructionsMap[v.Op]++
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
