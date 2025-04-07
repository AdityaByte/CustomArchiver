package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Here what we have to do
// We want to build a custom archiver
// So the architecture is like that
/*
We get a file we read that
// Structure - seralization
// Now the structure looks like
Filename
Filesize
DELIMETTER - ::END-METADATA::
DATA
DELIMETTER - ::END-FILE::
// Here new file with the same structure
*/

func main() {

	var archiveFile string
	var unarchiveFile string
	var archiveFileName string

	flag.StringVar(&archiveFile, "archive", "", "Files to archive")
	flag.StringVar(&unarchiveFile, "unarchive", "", "File to unarchive")
	flag.StringVar(&archiveFileName, "o", "data", "output file name")

	flag.Parse() // It's important to parse.

	if archiveFile != "" {
		if err := archive(&archiveFile, &archiveFileName); err != nil {
			log.Fatal(err)
		}
	}

	if unarchiveFile != "" {
		if err := unArchive(&unarchiveFile); err != nil {
			if err == io.EOF {
				log.Println("File unarchived successfully")
				return
			}
			log.Fatalf("ERROR: %v", err)
		}
	}

}

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

// In this function we want to unarchive the file like we want to read each and every bytes of it till we get
// those things that we want.
func unArchive(filename *string) error {

	if filename == nil {
		return fmt.Errorf("ERROR: File is empty")
	}

	file, err := os.Open(*filename)
	if err != nil {
		return fmt.Errorf("Failed to open the file name: %s : %v", file.Name(), err)
	}

	defer file.Close()

	if !strings.HasSuffix(file.Name(), ".adzip") {
		return fmt.Errorf("We Can't unarchive that file format", file.Name())
	}

	data, err := os.ReadFile(file.Name())
	if err != nil {
		return fmt.Errorf("Failed to read the file %v", err)
	}

	log.Println("Data length:", len(data))

	// Since we have many files inside an archive and we have to retrive it
	// from that bundle so we have to go through a for infinite loop here.

	// Here we have to define a reader which reads the data line by line
	reader := bufio.NewReader(bytes.NewReader(data))

	// Structure
	/*	filename \n
		filesize \n
		::END-METADATA::\n
		[]byte data\n
		::END-FILE::
	*/

	for {
		filename, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return io.EOF
			}
			return fmt.Errorf("Failed to read the file name: %v", err)
		}

		trimText(&filename)
		log.Println("Filename:", filename)

		filesize, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("Failed to read the file size: %v", err)
		}

		trimText(&filesize)
		actualFileSize, err := strconv.Atoi(filesize)

		if err != nil {
			return fmt.Errorf("Failed to parse the raw file size to integer: %v", err)
		}

		fmt.Println("File size:", actualFileSize)

		endMetaDataDellimeter, err := reader.ReadString('\n')
		trimText(&endMetaDataDellimeter)
		if err != nil || endMetaDataDellimeter != "::END-METADATA::" {
			return fmt.Errorf("Failed to read the metadata", err)
		}

		myData := make([]byte, actualFileSize)

		// readedData, err := reader.Read(myData) // Since the reader.Read has no guarentee that it can read the full data
		// Instead of this use this io.ReadFull function.

		readedData, err := io.ReadFull(reader, myData)
		if err != nil {
			return fmt.Errorf("Failed to read the actual data: %v", err)
		}

		if readedData != actualFileSize {
			return fmt.Errorf("Failed to read the actual data size read this much: %d", readedData)
		}

		if _, err := reader.ReadString('\n'); err != nil {
			return fmt.Errorf("Failed to read the new line delimetter: %v", err)
		}

		endFileDelimetter, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("Failed to percieve the end file delimetter: %v", err)
		}

		trimText(&endFileDelimetter)
		fmt.Println("End file delimetter:", endFileDelimetter)
		if endFileDelimetter != "::END-FILE::" {
			return fmt.Errorf("Failed to get the end file delimetter")
		}

		destination := fmt.Sprintf("storage/")

		if _, err := os.Stat(destination); os.IsNotExist(err) {
			if err := os.Mkdir(destination, os.ModePerm); err != nil {
				return fmt.Errorf("Failed to create the directory: %v", err)
			} else {
				log.Println("Directory created successfully.")
			}
		} else {
			log.Println("Directory already exists..")
		}

		myFile, err := os.Create(fmt.Sprintf("%s%s", destination, filename))
		if err != nil {
			return fmt.Errorf("Failed to create the file: %v", err)
		}

		fmt.Println("Filename:", myFile.Name())

		_, err = myFile.Write(myData)
		if err != nil {
			return fmt.Errorf("Failed to write the file: %v", err)
		}
	}
}

func trimText(str *string) {
	*str = strings.TrimSpace(*str)
}
