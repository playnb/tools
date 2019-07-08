package main

import (
	"os"
	"bufio"
	"io"
	"fmt"
	"strings"
	"flag"
	"regexp"
)

var InputFileName = `C:\ltp\EgretProjects\Boom\tools\msg\src\protocol.d.ts`
var OuputFileName = `C:\ltp\EgretProjects\Boom\tools\msg\src\msg.ts`

func readLine(fileName string, handler func(string)) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		//		line = strings.TrimSpace(line)
		handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			fmt.Println(err)
			return err
		}
	}
	f.Close()
	return nil
}

func main() {
	flag.StringVar(&InputFileName, "in", InputFileName, "输入")
	flag.StringVar(&OuputFileName, "out", OuputFileName, "输出")
	flag.Parse()
	_lines := ""
	handler := func(line string) {
		//TODO 以后用正则表达式匹配
		if strings.Contains(line, "constructor") {
		} else if strings.Contains(line, "public static") {
		} else if strings.Contains(line, "public toJSON") {
		} else {
			line = strings.Replace(line, "$protobuf", "protobuf", -1)
			line = strings.Replace(line, "export namespace", "namespace", -1)
			line = strings.Replace(line, "class", "export class", -1)

			{
				pat := `\bLong\b`
				re, _ := regexp.Compile(pat)
				line = re.ReplaceAllString(line, "protobuf.Long")
			}
			{
				pat := `\bmsg\.`
				re, _ := regexp.Compile(pat)
				line = re.ReplaceAllString(line, "")
			}

			_lines = _lines + line
		}
	}
	readLine(InputFileName, handler)

	f, err := os.OpenFile(OuputFileName, os.O_CREATE+os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("[OpenFile] " + err.Error())
	} else {
		fmt.Println("[In]  "+InputFileName)
		fmt.Println("[Out] "+OuputFileName)

		f.WriteString(_lines)
		f.Close()
	}
}
