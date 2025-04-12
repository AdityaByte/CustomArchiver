package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
)

func TestUnarchive(t *testing.T) {
	const (
		file1name = "test1.txt"
		file2name = "test2.txt"
		content1  = "aditya pawar"
		content2  = "golang developer"
	)

	createFileTest(t, file1name, content1)
	defer func() {
		if err := os.Remove(file1name); err != nil {
			t.Fatalf("Failed to remove the file: %s", file1name)
		} else {
			t.Logf("File removed successfully %s", file1name)
		}
	}()
	createFileTest(t, file2name, content2)
	defer func() {
		if err := os.Remove(file2name); err != nil {
			t.Fatalf("Failed to remove the file: %s", file2name)
		} else {
			t.Logf("File removed successfully %s", file2name)
		}
	}()

	fileArgs := fmt.Sprintf("%s %s", file1name, file2name)
	outputFilename := "test"

	if err := archive(&fileArgs, &outputFilename); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := os.Remove(outputFilename); err != nil {
			t.Fatalf("Failed to remove the file: %s,  ERROR: %v", outputFilename, err)
		} else {
			t.Logf("File removed successfully: %s", outputFilename)
		}
	}()

	// Here the main logic comes in we have to unarchive the file and check the data is correct or not.

	if err := unArchive(&outputFilename); err != nil {
		if err == io.EOF {
			t.Log("Got the EOF")
		} else {
			t.Fatal(err)
		}
	}

	destination := "storage"

	defer func() {
		if err := os.RemoveAll(destination); err != nil {
			t.Fatalf("Failed to remove the directory: %s, ERROR: %v", destination, err)
		} else {
			t.Logf("Directory removed successfully %s", destination)
		}
	} ()

	wg := sync.WaitGroup{} // created an instance of the waitgroup.
	wg.Add(2)

	errorChan := make(chan error, 2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		file1Data, err := os.ReadFile(fmt.Sprintf("%s/%s", destination, file1name))
		if err != nil {
			errorChan <- fmt.Errorf("Failed to read the filedata, %s, ERROR: %v", file1name, err)
			return
		}

		if len(string(file1Data)) != len(content1) {
			errorChan <- fmt.Errorf("Corrupted data: %s", file1name)
			return
		}

		errorChan <- nil
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		file2Data, err := os.ReadFile(fmt.Sprintf("%s/%s", destination, file2name))
		if err != nil {
			errorChan <- fmt.Errorf("Failed to read the filedata, %s, ERROR: %v", file2name, err)
			return
		}

		if len(string(file2Data)) != len(content2) {
			errorChan <- fmt.Errorf("Corrupted data: %s", file2name)
			return
		}

		errorChan <- nil
	}(&wg)

	wg.Wait()
	close(errorChan)

	for err := range errorChan {
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Unarchived Test Passed")

}
