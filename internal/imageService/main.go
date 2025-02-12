package imageService

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/godoes/printers"
	"github.com/golang/freetype/truetype"
	"github.com/oklog/ulid/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

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
		println(err.Error())
		return err
	}
	// Parse the font
	drawFont, err := truetype.Parse(fontBytes)
	if err != nil {
		println(err.Error())
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
	err = drawText(fnt, rgba, x, y+20, description)
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

func SaveImage(rgba *image.RGBA, outputPath string) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		println("failed to create output file:" + err.Error())
		return err
	}
	defer outFile.Close()

	err = png.Encode(outFile, rgba)
	if err != nil {
		println("failed to encode image:" + err.Error())
		return err
	}
	log.Printf("image saved successfully to %s", outputPath)
	return nil
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

	finalPath := generateNewName("v1", path)
	err = SaveImage(image, finalPath)
	if err != nil {
		return "", err
	}

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
		println(err.Error())
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
		println(err.Error())
		return err
	}
	println("Document started")
	err = p.StartPage()
	if err != nil {
		println(err.Error())
		return err
	}

	println("Page started")
	_, err = p.Write(data)
	if err != nil {
		println(err.Error())
		return err
	}

	err = p.EndPage()
	if err != nil {
		println(err.Error())
		return err
	}

	println("Page ended")
	err = p.EndDocument()
	if err != nil {
		println(err.Error())
		return err
	}

	println("Document ended")
	log.Println("Print job sent successfully")
	return nil

}

func generateNewName(prefix string, path string) string {
	homeDir := UserHomeDir()
	dbDir := homeDir + "\\AppData\\Local\\Kiosk"
	pathParts := strings.Split(path, string(os.PathSeparator))
	fileName := pathParts[len(pathParts)-1]
	newName := fmt.Sprintf(dbDir+"/%s-%s-%s", prefix, ulid.Make().String(), fileName)
	return newName
}
