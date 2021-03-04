package bimage

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/snail007/resize"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	typeJPEG uint = iota + 1
	typePNG
	typeGIF
	typeBMP
)

var (
	//go:embed simkai.ttf
	fontBytes []byte
	font      *truetype.Font
)

func init() {
	font, _ = freetype.ParseFont(fontBytes)
}

func CompressTo(src, dst string, level, width, height uint) (err error) {
	src, _ = filepath.Abs(src)
	dst, _ = filepath.Abs(dst)
	data, err := CompressFile(src, level, width, height)
	if err != nil {
		return
	}
	return ioutil.WriteFile(dst, data, 0755)
}

func CompressFile(file string, level, width, height uint) (data []byte, err error) {
	srcData, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	return Compress(srcData, level, width, height)
}

func Compress(src []byte, level, width, height uint) (data []byte, err error) {
	_, img, gifObj, err := getSupportedImage(src)
	if err != nil {
		return
	}
	if level > 10 || level == 0 {
		err = fmt.Errorf("error level, must be 1-10")
		return
	}
	imgWidth := uint(0)
	imgHeight := uint(0)
	if gifObj != nil {
		imgWidth = uint(gifObj.Config.Width)
		imgHeight = uint(gifObj.Config.Height)
	} else {
		imgWidth = uint(img.Bounds().Dx())
		imgHeight = uint(img.Bounds().Dy())
	}
	if width > imgWidth {
		width = imgWidth
	}
	if height > imgHeight {
		height = imgHeight
	}

	buf := new(bytes.Buffer)
	if img == nil && gifObj != nil {
		newGif := &gif.GIF{
			Image:           nil,
			Delay:           gifObj.Delay,
			LoopCount:       0,
			Disposal:        gifObj.Disposal,
			Config:          image.Config{},
			BackgroundIndex: gifObj.BackgroundIndex,
		}
		b := image.Rect(0, 0, gifObj.Config.Width, gifObj.Config.Height)
		for _, v := range gifObj.Image {
			newImg := resize.Resize(width, height, v, resize.NearestNeighbor)
			bufGif := new(bytes.Buffer)
			err = jpeg.Encode(bufGif, newImg, &jpeg.Options{
				Quality: 100 / 10 * int(10-level),
			})
			if err != nil {
				return
			}
			newGifImg, _ := jpeg.Decode(bufGif)
			p := image.NewPaletted(b, palette.Plan9)
			draw.Draw(p, b, newGifImg, image.Point{}, draw.Over)
			newGif.Image = append(newGif.Image, p)
		}
		err = gif.EncodeAll(buf, newGif)
	} else {
		newImg := resize.Resize(width, height, img, resize.NearestNeighbor)
		err = jpeg.Encode(buf, newImg, &jpeg.Options{
			Quality: 100 / 10 * int(10-level),
		})
	}
	if err != nil {
		return
	}
	data = buf.Bytes()
	return
}

func IsSupported(file string) bool {
	src, err := os.Open(file)
	if err != nil {
		return false
	}
	return IsSupportedByFile(src)
}

func IsSupportedByFile(src *os.File) bool {
	info, _ := src.Stat()
	data := make([]byte, info.Size())
	src.Read(data)
	return IsSupportedByBytes(data)
}

func IsSupportedByFormFile(file *multipart.FileHeader) bool {
	src, err := file.Open()
	if err != nil {
		return false
	}
	defer src.Close()
	data := make([]byte, file.Size)
	src.Read(data)
	return IsSupportedByBytes(data)
}

func IsSupportedByBytes(data []byte) bool {
	_, _, _, err := getSupportedImage(data)
	return err == nil
}
func TextMaskByFile(imageSrc, txt string) (err error) {
	imageBytes, err := ioutil.ReadFile(imageSrc)
	if err != nil {
		return
	}
	bs, err := TextMask(imageBytes, txt)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(imageSrc, bs, 0644)
	return
}
func TextMask(imageBytes []byte, txt string) (bs []byte, err error) {
	_, imgSrc, _, err := getSupportedImage(imageBytes)
	if err != nil {
		return
	}
	imgRGB := image.NewNRGBA(imgSrc.Bounds())
	for y := 0; y < imgSrc.Bounds().Dy(); y++ {
		for x := 0; x < imgSrc.Bounds().Dx(); x++ {
			imgRGB.Set(x, y, imgSrc.At(x, y))
		}
	}
	ctx := freetype.NewContext()
	ctx.SetDPI(100)
	ctx.SetFont(font)
	ctx.SetFontSize(12)
	ctx.SetClip(imgSrc.Bounds())
	ctx.SetDst(imgRGB)
	ctx.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 255, B: 255, A: 255}))
	pt := freetype.Pt(10, imgSrc.Bounds().Dy()-10)
	ctx.DrawString(txt, pt)

	ctx.SetSrc(image.NewUniform(color.RGBA{R: 0, G: 0, B: 0, A: 190}))
	pt = freetype.Pt(10, 20)
	ctx.DrawString(txt, pt)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, imgRGB, &jpeg.Options{Quality: 100})
	if err != nil {
		return
	}
	return buf.Bytes(), nil
}

func getSupportedImage(data []byte) (typ uint, img image.Image, gifObj *gif.GIF, err error) {
	l := len(data)
	if l > 512 {
		l = 512
	}
	d := data[:l]
	switch http.DetectContentType(d) {
	case "image/bmp":
		typ = typeBMP
		img, err = bmp.Decode(bytes.NewReader(data))
	case "image/jpeg":
		typ = typeJPEG
		img, err = jpeg.Decode(bytes.NewReader(data))
	case "image/gif":
		gifObj, err = gif.DecodeAll(bytes.NewReader(data))
		if err == nil {
			if len(gifObj.Image) == 1 {
				typ = typeGIF
				img = gifObj.Image[0]
			} else {
				err = fmt.Errorf("image format not supported")
			}
		}
	case "image/png":
		typ = typePNG
		img, err = png.Decode(bytes.NewReader(data))
	default:
		err = fmt.Errorf("image format not supported")
	}
	return
}
