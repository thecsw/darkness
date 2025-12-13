package reze

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/fogleman/gg"
	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/rei"
	"golang.org/x/image/webp"
)

// PreviewGenerator is a struct that generates a preview image.
type PreviewGenerator struct {
	_ struct{}

	// titleFont is the path to the font used for the title.
	titleFont string
	// websiteNameFont is the path to the font used for the website name.
	websiteNameFont string
	// websiteTimeFont is the path to the font used for the website time.
	websiteTimeFont string

	// width and height are the width and height of the preview image.
	width  float64
	height float64

	// backgroundColor is the background color of the preview image.
	backgroundColor string

	// avatarFile is the path to the avatar file.
	avatarFile string

	// avatarProperlySizedFile is the path to the resized avatar file.
	avatarProperlySizedFile string

	// avatarReadableConverted is a png version of the avatar to read.
	avatarReadableConverted string
}

// InitPreviewGenerator initializes a PreviewGenerator.
func InitPreviewGenerator(
	TitleFont string,
	NameFont string,
	TimeFont string,
	Width int,
	Height int,
	BackgroundColor string,
	AvatarFile string,
) PreviewGenerator {
	// Create the preview generator.
	p := PreviewGenerator{
		titleFont:       TitleFont,
		websiteNameFont: NameFont,
		websiteTimeFont: TimeFont,
		width:           float64(Width),
		height:          float64(Height),
		backgroundColor: BackgroundColor,
		avatarFile:      AvatarFile,
	}

	target := rei.Must(os.CreateTemp("", "reze_page_preview.png"))
	defer func(target *os.File) {
		err := target.Close()
		if err != nil {
			logger.Error("closing temporary file", "loc", target.Name(), "err", err)
		}
	}(target)
	p.avatarProperlySizedFile = target.Name()
	avatarFileReadable, shouldDelete := convertWebpToPNG(AvatarFile)
	if shouldDelete {
		p.avatarReadableConverted = avatarFileReadable
	}
	avatarOriginal := rei.Must(os.Open(filepath.Clean(avatarFileReadable)))
	defer func(avatarOriginal *os.File) {
		err := avatarOriginal.Close()
		if err != nil {
			logger.Error("closing avatar file", "loc", AvatarFile, "err", err)
		}
	}(avatarOriginal)
	originalDecoded, _ := rei.Must2(image.Decode(avatarOriginal))
	newWidth := p.calculateAvatarSize()
	newHeight := PreserveImageHeightRatio(originalDecoded, newWidth)
	resized := transform.Resize(originalDecoded, newWidth, newHeight, transform.Linear)
	rei.Try(imgio.PNGEncoder()(target, resized))
	// Rescale the avatar if needed.
	return p
}

const (
	defaultFg = "#2b2b2b"
)

// Generate generates a preview image and returns it as an io.Reader.
func (p PreviewGenerator) Generate(
	Title, Name, Time, ColorBg, ColorFg string,
) (io.Reader, error) {
	bg, fg := p.backgroundColor, defaultFg
	if len(ColorBg) > 0 {
		bg = ColorBg
	}
	if len(ColorFg) > 0 {
		fg = ColorFg
	}

	// Create the image context.
	dc := gg.NewContext(int(p.width), int(p.height))
	// Set the background color.
	dc.SetHexColor(bg)
	dc.Clear()

	// Set the font color.
	dc.SetHexColor(fg)

	// Let's try to dynamically size the font
	titleLength := float64(len(Title))
	titleFactor := float64(1)
	if titleLength > 37 {
		titleFactor = (37.0 / titleLength)
	}

	// Calculate the title offsets.
	titleFontSize := p.width * (0.115 * titleFactor)
	titleOffsetX := p.width * 0.125
	titleOffsetY := p.height * 0.38
	titleWidth := p.width * 0.75
	titleLineSpacing := 1.4
	titleAlign := gg.AlignRight

	// Calculate the website card offsets.
	websiteCardOffsetX := p.width * 0.14
	websiteCardOffsetY := p.height * 0.12

	// Calculate the avatar offsets.
	websiteCardAvatarSize := p.width * 0.125 // should be around 284px
	websiteCardAvatarOffsetX := websiteCardOffsetX + 0
	websiteCardAvatarOffsetY := websiteCardOffsetY + 0

	// Calculate the card title and time offsets.
	websiteCardAvatarToTextOffsetX := p.width * 0.016
	websiteCardTitleOffsetDiffY := p.height * 0.086
	websiteCartTimeOffsetDiffY := p.height * 0.15

	// Calculate the title offset.
	websiteCardTitleSize := titleFontSize * 0.37
	websiteCardTitleOffsetX := websiteCardOffsetX + websiteCardAvatarSize + websiteCardAvatarToTextOffsetX
	websiteCardTitleOffsetY := websiteCardOffsetY + websiteCardTitleOffsetDiffY

	// Calculate the time card offset.
	websiteCardTimeSize := titleFontSize * 0.25
	websiteCardTimeOffsetX := websiteCardAvatarOffsetX + websiteCardAvatarSize + websiteCardAvatarToTextOffsetX
	websiteCardTimeOffsetY := websiteCardAvatarOffsetY + websiteCartTimeOffsetDiffY

	// Draw the title.
	rei.Try(dc.LoadFontFace(p.titleFont, titleFontSize))
	dc.DrawStringWrapped(yunyun.FancyText(Title), titleOffsetX, titleOffsetY, 0, 0, titleWidth, titleLineSpacing, titleAlign)

	// Draw the website card.
	rei.Try(dc.LoadFontFace(p.websiteNameFont, websiteCardTitleSize))
	dc.DrawStringAnchored(yunyun.FancyText(Name), websiteCardTitleOffsetX, websiteCardTitleOffsetY, 0, 0)

	// Draw the timestamp of the page.
	rei.Try(dc.LoadFontFace(p.websiteTimeFont, websiteCardTimeSize))
	dc.DrawStringAnchored(Time, websiteCardTimeOffsetX, websiteCardTimeOffsetY, 0, 0)

	// Draw the avatar.
	im := rei.Must(gg.LoadPNG(p.avatarProperlySizedFile))
	dc.DrawImageAnchored(im, int(websiteCardAvatarOffsetX), int(websiteCardAvatarOffsetY), 0, 0)

	// Push the image to a buffer and return it.
	buf := new(bytes.Buffer)
	if err := dc.EncodePNG(buf); err != nil {
		return nil, fmt.Errorf("encoding png image to buffer: %v", err)
	}
	return buf, nil
}

// calculateAvatarSize calculates the size of the avatar based on the width of the preview.
func (p PreviewGenerator) calculateAvatarSize() int {
	return int(math.Floor(p.width * 0.11))
}

// Close removes the resized avatar file.
func (p PreviewGenerator) Close() error {
	if len(p.avatarReadableConverted) > 0 {
		if err := os.Remove(p.avatarReadableConverted); err != nil {
			return fmt.Errorf("removing readable avatar %s: %v", p.avatarReadableConverted, err)
		}
	}
	if err := os.Remove(p.avatarProperlySizedFile); err != nil {
		return fmt.Errorf("removing properly sized avatar %s: %v", p.avatarProperlySizedFile, err)
	}
	return nil
}

// SaveJpg saves a jpg image from an io.Reader.
func SaveJpg(reader io.Reader, filename string) error {
	im, _, err := image.Decode(reader)
	if err != nil {
		return fmt.Errorf("decoding image reader: %v", err)
	}
	target, err := os.Create(filepath.Clean(filename))
	if err != nil {
		return fmt.Errorf("creating file %s: %v", filename, err)
	}
	if err := imgio.JPEGEncoder(100)(target, im); err != nil {
		return fmt.Errorf("encoding to jpeg: %v", err)
	}
	if err := target.Close(); err != nil {
		return fmt.Errorf("closing file %s: %v", filename, err)
	}
	return nil
}

// convertWebpToPNG converts a webp image to a png image.
func convertWebpToPNG(filename string) (string, bool) {
	if strings.HasSuffix(filename, ".webp") {
		source := rei.Must(os.Open(filepath.Clean(filename)))
		img := rei.Must(webp.Decode(source))

		targetFilename := strings.ReplaceAll(filename, ".webp", ".png")
		targetFile := rei.Must(os.Create(filepath.Clean(targetFilename)))

		rei.Try(png.Encode(targetFile, img))
		rei.Try(targetFile.Close())
		rei.Try(source.Close())
		return targetFilename, true
	}
	return filename, false
}
