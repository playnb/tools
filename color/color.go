package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math"
	"os"
)

func Color2StandardRGB(c color.Color) (sR, sG, sB float64) {
	r, g, b, _ := c.RGBA()
	switch c.(type) {
	case color.RGBA:
		r = (r >> 8)
		g = (g >> 8)
		b = (b >> 8)
	case color.NRGBA:
		r = (r >> 8)
		g = (g >> 8)
		b = (b >> 8)
	}
	sR = float64(r)
	sG = float64(g)
	sB = float64(b)
	return sR, sG, sB
}

func StandardRGB2XYZ(sR, sG, sB float64) (X, Y, Z float64) {
	var_R := (float64(sR) / 255.0)
	var_G := (float64(sG) / 255.0)
	var_B := (float64(sB) / 255.0)

	if var_R > 0.04045 {
		var_R = math.Pow((var_R+0.055)/1.055, 2.4)
	} else {
		var_R = var_R / 12.92
	}
	if var_G > 0.04045 {
		var_G = math.Pow((var_G+0.055)/1.055, 2.4)
	} else {
		var_G = var_G / 12.92
	}
	if var_B > 0.04045 {
		var_B = math.Pow((var_B+0.055)/1.055, 2.4)
	} else {
		var_B = var_B / 12.92
	}

	var_R = var_R * 100
	var_G = var_G * 100
	var_B = var_B * 100

	X = var_R*0.4124 + var_G*0.3576 + var_B*0.1805
	Y = var_R*0.2126 + var_G*0.7152 + var_B*0.0722
	Z = var_R*0.0193 + var_G*0.1192 + var_B*0.9505

	return X, Y, Z
}

func XYZ2StandardRGB(X, Y, Z float64) (sR, sG, sB float64) {
	var_X := X / 100
	var_Y := Y / 100
	var_Z := Z / 100

	var_R := var_X*3.2406 + var_Y*-1.5372 + var_Z*-0.4986
	var_G := var_X*-0.9689 + var_Y*1.8758 + var_Z*0.0415
	var_B := var_X*0.0557 + var_Y*-0.2040 + var_Z*1.0570

	if var_R > 0.0031308 {
		var_R = 1.055*(math.Pow(var_R, 1/2.4)) - 0.055
	} else {
		var_R = 12.92 * var_R
	}
	if var_G > 0.0031308 {
		var_G = 1.055*(math.Pow(var_G, 1/2.4)) - 0.055
	} else {
		var_G = 12.92 * var_G
	}
	if var_B > 0.0031308 {
		var_B = 1.055*(math.Pow(var_B, 1/2.4)) - 0.055
	} else {
		var_B = 12.92 * var_B
	}

	sR = var_R * 255
	sG = var_G * 255
	sB = var_B * 255

	return sR, sG, sB
}

var ReferenceX = float64(95.047)
var ReferenceY = float64(100.0)
var ReferenceZ = float64(108.883)

func rad2deg(r float64) float64 {
	return r
}

func CieLab2Hue(var_a, var_b float64) float64 {
	var_bias := float64(0)
	if (var_a >= 0 && var_b == 0) {
		return 0
	}
	if (var_a < 0 && var_b == 0) {
		return 180
	}
	if (var_a == 0 && var_b > 0) {
		return 90
	}
	if (var_a == 0 && var_b < 0) {
		return 270
	}
	if (var_a > 0 && var_b > 0) {
		var_bias = 0
	}
	if (var_a < 0) {
		var_bias = 180
	}
	if (var_a > 0 && var_b < 0) {
		var_bias = 360
	}
	return (rad2deg(math.Atan(var_b/var_a)) + var_bias)
}

func XYZ2HunterLab(X, Y, Z float64) (L, a, b float64) {
	var_Ka := (175.0 / 198.04) * (ReferenceY + ReferenceX)
	var_Kb := (70.0 / 218.11) * (ReferenceY + ReferenceZ)

	L = 100.0 * math.Sqrt(Y/ReferenceY)
	a = var_Ka * (((X / ReferenceX) - (Y / ReferenceY)) / math.Sqrt(Y/ReferenceY))
	b = var_Kb * (((Y / ReferenceY) - (Z / ReferenceZ)) / math.Sqrt(Y/ReferenceY))
	return L, a, b
}

func HunterLab2XYZ(L, a, b float64) (X, Y, Z float64) {
	var_Ka := (175.0 / 198.04) * (ReferenceY + ReferenceX)
	var_Kb := (70.0 / 218.11) * (ReferenceY + ReferenceZ)

	Y = (math.Pow(L/ReferenceY, 2)) * 100.0
	X = (a/var_Ka*math.Sqrt(Y/ReferenceY) + (Y / ReferenceY)) * ReferenceX
	Z = - (b/var_Kb*math.Sqrt(Y/ReferenceY) - (Y / ReferenceY)) * ReferenceZ

	return X, Y, Z
}

func XYZ2CIELab(X, Y, Z float64) (L, a, b float64) {
	var_X := X / ReferenceX
	var_Y := Y / ReferenceY
	var_Z := Z / ReferenceZ

	if var_X > 0.008856 {
		var_X = math.Pow(var_X, 1.0/3.0)
	} else {
		var_X = (7.787 * var_X) + (16 / 116)
	}
	if var_Y > 0.008856 {
		var_Y = math.Pow(var_Y, 1.0/3.0)
	} else {
		var_Y = (7.787 * var_Y) + (16 / 116)
	}
	if var_Z > 0.008856 {
		var_Z = math.Pow(var_Z, 1.0/3.0)
	} else {
		var_Z = (7.787 * var_Z) + (16 / 116)
	}

	L = (116 * var_Y) - 16
	a = 500 * (var_X - var_Y)
	b = 200 * (var_Y - var_Z)
	return L, a, b
}

func CIELab2XYZ(L, a, b float64) (X, Y, Z float64) {
	var_Y := (L * + 16) / 116
	var_X := a/500 + var_Y
	var_Z := var_Y - b/200

	if math.Pow(var_Y, 3) > 0.008856 {
		var_Y = math.Pow(var_Y, 3)
	} else {
		var_Y = (var_Y - 16/116) / 7.787
	}
	if math.Pow(var_X, 3) > 0.008856 {
		var_X = math.Pow(var_X, 3)
	} else {
		var_X = (var_X - 16/116) / 7.787
	}
	if math.Pow(var_Z, 3) > 0.008856 {
		var_Z = math.Pow(var_Z, 3)
	} else {
		var_Z = (var_Z - 16/116) / 7.787
	}

	X = var_X * ReferenceX
	Y = var_Y * ReferenceY
	Z = var_Z * ReferenceZ

	return X, Y, Z
}

func RGB2HSL(R, G, B float64) (H, S, L float64) {
	var_R := (R / 255)
	var_G := (G / 255)
	var_B := (B / 255)

	var_Min := math.Min(var_R, math.Min(var_G, var_B)) //Min. value of RGB
	var_Max := math.Max(var_R, math.Max(var_G, var_B)) //Max. value of RGB
	del_Max := var_Max - var_Min                       //Delta RGB value

	L = (var_Max + var_Min) / 2

	if del_Max == 0 { //This is a gray, no chroma...
		H = 0
		S = 0
	} else { //Chromatic data...
		if L < 0.5 {
			S = del_Max / (var_Max + var_Min)
		} else {
			S = del_Max / (2 - var_Max - var_Min)
		}

		del_R := (((var_Max - var_R) / 6) + (del_Max / 2)) / del_Max
		del_G := (((var_Max - var_G) / 6) + (del_Max / 2)) / del_Max
		del_B := (((var_Max - var_B) / 6) + (del_Max / 2)) / del_Max

		if var_R == var_Max {
			H = del_B - del_G
		} else if var_G == var_Max {
			H = (1 / 3) + del_R - del_B
		} else if var_B == var_Max {
			H = (2 / 3) + del_G - del_R
		}

		if H < 0 {
			H += 1
		}
		if H > 1 {
			H -= 1
		}
	}
	return H, S, L
}

func Hue2RGB(v1, v2, vH float64) float64 { //Function Hue_2_RGB
	if vH < 0 {
		vH += 1
	}
	if vH > 1 {
		vH -= 1
	}
	if (6 * vH) < 1 {
		return v1 + (v2-v1)*6*vH
	}
	if (2 * vH) < 1 {
		return v2
	}
	if (3 * vH) < 2 {
		return v1 + (v2-v1)*((2/3)-vH)*6
	}
	return v1
}

func HSL2RGB(H, S, L float64) (R, G, B float64) {
	var_1 := float64(0)
	var_2 := float64(0)
	if S == 0 {

		R = L * 255
		G = L * 255
		B = L * 255
	} else {
		if (L < 0.5) {
			var_2 = L * (1 + S)
		} else {
			var_2 = (L + S) - (S * L)
		}

		var_1 = 2*L - var_2

		R = 255 * Hue2RGB(var_1, var_2, H+0.333333333)
		G = 255 * Hue2RGB(var_1, var_2, H)
		B = 255 * Hue2RGB(var_1, var_2, H-0.333333333)
	}
	return R, G, B
}

func RGB2HSV(R, G, B float64) (H, S, V float64) {
	var_R := (R / 255)
	var_G := (G / 255)
	var_B := (B / 255)

	var_Min := math.Min(var_R, math.Min(var_G, var_B)) //Min. value of RGB
	var_Max := math.Max(var_R, math.Max(var_G, var_B)) //Max. value of RGB
	del_Max := var_Max - var_Min                       //Delta RGB value

	V = var_Max

	if del_Max == 0 { //This is a gray, no chroma...

		H = 0
		S = 0
	} else { //Chromatic data...

		S = del_Max / var_Max

		del_R := (((var_Max - var_R) / 6) + (del_Max / 2)) / del_Max
		del_G := (((var_Max - var_G) / 6) + (del_Max / 2)) / del_Max
		del_B := (((var_Max - var_B) / 6) + (del_Max / 2)) / del_Max

		if var_R == var_Max {
			H = del_B - del_G
		} else if var_G == var_Max {
			H = (1 / 3) + del_R - del_B
		} else if var_B == var_Max {
			H = (2 / 3) + del_G - del_R
		}

		if H < 0 {
			H += 1
		}
		if H > 1 {
			H -= 1
		}
	}
	return H, S, V
}

func GetMainColor(img image.Image, rect image.Rectangle) color.Color {
	threshold := float64(0.1)
	sumHue := float64(0)
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			c := img.At(x, y)
			h, _, _ := RGB2HSV(Color2StandardRGB(c))
			sumHue += h
		}
	}
	avgHue := sumHue / float64((rect.Max.X-rect.Min.X)*(rect.Max.Y-rect.Min.Y))
	colors := make([]color.Color, 0, 0)
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			c := img.At(x, y)
			h, _, _ := RGB2HSV(Color2StandardRGB(c))
			if math.Abs(h-avgHue) > threshold {
				colors = append(colors, c)
			}
		}
	}
	if len(colors) == 0 {
		return color.Black
	}
	R := uint32(0)
	G := uint32(0)
	B := uint32(0)
	for _, v := range colors {
		r, g, b, _ := v.RGBA()
		R += r
		G += g
		B += b
	}

	R = R / uint32(len(colors))
	G = G / uint32(len(colors))
	B = B / uint32(len(colors))
	return &color.RGBA{uint8(R), uint8(G), uint8(B), 255}
}

func GetMainColor2(img image.Image, rect image.Rectangle) color.Color {
	R := uint32(0)
	G := uint32(0)
	B := uint32(0)
	count := uint32(0)
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			c := img.At(x, y)
			r, g, b := Color2StandardRGB(c)
			R += uint32(r)
			G += uint32(g)
			B += uint32(b)
			count++
		}
	}
	R = R / count
	G = G / count
	B = B / count
	return &color.RGBA{uint8(R), uint8(G), uint8(B), 255}
}

func GetMainColor3(img image.Image, rect image.Rectangle) color.Color {
	sumH := float64(0)
	sumS := float64(0)
	sumL := float64(0)
	count := float64(0)
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			c := img.At(x, y)
			h, s, l := RGB2HSL(Color2StandardRGB(c))
			sumH += h
			sumL += l
			sumS += s
			count += 1
		}
	}
	r, g, b := HSL2RGB(sumH/count, sumS/count, sumL/count)
	return &color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}

var dirStr = `C:\Users\liutp\Pictures\`

func main() {
	f, e := os.Open(dirStr + `微信图片_20181010175044.png`)
	if e != nil {
		fmt.Println(e)
	}
	img, _, e := image.Decode(f)
	if e != nil {
		fmt.Println(e)
	}

	f.Close()
	c := GetMainColor2(img, img.Bounds())
	fmt.Println(c)

	r := image.Rect(0, 0, 100, 100)
	m := image.NewRGBA(r)
	for x := 0; x < r.Max.X; x++ {
		for y := 0; y < r.Max.Y; y++ {
			m.Set(x, y, c)
		}
	}

	choice := 1

	if choice == 1 {
		dstImg := image.NewRGBA(img.Bounds())
		//draw.DrawMask(dstImg, img.Bounds(), img, image.Pt(0, 0), m, image.Pt(0, 0), draw.Src)
		draw.Draw(dstImg, img.Bounds(), img, image.Pt(0, 0), draw.Src)
		draw.Draw(dstImg, m.Bounds(), m, image.Pt(0, 0), draw.Src)
		out, _ := os.Create(dirStr + `1.png`)
		png.Encode(out, dstImg)
		out.Close()
	}
	if choice == 2 {
		dstImg := resize.Resize(1, 1, img, resize.Lanczos3)
		out, _ := os.Create(dirStr + `1.png`)
		png.Encode(out, dstImg)
		out.Close()
	}

	if choice == 3 {
		scale := 10
		mX := img.Bounds().Max.X / scale
		mY := img.Bounds().Max.Y / scale
		dstImg := image.NewRGBA(img.Bounds())
		for i := 0; i < mX; i++ {
			for j := 0; j < mY; j++ {
				r := image.Rect(i*scale, j*scale, (i+1)*scale, (j+1)*scale)
				c := GetMainColor2(img, r)

				for x := 0; x < scale; x++ {
					for y := 0; y < scale; y++ {
						dstImg.Set(i*scale+x, j*scale+y, c)
					}
				}
			}
		}
		out, _ := os.Create(`D:\Users\liutianpeng\Pictures\风格化的图片\1.png`)
		png.Encode(out, dstImg)
		out.Close()
	}
}
