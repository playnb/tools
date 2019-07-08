package main

import (
	"cell/common/mustang/util"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var _declare_template = `
message $PROTO_NAME$_Cell{
	optional uint64 UserID = 1;
	optional $PROTO_NAME$ CmdData = 2;
}

`

var debug = false
func main() {
	path := ""
	flag.StringVar(&path, "path", "", "path")
	flag.BoolVar(&debug, "debug", false, "path")
	flag.Parse()
	//base_dir = `C:\project_0\`

	_import_files := ""
	_declare_str := ""
	util.GetFilelist(path, func(fileName string) {
		if strings.HasSuffix(fileName, ".proto") && !strings.HasSuffix(fileName, "_inner.proto") {
			has_msg := false
			util.ReadFileLine(fileName, func(line string) {
				if strings.Contains(line, "message") && strings.Contains(line, "MSG_") {
					line = strings.Trim(line, "{}")
					ss := strings.Split(line, " ")
					for _, str := range ss {
						if strings.Contains(str, "MSG_") {
							has_msg = true
							_declare_str = _declare_str + strings.Replace(_declare_template, "$PROTO_NAME$", str, -1)
						}
					}
				}
			})
			if has_msg {
				_import_files = _import_files + `import "` + filepath.Base(fileName) + `";`
				_import_files = _import_files + "\n"
			}
		}
	})

	proto := `syntax = "proto2";
	
package msg;
	
` + _import_files + `
` + _declare_str

	proto_path := path + `\forward_inner.proto`
	if debug {
		fmt.Println("生成文件: " + proto_path)
	}
	f, err := os.OpenFile(proto_path, os.O_CREATE+os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("[OpenFile] " + err.Error())
	}
	f.WriteString(proto)
}
