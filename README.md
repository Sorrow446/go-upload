# go-upload
File uploader with support for multiple hosts and progress reporting written in Go. Large file-friendly.
![](https://i.imgur.com/Mtfn3pu.png)
![](https://i.imgur.com/D2a7NLu.png)   
[Windows, Linux, macOS and Android binaries](https://github.com/Sorrow446/go-upload/releases)

## Usage
Upload single file to anonfiles:   
`go-ul_x64.exe anonfiles -f G:\file.bin`

Upload two files to anonfiles and catbox and write output template:   
`go-ul_x64.exe anonfiles catbox -f G:\file.bin G:\file2.bin -o urls.txt`

Upload all files in `G:\stuff` to zippyshare recursively with a 500 kB/s limit and write output template:   
`go-ul_x64.exe zippyshare -d G:\stuff -r -o urls.txt -l 0.5`

Upload a single file to FTP server to /x/y/ and overwrite it if it already exists.   
`go-ul_x64.exe ftp -f G:\file.bin -U ftp://myusername:mypassword@ftp.server.com:21/x/y/ -O`

```
Usage: go-ul_x64.exe  [--outpath OUTPATH] [--wipe] [--files FILES] [--private] [--template TEMPLATE] [--overwrite] [--user USER] [--directories DIRECTORIES] [--recursive] [--speedlimit SPEEDLIMIT] HOSTS [HOSTS ...]

Positional arguments:
  HOSTS                  Which hosts to upload to.

Options:
  --outpath OUTPATH, -o OUTPATH
                         Path of text file to write template to. It will be created if it doesn't already exist.
  --wipe, -w             Wipe output text file on startup.
  --files FILES, -f FILES
  --private, -P          *Set upload as private.
  --template TEMPLATE, -t TEMPLATE
                         Output text file template. Vars: filename, filepath, fileUrl [default: # {{.filename}}\n{{.fileUrl}}\n]
  --overwrite, -O        *Overwrite file on host if it already exists.
  --user USER, -u USER   *User form for FTP. Folders will be created recursively if they don't already exist.
  --directories DIRECTORIES, -d DIRECTORIES
  --recursive, -r        Include subdirectories.
  --speedlimit SPEEDLIMIT, -l SPEEDLIMIT
                         Upload speed limit in megabytes. Example: 0.5 = 500 kB/s, 1 = 1 MB/s, 1.5 = 1.5 MB/s. [default: -1]
  --help, -h             display this help and exit
```
\* = Not supported for all hosts.

### Template

Default: `# {{.filename}}\n{{.fileUrl}}\n`    
Output with the default template:
```
# 2.jpg
https://anonfiles.com/Hde2H4F5ue/2_jpg
```
Vars: filename, filePath, fileUrl

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
|[WeTransfer](https://wetransfer.com/)|Wetransfer
|[zippyshare](https://www.zippyshare.com/)|zippyshare

Host arguments are case insensitive.
