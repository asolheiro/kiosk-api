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

	"github.com/godoes/printers"
	"github.com/golang/freetype/truetype"
	"github.com/oklog/ulid/v2"
	"golang.org/x/image/font"
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

func UpdateImage(rgba *image.RGBA, x, y int, title string, description string, fontZise int) error {
	fontPath := "ARIALBD.ttf"
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// Parse the font
	drawFont, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Create the font drawer
	fnt := truetype.NewFace(drawFont, &truetype.Options{
		Size:    float64(fontZise),
		DPI:     72,
		Hinting: font.HintingFull,
	})

	err = drawText(fnt, rgba, x, y, title)
	if err != nil {
		fmt.Printf("error drawing text: %v", err)
		return err
	}

	// Draw description
	err = drawText(fnt, rgba, x, y+50, description)
	if err != nil {
		fmt.Printf("error drawing text: %v", err)
		return err
	}
	return nil
}

func drawText(fnt font.Face, rgba *image.RGBA, x, y int, label string) error {
	col := color.RGBA{0, 0, 0, 255}
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(col),
		Face: fnt,
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

func DefaultPrint(path string, title string, description string, x, y, fontZise int) (string, error) {
	image, err := GetTemplateImage(path)
	if err != nil {
		return "", err
	}
	err = UpdateImage(image, x, y, title, description, fontZise)
	if err != nil {

		return "", err
	}

	if _, err := os.Stat("./temp"); os.IsNotExist(err) {
		err = os.Mkdir("./temp", os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("failed to create temp directory: %v", err)
		}
	}

	finalPath := generateNewName("v1", path)
	SaveImage(image, finalPath)

	return finalPath, nil
}

func Print(filePath string, printerName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// Read file data
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	p, err := printers.Open(printerName)
	if err != nil {
		return err
	}

	defer func() {
		_ = p.Close()
	}()
	println("Printer opened")
	err = p.StartDocument("Print Job", "RAW")
	if err != nil {
		log.Fatal(err)
	}
	println("Document started")
	err = p.StartPage()
	if err != nil {
		log.Fatal(err)
	}

	println("Page started")
	_, err = p.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	err = p.EndPage()
	if err != nil {
		log.Fatal(err)
	}

	println("Page ended")
	err = p.EndDocument()
	if err != nil {
		log.Fatal(err)
	}

	println("Document ended")
	log.Println("Print job sent successfully")
	return nil

}

func generateNewName(prefix string, path string) string {
	pathParts := strings.Split(path, string(os.PathSeparator))
	fileName := pathParts[len(pathParts)-1]
	newName := fmt.Sprintf("./temp/%s-%s-%s", prefix, ulid.Make().String(), fileName)
	return newName
}
