package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	MessageHeader = "message-id:"
)

var (
	directoryFlag   = flag.String("d", "", "maildir to be scanned")
	noOperationFlag = flag.Bool("n", false, "scan the maildir but do not delete any files")
	summarizeFlag   = flag.Bool("s", false, "print a summary instead of the file paths of the duplicates")

	mails      int
	duplicates map[string][]string
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage for %s\n", filepath.Base(os.Args[0]))
	flag.PrintDefaults()
}

func traverse(path string, info fs.FileInfo, err error) error {
	if err != nil || info.IsDir() {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		if len(line) < 2 {
			continue
		}
		if strings.ToLower(line[0]) == MessageHeader {
			duplicates[line[1]] = append(duplicates[line[1]], path)
			mails++
			return nil
		}
	}
	// return fmt.Errorf("Missing header in %s", path)
	return nil
}

func main() {
	flag.Parse()

	_, err := os.Stat(*directoryFlag)
	if err != nil {
		if len(flag.Args()) > 0 {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		usage()
		os.Exit(1)
	}

	duplicates = make(map[string][]string)

	if err := filepath.Walk(*directoryFlag, traverse); err != nil {
		panic(err)
	}

	if *summarizeFlag {
		fmt.Printf("%d mail(s)\n", mails)
		fmt.Printf("%d duplicate(s)\n", mails-len(duplicates))
	}

	for _, dups := range duplicates {
		if len(dups) < 2 {
			continue
		}
		for j, dup := range dups {
			if j < 1 {
				continue
			}
			if !*summarizeFlag {
				fmt.Println(dup)
			}
			if !*noOperationFlag {
				if err := os.Remove(dup); err != nil {
					panic(err)
				}
			}
		}
	}
}
