package main

import (
	"fmt"
	"os"
	"testing"
)

func createFileTest(t *testing.T, filename, content string) {

	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create file name: %s, ERROR: %v", filename, err)
	}

	n, err := file.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write content to file: %s, ERROR: %v", filename, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.Fatalf("Failed to close the file: %s", filename)
		} else {
			t.Logf("File closed successfully: %s", filename)
		}
	}()

	if n != len(content) {
		t.Fatalf("Incomplete write to file: %s", filename)
	}

	t.Logf("File created successfully: %s\n", filename)

}

func TestArchive(t *testing.T) {

	// Dummy Data:
	const (
		file1name    = "test1.txt"
		file2name    = "test2.txt"
		file1content = "This is the file1 - @author adityabyte"
		file2content = "This is the file2 - @author adityapawar"
	)

	createFileTest(t, file1name, file1content) // Note: &testing.T by doing this we are creating a new instance of test.
	defer func() {
		if err := os.Remove(file1name); err != nil {
			t.Fatalf("Failed to remove the file: %s", file1name)
		} else {
			t.Logf("File removed successfully: %s", file1name)
		}
	}()
	createFileTest(t, file2name, file2content)
	defer func() {
		if err := os.Remove(file2name); err != nil {
			t.Fatalf("Failed to remove the file: %s", file2name)
		} else {
			t.Logf("File removed successfully: %s", file2name)
		}
	}()

	fileArg := fmt.Sprintf("%s %s", file1name, file2name)
	outputFilename := "test"
	defer func() {
		if err := os.Remove(outputFilename); err != nil {
			t.Fatalf("Failed to remove the file: %s", outputFilename)
		} else {
			t.Logf("File removed successfully: %s", outputFilename)
		}
	}()

	if err := archive(&fileArg, &outputFilename); err != nil {
		t.Fatal(err)
	}

	// We don't have to explicitly add .adzip cause the string changed when it gets passed to the archive function.
	if _, err := os.Stat(outputFilename); os.IsNotExist(err) {
		t.Fatalf("Expected file %s.adzip not found", outputFilename)
	}

	t.Log("Archived Test Passed")
}
