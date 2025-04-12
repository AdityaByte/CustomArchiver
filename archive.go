package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func archive(archiveFile *string, archiveFileName *string) error {
	var buf bytes.Buffer

	trimText(archiveFile) // If the user enter any leading or trailing space mistakenly it will remove that.
	var args []string = strings.Split(*archiveFile, " ")

	if len(args) == 0 {
		return fmt.Errorf("No file provided")
	}

	fmt.Println("Length of args:", len(args))

	for index, filename := range args {

		fmt.Println("Filename:", filename)

		data, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("Failed to read the file %s : %v", filename, err)
		}

		name := filepath.Base(filename)

		buf.Write([]byte(fmt.Sprintf("%s\n%d\n%s\n", name, len(data), "::END-METADATA::")))
		buf.Write(data)
		buf.Write([]byte("\n::END-FILE::\n"))

		log.Printf("Task has been done %d\n", index)
	}

	createArchive(buf.Bytes(), archiveFileName)
	return nil
}

func createArchive(data []byte, filename *string) error {
	*filename = fmt.Sprintf("%s.adzip", *filename)
	file, err := os.Create(*filename) // By doing this we opened the file so we don't need to open it more via WriteFile()
	if err != nil {
		return fmt.Errorf("Failed to create the file %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("Failed to write in the file: %v", err)
	}

	return nil
}
