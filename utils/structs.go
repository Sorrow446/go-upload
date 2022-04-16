package utils

type Args struct {
	Hosts       []string `arg:"positional, required" help:"Which hosts to upload to."`
	OutPath     string   `arg:"-o" help:"Path of text file to write template to. It will be created if it doesn't already exist."`
	Wipe        bool     `arg:"-w" help:"Wipe output text file on startup."`
	Files       []string `arg:"-f, help:"Paths of files to upload."`
	Private     bool     `arg:"-P" help:"*Set upload as private."`
	Template    string   `arg:"-t" default:"# {{.filename}}\n{{.fileUrl}}\n" help:"Output text file template. Vars: filename, filepath, fileUrl"`
	Overwrite   bool     `arg:"-O" help:"*Overwrite file on host if it already exists."`
	User        string   `arg:"-u" help:"*User form for FTP. Folders will be created recursively if they don't already exist."`
	Directories []string `arg:"-d, help:"Paths of folders to upload."`
	Recursive   bool     `arg:"-r" help:"Include subdirectories."`
	SpeedLimit  float64  `arg:"-l" default:"-1" help:"Upload speed limit in megabytes. Example: 0.5 = 500 kB/s, 1 = 1 MB/s, 1.5 = 1.5 MB/s."`
	ByteLimit   int64    `arg:"-"`
}

type myTransport struct{}

type WriteCounter struct {
	Total      int64
	TotalStr   string
	Uploaded   int64
	Percentage int
	StartTime  int64
}

type TemplateEscPair struct {
	From []byte
	To   []byte
}
