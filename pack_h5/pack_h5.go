package main

import (
	"bytes"
	"cell/common/mustang/util"
	"cell/common/mustang/worker"
	"encoding/hex"
	"flag"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync/atomic"
	"time"
)

var MTCT_BASE = `C:\Program Files\ARM\Mali Developer Tools\Mali Texture Compression Tool v4.3.0\bin`
var convertKTX = false
var baseDir = ``
var fileCount = int32(0)
var maxFileCount = int32(0)

func crc32File(f string) {
	ieee := crc32.NewIEEE()
	b, e := ioutil.ReadFile(f)
	if e == nil {
		ieee.Write(b)
		code := hex.EncodeToString(ieee.Sum(nil))
		newFileName := baseDir + `_P` + strings.TrimPrefix(f, baseDir)
		nd, nf := path.Split(newFileName)
		nd = strings.ReplaceAll(nd, `//`, `\`)
		nd = strings.ReplaceAll(nd, `/`, `\`)
		nd = strings.TrimSuffix(nd, `\`)
		os.MkdirAll(nd, os.ModePerm)
		ss := strings.Split(nf, ".")
		nf = nd + `\`
		for i := 0; i < len(ss); i++ {
			nf += ss[i]
			if i == len(ss)-1 {

			} else {
				nf += "."
			}
			if i == len(ss)-2 {
				nf += code
				nf += "."
			}
		}
		ioutil.WriteFile(nf, b, os.ModePerm)

		if convertKTX && (strings.HasSuffix(nf, ".png") || strings.HasSuffix(nf, ".jpg")) {
			var out bytes.Buffer
			fmt.Println(nf)
			fmt.Println(nd)

			//c := exec.Command(MTCT, nf, nd, "-ktx")
			{
				c := exec.Command("convert", nf, nf+".ppm")
				c.Stderr = &out
				c.Run()
			}
			{
				c := exec.Command("etcpack", nf+".ppm", nd, "-ktx")
				c.Stderr = &out
				c.Run()
			}
			fmt.Println(out.String())
		} else {
			fmt.Println(nf)
		}

		atomic.AddInt32(&fileCount, 1)
		fmt.Printf("%d/%d(%f)\n", fileCount, maxFileCount, float32(fileCount)/float32(maxFileCount))
	} else {
		panic(e)
	}
}

func main() {
	flag.StringVar(&baseDir, "base", ``, "打包的目标目录")
	flag.Parse()
	t1 := time.Now()
	wg := worker.NewWaitAndLimit(16)
	fileList := util.GetAllFiles(baseDir)
	maxFileCount = int32(len(fileList))
	for _, f := range fileList {
		wg.Add()
		go func(file string) {
			defer wg.Done()
			crc32File(file)
		}(f)
	}
	wg.Wait()

	fmt.Printf("消耗时间:%ds \n", time.Now().Unix()-t1.Unix())
}
