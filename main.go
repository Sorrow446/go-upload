package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"main/hosts/anonfiles"
	"main/hosts/catbox"
	"main/hosts/fileio"
	"main/hosts/filemail"
	"main/hosts/ftp"
	"main/hosts/gofile"
	"main/hosts/krakenfiles"
	"main/hosts/letsupload"
	"main/hosts/megaup"
	"main/hosts/mixdrop"
	"main/hosts/pixeldrain"
	"main/hosts/racaty"
	"main/hosts/transfersh"
	"main/hosts/uguu"
	"main/hosts/wetransfer"
	"main/hosts/workupload"
	"main/hosts/zippyshare"
	"main/utils"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/alexflint/go-arg"
	"github.com/dustin/go-humanize"
)

const megabyte = 1000000

var (
	funcMap = map[string]func(*utils.Args, string) (string, error){
		"anonfiles":   anonfiles.Run,
		"catbox":      catbox.Run,
		"fileio":      fileio.Run,
		"filemail":    filemail.Run,
		"ftp":         ftp.Run,
		"gofile":      gofile.Run,
		"krakenfiles": krakenfiles.Run,
		"letsupload":  letsupload.Run,
		"megaup":      megaup.Run,
		"mixdrop":     mixdrop.Run,
		"pixeldrain":  pixeldrain.Run,
		"racaty":      racaty.Run,
		"transfersh":  transfersh.Run,
		"uguu":        uguu.Run,
		"wetransfer":  wetransfer.Run,
		"zippyshare":  zippyshare.Run,
		"workupload":  workupload.Run,
	}
	templateEscPairs = []utils.TemplateEscPair{
		// Newline
		{From: []byte{'\x5C', '\x6E'}, To: []byte{'\x0A'}},
		// Tab
		{From: []byte{'\x5C', '\x74'}, To: []byte{'\x09'}},
	}
)

func populateDirs(path string) ([]string, error) {
	var paths []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if !f.IsDir() {
			filePath := filepath.Join(path, f.Name())
			paths = append(paths, filePath)
		}
	}
	return paths, nil
}

func populateDirsRec(srcPath string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	return dirs, err
}

func checkExists(path string, isDir bool) (bool, error) {
	f, err := os.Stat(path)
	if err == nil {
		if isDir {
			return f.IsDir(), nil
		} else {
			return !f.IsDir(), nil
		}
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func processDirs(args *utils.Args) error {
	var (
		allDirs  []string
		popPaths []string
	)
	for _, dir := range args.Directories {
		exists, err := checkExists(dir, true)
		if err != nil {
			return err
		}
		if exists {
			if !foldContains(allDirs, dir) {
				allDirs = append(allDirs, dir)
				if args.Recursive {
					popPaths, err = populateDirsRec(dir)
				} else {
					popPaths, err = populateDirs(dir)
				}
				if err != nil {
					return err
				}
				args.Files = append(args.Files, popPaths...)
			} else {
				fmt.Println("Filtered duplicate directory:", dir)
			}

		} else {
			fmt.Println("Filtered non-existent directory:", dir)
		}
	}
	return nil
}

func foldContains(arr []string, value string) bool {
	for _, item := range arr {
		if strings.EqualFold(item, value) {
			return true
		}
	}
	return false
}

func filterHosts(hosts []string) []string {
	var filteredHosts []string
	for _, host := range hosts {
		if !foldContains(filteredHosts, host) {
			filteredHosts = append(filteredHosts, host)
		}
	}
	return filteredHosts
}

func filterPaths(paths []string) ([]string, error) {
	var filteredPaths []string
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for _, path := range paths {
		if !filepath.IsAbs(path) {
			path = filepath.Join(wd, path)
		}
		exists, err := checkExists(path, false)
		if err != nil {
			return nil, err
		}
		if exists {
			if !foldContains(filteredPaths, path) {
				filteredPaths = append(filteredPaths, path)
			} else {
				fmt.Println("Filtered duplicate file:", path)
			}
		} else {
			fmt.Println("Filtered non-existent file:", path)
		}
	}
	return filteredPaths, nil
}

func parseArgs() (*utils.Args, error) {
	var args utils.Args
	arg.MustParse(&args)
	if args.SpeedLimit != -1 && args.SpeedLimit <= 0 {
		return nil, errors.New("Invalid speed limit.")
	}
	if len(args.Files) == 0 && len(args.Directories) == 0 {
		return nil, errors.New("File path and/or directory required.")
	}
	args.ByteLimit = int64(megabyte * args.SpeedLimit)
	if args.SpeedLimit != -1 {
		fmt.Printf("Upload speed limiting is active, limit: %s/s.\n",
			humanize.Bytes(uint64(args.ByteLimit)))
	}
	if len(args.Directories) > 0 {
		err := processDirs(&args)
		if err != nil {
			return nil, err
		}
	}
	paths, err := filterPaths(args.Files)
	if err != nil {
		errString := fmt.Sprintf("Failed to filter paths.\n%s", err)
		return nil, errors.New(errString)
	}
	if len(paths) == 0 {
		return nil, errors.New("All files were filtered.")
	}
	hosts := filterHosts(args.Hosts)
	args.Hosts = hosts
	args.Files = paths
	return &args, nil
}

func escapeTemplate(template []byte) []byte {
	var escaped []byte
	for i, pair := range templateEscPairs {
		if i != 0 {
			template = escaped
		}
		escaped = bytes.ReplaceAll(template, pair.From, pair.To)
	}
	return escaped
}

func parseTemplate(templateText string, meta map[string]string) []byte {
	var buffer bytes.Buffer
	for {
		err := template.Must(template.New("").Parse(templateText)).Execute(&buffer, meta)
		if err == nil {
			break
		}
		fmt.Println("Failed to parse template. Default will be used instead.")
		templateText = "# {{.filename}}\n{{.fileUrl}}\n"
		buffer.Reset()
	}
	return escapeTemplate(buffer.Bytes())
}

func writeTxt(path, filePath, fileUrl, templateText string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	meta := map[string]string{
		"filename": filepath.Base(filePath),
		"filePath": filePath,
		"fileUrl":  fileUrl,
	}
	parsed := parseTemplate(templateText, meta)
	_, err = f.Write(parsed)
	f.Close()
	return err
}

func outSetup(path string, wipe bool) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	if wipe {
		err = f.Truncate(0)
		if err != nil {
			return err
		}
	}
	return nil
}

func outSetupJob(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	jobs := &utils.UploadJobs{
		Jobs: []utils.UploadJob{},
	}

	m, err := json.MarshalIndent(&jobs, "", "\t")
	if err != nil {
		return err
	}
	_, err = f.Write(m)
	return err
}

func writeJob(jobPath, _url, host, filePath string, jobErr error) error {
	var (
		ok      = true
		errText string
	)
	if jobErr != nil {
		ok = false
		errText = jobErr.Error()
	}
	job := &utils.UploadJob{
		URL:       _url,
		Host:      host,
		Filename:  filepath.Base(filePath),
		FilePath:  filePath,
		Ok:        ok,
		ErrorText: errText,
	}
	data, err := ioutil.ReadFile(jobPath)
	if err != nil {
		return err
	}
	var jobs utils.UploadJobs
	err = json.Unmarshal(data, &jobs)
	if err != nil {
		return err
	}
	jobs.Jobs = append(jobs.Jobs, *job)
	m, err := json.MarshalIndent(&jobs, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(jobPath, m, 0755)
	return err
}

func main() {
	args, err := parseArgs()
	if err != nil {
		panic(err)
	}
	outPath := args.OutPath
	if outPath != "" {
		err := outSetup(outPath, args.Wipe)
		if err != nil {
			panic(err)
		}
	}
	if args.JobOutPath != "" {
		err := outSetupJob(args.JobOutPath)
		if err != nil {
			panic(err)
		}
	}
	for i, host := range args.Hosts {
		lowerHost := strings.ToLower(host)
		hostFunc, ok := funcMap[lowerHost]
		if !ok {
			fmt.Println("Invalid host:", host)
			continue
		}
		if i != 0 {
			fmt.Println("")
		}
		fmt.Println("--" + lowerHost + "--")
		pathTotal := len(args.Files)
		for num, path := range args.Files {
			fmt.Printf("File %d of %d:\n", num+1, pathTotal)
			fmt.Println(path)
			fileUrl, err := hostFunc(args, path)
			if args.JobOutPath != "" {
				jobErr := writeJob(args.JobOutPath, fileUrl, host, path, err)
				if jobErr != nil {
					// Intentional.
					panic(jobErr)
				}
			}
			if err != nil {
				fmt.Println("Upload failed.\n" + err.Error())
				continue
			}
			fmt.Println(fileUrl)
			if outPath != "" {
				err = writeTxt(outPath, path, fileUrl, args.Template)
				if err != nil {
					fmt.Println("Failed to write to output text file.\n" + err.Error())
				}
			}
		}
	}
}
