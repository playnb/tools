package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func getFilelist(path string, fileFunc func(string)) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		_, fileName := filepath.Split(path)
		fileFunc(fileName)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

func splitFileName(name string) int {
	ss1 := strings.Split(name, ".")
	ss2 := strings.Split(ss1[0], "_")
	if len(ss2) > 1 {
		s := ss2[len(ss2)-1]
		n, _ := strconv.Atoi(s)
		return n
	}
	return 0
}

func mergeImg(s string, dstImg draw.Image, x, y int) {
	file, err := os.Open(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	img, err := jpeg.Decode(file) //解码
	if err != nil {
		fmt.Println(err)
		return
	}

	draw.Draw(dstImg, image.Rect(x, y, x+512, y+512), img, image.Pt(0, 0), draw.Src)
}

/*
合并： --merge=true --split_dir=C:\ltp\code\mmo\client\resource\assets\sprites\map\ --one_file=C:\ltp\code\mmo\common\maps\ --name=map_4 --width=7 --height=6

*/
func main() {
	pflag.String("split_dir", "", "源文件夹")
	pflag.String("one_file", "", "目标文件")
	pflag.String("name", "", "地图名")
	pflag.Int("size", 512, "tile的尺寸")
	pflag.Int("width", 6, "宽度格子数")
	pflag.Int("height", 6, "高度格子数")
	pflag.Bool("merge", true, "合并/拆分")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	oneFilePath := viper.GetString("one_file") + "/" + viper.GetString("name") + "/" + viper.GetString("name") + ".jpg"
	splitDir := viper.GetString("split_dir") + "/" + viper.GetString("name")
	mapName := viper.GetString("name")

	if viper.GetBool("merge") {
		files := make([]string, 0)
		getFilelist(splitDir, func(s string) {
			files = append(files, s)
		})

		sort.Slice(files, func(i, j int) bool {
			return splitFileName(files[i]) < splitFileName(files[j])
		})

		width := viper.GetInt("width") * viper.GetInt("size")
		height := viper.GetInt("height") * viper.GetInt("size")
		dstImg := image.NewNRGBA(image.Rect(0, 0, width, height))
		i := 0
		j := 0
		for _, f := range files {
			mergeImg(splitDir+"/"+f, dstImg, j*512, i*512)
			j++
			if j >= viper.GetInt("width") {
				j = 0
				i++
			}
		}

		os.MkdirAll(path.Dir(oneFilePath), os.ModePerm)
		file1, err := os.Create(oneFilePath)
		if err != nil {
			fmt.Println(err)
		}
		jpeg.Encode(file1, dstImg, &jpeg.Options{Quality: 100})
		file1.Close()
	} else {
		var srcImg image.Image
		var err error

		file, err := os.Open(oneFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		srcImg, err = jpeg.Decode(file) //解码
		if err != nil {
			fmt.Println(err)
			return
		}

		os.MkdirAll(splitDir, os.ModePerm)
		for i := 0; i < viper.GetInt("width"); i++ {
			for j := 0; j < viper.GetInt("height"); j++ {
				fileName := fmt.Sprintf(mapName+"_%d_%d.jpg", i, j)
				out, err := os.Create(splitDir + "/" + fileName)
				if err != nil {
					fmt.Println(err)
					return
				}
				dstImg := image.NewNRGBA(image.Rect(0, 0, viper.GetInt("size"), viper.GetInt("size")))

				draw.Draw(dstImg, dstImg.Bounds(), srcImg, image.Pt(j*512, i*512), draw.Src)

				jpeg.Encode(out, dstImg, &jpeg.Options{Quality: 100})
				out.Close()
			}
		}
	}
}
