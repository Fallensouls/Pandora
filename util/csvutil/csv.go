package csvutil

import (
	"encoding/csv"
	"io"
	"os"
)

func GetTestData(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data [][]string

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		data = append(data, record)
	}
	return data, nil
}
