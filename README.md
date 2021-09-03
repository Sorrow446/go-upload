# go-upload
File uploader with support for multiple hosts and progress reporting written in Go.

## Usage
Upload single file to anonfiles:   
`go-ul_x64.exe anonfiles -f G:\file.bin`

Upload two files to anonfiles and catbox and write output template:   
`go-ul_x64.exe anonfiles catbox -f G:\file.bin G:\file2.bin -o urls.txt`

Upload a single file to FTP server to /x/y/ and overwrite it if it already exists.   
`go-ul_x64.exe ftp -f G:\file.bin -U ftp://myusername:mypassword@ftp.server.com:21/x/y/ -O`

```
Usage: go-ul_x64.exe [--outpath OUTPATH] [--wipe] --files FILES [--private] [--template TEMPLATE] [--overwrite] [--user USER] HOSTS [HOSTS ...]

Positional arguments:
  HOSTS                  Which hosts to upload to.

Options:
  --outpath OUTPATH, -o OUTPATH
                         Path of text file to write template to. It will be created if it doesn't already exist.
  --wipe, -w             Wipe output text file on startup.
  --files FILES, -f FILES
                         Paths of files to upload.
  --private, -P          *Make upload private.
  --template TEMPLATE, -t TEMPLATE
                         Output text file template. Vars: filename, filepath, fileUrl [default: # {{.filename}}\n{{.fileUrl}}]
  --overwrite, -O        *Overwrite file on host if it already exists.
  --user USER, -u USER   *User form for FTP. Folders will be created recursively if they don't already exist.
  --help, -h             display this help and exit
```
\* = Not supported for all hosts.

## Supported hosts
|Host|Argument|
| --- | --- |
|[anonfiles](https://anonfiles.com/)|anonfiles
|[Catbox](https://catbox.moe/)|catbox
|[file.io](https://www.file.io/)|fileio
|[Filemail](https://www.filemail.com/)|filemail
|FTP|ftp
|[Gofile](https://gofile.io/)|gofile
|[MegaUp](https://megaup.net/)|megaup
|[pixeldrain](https://pixeldrain.com/)|pixeldrain
|[Uguu](https://uguu.se/)|uguu
|[zippyshare](https://www.zippyshare.com/)|zippyshare

Host arguments are case insensitive.
