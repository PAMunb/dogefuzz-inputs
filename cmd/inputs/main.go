package main

import (
	"context"
	"os"
	"strings"

	"github.com/dogefuzz/inputs/pkg/common"
	"github.com/dogefuzz/inputs/pkg/solc"
	"github.com/dogefuzz/inputs/pkg/vandal"
)

func main() {
	contractSlice := make([]common.ContractInfo, 0)
	data := make([][]string, 0)

	data = common.ReadCsvFile(os.Args[1])

	contractSlice = createContractSlice(data)

	for i := 0; i < len(contractSlice); i++ {
		contractContent, err := readContractContent(contractSlice[i].Name)

		if err == nil {
			contractName := strings.Split(contractSlice[i].Name, ".")

			blocks, predecessors, criticalInstructions := getNumberOfBlocksAndCriticalInstructions(contractName[0], contractContent)

			setNumberOfBlocks(&contractSlice[i], blocks)
			setNumberOfPredecessors(&contractSlice[i], predecessors)
			setNumberOfCriticalInstructions(&contractSlice[i], criticalInstructions)
		}
	}
}

func setNumberOfBlocks(contract *common.ContractInfo, numberOfBlocks int) {
	contract.NumberOfBlocks = numberOfBlocks
}

func setNumberOfPredecessors(contract *common.ContractInfo, numberOfPredecessors int) {
	contract.NumberOfPredecessors = numberOfPredecessors
}

func setNumberOfCriticalInstructions(contract *common.ContractInfo, numberOfCriticalInstructions int) {
	contract.NumberOfCriticalInstructions = numberOfCriticalInstructions
}

func getNumberOfBlocksAndCriticalInstructions(contractName string, contractContent string) (int, int, int) {
	var nos, arestas, instrucoes int

	compiler := solc.NewSolidityCompiler("/tmp/dogefuzz/")
	name := strings.Split(contractName, ".")
	contract, _ := compiler.CompileSource(name[0], contractContent)

	c := vandal.NewVandalClient("http://localhost:5005")
	blocks, _, _ := c.Decompile(context.Background(), contract.CompiledCode)

	nos = len(blocks)
	for _, block := range blocks {
		arestas += len(block.Predecessors)
		for _, v := range block.Instructions {
			switch v.Op {
			case "CALL":
				instrucoes++
			case "SELFDESTRUCT":
				instrucoes++
			case "CALLCODE":
				instrucoes++
			case "DELEGATECALL":
				instrucoes++
			}
		}
	}

	return nos, arestas, instrucoes
}

func readContractContent(contractName string) (string, error) {
	data, err := os.ReadFile("contracts/" + contractName)

	return string(data), err
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
