package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var _msg_user_index = 1
var _msg_inner_index = 30001

var _lua_code = `

--服务器与客户端协议号对应表;
PROTO_MSG = {
	$LUA_DEF$
}`

var _template_ts_code = `
//服务器与客户端协议号对应表
class Protocol {
    //public static Root: protobuf.Root = new protobuf.Root()
	//public static All = {}
    //public static MSG_LoginRequest = new ProtocolType()
	$TS_DEF$

    public static LoadAllProtobuf(handleMgr: core.IMessageHandleMgr) {
    	//handleMgr.RegisterEncodeDecode(Protocol.MSG_RoomUserInfo_Index, msg.MSG_RoomUserInfo.decode, msg.MSG_RoomUserInfo.encode)
        //Protocol.MSG_LoginRequest.Index = 1
        //Protocol.MSG_LoginRequest.Name = "msg.MSG_LoginRequest"
        //Protocol.MSG_LoginRequest.Proto = Protocol.Root.lookupType(Protocol.MSG_LoginRequest.Name)
		$TS_INDEX$
    }
}
`
var _template_ts_def = `	public static $MSG_NAME$_Index = $MSG_INDEX$
`
var _template_ts_index = `
	handleMgr.RegisterEncodeDecode(Protocol.$MSG_NAME$_Index, msg.$MSG_NAME$.decode, msg.$MSG_NAME$.encode)
`

var _template_cs_code = `
using System;
using System.Collections;
using System.Collections.Generic;
using System.Runtime.Serialization.Formatters.Binary;
using System.IO;

public class MsgIndex
{$CS_Index$
    public static bool Handle(int index, byte[] data, int offset, int length)
    {
        return true;
    }
	public static void InitStatic()
    {
    }
}
`

/*
var _template_cs_code = `
using System;
using System.Collections;
using System.Collections.Generic;
using Msg;
using System.Runtime.Serialization.Formatters.Binary;
using System.IO;

using pb = Google.Protobuf;

//服务器与客户端协议号对应表
public class MsgIndex
{$CS_Index$

    static MsgIndex()
    {
    }

    internal static class GetValueImpl
    {
        private static int DefImpl<T>() => default(int);
        private static int IntRet() => int.MaxValue;

        internal static class Specializer<T>
        {
            internal static Func<int> Fun;
            internal static int Call() => null != Fun ? Fun() : DefImpl<T>();
        }

        static GetValueImpl()
        {
        }
    }

    public static void InitStatic()
    {$CS_Def$
    }

    public delegate void HandleMessage<T>(T m);

    internal static class HandleImpl
    {
        internal static class Specializer<T>
        {
            internal static HandleMessage<T> handles;
            internal static Func<T, bool> Fun;
            internal static bool Handle(T m)
            {
                if (handles != null)
                {
                    handles.Invoke(m);
                    return true;
                }
                return false;
            }
        }
    }

    public static int GetValue<T>(T t) => GetValueImpl.Specializer<T>.Call();

    public static bool Handle<T>(T m) 
    {
        return HandleImpl.Specializer<T>.Handle(m);
    }

    public static void AddHandle<T>(HandleMessage<T> handle) 
    {
        HandleImpl.Specializer<T>.handles += handle;
    }

    public static bool Handle(int index, byte[] data, int offset, int length)
    {
        switch (index)
        {$CS_Handle$
            default:
                return false;
        }
    }
}
`
*/

/*
	private static _$MSG_NAME$:ProtocolType<msg.$MSG_NAME$> = null
	public static get $MSG_NAME$(): ProtocolType<msg.$MSG_NAME$> {
		if (Protocol._$MSG_NAME$ == null) {
			Protocol._$MSG_NAME$ = new ProtocolType<msg.$MSG_NAME$>()
			Protocol._$MSG_NAME$.Index = $MSG_INDEX$
			Protocol._$MSG_NAME$.Name = "msg.$MSG_NAME$"
			Protocol._$MSG_NAME$.Proto = Protocol.Root.lookupType(Protocol._$MSG_NAME$.Name)
			Protocol.All[Protocol._$MSG_NAME$.Index] = Protocol._$MSG_NAME$
		}
		return Protocol._$MSG_NAME$
	}
*/
/*
Protocol.$MSG_NAME$.Index
*/

var _msg_name_index = make(map[string]int)

func readLine(fileName string, handler func(string)) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

type msgIndexDef struct {
	Package  string
	MsgName  string
	MsgIndex int
}

type MsgIndex map[int]*msgIndexDef

func (mi *MsgIndex) Get(index int) *msgIndexDef {
	if md, ok := (*mi)[index]; ok {
		return md
	}
	return nil
}

var index = make(MsgIndex)

func getFilelist(path string) bool {
	/*
	_const_str := ""
	_init_str := ""
	_know_msg_str := ""
	_lua_def := ""
	_ts_def := ""
	_ts_index := ""
	_cs_index := ""
	_cs_def := ""
	_cs_handle := ""
*/
	_gen_code := false

	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".proto") {
			//_data.proto表示数据类型定义，里面的东西都不做Index
			if strings.HasSuffix(path, "_data.proto") {
				return nil
			}
			//fmt.Println("Process: "+path)
			_inner := strings.HasSuffix(path, "_inner.proto")
			_package := ""

			handler := func(line string) {
				//TODO 以后用正则表达式匹配
				if strings.Contains(line, "package") {
					_package = strings.Trim(line, "package")
					_package = strings.Trim(_package, " ;")
					//fmt.Println(_package)
				}
				if strings.Contains(line, "message") {
					if strings.Contains(line, "MSG_") {
						_msg_name := strings.Replace(line, "message", "", -1)
						_msg_name = strings.Trim(_msg_name, " ;{")

						_msg_index := 0
						_find_index := false

						_msg_index, _find_index = _msg_name_index[_msg_name]
						if _find_index == false {
							if _inner {
								_msg_index = _msg_inner_index
							} else {
								_msg_index = _msg_user_index
							}
						}

						index[_msg_index] = &msgIndexDef{
							Package:  _package,
							MsgIndex: _msg_index,
							MsgName:  _msg_name,
						}

						/*
						str := "Index_$PACK$_$MSG_NAME$ = $MSG_INDEX$"
						str = strings.Replace(str, "$PACK$", _package, -1)
						str = strings.Replace(str, "$MSG_NAME$", _msg_name, -1)
						str = strings.Replace(str, "$MSG_INDEX$", strconv.Itoa(_msg_index), -1)
						_const_str = _const_str + "\n" + str

						str = `msg_name_index["$PACK$.$MSG_NAME$"] = &MsgIndex{Index: Index_$PACK$_$MSG_NAME$, MsgType: reflect.TypeOf(&$PACK$.$MSG_NAME${})}`
						str = strings.Replace(str, "$PACK$", _package, -1)
						str = strings.Replace(str, "$MSG_NAME$", _msg_name, -1)
						_init_str = _init_str + "\n" + str

						str = "host.KnowMsg(&$PACK$.$MSG_NAME${})"
						str = strings.Replace(str, "$PACK$", _package, -1)
						str = strings.Replace(str, "$MSG_NAME$", _msg_name, -1)
						_know_msg_str = _know_msg_str + "\n" + str

						str = `$MSG_NAME$ = { messageID = $MSG_INDEX$, name = "$PACK$.$MSG_NAME$" },`
						str = strings.Replace(str, "$PACK$", _package, -1)
						str = strings.Replace(str, "$MSG_NAME$", _msg_name, -1)
						str = strings.Replace(str, "$MSG_INDEX$", strconv.Itoa(_msg_index), -1)
						_lua_def = _lua_def + "\n" + str

						str = `
    public const int $MSG_NAME$ = $MSG_INDEX$;
    //private static int Index_$MSG_NAME$() => MsgIndex.$MSG_NAME$;`
						str = strings.Replace(str, "$PACK$", _package, -1)
						str = strings.Replace(str, "$MSG_NAME$", _msg_name, -1)
						str = strings.Replace(str, "$MSG_INDEX$", strconv.Itoa(_msg_index), -1)
						_cs_index = _cs_index + str

						str = `
        GetValueImpl.Specializer<Msg.$MSG_NAME$>.Fun = () => { return MsgIndex.$MSG_NAME$; };`
						str = strings.Replace(str, "$PACK$", _package, -1)
						str = strings.Replace(str, "$MSG_NAME$", _msg_name, -1)
						str = strings.Replace(str, "$MSG_INDEX$", strconv.Itoa(_msg_index), -1)
						_cs_def = _cs_def + str

						str = `
            case MsgIndex.$MSG_NAME$:
                {
                    Msg.$MSG_NAME$ m = Msg.$MSG_NAME$.Parser.ParseFrom(data, offset, length);
                    return HandleImpl.Specializer<Msg.$MSG_NAME$>.Handle(m);
                }`
						str = strings.Replace(str, "$PACK$", _package, -1)
						str = strings.Replace(str, "$MSG_NAME$", _msg_name, -1)
						str = strings.Replace(str, "$MSG_INDEX$", strconv.Itoa(_msg_index), -1)
						_cs_handle = _cs_handle + str

						if _msg_index < 30000 {
							_ts_def += strings.Replace(strings.Replace(_template_ts_def, "$MSG_INDEX$", strconv.Itoa(_msg_index), -1), "$MSG_NAME$", _msg_name, -1)
							_ts_index += strings.Replace(strings.Replace(_template_ts_index, "$MSG_INDEX$", strconv.Itoa(_msg_index), -1), "$MSG_NAME$", _msg_name, -1)
						}

						*/

						if _find_index == false {
							if _inner {
								_msg_inner_index++
							} else {
								_msg_user_index++
							}
						}
						_gen_code = true
					}
				}
			}
			readLine(path, handler)

		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
		return false
	}
	return _gen_code
	/*
		if _gen_code {
			f, err := os.OpenFile(go_path, os.O_CREATE+os.O_TRUNC, 0666)
			if err != nil {
				fmt.Println("[OpenFile] " + err.Error())
			}
			f.WriteString(genGoCode())

			str = strings.Replace(_lua_code, "$LUA_DEF$", _lua_def, -1)
			f, err = os.OpenFile(lua_path, os.O_CREATE+os.O_TRUNC, 0666)
			if err != nil {
				fmt.Println("[OpenFile] " + err.Error())
			}
			f.WriteString(str)

			str = strings.Replace(_template_ts_code, "$TS_DEF$", _ts_def, -1)
			str = strings.Replace(str, "$TS_INDEX$", _ts_index, -1)
			f, err = os.OpenFile(ts_path, os.O_CREATE+os.O_TRUNC, 0666)
			if err != nil {
				fmt.Println("[OpenFile] " + err.Error())
			}

			str = strings.Replace(_template_cs_code, "$CS_Index$", _cs_index, -1)
			str = strings.Replace(str, "$CS_Def$", _cs_def, -1)
			str = strings.Replace(str, "$CS_Handle$", _cs_handle, -1)
			f, err = os.OpenFile(cs_path, os.O_CREATE+os.O_TRUNC, 0666)
			if err != nil {
				fmt.Println("[OpenFile] " + err.Error())
			}

			f.WriteString(str)

		}
	*/
}

func saveToFile(fileName string, data string) {
	f, err := os.OpenFile(fileName, os.O_CREATE+os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("[SaveToFile] Error " + err.Error())
	} else {
		f.WriteString(data)
	}
}

func genGoCode() string {

	var _go_code = `package protocol
			
import (
	"cell/common/protocol/msg"
	"reflect"
)

/*
自動生成的代碼，不要手動修改
*/

type IMsgLibrary interface {
	KnowMsg(sample interface{})
}

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
}

func KnowAllMsg(host IMsgLibrary) {
	$KNOW_MSG$
}
`
	_const_str := ""
	_init_str := ""
	_know_msg_str := ""
	for _, v := range index {
		str := "Index_$PACK$_$MSG_NAME$ = $MSG_INDEX$"
		str = strings.Replace(str, "$PACK$", v.Package, -1)
		str = strings.Replace(str, "$MSG_NAME$", v.MsgName, -1)
		str = strings.Replace(str, "$MSG_INDEX$", strconv.Itoa(v.MsgIndex), -1)
		_const_str = _const_str + "\n" + str

		str = `msg_name_index["$PACK$.$MSG_NAME$"] = &MsgIndex{Index: Index_$PACK$_$MSG_NAME$, MsgType: reflect.TypeOf(&$PACK$.$MSG_NAME${})}`
		str = strings.Replace(str, "$PACK$", v.Package, -1)
		str = strings.Replace(str, "$MSG_NAME$", v.MsgName, -1)
		_init_str = _init_str + "\n" + str

		str = "host.KnowMsg(&$PACK$.$MSG_NAME${})"
		str = strings.Replace(str, "$PACK$", v.Package, -1)
		str = strings.Replace(str, "$MSG_NAME$", v.MsgName, -1)
		_know_msg_str = _know_msg_str + "\n" + str
	}

	str := strings.Replace(_go_code, "$CONST$", _const_str, -1)
	str = strings.Replace(str, "$INIT$", _init_str, -1)
	str = strings.Replace(str, "$KNOW_MSG$", _know_msg_str, -1)

	return str
}

func genDartCode() string {
	code := "class MessageIndex {"
	for _, v := range index {
		if v.MsgIndex > 30000 {
			continue
		}
		str := "	static const int $MSG_NAME$ = $MSG_INDEX$;"
		str = strings.Replace(str, "$PACK$", v.Package, -1)
		str = strings.Replace(str, "$MSG_NAME$", v.MsgName, -1)
		str = strings.Replace(str, "$MSG_INDEX$", strconv.Itoa(v.MsgIndex), -1)
		code = code + "\n" + str
	}
	code += "\n}\n"
	return code
}

var check_index = make(map[int]string)

func main() {
	fmt.Println("自动生成编号")
	path := ""
	go_path := ""
	lua_path := ""
	ts_path := ""
	cs_path := ""
	dart_path := ""

	flag.StringVar(&path, "path", "", "path")
	flag.StringVar(&go_path, "go_path", ".", "go_path")
	flag.StringVar(&lua_path, "lua_path", ".", "lua_path")
	flag.StringVar(&ts_path, "ts_path", ".", "ts_path")
	flag.StringVar(&cs_path, "cs_path", ".", "cs_path")
	flag.StringVar(&dart_path, "dart_path", ".", "dart_path")
	flag.Parse()

	go_path = go_path + "/msg_index.go"
	lua_path = lua_path + "/PROTO_MSG.lua"
	ts_path = ts_path + "/msg_index.ts"
	cs_path = cs_path + "/msg_index.cs"
	dart_path = dart_path + "/msg_index.dart"

	readLine(go_path, func(line string) {
		line = strings.Trim(line, " ")
		if strings.HasPrefix(line, "Index_MAX") {
			return
		}
		if strings.HasPrefix(line, "Index_msg_") {
			ss1 := strings.Split(line, "=")
			ss2 := strings.Split(strings.Trim(ss1[0], " "), "Index_msg_")

			_index, _ := strconv.Atoi(strings.Trim(ss1[1], " "))
			_name := strings.Trim(ss2[1], " ")
			//fmt.Println(">" + _name + ":" + strconv.Itoa(_index))
			if _n, ok := check_index[_index]; ok {
				fmt.Printf("============> ERROR \n")
				fmt.Printf("============> %s和%s的Index冲突了  %d \n", _name, _n, _index)
				fmt.Printf("----------------------------------------------- \n")
				return
			}
			check_index[_index] = _name
			_msg_name_index[_name] = _index
			if _index > 30000 {
				if _msg_inner_index <= _index {
					_msg_inner_index = _index + 1
					//fmt.Println("[INNER]" + _name + ":" + strconv.Itoa(_index))
				}
			} else {
				if _msg_user_index <= _index {
					_msg_user_index = _index + 1
					//fmt.Println("[USER]" + _name + ":" + strconv.Itoa(_index))
				}
			}
		}
	})

	//fmt.Println("Read: " + path)
	//fmt.Println("Write: " + go_path)

	getFilelist(path)

	{
		saveToFile(go_path, genGoCode())
		saveToFile(dart_path, genDartCode())
	}

	go_root := os.Getenv("GOROOT")
	cmd := exec.Command(go_root+`\bin\gofmt.exe`, "-w", go_path)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Waiting for command to finish...")
	//fmt.Println(cmd.Args)
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Command finished with error: %v \n", err)
	}
	//fmt.Println(out.String())
	//fmt.Println("")
}
