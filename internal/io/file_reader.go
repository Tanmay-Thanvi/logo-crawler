package io

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// ReadPublishers reads publishers from a file
func ReadPublishers(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var publishers []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		publisher := strings.TrimSpace(scanner.Text())
		if publisher != "" {
			publishers = append(publishers, publisher)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	return publishers, nil
}
