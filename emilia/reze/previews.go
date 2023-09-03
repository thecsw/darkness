package reze

import (
	"bytes"
	"fmt"
	"golang.org/x/image/webp"
	"image/png"
	"io"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/thecsw/rei"
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

var (
	previewAvatarGenerateOnce sync.Once
)

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

	previewAvatarGenerateOnce.Do(func() {
		target := rei.Must(os.CreateTemp("", "reze_page_preview.png"))
		defer func(target *os.File) {
			err := target.Close()
			if err != nil {
				logger.Error("closing temporary file", "loc", target.Name(), "err", err)
			}
		}(target)
		p.avatarProperlySizedFile = target.Name()
		avatarFileReadable, shouldDelete := convertToReadable(AvatarFile)
		if shouldDelete {
			p.avatarReadableConverted = avatarFileReadable
		}
		avatarOriginal := rei.Must(os.Open(avatarFileReadable))
		defer func(avatarOriginal *os.File) {
			err := avatarOriginal.Close()
			if err != nil {
				logger.Error("closing avatar file", "loc", AvatarFile, "err", err)
			}
		}(avatarOriginal)
		originalDecoded := rei.Must(imaging.Decode(avatarOriginal, imaging.AutoOrientation(false)))
		resized := imaging.Resize(originalDecoded, p.calculateAvatarSize(), 0, imaging.Lanczos)
		rei.Try(imaging.Encode(target, resized, imaging.PNG))
	})

	// Rescale the avatar if needed.
	return p
}

// Generate generates a preview image and returns it as an io.Reader.
func (p PreviewGenerator) Generate(
	Title, Name, Time string,
) (io.Reader, error) {
	// Create the image context.
	dc := gg.NewContext(int(p.width), int(p.height))
	// Set the background color.
	dc.SetHexColor(p.backgroundColor)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	// Calculate the title offsets.
	titleFontSize := p.width * 0.115
	titleOffsetX := p.width * 0.125
	titleOffsetY := p.height * 0.444
	titleWidth := p.width * 0.748
	titleLineSpacing := 1.4
	titleAlign := gg.Align(2)

	// Calculate the website card offsets.
	websiteCardOffsetX := p.width * 0.1418439
	websiteCardOffsetY := p.height * 0.1604278

	// Calculate the avatar offsets.
	websiteCardAvatarSize := p.width * 0.1258865 // should be around 284px
	websiteCardAvatarOffsetX := websiteCardOffsetX + 0
	websiteCardAvatarOffsetY := websiteCardOffsetY + 0

	// Calculate the card title and time offsets.
	websiteCardAvatarToTextOffsetX := p.width * 0.015957
	websiteCardTitleOffsetDiffY := p.height * 0.08622
	websiteCartTimeOffsetDiffY := p.height * 0.1504010

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
	dc.DrawStringWrapped(Title, titleOffsetX, titleOffsetY, 0, 0, titleWidth, titleLineSpacing, titleAlign)

	// Draw the website card.
	rei.Try(dc.LoadFontFace(p.websiteNameFont, websiteCardTitleSize))
	dc.DrawStringAnchored(Name, websiteCardTitleOffsetX, websiteCardTitleOffsetY, 0, 0)

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
	return int(math.Floor(p.width * 0.120))
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
	im, err := imaging.Decode(reader)
	if err != nil {
		return fmt.Errorf("decoding image reader: %v", err)
	}
	target, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating file %s: %v", filename, err)
	}
	if err := imaging.Encode(target, im, imaging.JPEG); err != nil {
		return fmt.Errorf("encoding to jpeg: %v", err)
	}
	if err := target.Close(); err != nil {
		return fmt.Errorf("closing file %s: %v", filename, err)
	}
	return nil
}

func convertToReadable(filename string) (string, bool) {
	if strings.HasSuffix(filename, ".webp") {
		source := rei.Must(os.Open(filename))
		img := rei.Must(webp.Decode(source))

		targetFilename := strings.ReplaceAll(filename, ".webp", ".png")
		targetFile := rei.Must(os.Create(targetFilename))

		rei.Try(png.Encode(targetFile, img))
		rei.Try(targetFile.Close())
		rei.Try(source.Close())
		return targetFilename, true
	}
	return filename, false
}
