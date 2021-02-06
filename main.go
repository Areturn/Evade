package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var opts struct {
	InputFile      string `short:"i" long:"input-file" description:"Input source file" required:"true" value-name:"File"`
	OutputDir      string `short:"o" long:"output-dir" description:"Output directory" value-name:"Dir" default:"output"`
	FileNamePrefix string `short:"p" long:"filename-prefix" description:"Output file name prefix" value-name:"Prefix" default-mask:"InputFileName"`
	Size           int    `short:"s" long:"size" description:"Shard size" default:"4096" value-name:"Byte"`
	DisableAppend  bool   `long:"disable-append" description:"Disable append mode"`
}

func InitFlag() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.Usage = fmt.Sprintf("<-i filename> [-o output] [-p prefix] [-s 4096] [--disable-append]")
	_, err := parser.ParseArgs(os.Args)
	if flags.WroteHelp(err) {
		os.Exit(0)
	} else if err != nil {
		parser.WriteHelp(os.Stdin)
		os.Exit(1)
	}
}

func CopyFile(dstFileName string, srcFileName string) (written int64, err error) {
	srcFile, err := os.Open(srcFileName)
	if err != nil {
		return
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dstFileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

func CreateDir(dir string) (dirPath string, err error) {
	var stat os.FileInfo
	var path string
	if !regexp.MustCompile(`^(/|[A-Z]:\\)`).MatchString(dir) {
		path, err = AsbPath()
		if err != nil {
			return
		}
		dir = path + "/" + dir
	}
	if stat, err = os.Stat(dir); err != nil && os.IsExist(err) {
		return
	} else if os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return
		}
	} else if !stat.IsDir() {
		err = fmt.Errorf("目录'%s'已存在,且是个文件,请检查.", dir)
		return
	}
	return dir, nil
}

func AsbPath() (path string, err error) {
	// 获取可执行文件相对于当前工作目录的相对路径
	path = filepath.Dir(os.Args[0])
	// 根据相对路径获取可执行文件的绝对路径
	path, err = filepath.Abs(path)
	return
}

func main() {
	var dirPath string
	var inputFile *os.File
	InitFlag()
	dirPath, err := CreateDir(opts.OutputDir)
	checkErr(err)

	inputFile, err = os.Open(opts.InputFile)
	checkErr(err)
	defer inputFile.Close()

	info, err := inputFile.Stat()
	checkErr(err)

	split := strings.Split(info.Name(), ".")
	var suffix, prefix string
	l := len(split)
	if l >= 2 {
		suffix = split[l-1]
		prefix = strings.Join(split[:l-1], ".")
	} else {
		prefix = split[0]
	}

	if opts.FileNamePrefix != "" {
		prefix = opts.FileNamePrefix
	}

	var buf = make([]byte, opts.Size)
	var size int
	var filename string
	for {
		n, err := inputFile.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			} else {
				checkErr(err)
			}
		}
		oldfilename := filename
		oldsize := size
		size += n
		if suffix != "" && !opts.DisableAppend {
			filename = fmt.Sprintf("%s/%s_%d.%s", dirPath, prefix, size, suffix)
		} else if suffix != "" && opts.DisableAppend {
			filename = fmt.Sprintf("%s/%s_%d_%d.%s", dirPath, prefix, oldsize, size, suffix)
		} else if suffix == "" && opts.DisableAppend {
			filename = fmt.Sprintf("%s/%s_%d_%d", dirPath, prefix, oldsize, size)
		} else {
			filename = fmt.Sprintf("%s/%s_%d", dirPath, prefix, size)
		}

		if oldfilename != "" && !opts.DisableAppend {
			_, err = CopyFile(filename, oldfilename)
			checkErr(err)
		}
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		checkErr(err)
		defer file.Close()
		_, err = file.Write(buf[:n])
		checkErr(err)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
