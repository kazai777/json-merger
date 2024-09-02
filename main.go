package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	input := getInput()
	files := getFiles(input)
	if len(files) == 0 {
		fmt.Println("No JSON files found.")
		return
	}

	outputFileName := getOutputFileName()

	mergedData := mergeFiles(files)

	err := writeOutputFile(outputFileName, mergedData)
	if err != nil {
		fmt.Printf("Error writing output file %s: %v\n", outputFileName, err)
		return
	}

	fmt.Printf("JSON files successfully merged into %s\n", outputFileName)
}

// getInput prompts the user for input and returns it
func getInput() string {
	var input string
	fmt.Println("Please enter the JSON files to merge (separated by spaces), or a directory, or '.' for the current directory:")
	fmt.Scanln(&input)
	return input
}

// getFiles returns a list of files based on the user input
func getFiles(input string) []string {
	if input == "." {
		files, _ := filepath.Glob("*.json")
		return files
	}

	fileInfo, err := os.Stat(input)
	if err == nil && fileInfo.IsDir() {
		files, _ := filepath.Glob(filepath.Join(input, "*.json"))
		return files
	}

	return strings.Split(input, " ")
}

// getOutputFileName prompts the user for the output file name and ensures it has a .json extension
func getOutputFileName() string {
	var outputFileName string
	fmt.Println("Please enter the name of the output JSON file:")
	fmt.Scanln(&outputFileName)

	if !strings.HasSuffix(outputFileName, ".json") {
		outputFileName += ".json"
	}

	return outputFileName
}

// mergeFiles reads and merges JSON data from the provided files
func mergeFiles(files []string) []interface{} {
	var mergedData []interface{}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file, err)
			continue
		}

		var jsonData interface{}
		err = json.Unmarshal(data, &jsonData)
		if err != nil {
			fmt.Printf("Error parsing JSON from file %s: %v\n", file, err)
			continue
		}

		switch jd := jsonData.(type) {
		case []interface{}:
			mergedData = append(mergedData, jd...)
		case map[string]interface{}:
			mergedData = append(mergedData, jd)
		default:
			fmt.Printf("Unexpected JSON content type in file %s\n", file)
		}
	}

	return mergedData
}

// writeOutputFile writes the merged JSON data to the specified output file
func writeOutputFile(outputFileName string, mergedData []interface{}) error {
	mergedJson, err := json.MarshalIndent(mergedData, "", "  ")
	if err != nil {
		return fmt.Errorf("error generating output JSON: %v", err)
	}

	err = os.WriteFile(outputFileName, mergedJson, 0o644)
	if err != nil {
		return fmt.Errorf("error writing output file %s: %v", outputFileName, err)
	}

	return nil
}
