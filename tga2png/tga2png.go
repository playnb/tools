package main

import (
	"bytes"
	"cell/common/mustang/worker"
	"fmt"
	"github.com/andybons/gogif"
	"github.com/ftrvxmtrx/tga"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"image"
	"image/gif"
	"image/png"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
)

func GetAllFiles(pathname string) []string {
	var fileList []string

	rd, _ := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			fileList = append(fileList, GetAllFiles(pathname+fi.Name()+`\`)...)
		} else {
			fileList = append(fileList, pathname+fi.Name())
		}
	}
	return fileList
}

func getFileList(path string) []string {
	var fileList []string
	err := filepath.Walk(path, func(p string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if path == p {
			return nil
		}
		if f.IsDir() {
			fileList = append(fileList, getFileList(p+`\`)...)
			return nil
		} else {
			fileList = append(fileList, p)
			return nil
		}
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	return fileList
}

/*
PNG压缩
http://nullice.com/limitPNG/#download

body:
00000	4*8 待机
00002	16*8 跑步
00003	8*8 攻击普通
00004	8*8 攻击刀2
00005	8*8 攻击刀3
00006	8*8 攻击杖1
00007	8*8 攻击杖2
00008	8*8 攻击琴1
00009	8*8 攻击琴2
00010	8*8 攻击弓1
00011	8*8 攻击弓2
00012	7*8 待机  马上
00013	7*8 跑步  马上
00014	4*8 采集
00015	6*8 死亡
00016	3*8 打坐
00017	1*8 轻功启动
00018	8*8 攻击X
00019	1*8 轻功中
00020	11*8 后空翻
00021	11*8 侧空翻
00022	8*8 攻击杖  马上
00023	8*8 死亡击飞
00032	4*8 待机  马上
00033	7*8 跑步  马上
00034	8*8 攻击刀  马上
00035	8*8 攻击琴  马上
00036	8*8 攻击弓  马上
00037	8*8 死亡  马上
00038	11*8 抬前蹄  马上
*/

type ImgDesc struct {
	Name  string
	Frame int
	Dir   int
	Ride  bool
	Desc  string
}

var AllDesc = []*ImgDesc{
	{"00000", 4, 8, false, "待机"},
	{"00002", 16, 8, false, "跑步"},
	{"00003", 8, 8, false, "攻击普通"},
	{"00004", 8, 8, false, "攻击刀1"},
	{"00005", 8, 8, false, "攻击刀2"},
	{"00006", 8, 8, false, "攻击杖1"},
	{"00007", 8, 8, false, "攻击杖2"},
	{"00008", 8, 8, false, "攻击琴1"},
	{"00009", 8, 8, false, "攻击琴2"},
	{"00010", 8, 8, false, "攻击弓1"},
	{"00011", 8, 8, false, "攻击弓2"},
	{"00012", 7, 8, true, "待机"},
	{"00013", 7, 8, true, "跑步"},
	{"00014", 4, 8, false, "采集"},
	{"00015", 6, 8, false, "死亡"},
	{"00016", 3, 8, false, "打坐"},
	{"00017", 1, 8, false, "轻功启动"},
	{"00018", 8, 8, false, "攻击X"},
	{"00019", 1, 8, false, "轻功中"},
	{"00020", 11, 8, false, "后空翻"},
	{"00021", 11, 8, false, "侧空翻"},
	{"00022", 8, 8, true, "攻击杖"},
	{"00023", 8, 8, false, "死亡击飞"},
	{"00032", 4, 8, true, "待机"},
	{"00033", 7, 8, true, "跑步"},
	{"00034", 8, 8, true, "攻击刀"},
	{"00035", 8, 8, true, "攻击琴"},
	{"00036", 8, 8, true, "攻击弓"},
	{"00037", 8, 8, true, "死亡"},
	{"00038", 11, 8, true, "抬前蹄"},
}

func findDesc(name string) *ImgDesc {
	for _, v := range AllDesc {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func splitDir(v string) (string, string, string) {
	dNameOut := ""
	fName := ""
	actionName := ""
	roleName := ""
	dNameOut, fName = filepath.Split(v)
	dNameOut = strings.TrimSuffix(dNameOut, `\`)
	dNameOut = strings.TrimSuffix(dNameOut, `/`)
	dNameOut, actionName = filepath.Split(dNameOut)
	dNameOut = strings.TrimSuffix(dNameOut, `\`)
	dNameOut = strings.TrimSuffix(dNameOut, `/`)
	dNameOut, roleName = filepath.Split(dNameOut)
	return fName, actionName, roleName
}

func readImage(f string) image.Image {
	fileIn, err := os.Open(f)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	defer fileIn.Close()
	img, err := tga.Decode(fileIn)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	return img
}

func shortString(str string) string {
	n, _ := strconv.Atoi(str)
	return strconv.FormatInt(int64(n), 10)
}

//一个动画
type Animation struct {
	Img []image.Image
}

//一个动作
type Action struct {
	ActionName string
	RoleName   string
	Animation  []*Animation //8方向
}

var animations = make(map[string]*Action)

func main() {
	kindStr := "wing"
	kindStrShort := "WI"
	pflag.String("in_dir", `C:\ltp\code\h5\zt2\wing\`, "源文件夹")
	pflag.String("out_dir", `C:\ltp\code\h5\zt2_png\wing\`, "目标文件夹")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	inDir := viper.GetString("in_dir")
	outDir := viper.GetString("out_dir")
	//files := getFileList(inDir)
	files := GetAllFiles(inDir)
	totalCount := int32(len(files))
	doneCount := int32(0)
	wl := worker.NewWaitAndLimit(10)
	// os.Exit(11)

	for _, v := range files {
		_, actionName, roleName := splitDir(v)
		animations[actionName+"_"+roleName] = &Action{}
		animations[actionName+"_"+roleName].ActionName = shortString(actionName)
		animations[actionName+"_"+roleName].RoleName = shortString(roleName)
	}

	loadImage := func(v string) {
		_, actionName, roleName := splitDir(v)
		action := animations[actionName+"_"+roleName]
		var animation *Animation

		/*
				desc := findDesc(actionName)
			if len(action.Animation) == 0 {
				animation = &Animation{}
				action.Animation = append(action.Animation, animation)
			} else {
				animation = action.Animation[len(action.Animation)-1]
			}
			if len(animation.Img) == desc.Frame {
				animation = &Animation{}
				action.Animation = append(action.Animation, animation)
			}
		*/
		if len(action.Animation) == 0 {
			animation = &Animation{}
			action.Animation = append(action.Animation, animation)
		} else {
			animation = action.Animation[0]
		}

		animation.Img = append(animation.Img, readImage(v))

		doneCount = atomic.AddInt32(&doneCount, 1)
		fmt.Printf("RAED %s	%d%%\n", filepath.Dir(v), doneCount*100/totalCount)
	}
	for _, v := range files {
		loadImage(v)
	}
	wl.Wait()

	for _, action := range animations {
		if len(action.Animation) == 1 {
			if len(action.Animation[0].Img)%8 != 0 {
				panic("不是8方向???")
			}
			frameCount := len(action.Animation[0].Img) / 8
			as := action.Animation[0].Img
			action.Animation = nil
			animation := &Animation{}
			action.Animation = append(action.Animation, animation)
			for i := 0; i < len(as); i++ {
				if frameCount == len(animation.Img) {
					animation = &Animation{}
					action.Animation = append(action.Animation, animation)
				}
				animation.Img = append(animation.Img, as[i])
			}
		}
	}

	doneCount = 0
	var batch []string
	outGif := false
	outputAction := func(action *Action) {
		for d, animation := range action.Animation {
			if d >= 5 {
				continue
			}
			if outGif {
				s := fmt.Sprintf("%s/%s/%s/%d.gif", outDir, action.RoleName, action.ActionName, d)
				nframes := len(animation.Img)
				delay := 8
				anim := gif.GIF{LoopCount: nframes}
				for _, img := range animation.Img {
					paletted := image.NewPaletted(img.Bounds(), nil)
					quantizer := gogif.MedianCutQuantizer{NumColor: 255}
					quantizer.Quantize(paletted, img.Bounds(), img, image.ZP)
					/*
						paletted := image.NewPaletted(img.Bounds(), palette.Plan9)
						draw.FloydSteinberg.Draw(paletted, img.Bounds(), img, image.ZP)
					*/
					anim.Delay = append(anim.Delay, delay)
					anim.Image = append(anim.Image, paletted)
					atomic.AddInt32(&doneCount, 1)
				}
				os.MkdirAll(filepath.Dir(s), os.ModePerm)
				f, _ := os.Create(s)
				gif.EncodeAll(f, &anim)
				f.Close()
				fmt.Printf("WRITE %s	%d%%\n", filepath.Dir(s), doneCount*100/totalCount)

			} else {
				for k, img := range animation.Img {
					s := fmt.Sprintf("%s/%s/%s/%d/", outDir, action.RoleName, action.ActionName, d)
					s = s + fmt.Sprintf("%03d.png", k)

					s = fmt.Sprintf("%s/%s/%s/%d%03d.png", outDir, action.RoleName, action.ActionName, d, k)

					buf := bytes.NewBuffer(nil)
					err := png.Encode(buf, img)

					err = os.MkdirAll(filepath.Dir(s), os.ModePerm)
					if err != nil {
						fmt.Errorf(err.Error())
					}
					fileOut, err := os.Create(s)
					if err != nil {
						fmt.Errorf(err.Error())
					}
					fileOut.Write(buf.Bytes())
					if err != nil {
						fmt.Errorf(err.Error())
					}
					fileOut.Close()

					doneCount = atomic.AddInt32(&doneCount, 1)
					fmt.Printf("WRITE %s	%d%%\n", filepath.Dir(s), doneCount*100/totalCount)
				}
				//magick -delay 16 -loop 0 sourcePath gifPath
				if false {
					sourcePath := fmt.Sprintf("%s%s\\%s\\%d\\*.png", outDir, action.RoleName, action.ActionName, d)
					gifPath := fmt.Sprintf("%s%s\\%s\\%d.gif", outDir, action.RoleName, action.ActionName, d)
					err := exec.Command("magick -delay 16 -loop 0", sourcePath, gifPath).Run()
					if err != nil {
						fmt.Errorf(err.Error())
					}
				} else {
					/*
						sourcePath := fmt.Sprintf("%s%s\\%s\\", outDir, action.RoleName, action.ActionName)
						jsonPath := fmt.Sprintf("%s%s\\%s.json", outDir, action.RoleName, action.ActionName)
						err := exec.Command(`D:\Program Files\Egret\TextureMerger\TextureMerger`,
							`-p`,
							sourcePath,
							`-o`,
							jsonPath).Run()
						if err != nil {
							fmt.Println(err.Error())
						}
					*/
				}
			}
		}

		sourcePath := fmt.Sprintf("%s%s\\%s\\", outDir, action.RoleName, action.ActionName)
		jsonPath := fmt.Sprintf("%%OUT%%\\%s\\%s_%s_%s.json", action.RoleName, kindStrShort, action.RoleName, action.ActionName)
		cmdLine := "TextureMerger -p " + sourcePath + " -o " + jsonPath
		//fmt.Println(cmdLine)
		batch = append(batch, cmdLine)
	}

	for _, action := range animations {
		wl.Add()
		go func(a *Action) {
			defer wl.Done()
			outputAction(a)
		}(action)
	}
	wl.Wait()

	f, _ := os.Create(outDir + "texture.bat")
	f.WriteString(`SET OUT=C:\ltp\code\h5\zt2_png\sprites\` + kindStr + "\n")
	f.WriteString(`SET PATH=C:\Program Files\Egret\TextureMerger;%PATH%` + "\n")
	for _, c := range batch {
		f.WriteString(c + "\n\r")
	}
	f.Close()
	return

	convert := func(f string) {
		if filepath.Ext(f) != ".tga" {
			return
		}
		s := strings.TrimPrefix(f, inDir)
		s = strings.TrimSuffix(s, ".tga")
		s += ".png"
		s = outDir + s

		if false {
			dNameOut, fNameOut := filepath.Split(s)
			dNameOut = strings.TrimSuffix(dNameOut, `\`)
			dNameOut = strings.TrimSuffix(dNameOut, `/`)
			//fmt.Println("dNameOut: " + dNameOut)
			//fmt.Println("fNameOut: " + fNameOut)
			dNameOut, listDirName := filepath.Split(dNameOut)
			//fmt.Println("dNameOut: " + dNameOut)
			//fmt.Println("listDirName: " + listDirName)

			s = dNameOut + "/" + listDirName + "_" + fNameOut
		}

		doneCount = atomic.AddInt32(&doneCount, 1)
		fmt.Printf("%s	%d%%\n", filepath.Dir(s), doneCount*100/totalCount)

		err := os.MkdirAll(filepath.Dir(s), os.ModePerm)
		if err != nil {
			fmt.Errorf(err.Error())
		}
		fileIn, err := os.Open(f)
		if err != nil {
			fmt.Errorf(err.Error())
		}
		defer fileIn.Close()
		img, err := tga.Decode(fileIn)
		if err != nil {
			fmt.Errorf(err.Error())
		}

		fileOut, err := os.Create(s)
		if err != nil {
			fmt.Errorf(err.Error())
		}
		defer fileOut.Close()

		buf := bytes.NewBuffer(nil)
		err = png.Encode(buf, img)

		fileOut.Write(buf.Bytes())
		if err != nil {
			fmt.Errorf(err.Error())
		}
	}

	for _, v := range files {
		wl.Add()
		func() {
			defer wl.Done()
			convert(v)
		}()
	}
	wl.Wait()
}
