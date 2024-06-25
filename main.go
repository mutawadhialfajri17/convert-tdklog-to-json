package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"convert-metadata-log-to-json/utils"
)

func main() {
	// read metadata.log
	metadataByte, err := os.ReadFile(utils.LogFilename)
	if err != nil {
		utils.PrintLog(err)
	}
	metadataStr := string(metadataByte)

	// check is metadata full
	if strings.Contains(metadataStr, "Show more") { // assume the request in metadata doesnt contain data `Show more`
		utils.PrintLog(errors.New("your metadata is not full"))
	}

	nestedMap := make(map[string]interface{})
	lines := strings.Split(metadataStr, "\n")
	for _, line := range lines {
		// split the line into key value
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			utils.PrintLog(fmt.Errorf("invalid input format: %s", line))
		}
		keyPath := parts[0]
		value := parts[1]

		// split the key path into individual keys
		keys := strings.Split(keyPath, ".")

		// remove "metadata" if exist
		if keys[0] == "metadata" {
			keys = keys[1:]
		}

		// traverse the nested map to find the correct place to insert the value
		currentLevel := nestedMap
		for i, key := range keys {
			// check if it the last key
			if i == len(keys)-1 {
				// try to convert the value to a boolean or a number if applicable
				if value == "true" {
					currentLevel[key] = true
				} else if value == "false" {
					currentLevel[key] = false
				} else if num, err := strconv.ParseFloat(value, 64); err == nil {
					currentLevel[key] = num
				} else if json.Valid([]byte(value)) { // also, check if the value is already a JSON array or object
					var jsonValue interface{}
					err := json.Unmarshal([]byte(value), &jsonValue)
					if err != nil {
						utils.PrintLog(err)
					}
					currentLevel[key] = jsonValue
				} else {
					currentLevel[key] = value
				}
			} else {
				// ensure the next level is a map
				if _, exists := currentLevel[key]; !exists {
					currentLevel[key] = make(map[string]interface{})
				}
				currentLevel = currentLevel[key].(map[string]interface{})
			}
		}
	}

	result, err := json.Marshal(nestedMap)
	if err != nil {
		utils.PrintLog(err)
	}
	fmt.Println(string(result))
}
