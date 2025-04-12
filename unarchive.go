package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// In this function we want to unarchive the file like we want to read each and every bytes of it till we get
// those things that we want.
func unArchive(filename *string) error {

	if filename == nil {
		return fmt.Errorf("ERROR: File is empty")
	}

	file, err := os.Open(*filename)
	myFileName := file.Name()
	if err != nil {
		return fmt.Errorf("Failed to open the file name: %s : %v", myFileName, err)
	}

	defer file.Close()

	if !strings.HasSuffix(file.Name(), ".adzip") {
		return fmt.Errorf("We Can't unarchive that file format %s", myFileName)
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
			return fmt.Errorf("Failed to read the metadata %v", err)
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
		defer myFile.Close()

		fmt.Println("Filename:", myFile.Name())

		_, err = myFile.Write(myData)
		if err != nil {
			return fmt.Errorf("Failed to write the file: %v", err)
		}
	}
}
