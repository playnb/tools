package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Cell struct {
	VarName    string
	ClientType string
	ServerType string
	Values     []string
}

type TypeFunc struct {
	Description string
	Val         func(string) interface{}
	Default     func() interface{}
}

var types = make(map[string]*TypeFunc)

func hasType(t string) bool {
	if len(t) == 0 {
		return true
	}
	if _, ok := types[t]; ok {
		return true
	}
	return false
}

func init() {
	types["string"] = &TypeFunc{
		Val: func(s string) interface{} {
			return s
		},
		Default: func() interface{} {
			return ""
		},
		Description: "字符串类型",
	}

	types["int"] = &TypeFunc{
		Val: func(s string) interface{} {
			n, _ := strconv.ParseInt(s, 10, 64)
			return n
		},
		Default: func() interface{} {
			return 0
		},
		Description: "整数类型",
	}

	types["float"] = &TypeFunc{
		Val: func(s string) interface{} {
			n, _ := strconv.ParseFloat(s, 64)
			return n
		},
		Default: func() interface{} {
			return 0
		},
		Description: "浮点数类型",
	}

	types["array<int>"] = &TypeFunc{
		Val: func(s string) interface{} {
			a := make([]int64, 0)
			if len(s) > 0 {
				ss := strings.Split(s, ",")
				for _, v := range ss {
					n, _ := strconv.ParseInt(v, 10, 64)
					a = append(a, n)
				}
			}
			return a
		},
		Default: func() interface{} {
			a := make([]int64, 0)
			return a
		},
		Description: "整数数组类型",
	}

	types["array<float>"] = &TypeFunc{
		Val: func(s string) interface{} {
			a := make([]float64, 0)
			if len(s) > 0 {
				ss := strings.Split(s, ",")
				for _, v := range ss {
					n, _ := strconv.ParseFloat(v, 64)
					a = append(a, n)
				}
			}
			return a
		},
		Default: func() interface{} {
			a := make([]float64, 0)
			return a
		},
		Description: "浮点数数组类型",
	}

	types["array<string>"] = &TypeFunc{
		Val: func(s string) interface{} {
			if len(s) > 0 {
				return strings.Split(s, ",")
			}
			return make([]string, 0)
		},
		Default: func() interface{} {
			return make([]string, 0)
		},
		Description: "字符串数组类型",
	}
}

func main() {
	pflag.String("excel", "", "目标excel文件")
	pflag.String("server", "", "服务器json文件目录")
	pflag.String("client", "", "客户端json文件目录")
	pflag.String("help", "", "客户端json文件目录")
	pflag.Bool("types", false, "支持的类型")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if viper.GetBool("types") {
		fmt.Println("支持的类型有:")
		for k, v := range types {
			fmt.Printf("%s: %s\n", k, v.Description)
		}
		return
	}

	excelFile := viper.GetString("excel")
	xlFile, err := xlsx.OpenFile(excelFile)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	maxCount := 0

	fmt.Println("读取Excel文件: " + excelFile)
	for _, sheet := range xlFile.Sheets {
		if !strings.HasPrefix(sheet.Name, "(pack)") {
			continue
		}
		cells := make(map[int]*Cell)
		fileName := strings.TrimPrefix(sheet.Name, "(pack)")

		for rowIndex, row := range sheet.Rows {
			switch rowIndex {
			case 0: //变量名
				for colIndex, cell := range row.Cells {
					cells[colIndex] = &Cell{VarName: cell.Value}
				}
			case 1: //客户端类型
				for colIndex, cell := range row.Cells {
					cells[colIndex].ClientType = strings.Trim(cell.Value, " ")
					if !hasType(cell.Value) {
						log.Fatalln(excelFile + " 未知的客户端类型=>  " + cells[colIndex].VarName + ":" + cell.Value)
					}
				}
			case 2: //服务端类型
				for colIndex, cell := range row.Cells {
					cells[colIndex].ServerType = strings.Trim(cell.Value, " ")
					if !hasType(cell.Value) {
						log.Fatalln(excelFile + " 未知的服务端类型=>  " + cells[colIndex].VarName + ":" + cell.Value)
					}
				}
			case 3: //策划用的字段名
			default:
				for colIndex, cell := range row.Cells {
					cells[colIndex].Values = append(cells[colIndex].Values, cell.Value)
					if maxCount < len(cells[colIndex].Values) {
						maxCount = len(cells[colIndex].Values)
					}
				}
			}
		}
		fileName += ".json"
		fmt.Println("生成json数据文件: " + fileName)
		{
			vals := make([]map[string]interface{}, maxCount)
			for k, _ := range vals {
				vals[k] = make(map[string]interface{})
			}
			hasValue := false
			for _, c := range cells {
				for i := 0; i < maxCount; i++ {
					valType := c.ServerType
					if len(valType) > 0 {
						if i < len(c.Values) {
							vals[i][c.VarName] = types[valType].Val(c.Values[i])
						} else {
							vals[i][c.VarName] = types[valType].Default()
						}
						hasValue = true
					}
				}
			}
			if hasValue {
				str, _ := jsoniter.MarshalToString(vals)
				os.MkdirAll(viper.GetString("server"), os.ModePerm)
				ioutil.WriteFile(viper.GetString("server")+"/"+fileName, []byte(str), os.ModePerm)
			}
		}
		{
			vals := make([]map[string]interface{}, maxCount)
			for k, _ := range vals {
				vals[k] = make(map[string]interface{})
			}
			hasValue := false
			for _, c := range cells {
				for i := 0; i < maxCount; i++ {
					valType := c.ClientType
					if len(valType) > 0 {
						if i < len(c.Values) {
							vals[i][c.VarName] = types[valType].Val(c.Values[i])
						} else {
							vals[i][c.VarName] = types[valType].Default()
						}
						hasValue = true
					}
				}
			}
			if hasValue {
				str, _ := jsoniter.MarshalToString(vals)
				os.MkdirAll(viper.GetString("client"), os.ModePerm)
				ioutil.WriteFile(viper.GetString("client")+"/"+fileName, []byte(str), os.ModePerm)
			}
		}
	}

}
