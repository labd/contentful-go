package testutil

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

func readTestData(fileName string) string {
	path := "../../testdata/" + fileName
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return string(content)
}

func ModelFromTestData(fileName string, value any) error {
	content := readTestData(fileName)

	err := json.NewDecoder(strings.NewReader(content)).Decode(&value)
	if err != nil {
		return err
	}

	return nil
}
