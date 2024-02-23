# dupmail

dupmail scans a directory that (presumably) contains offline e-mail
files. The program deems a mail a duplicate by comparing the Message-ID header
and in case there isn't one the mail is simply skipped.

These mails are typically fetched by a program such as OfflineIMAP,
mbsync or fdm and should work nicely with any of them.

On exit dupmail will print the file paths of all the duplicates.

## Installation

    $ go build -o dupmail dupmail.go
    # cp dupmail /usr/local/bin/

## Usage

`dupmail [-d maildir] [-v] [-n] [-s]`

* `-d`	maildir to be scanned
* `-n`	scan the maildir but do not delete any files
* `-s`	print a summary instead of the file paths of the duplicates
* `-v`  increase verbosity. Dupmail will print out the duplicates
in a very script friendly "\<Message-ID> \<TAB> \<File Path>" format.
