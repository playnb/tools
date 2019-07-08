package main

import (
	"cell/common/mustang/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type SpriteSheets struct {
	File   string                 `json:"file"`
	Frames map[string]SpriteFrame `json:"frames"`
}

type SpriteFrame struct {
	X       int `json:"x"`
	Y       int `json:"y"`
	H       int `json:"h"`
	W       int `json:"w"`
	OffX    int `json:"offX"`
	OffY    int `json:"offY"`
	SourceW int `json:"sourceW"`
	SourceH int `json:"sourceH"`
}

type MoveClip struct {
	MC  map[string]MoveClipMC  `json:"mc"`
	Res map[string]MoveClipRes `json:"res"`
}

type MoveClipMC struct {
	FrameRate int                `json:"frameRate"`
	Labels    [] MoveClipMCLabel `json:"labels"`
	Frames    [] MoveClipMCFrame `json:"frames"`
}
type MoveClipMCLabel struct {
	Name  string `json:"name"`
	Frame int    `json:"frame"`
	End   int    `json:"end"`
}
type MoveClipMCFrame struct {
	Res string `json:"res"`
	X   int    `json:"x"`
	Y   int    `json:"y"`
}
type MoveClipRes struct {
	X int `json:"x"`
	Y int `json:"y"`
	H int `json:"h"`
	W int `json:"w"`
}

func convertMc(filePath string) {
	SS := &SpriteSheets{}
	data, _ := ioutil.ReadFile(filePath)
	err := json.Unmarshal(data, SS)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(SS)

	{
		mName := strings.TrimSuffix(SS.File, ".png")

		mc := &MoveClip{}
		mc.Res = make(map[string]MoveClipRes)
		mc.MC = make(map[string]MoveClipMC)
		for k, v := range SS.Frames {
			mc.Res[mName + "_" + k] = MoveClipRes{
				X: v.X,
				Y: v.Y,
				H: v.H,
				W: v.W,
			}
		}
		motion := MoveClipMC{}
		motion.FrameRate = 5
		for k, v := range SS.Frames {
			frame := MoveClipMCFrame{}
			frame.Res = mName + "_" + k
			frame.X = v.OffX - v.SourceW/2
			frame.Y = v.OffY - v.SourceH/2
			motion.Frames = append(motion.Frames, frame)
		}
		sort.Slice(motion.Frames, func(i, j int) bool {
			return motion.Frames[i].Res < motion.Frames[j].Res
		})
		index := 1
		fc := len(motion.Frames) / 5
		for i := 1; i <= 5; i++ {
			motion.Labels = append(motion.Labels, MoveClipMCLabel{
				Name:  fmt.Sprintf("%d", i),
				Frame: index,
				End:   index + fc - 1,
			})
			index = index + fc
		}
		mc.MC[mName] = motion

		data, _ := json.Marshal(mc)
		ioutil.WriteFile(filePath, data, os.ModePerm)
	}
}

func main() {
	files := util.GetAllFiles(`C:\ltp\code\mmo\server\src\cell\common\tools\mc\sprites\`)
	for _, v := range files {
		if strings.HasSuffix(v, ".json") {
			convertMc(v)
		}
	}
}
