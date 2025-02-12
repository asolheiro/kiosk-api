package imageService

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/oklog/ulid/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func maxCharsThatFit(text string, face font.Face, maxWidth int) int {
	d := &font.Drawer{
		Face: face,
	}

	for i := 1; i <= len(text); i++ {
		width := d.MeasureString(text[:i])
		if width.Ceil() > maxWidth {
			return i - 1
		}
	}
	return len(text)
}

func loadFont(fontPath string, fontSize float64) (font.Face, error) {
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}

	parsedFont, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	return face, nil
}

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

func parseTextToPrinter(text string, face font.Face, maxWidth int) string {
	maxCharLimiter := maxCharsThatFit(text, face, maxWidth)

	return text[:maxCharLimiter]

}

func UpdateImage(rgba *image.RGBA, x, y int, title string, description string, fontZise int, maxWidth int) error {
	fontPath := "ARIALBD.ttf"
	face, err := loadFont(fontPath, float64(fontZise))

	parsedTitle := parseTextToPrinter(title, face, maxWidth)
	err = drawText(face, rgba, x, y, parsedTitle)
	if err != nil {
		fmt.Printf("error drawing text: %v", err)
		return err
	}

	// Draw description
	parsedDescription := parseTextToPrinter(description, face, maxWidth)
	err = drawText(face, rgba, x, y+20, parsedDescription)
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

func DefaultPrint(path string, title string, description string, x, y, fontZise, widthLimit int) (string, error) {
	image, err := GetTemplateImage(path)
	if err != nil {
		return "", err
	}
	err = UpdateImage(image, x, y, title, description, fontZise, widthLimit)
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

func sendToPrinter(filePath string, printerName string) error {
	// Use the `print` command on Windows to send the file to the default printer
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("print", "/D:"+"\""+printerName+"\"", filePath)
	} else {
		cmd = exec.Command("lp", "-d "+"\""+printerName+"\"", filePath)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to print: %v, output: %s", err, output)
	}

	fmt.Println("Print job sent successfully")
	return nil
}

func Print(filePath string, printerName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	file.Close()

	err = sendToPrinter(filePath, printerName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
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
