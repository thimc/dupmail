package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	dirpath = flag.String("d", "", "maildir to be scanned")
	noop    = flag.Bool("n", false, "scan the maildir but do not delete any files")

	messageHeader = "message-id:"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage for %s\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	var (
		mail        = make(map[string]bool)
		start       = time.Now()
		mails, dups int
	)
	if *dirpath == "" {
		flag.Usage()
		os.Exit(1)
	}
	if _, err := os.Stat(*dirpath); err != nil {
		fmt.Fprintf(os.Stderr, "stat %q: %s\n", *dirpath, err)
		os.Exit(1)
	}
	if err := filepath.Walk(*dirpath, func(path string, info fs.FileInfo, err error) error {
		if !info.Mode().IsRegular() {
			return err
		}
		mails++
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		r := bufio.NewReader(f)
		for {
			ln, err := r.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if !strings.HasPrefix(strings.ToLower(ln), messageHeader) {
				continue
			}
			var (
				parts = strings.Split(ln[:len(ln)-1], " ")
				msgid = strings.Trim(strings.Join(parts[1:], ""), "<>")
			)
			if !mail[msgid] {
				mail[msgid] = true
				continue
			}
			dups++
			fmt.Printf("%s\t%s\n", path, msgid)
			if *noop {
				continue
			}
			if err := os.Remove(path); err != nil {
				return err
			}
			return nil
		}
		return nil
	}); err != nil {
		fmt.Fprintf(os.Stderr, "walk %q: %s\n", *dirpath, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "processed %d mail(s) in %.2fs. ", mails, time.Since(start).Seconds())
	if dups > 0 {
		fmt.Fprintf(os.Stderr, "%d duplicate(s) found.", dups)
	}
	fmt.Fprintln(os.Stderr)
}
