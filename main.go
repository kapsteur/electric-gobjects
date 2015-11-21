package main

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	dpi    = 72.0
	size   = 20.0
	height = 1920
	width  = 1080
)

type Bit struct {
	X   int
	Y   int
	Val string
}

func main() {
	now := fmt.Sprintf("%d", time.Now().Unix())
	os.Mkdir(now, 0777)

	//Load font file
	b, err := ioutil.ReadFile("Arial.ttf")
	if err != nil {
		log.Println(err)
		return
	}

	//Parse font file
	ttf, err := truetype.Parse(b)
	if err != nil {
		log.Println(err)
		return
	}

	//Create Font.Face from font
	face := truetype.NewFace(ttf, &truetype.Options{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})

	//Create initial matrix
	m := make([]Bit, 0)
	for y := 17; y <= height; y = y + int(size) {
		for x := 4; x <= width; x = x + int(size) {
			val := "1"
			rand.Seed(time.Now().UnixNano())
			if rand.Intn(2) == 1 {
				val = "0"
			}

			m = append(m, Bit{X: x, Y: y, Val: val})
		}
	}

	//Create next matrix
	mm := make([][]Bit, 0)
	mm = append(mm, m)

	//Get a random list
	rand.Seed(time.Now().UnixNano())
	randBit := rand.Perm(len(m))

	//Create each next matrix
	for i := 1; i < len(randBit); i++ {

		tmpM := make([]Bit, len(m))

		//Copy last matrix
		copy(tmpM, mm[i-1])

		//Update one bit
		tmpBit := tmpM[randBit[i]]
		if tmpBit.Val == "0" {
			tmpBit.Val = "1"
		} else {
			tmpBit.Val = "0"
		}

		//Save bit in matrix
		tmpM[randBit[i]] = tmpBit

		//Append new matrix
		mm = append(mm, tmpM)

	}

	//Create all matrix images
	for i := 0; i < len(randBit); i++ {

		if i%100 == 0 {
			log.Printf("i:%d/%d", i, len(randBit))
		}

		//Create template images
		src := image.NewNRGBA(image.Rect(0, 0, width, height))
		dst := image.NewNRGBA(image.Rect(0, 0, width, height))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				src.Set(x, y, color.NRGBA{uint8(0), uint8(0), uint8(0), 255})
				dst.Set(x, y, color.NRGBA{uint8(255), uint8(255), uint8(255), 255})
			}
		}

		//Create file
		f, err := os.OpenFile(fmt.Sprintf("%s/%d.png", now, i), os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("OpenFile - Err:%s", err)
		}

		//Draw the bit
		for _, val := range mm[i] {
			p := fixed.P(val.X, val.Y)
			d := font.Drawer{Dst: dst, Src: src, Face: face, Dot: p}
			d.DrawString(val.Val)
		}

		//Save the file
		if err = png.Encode(f, dst); err != nil {
			log.Fatalf("Encode0 - Err:%s", err)
		}

	}
}
