package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type enumArr []string

type NumNum struct {
	Directory string              `json:"path"`
	Enums     map[string][]string `json:"enums"`
}

func createLine(s string) string {
	return fmt.Sprintf("\t%s = \"%s\",", s, s)
}

func main() {
	filename := "enums.json"

	fileData, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err.Error())
		return
	}

	var result NumNum
	if err = json.Unmarshal(fileData, &result); err != nil {
		fmt.Printf("Error parsing json: %s\n", err.Error())
		return
	}

	if err := os.MkdirAll(result.Directory, 0666); err != nil {
		fmt.Printf("Error creating directory %s\n", result.Directory)
		return
	}

	for enumName, enumArr := range result.Enums {
		filePath := fmt.Sprintf("%s/%s.luau", result.Directory, enumName)

		file, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Error creating file %s; skipping\n", filePath)
			continue
		}
		defer file.Close()

		var tblOutput string
		typeOutput := fmt.Sprintf("type %s = ", enumName)
		for i := range enumArr {
			elem := enumArr[i]
			if i > 0 {
				tblOutput += "\n"
				typeOutput += " | "
			}
			tblOutput += createLine(elem)
			typeOutput += fmt.Sprintf("\"%s\"", elem)
		}

		output := fmt.Sprintf("%s\n\nreturn table.freeze {\n%s\n}", typeOutput, tblOutput)
		if _, err := file.WriteString(output); err != nil {
			fmt.Printf("Error writing to file %s; skipping\n", filePath)
			continue
		}
	}
}
