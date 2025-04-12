package main

import (
	"flag"
	"io"
	"log"
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

func trimText(str *string) {
	*str = strings.TrimSpace(*str)
}
