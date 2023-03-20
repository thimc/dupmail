package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type MailFile struct {
	filepath, msgid string
}

var (
	MessageHeader   = "Message-ID:"
	re              *regexp.Regexp
	mails           []MailFile
	duplicates      int
	directoryFlag   = flag.String("d", "", "directory to be scanned")
	verboseFlag     = flag.Bool("v", false, "increase verbosity")
	keepFilesFlag   = flag.Bool("n", false, "don't remove any files, just display the statistics")
	exitOnErrorFlag = flag.Bool("e", false, "exit on error or if missing mail header")
)

func init() {
	re = regexp.MustCompile(`<([^>]+)>`)
}

func duplicateMail(messageId string) (MailFile, bool) {
	for _, k := range mails {
		if k.msgid == messageId {
			return k, true
		}
	}
	return MailFile{}, false
}

func handleDirectory(dirPath string) {
	err := filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}

			f, err := os.Open(path)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			messageId := ""

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				match := strings.Contains(scanner.Text(), MessageHeader)
				if match {
					messageId = fmt.Sprint(re.FindAllString(scanner.Text(), 1))
				}
			}
			// TODO: The way we determine if the file has a Message-ID header
			// is to check whether it's length is greater than 3 bytes.
			// In my case I have had mails where there's a field for the ID
			// itself but there isn't actually a *real* value there.
			if len(messageId) < 4 {
				if *exitOnErrorFlag {
					log.Printf("Missing header: %v\n", path)
					os.Exit(1)
				}
				return nil
			}
			m := MailFile{filepath: path, msgid: messageId}
			if err := scanner.Err(); err != nil {
				panic(err)
			}

			if match, duplicate := duplicateMail(m.msgid); !duplicate {
				mails = append(mails, m)
				return nil
			} else {
				duplicates++
				if *verboseFlag {
					log.Printf("Duplicate header: %v %v", match.msgid, match.filepath)
				}
				if *keepFilesFlag == false {
					err = os.Remove(m.filepath)
					if err != nil {
						panic(err)
					}
				}
			}

			return nil
		})
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	if _, err := os.Stat(*directoryFlag); err != nil {
		panic(err)
	}
	handleDirectory(*directoryFlag)
	fmt.Printf("Found %d mails (%d duplicates)\n", len(mails), duplicates)
}
