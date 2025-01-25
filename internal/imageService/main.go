package imageService

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/oklog/ulid/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func GetTemplateImage(path string) (*image.RGBA, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer imgFile.Close()

	img, err := png.Decode(imgFile)
	if err != nil {
		return nil, err
	}
	rgba := image.NewRGBA(img.Bounds())

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
	return rgba, nil
}

func UpdateImage(rgba *image.RGBA, x, y int, title string, description string) error {
	// Draw title
	err := drawText(rgba, x, y, title)
	if err != nil {
		fmt.Printf("error drawing text: %v", err)
		return err
	}

	// Draw description
	err = drawText(rgba, x, y+50, description)
	if err != nil {
		fmt.Printf("error drawing text: %v", err)
		return err
	}
	return nil
}

func drawText(rgba *image.RGBA, x, y int, label string) error {
	col := color.RGBA{0, 0, 0, 255}
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
	return nil
}

func SaveImage(rgba *image.RGBA, outputPath string) {
	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, rgba)
	if err != nil {
		log.Fatalf("failed to encode image: %v", err)
	}
	log.Printf("image saved successfully to %s", outputPath)
}

func DefaultPrint(path string, title string, description string, x, y int) error {
	image, err := GetTemplateImage(path)
	if err != nil {
		return err
	}
	err = UpdateImage(image, x, y, title, description)
	if err != nil {

		return err
	}

	SaveImage(image, generateNewName("v1", path))

	return nil
}

func CustomPrint(path string, title string, description string, x, y int) error {
	im, err := gg.LoadImage(path)
	if err != nil {
		return err
	}

	dc := gg.NewContext(im.Bounds().Dx(), im.Bounds().Dy())
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("Roboto-VariableFont_wdth,wght.ttf", 14); err != nil {
		return err
	}
	dc.DrawImage(im, 0, 0)
	dc.DrawStringAnchored(title, float64(x), float64(y), 0.5, 0.5)

	dc.DrawStringAnchored(description, float64(x), float64(y+50), 0.5, 0.5)
	dc.Clip()
	dc.SavePNG(generateNewName("v2", path))
	return nil
}

func generateNewName(prefix string, path string) string {
	newName := path
	pathParts := strings.Split(path, "/")
	fileName := pathParts[len(pathParts)-1]
	newName = strings.Join(append(pathParts[:len(pathParts)-1], prefix+" "+ulid.Make().String()+"-"+fileName), "/")
	return newName
}
