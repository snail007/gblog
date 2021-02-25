package bimage

import (
	"bytes"
	"fmt"
	"github.com/snail007/resize"
	"golang.org/x/image/bmp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

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
	_, img, err := getSupportedImage(src)
	if err != nil {
		return
	}
	if level > 10 || level == 0 {
		err = fmt.Errorf("error level, must be 1-10")
		return
	}
	w := uint(img.Bounds().Dx())
	if width > w {
		width = w
	}
	h := uint(img.Bounds().Dy())
	if height > h {
		height = h
	}
	buf := new(bytes.Buffer)
	newImg := resize.Resize(width, height, img, resize.NearestNeighbor)
	err = jpeg.Encode(buf, newImg, &jpeg.Options{
		Quality: 100 / 10 * int(10-level),
	})
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
	data := make([]byte, 512)
	src.Read(data)
	return IsSupportedByBytes(data)
}

func IsSupportedByFormFile(file *multipart.FileHeader) bool {
	src, err := file.Open()
	if err != nil {
		return false
	}
	defer src.Close()
	data := make([]byte, 512)
	src.Read(data)
	return IsSupportedByBytes(data)
}

func IsSupportedByBytes(data []byte) bool {
	_, _, err := getSupportedImage(data)
	return err == nil
}

func getSupportedImage(data []byte) (typ uint, img image.Image, err error) {
	l := len(data)
	if l > 512 {
		l = 512
	}
	d := data[:l]
	switch http.DetectContentType(d) {
	case "image/bmp":
		img, _ = bmp.Decode(bytes.NewReader(data))
	case "image/jpeg":
		img, _ = jpeg.Decode(bytes.NewReader(data))
	case "image/gif":
		gifs, e := gif.DecodeAll(bytes.NewReader(data))
		if e != nil {
			err = e
			return
		}
		if len(gifs.Image) > 1 {
			err = fmt.Errorf("image format not supported")
			return
		}
		img = gifs.Image[0]
	case "image/png":
		img, _ = png.Decode(bytes.NewReader(data))
	default:
		err = fmt.Errorf("image format not supported")
	}
	return
}
