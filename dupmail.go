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
	verboseFlag     = flag.Bool("v", false, "increase verbosity")
	exitOnErrorFlag = flag.Bool("e", false, "exit once a mail without a message-id has been found")

	mails          int
	duplicates     map[string][]string
	missingHeaders []string
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

	var messageHeader string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		if len(line) < 2 {
			continue
		}
		messageHeader = strings.Join(line[1:], " ")
		if strings.ToLower(line[0]) == MessageHeader {
			duplicates[messageHeader] = append(duplicates[messageHeader], path)
			mails++
			return nil
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if *verboseFlag {
		if !*summarizeFlag {
			fmt.Fprintf(os.Stderr, "warning: missing message-id header %s\n", path)
		}
		missingHeaders = append(missingHeaders, path)
	}
	if *exitOnErrorFlag {
		return fmt.Errorf("Missing header in %s", path)
	}
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
		fmt.Fprintf(os.Stderr, "%d mail(s)\n", mails)
		fmt.Fprintf(os.Stderr, "%d duplicate(s)\n", mails-len(duplicates))
		if *verboseFlag {
			fmt.Fprintf(os.Stderr, "%d mail(s) are missing the \"%s\" header\n", len(missingHeaders), MessageHeader[:len(MessageHeader)-1])
		}
	}

	for msgid, dups := range duplicates {
		if len(dups) < 2 {
			continue
		}
		for n, dup := range dups {
			if !*verboseFlag && n < 1 {
				continue
			}
			if !*summarizeFlag {
				if *verboseFlag {
					fmt.Printf("%s\t", msgid)
				}
				fmt.Println(dup)
			}
			if !*noOperationFlag {
				if err := os.Remove(dup); err != nil {
					panic(err)
				}
			}
		}
	}

	if len(duplicates) > 1 {
		os.Exit(1)
	}
}
