# A naive reimplementation of the ls command line tool

This is just a fun learning project. Don't use this code for anything serious

## Features

 - Shows all files in a directory, either the current working directory or a specified one.
 - Alphabetic ordering
 - Show details (file size): -l
 - Sort by file size: -S
 - Show sizes in human readable format: -h
 - Reverse sort order: -r

## Usage examples

### List current directory

    ls2

### List a specified directory

    ls2 /etc
    
### Show details and sort by size

    ls2 -l -S
