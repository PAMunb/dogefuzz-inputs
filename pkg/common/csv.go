package common

import (
	"encoding/csv"
	"os"
)

func ReadCsvFile(pathFile string) [][]string {
	f, err := os.Open(pathFile)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	return data
}
