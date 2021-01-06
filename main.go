package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/integrii/flaggy"
)

var (
	input  string
	output string

	compact bool
	invert  bool

	red   uint32 = 127
	green uint32 = 127
	blue  uint32 = 127
	alpha uint32 = 127
)

func validColour(c color.Color) bool {
	r, g, b, a := c.RGBA()
	return r >= red && g >= green && b >= blue && a >= alpha
}

func processImage(img image.Image, out *os.File) {
	out.WriteString(fmt.Sprintf("unsigned char toersten[%v][%v] = {\n", img.Bounds().Dy(), img.Bounds().Dx()))

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		out.WriteString("\t{")
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			if valid := validColour(img.At(x, y)); valid && !invert {
				out.WriteString("1")
			} else if !valid && invert {
				out.WriteString("1")
			} else {
				out.WriteString("0")
			}

			if x < img.Bounds().Max.X-1 {
				out.WriteString(", ")
			}
		}
		out.WriteString("}")
		if y < img.Bounds().Max.Y-1 {
			out.WriteString(", ")
		}
		out.WriteString("\n")
	}

	out.WriteString("};")
}

func processCompact(img image.Image, out *os.File) {
	out.WriteString(fmt.Sprintf("unsigned char toersten[%v][%v] = {\n", img.Bounds().Dy(), img.Bounds().Dx()/8))

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		b := uint8(0)
		shift := 0
		out.WriteString("\t{")
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			b <<= 1
			shift++
			if valid := validColour(img.At(x, y)); valid && !invert {
				b++
			} else if !valid && invert {
				b++
			}

			if shift == 8 {
				out.WriteString(fmt.Sprintf("%#02x", b))
				if x < img.Bounds().Max.X-1 {
					out.WriteString(", ")
				}
				b = 0
				shift = 0
			}
		}
		out.WriteString("}")
		if y < img.Bounds().Max.Y-1 {
			out.WriteString(", ")
		}
		out.WriteString("\n")
	}

	out.WriteString("};")
}

func main() {
	flaggy.AddPositionalValue(&input, "Input", 1, true, "Input file")
	flaggy.AddPositionalValue(&output, "Output", 2, true, "Output file")

	flaggy.Bool(&compact, "c", "compact", "8 pixels are packed into one byte.")
	flaggy.Bool(&invert, "i", "invert", "Invert all bits. 0xff -> 0x00")

	flaggy.UInt32(&red, "", "red", "Red threshold. Default 127")
	flaggy.UInt32(&green, "", "green", "Green threshold. Default 127")
	flaggy.UInt32(&blue, "", "blue", "Blue threshold. Default 127")
	flaggy.UInt32(&alpha, "", "alpha", "Alpha threshold. Default 127")

	flaggy.Parse()

	file, err := os.Open(input)
	if err != nil {
		log.Fatalf("Error opening input file: %v", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		log.Fatalf("Error decoding image: %v", err)
	}

	out, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatalf("Error opening output file: %v", err)
	}
	defer file.Close()

	if compact {
		processCompact(img, out)
	} else {
		processImage(img, out)
	}
}
