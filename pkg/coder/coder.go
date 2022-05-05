package coder

import (
	"fmt"
	"image"
	"sync"

	"math"
)

const (
	RED   = 0
	GREEN = 1
	BLUE  = 2
)

type Coder struct {
	rgbBitMap [][][]uint32
	width     uint32
	height    uint32
	raport    []string
}

func Coder_createCoder(bitmap image.Image) *Coder {
	coder := &Coder{}

	//PLACEHOLDER
	fmt.Println(bitmap.Bounds().Max.Y + 1)
	fmt.Println(bitmap.Bounds().Max.X + 1)

	coder.rgbBitMap = make([][][]uint32, bitmap.Bounds().Max.Y+1)
	coder.raport = make([]string, 9)
	for i := range coder.rgbBitMap {

		coder.rgbBitMap[i] = make([][]uint32, bitmap.Bounds().Max.X+1)

		for j := range coder.rgbBitMap[i] {
			coder.rgbBitMap[i][j] = []uint32{0, 0, 0}
		}

	}

	for i := 0; i < bitmap.Bounds().Max.Y; i++ {
		for j := 0; j < bitmap.Bounds().Max.X; j++ {
			r, g, b, _ := bitmap.At(i, j).RGBA()

			coder.rgbBitMap[i+1][j+1][RED] = r / 256
			coder.rgbBitMap[i+1][j+1][BLUE] = b / 256
			coder.rgbBitMap[i+1][j+1][GREEN] = g / 256
		}

	}

	coder.width = uint32(bitmap.Bounds().Max.X)
	coder.height = uint32(bitmap.Bounds().Max.Y)
	return coder
}

func (c *Coder) Coder_run() {
	var wg sync.WaitGroup
	predicantsTypes := []func(x, y uint32) (uint32, uint32, uint32){c.type0, c.type1, c.type2, c.type3, c.type4, c.type5, c.type6, c.type7, c.type8}

	for i := range predicantsTypes {

		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()

			c.counter(predicantsTypes[i], i)
		}()
	}

	wg.Wait()
	c.printRaport()
}

func (c *Coder) printRaport() {
	fmt.Println("TYPE 8 is basic file")
	for _, s := range c.raport {
		fmt.Println(s)
	}
	fmt.Println("TYPE 8 is basic file")

}

func (c *Coder) counter(predicate func(uint32, uint32) (uint32, uint32, uint32), t int) {
	rgb := [][]uint32{make([]uint32, 256), make([]uint32, 256), make([]uint32, 256)}
	probsRed := make([]float64, 256)
	probsGreen := make([]float64, 256)
	probsBlue := make([]float64, 256)
	probs := make([]float64, 256)
	all := c.height * c.width
	for i := uint32(0); i < c.height; i++ {

		for j := uint32(0); j < c.width; j++ {
			redVal, greenVal, blueVal := predicate(j+1, i+1)
			rgb[RED][subMod256(c.rgbBitMap[i+1][j+1][RED], redVal)]++
			rgb[GREEN][subMod256(c.rgbBitMap[i+1][j+1][GREEN], greenVal)]++
			rgb[BLUE][subMod256(c.rgbBitMap[i+1][j+1][BLUE], blueVal)]++
		}
	}

	for i := range rgb[RED] {
		probsRed[i] = float64(rgb[RED][i]) / float64(all)
		probsGreen[i] = float64(rgb[GREEN][i]) / float64(all)
		probsBlue[i] = float64(rgb[BLUE][i]) / float64(all)
		probs[i] = float64(rgb[RED][i]+rgb[GREEN][i]+rgb[BLUE][i]) / float64(all*3)
	}
	HR, HG, HB, H := calcEntroptyRGB(probsRed, probsGreen, probsBlue, probs)
	raport := fmt.Sprintln("ENTROPY RED TYPE", t, "::", HR)
	raport += fmt.Sprintln("ENTROPY GREEN TYPE", t, "::", HG)
	raport += fmt.Sprintln("ENTROPY BLUE TYPE", t, "::", HB)
	raport += fmt.Sprintln("ENTROPY TYPE", t, "::", H)
	raport += fmt.Sprintln()

	c.raport[t] = raport
}

func (c *Coder) type8(x, y uint32) (uint32, uint32, uint32) {
	return 0, 0, 0
}

func (c *Coder) type0(x, y uint32) (uint32, uint32, uint32) {
	return c.rgbBitMap[y][x-1][RED], c.rgbBitMap[y][x-1][GREEN], c.rgbBitMap[y][x-1][BLUE]
}

func (c *Coder) type1(x, y uint32) (uint32, uint32, uint32) {
	return c.rgbBitMap[y-1][x][RED], c.rgbBitMap[y-1][x][GREEN], c.rgbBitMap[y-1][x][BLUE]
}

func (c *Coder) type2(x, y uint32) (uint32, uint32, uint32) {
	return c.rgbBitMap[y-1][x-1][RED], c.rgbBitMap[y-1][x-1][GREEN], c.rgbBitMap[y-1][x-1][BLUE]
}

func (c *Coder) type3(x, y uint32) (uint32, uint32, uint32) {
	R := subMod256(addMod256(c.rgbBitMap[y-1][x][RED], c.rgbBitMap[y][x-1][RED]), c.rgbBitMap[y-1][x-1][RED])
	G := subMod256(addMod256(c.rgbBitMap[y-1][x][GREEN], c.rgbBitMap[y][x-1][GREEN]), c.rgbBitMap[y-1][x-1][GREEN])
	B := subMod256(addMod256(c.rgbBitMap[y-1][x][BLUE], c.rgbBitMap[y][x-1][BLUE]), c.rgbBitMap[y-1][x-1][BLUE])

	return R, G, B
}

func (c *Coder) type4(x, y uint32) (uint32, uint32, uint32) {
	R := addMod256(c.rgbBitMap[y-1][x][RED], divMod256(subMod256(c.rgbBitMap[y][x-1][RED], c.rgbBitMap[y-1][x-1][RED]), 2))
	G := addMod256(c.rgbBitMap[y-1][x][GREEN], divMod256(subMod256(c.rgbBitMap[y][x-1][GREEN], c.rgbBitMap[y-1][x-1][GREEN]), 2))
	B := addMod256(c.rgbBitMap[y-1][x][BLUE], divMod256(subMod256(c.rgbBitMap[y][x-1][BLUE], c.rgbBitMap[y-1][x-1][BLUE]), 2))

	return R, G, B
}

func (c *Coder) type5(x, y uint32) (uint32, uint32, uint32) {
	R := addMod256(c.rgbBitMap[y][x-1][RED], divMod256(subMod256(c.rgbBitMap[y-1][x][RED], c.rgbBitMap[y-1][x-1][RED]), 2))
	G := addMod256(c.rgbBitMap[y][x-1][GREEN], divMod256(subMod256(c.rgbBitMap[y-1][x][GREEN], c.rgbBitMap[y-1][x-1][GREEN]), 2))
	B := addMod256(c.rgbBitMap[y][x-1][BLUE], divMod256(subMod256(c.rgbBitMap[y-1][x][BLUE], c.rgbBitMap[y-1][x-1][BLUE]), 2))

	return R, G, B
}

func (c *Coder) type6(x, y uint32) (uint32, uint32, uint32) {
	R := divMod256(addMod256(c.rgbBitMap[y][x-1][RED], c.rgbBitMap[y-1][x][RED]), 2)
	G := divMod256(addMod256(c.rgbBitMap[y][x-1][GREEN], c.rgbBitMap[y-1][x][GREEN]), 2)
	B := divMod256(addMod256(c.rgbBitMap[y][x-1][BLUE], c.rgbBitMap[y-1][x][BLUE]), 2)

	return R, G, B
}

func (c *Coder) type7(x, y uint32) (uint32, uint32, uint32) {
	NWR := c.rgbBitMap[y-1][x-1][RED]
	WR := c.rgbBitMap[y][x-1][RED]
	NR := c.rgbBitMap[y-1][x][RED]

	NWB := c.rgbBitMap[y-1][x-1][BLUE]
	WB := c.rgbBitMap[y][x-1][BLUE]
	NB := c.rgbBitMap[y-1][x][BLUE]

	NWG := c.rgbBitMap[y-1][x-1][GREEN]
	WG := c.rgbBitMap[y][x-1][GREEN]
	NG := c.rgbBitMap[y-1][x][GREEN]

	return type7Helper(NWR, WR, NR), type7Helper(NWB, WB, NB), type7Helper(NWG, WG, NG)
}

func type7Helper(NW, W, N uint32) uint32 {
	if NW >= min(W, N) {
		return min(W, N)
	}
	if NW <= max(W, N) {
		return max(W, N)
	}
	return subMod256(addMod256(W, N), NW)
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func addMod256(a, b uint32) uint32 {
	return (a + b) % 256
}

func subMod256(a, b uint32) uint32 {
	return (a - b) % 256
}

func divMod256(a, b uint32) uint32 {
	return (a / b) % 256
}

func calcEntroptyRGB(probsRed, probsGreen, probsBlue, probs []float64) (float64, float64, float64, float64) {

	HR := 0.0
	HG := 0.0
	HB := 0.0
	H := 0.0

	for i := 0; i < 256; i++ {
		Pxr := probsRed[i]
		Pxg := probsGreen[i]
		Pxb := probsBlue[i]
		Px := probs[i]

		if Pxr != 0.0 {
			Ixr := -math.Log2(Pxr)
			HR += Pxr * Ixr
		}

		if Pxg != 0.0 {
			Ixg := -math.Log2(Pxg)
			HG += Pxg * Ixg
		}

		if Pxb != 0.0 {
			Ixb := -math.Log2(Pxb)
			HB += Pxb * Ixb
		}

		if Px != 0.0 {
			Ix := -math.Log2(Px)
			H += Px * Ix
		}
	}

	return HR, HG, HB, H
}
