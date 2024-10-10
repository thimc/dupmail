# dupmail

dupmail scans a directory that (presumably) contains offline e-mail
files which are usually fetched by programs such as offlineimap,
mbsync or fdm. The program deems a mail a duplicate by comparing
the Message-ID header and in case there isn't one the mail is simply
skipped.

On exit dupmail will print the message id and file paths of all the
duplicates.

## Installation

    # go build -o /usr/local/bin/dupmail main.go

## Usage

`dupmail [-d maildir] [-n]`

* `-d`	maildir to be scanned
* `-n`	scan the maildir but do not delete any files
