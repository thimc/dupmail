# dupmail

dupmail scans a directory that (presumably) contains offline e-mail
files. The program assumes that there's a Message-ID header stored
in each file.

These mails are typically fetched by a program such as OfflineIMAP,
mbsync or fdm and should work nicely with any of them.

On exit it summarizes how many mails it iterated over and how many
were duplicates.

## Installation

    $ go build -o dupmail main.go
    # cp dupmail /usr/local/bin/

## Usage

`dupmail [-d maildir] [-v] [-n] [-s]`

* `-d`	maildir to be scanned
* `-n`	scan the maildir but do not delete any files
* `-s`	print a summary instead of the file paths of the duplicates
