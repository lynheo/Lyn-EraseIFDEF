package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
)

type ProcInfo struct {
	deleteDefine   string
	backupOriginal bool
}

func parseCommandLine() ([]string, ProcInfo, bool) {
	fileName := flag.String("file", "", "FileName")
	fileNameRegExp := flag.String("rfile", "", "FileName RegExp")
	runPath := flag.String("path", "", "Running Path")
	deleteDefine := flag.String("define", "", "Delete Define(Required)")
	backupOriginal := flag.Bool("backup", false, "Backup Original File")

	var fileList []string
	var procInfo ProcInfo

	flag.Parse()
	if flag.NFlag() == 0 {
		flag.Usage()
		return fileList, procInfo, false
	}

	if *deleteDefine == "" {
		flag.Usage()
		return fileList, procInfo, false
	}
	fmt.Println(*runPath)
	procInfo.backupOriginal = *backupOriginal
	procInfo.deleteDefine = *deleteDefine

	if *fileNameRegExp == "" {
		fileList = append(fileList, filepath.Join(*runPath, *fileName))
	} else {
		rx, result := regexp.Compile(*fileNameRegExp)

		if result != nil {
			rx = nil
			flag.Usage()
			fmt.Println("Regexp Syntex Error")
			return fileList, procInfo, false
		} else {
			fmt.Println("File scan from ", *runPath)
			allFiles, _ := ioutil.ReadDir(*runPath)
			fmt.Println("FileCount : ", len(allFiles))

			for _, f := range allFiles {
				if rx.MatchString(f.Name()) == true {
					fileList = append(fileList, filepath.Join(*runPath, f.Name()))
					fmt.Println("ScanFile : ", f.Name(), " Ok")
				} else {
					fmt.Println("ScanFile : ", f.Name(), " Skip")
				}
			}
		}
	}

	return fileList, procInfo, true
}

func main() {
	runtime.GOMAXPROCS(1)

	fileList, procInfo, isSuccess := parseCommandLine()
	if isSuccess == false {
		return
	}

	fmt.Println("Define : ", procInfo.deleteDefine)
	wg := new(sync.WaitGroup)

	for _, fileName := range fileList {
		wg.Add(1)
		go func(wg *sync.WaitGroup, procInfo ProcInfo, fileName string) {
			fmt.Println("Start : ", fileName)
			doProcessSource(fileName, procInfo)
			fmt.Println("End : ", fileName)
			wg.Done()
		}(wg, procInfo, fileName)
	}

	wg.Wait()
}
