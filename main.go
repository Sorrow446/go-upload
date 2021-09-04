// Filter duplicate host args.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"main/hosts/anonfiles"
	"main/hosts/catbox"
	"main/hosts/fileio"
	"main/hosts/filemail"
	"main/hosts/ftp"
	"main/hosts/gofile"
	"main/hosts/megaup"
	"main/hosts/pixeldrain"
	"main/hosts/uguu"
	"main/hosts/zippyshare"
	"main/utils"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/alexflint/go-arg"
)

var funcMap = map[string]func(*utils.Args, string) (string, error){
	"anonfiles":  anonfiles.Run,
	"catbox":     catbox.Run,
	"fileio":     fileio.Run,
	"ftp":        ftp.Run,
	"gofile":     gofile.Run,
	"pixeldrain": pixeldrain.Run,
	"uguu":       uguu.Run,
	"zippyshare": zippyshare.Run,
	"megaup":     megaup.Run,
	"filemail":   filemail.Run,
}

func parseArgs() (*utils.Args, error) {
	var args utils.Args
	arg.MustParse(&args)
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
	for _, path := range paths {
		exists, err := fileExists(path)
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

func fileExists(path string) (bool, error) {
	f, err := os.Stat(path)
	if err == nil {
		return !f.IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func parseTemplate(templateText string, meta map[string]string) string {
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
	return buffer.String()
}

func writeTxt(path, filename, filepath, fileUrl, templateText string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	meta := map[string]string{
		"filename": filename,
		"filepath": filepath,
		"fileUrl":  fileUrl,
	}
	parsed := parseTemplate(templateText, meta)
	_, err = f.Write([]byte(parsed))
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
			filename := filepath.Base(path)
			fmt.Println(filename)
			fileUrl, err := hostFunc(args, path)
			if err != nil {
				fmt.Println("Upload failed.\n", err)
				continue
			}
			fmt.Println(fileUrl)
			if outPath != "" {
				err = writeTxt(outPath, filename, path, fileUrl, args.Template)
				if err != nil {
					fmt.Println("Failed to write to output text file.\n", err)
				}
			}
		}
	}
}
