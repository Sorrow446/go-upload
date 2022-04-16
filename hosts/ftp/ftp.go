package ftp

import (
	"errors"
	"fmt"
	"io"
	"main/utils"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/jlaffaye/ftp"
	_ftp "github.com/jlaffaye/ftp"
)

func auth(host, username, password string) (*_ftp.ServerConn, error) {
	c, err := _ftp.Dial(host, _ftp.DialWithTimeout(time.Second*30))
	if err != nil {
		return nil, err
	}
	err = c.Login(username, password)
	if err != nil {
		c.Quit()
		return nil, err
	}
	return c, nil
}

func fileExists(c *_ftp.ServerConn, path, filename string) (bool, error) {
	entries, err := c.List(path)
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		if entry.Type == _ftp.EntryTypeFile && entry.Name == filename {
			return true, nil
		}
	}
	return false, nil
}

func dirExists(c *ftp.ServerConn, directory string) (bool, error) {
	entries, err := c.List("")
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		if entry.Type == ftp.EntryTypeFolder && entry.Name == directory {
			return true, nil
		}
	}
	return false, nil
}

func upload(c *_ftp.ServerConn, path, filename string) error {
	defer fmt.Println("")
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	size := stat.Size()
	counter := &utils.WriteCounter{
		Total:     size,
		TotalStr:  humanize.Bytes(uint64(size)),
		StartTime: time.Now().UnixMilli(),
	}
	err = c.Stor(filename, io.TeeReader(f, counter))
	return err
}

func parseUrl(userString string) (*User, error) {
	if !strings.HasPrefix(userString, "ftp://") {
		userString = "ftp://" + userString
	}
	u, err := url.Parse(userString)
	if err != nil {
		return nil, err
	}
	host := u.Host
	path := u.Path
	userInfo := u.User
	username := userInfo.Username()
	password, _ := userInfo.Password()
	if host == "" {
		return nil, errors.New("Host required.")
	} else if username == "" {
		return nil, errors.New("Username required.")
	} else if password == "" {
		return nil, errors.New("Password required.")
	}
	if path == "/" {
		path = ""
	}
	user := &User{
		Host:     host,
		Username: username,
		Password: password,
		Path:     path,
	}
	return user, nil
}

func makeDirRecur(c *ftp.ServerConn, path string) error {
	path = strings.Trim(path, "/")
	splitPath := strings.Split(path, "/")
	for _, dir := range splitPath {
		exists, err := dirExists(c, dir)
		if err != nil {
			return err
		}
		if !exists {
			err := c.MakeDir(dir)
			if err != nil {
				return err
			}
		}
		err = c.ChangeDir(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

func Run(args *utils.Args, path string) (string, error) {
	filename := filepath.Base(path)
	if args.User == "" {
		return "", errors.New("User required (host, port, username and password).")
	}
	u, err := parseUrl(args.User)
	if err != nil {
		return "", err
	}
	c, err := auth(u.Host, u.Username, u.Password)
	if err != nil {
		return "", err
	}
	defer c.Quit()
	if u.Path != "" {
		err = makeDirRecur(c, u.Path)
		if err != nil {
			return "", err
		}
	}
	if !args.Overwrite {
		exists, err := fileExists(c, u.Path, filename)
		if err != nil {
			return "", err
		}
		if exists {
			return "", errors.New("File already exists on FTP. Use the -O flag to overwrite.")
		}
	}
	err = upload(c, path, filename)
	if err != nil {
		return "", err
	}
	outPath := u.Path + "/" + filename
	return outPath, nil
}
