package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type MsgMap struct {
	Name string `xml:"msg,attr"`
}

type MapConfig struct {
	RoomService []MsgMap `xml:"room_service"`
}

func main() {
	base_dir := ""
	flag.StringVar(&base_dir, "base_dir", `C:\code\`, "base_dir")
	flag.Parse()

	path := base_dir + `cell\common\trunk\branch_0`
	//outPath := base_dir + `server\trunk\branch_0\src\cell\net_service`
	fileName := path + `\msg_map.xml`

	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	config := &MapConfig{}
	err = xml.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}

	for _, msg := range config.RoomService {
		fmt.Println(msg.Name)
	}

	fmt.Println("[RoomService]: ", config.RoomService)
}
