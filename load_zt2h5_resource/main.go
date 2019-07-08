package main

import (
	"cell/common/mustang/worker"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var baseUrl = `http://192.168.116.221/zt2h5/resource/`

func loadFile(url string, file string) []byte {
	fmt.Println("LoadFile " + url + file)
	client := &http.Client{}
	client.Timeout = time.Second * 60 //设置超时时间
	resp, err := client.Get(url + file)
	if err != nil {
		panic(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
		return nil
	}
	return body
}

func main() {
	defer fmt.Println("exit...")
	index := make(map[string]string)
	err := json.Unmarshal(loadFile(baseUrl, `version.json`), &index)
	if err != nil {
		log.Println(err)
		return
	}

	count := 0
	limit := worker.NewWaitAndLimit(10)
	for k, v := range index {
		fileExt := path.Ext(k)
		limit.Add()
		func(k, v string) {
			defer limit.Done()
			var err error
			data := loadFile(baseUrl, strings.TrimSuffix(k, fileExt)+"."+v+fileExt)
			file := `zt2h5\` + k
			err = os.MkdirAll(path.Dir(file), os.ModeDir)
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile(file, data, os.ModePerm)
			if err != nil {
				panic(err)
			}

			f, err := os.Create(file)
			if err != nil {
				panic(err)
			}
			_, err = f.Write(data)
			if err != nil {
				panic(err)
			}
			f.Close()
			count++
		}(k, v)
	}
	limit.Wait()
	fmt.Println(count)
}
