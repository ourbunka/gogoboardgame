package utils

import (
	"bytes"
	"embed"
	"image"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

func LoadImage(path string, UseEmbeded bool, Build string, Resources embed.FS) (*ebiten.Image, *image.Alpha, error) {
	var imageByte []byte
	var err error
	if path == "" {
		if UseEmbeded {
			log.Println("====== ", Build, "==========")
			imageByte, err = Resources.ReadFile("resources/board/board_center.png")
		} else {
			log.Println("====== ", Build, "==========")
			imageByte, err = os.ReadFile("./resources/board/board_center.png")
		}
	} else {
		if UseEmbeded {
			log.Println("====== ", Build, "==========")
			path = strings.TrimPrefix(path, "./")
			imageByte, err = Resources.ReadFile(path)
		} else {
			log.Println("====== ", Build, "==========")
			imageByte, err = os.ReadFile(path)
		}

	}
	if err != nil {
		println("Failed to load test image, REASON :" + err.Error())
		return nil, nil, err
	}
	img, fileFormat, err := image.Decode(bytes.NewReader(imageByte))
	if err != nil {
		println(path)
		println("Failed to decode image : " + err.Error() + "  file format : " + fileFormat)
		return nil, nil, err
	}
	//println(fileFormat)
	ebitenImage := ebiten.NewImageFromImage(img)
	bound := img.Bounds()
	ebitenAlphaImage := image.NewAlpha(bound)

	for j := bound.Min.Y; j < bound.Max.Y; j++ {
		for i := bound.Min.X; i < bound.Max.X; i++ {
			ebitenAlphaImage.Set(i, j, img.At(i, j))
		}
	}
	return ebitenImage, ebitenAlphaImage, nil
}
