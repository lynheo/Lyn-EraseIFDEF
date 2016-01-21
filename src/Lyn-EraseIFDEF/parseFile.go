package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

const (
	Off = iota
	Use
	Remove
)

func readFileLines(fileName string) []string {
	var lines []string

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer file.Close()
	r := bufio.NewReader(file)

	for {
		s, _, e := r.ReadLine()
		if e != nil {
			break
		}

		lines = append(lines, string(s))
	}
	return lines
}

func copyFile(original, dest string) {
	r, err := os.Open(original)
	if err != nil {
		panic(err)
	}

	defer r.Close()

	w, err := os.Create(dest)
	if err != nil {
		panic(err)
	}

	defer w.Close()

	_, err = io.Copy(w, r)
	if err != nil {
		panic(err)
	}
}

func doProcessSource(fileName string, procInfo ProcInfo) {
	if procInfo.backupOriginal == true {
		backupFileName := fileName + ".bak"
		copyFile(fileName, backupFileName)
	}

	lines := readFileLines(fileName)

	status := Off
	var depthstack DepthStack

	output := make([]string, 0, len(lines))
	for _, line := range lines {
		originalLine := strings.TrimRightFunc(line, unicode.IsSpace)
		line = strings.TrimSpace(originalLine)

		if len(line) == 0 {
			output = append(output, line)
			continue
		}

		if line[0] == '#' {
			var preprocessor string
			var variable string
			fmt.Sscanf(line, "%s %s", &preprocessor, &variable)

			preprocessor = strings.TrimSpace(preprocessor)
			variable = strings.TrimSpace(variable)

			switch preprocessor {
			case "#ifdef":
				depthstack.Push(variable)
				if variable == procInfo.deleteDefine {
					status = Use
					continue
				}
			case "#ifndef":
				depthstack.Push(variable)
				if variable == procInfo.deleteDefine {
					status = Remove
					continue
				}
			case "#if":
				depthstack.Push(variable)
			case "#else":
				if depthstack.Top() == procInfo.deleteDefine {
					if status == Use {
						status = Remove
						continue
					} else if status == Remove {
						status = Use
					}
				}
			case "#endif":
				topValue := depthstack.Top()
				depthstack.Pop()
				if topValue == procInfo.deleteDefine {
					status = Off
					continue
				}
			}
		}

		switch status {
		case Use, Off:
			output = append(output, originalLine)
		}
	}

	file, err := os.Create(fileName)
	if err == nil {
		defer file.Close()
		for _, line := range output {
			file.WriteString(line)
			file.WriteString("\n")
		}
	} else {
		panic(err)
	}
}
