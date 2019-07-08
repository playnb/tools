package main

import (
	"bufio"
	//	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	//	"os/exec"
	"path/filepath"
	//	"strconv"
	"strings"
)

var _msg_index = 1
var _go_code = `package protocol
			
import (
	"cell/common/protocol/msg"
	"reflect"
)

const (
	$CONST$
	Index_MAX
)

type MsgIndex struct {
	Index   uint16
	MsgType reflect.Type
}

var msg_name_index map[string]*MsgIndex

func init() {
	msg_name_index = make(map[string]*MsgIndex)

	//初始化部分
	$INIT$
}`

var _lua_code = `

--服务器与客户端协议号对应表;
PROTO_MSG = {
	$LUA_DEF$
}`

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
			fmt.Println("Error: " + err.Error())
			return err
		}
	}
	f.Close()
	return nil
}

func getFilelist(path string) {
	err := filepath.Walk(path, func(fileName string, finfo os.FileInfo, err error) error {
		if finfo == nil {
			return err
		}
		if finfo.IsDir() {
			return nil
		}

		if strings.HasSuffix(fileName, ".pb.go") {

			_lines := ""
			handler := func(line string) {
				//TODO 以后用正则表达式匹配
				process := true
				for process {
					process = false
					if strings.Contains(line, "import") && strings.Contains(line, "msg1") {
						line = ""
					} else if strings.Contains(line, "import") && strings.Contains(line, "cell/common/protocol/.") {
						line = ""
					} else if strings.Contains(line, "msg1.") {
						line = strings.Replace(line, "msg1.", "", -1)
						process = true
					} else if strings.Contains(line, "import") && strings.Contains(line, "msg2") {
						line = ""
					} else if strings.Contains(line, "msg2.") {
						line = strings.Replace(line, "msg2.", "", -1)
						process = true
					} else if strings.Contains(line, "import") && strings.Contains(line, "msg3") {
						line = ""
					} else if strings.Contains(line, "msg3.") {
						line = strings.Replace(line, "msg3.", "", -1)
						process = true
					} else if strings.Contains(line, "import") && strings.Contains(line, "msg4") {
						line = ""
					} else if strings.Contains(line, "msg4.") {
						line = strings.Replace(line, "msg4.", "", -1)
						process = true
					} else if strings.Contains(line, "import") && strings.Contains(line, "msg5") {
						line = ""
					} else if strings.Contains(line, "msg5.") {
						line = strings.Replace(line, "msg5.", "", -1)
						process = true
					} else if strings.Contains(line, "import") && strings.Contains(line, "msg6") {
						line = ""
					} else if strings.Contains(line, "msg6.") {
						line = strings.Replace(line, "msg6.", "", -1)
						process = true
					} else if strings.Contains(line, "fileDescriptor0") {
						line = strings.Replace(line, "fileDescriptor0", "fileDescriptor"+strings.Split(finfo.Name(), ".")[0], -1)
					}
				}
				_lines = _lines + line
			}
			readLine(fileName, handler)

			//fmt.Println(_lines)

			//fmt.Println("处理文件:" + fileName)
			//err := os.Remove(fileName)
			//if err != nil {
			//	fmt.Println("[Remove] ", err.Error())
			//	return err
			//}
			f, err := os.OpenFile(fileName, os.O_CREATE+os.O_TRUNC, 0666)
			if err != nil {
				fmt.Println("[OpenFile] " + err.Error())
				return err
			}
			f.WriteString(_lines)
			f.Close()
		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

func main() {
	path := ""
	flag.StringVar(&path, "path", `C:\code\`, "项目的根目录")
	flag.Parse()

	getFilelist(path)

}
