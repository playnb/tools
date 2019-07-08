package main

import (
	"cell/conf"
	rdExcel "cell/conf/tbl"
	"cell/common/mustang/util"
	"flag"
	"fmt"
	"strings"
	"time"
)

func main() {
	flag.Set("config", "./")
	//flag.Set("config", `C:\code\doc`)
	flag.Parse()
	fmt.Println(conf.GetMe().ConfigDir)
	rdExcel.ReadTbl()
	rdExcel.LoadConfig()
	fmt.Println("读取Excel文件")

	fmt.Println("================================================\n\n")
	fmt.Println("================================================")
	//	allNames := make(map[string]bool)
	rdExcel.ForeachRandomNameFirstTbl(func(name string) {
		ok := true
		util.ReadFileLine(conf.GetMe().ConfigDir+"/Filter/BadWord.txt", func(line string) {
			if !ok {
				return
			}
			if len(strings.Split(name, line)[0]) != len(name) {
				fmt.Printf("前缀 %s 有问题(%s)\n", name, line)
				ok = false
			}
		})
	})

	rdExcel.ForeachRandomNameSecondTbl(func(name string) {
		ok := true
		util.ReadFileLine(conf.GetMe().ConfigDir+"/Filter/BadWord.txt", func(line string) {
			if !ok {
				return
			}
			if len(strings.Split(name, line)[0]) != len(name) {
				fmt.Printf("中间 %s 有问题(%s)\n", name, line)
				ok = false
			}
		})
	})

	rdExcel.ForeachRandomNameThirdTbl(func(name string) {
		ok := true
		util.ReadFileLine(conf.GetMe().ConfigDir+"/Filter/BadWord.txt", func(line string) {
			if !ok {
				return
			}
			if len(strings.Split(name, line)[0]) != len(name) {
				fmt.Printf("后缀 %s 有问题(%s)\n", name, line)
				ok = false
			}
		})
	})

	for {
		time.Sleep(time.Second)
	}
}
