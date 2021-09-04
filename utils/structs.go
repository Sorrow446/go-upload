package utils

type Args struct {
	Hosts     []string `arg:"positional, required" help:"Which hosts to upload to."`
	OutPath   string   `arg:"-o" help:"Path of text file to write template to. It will be created if it doesn't already exist."`
	Wipe      bool     `arg:"-w" help:"Wipe output text file on startup."`
	Files     []string `arg:"-f, required" help:"Paths of files to upload."`
	Private   bool     `arg:"-P" help:"*Set upload as private."`
	Template  string   `arg:"-t" default:"# {{.filename}}\n{{.fileUrl}}\n" help:"Output text file template. Vars: filename, filepath, fileUrl"`
	Overwrite bool     `arg:"-O" help:"*Overwrite file on host if it already exists."`
	User      string   `arg:"-U" help:"*User form for FTP. Folders will be created recursively if they don't already exist."`
}

type myTransport struct{}

type WriteCounter struct {
	Total      int64
	TotalStr   string
	Uploaded   int64
	Percentage int
}
