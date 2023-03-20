# dupmail

dupmail scans a directory that (presumably) contains offline e-mail files.\
The program assumes that there's a Message-ID header stored in each file.

These mails are typically fetched by a program such as OfflineIMAP, mbsync\
or fdm and should work nicely with any of them.

dupmail is made with a mindset that it should be ran as a cronjob.

On exit it summarizes how many mails it iterated over and how many\
were duplicates.

## Usage
`dupmail [-d directory] [-v] [-n] [-e]`
* `-d`	specifies the directory where the mails are stored
* `-v`	verbose output
* `-n`	don't remove any files
* `-e`	exit on error or if the mail doesn't contain a Message-ID header

