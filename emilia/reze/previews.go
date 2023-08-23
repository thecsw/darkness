package reze

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/thecsw/rei"
)

const (
	// resizedAvatarTempFile is the name of the resized avatar file.
	resizedAvatarTempFile = "/tmp/.resized_avatar_reze.png"
)

// PreviewGenerator is a struct that generates a preview image.
type PreviewGenerator struct {
	_ struct{}

	// TitleFont is the path to the font used for the title.
	TitleFont string
	// WebsiteNameFont is the path to the font used for the website name.
	WebsiteNameFont string
	// WebsiteTimeFont is the path to the font used for the website time.
	WebsiteTimeFont string

	// Width and Height are the width and height of the preview image.
	Width  float64
	Height float64

	// BackgroundColor is the background color of the preview image.
	BackgroundColor string

	// AvatarFile is the path to the avatar file.
	AvatarFile string

	// avatarProperlySizedFile is the path to the resized avatar file.
	avatarProperlySizedFile string
}

var (
	once sync.Once
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
		TitleFont:               TitleFont,
		WebsiteNameFont:         NameFont,
		WebsiteTimeFont:         TimeFont,
		Width:                   float64(Width),
		Height:                  float64(Height),
		BackgroundColor:         BackgroundColor,
		AvatarFile:              AvatarFile,
		avatarProperlySizedFile: resizedAvatarTempFile,
	}
	once.Do(p.rescaleAvatarAsNeeded)
	// Rescale the avatar if needed.
	return p
}

// Generate generates a preview image and returns it as an io.Reader.
func (p PreviewGenerator) Generate(
	Title, Name, Time string,
) (io.Reader, error) {
	// Create the image context.
	dc := gg.NewContext(int(p.Width), int(p.Height))
	// Set the background color.
	dc.SetHexColor(p.BackgroundColor)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	// Calculate the title offsets.
	titleFontSize := p.Width * 0.115
	titleOffsetX := p.Width * 0.125
	titleOffsetY := p.Height * 0.444
	titleWidth := p.Width * 0.748
	titleLineSpacing := 1.4
	titleAlign := gg.Align(2)

	// Calculate the website card offsets.
	websiteCardOffsetX := p.Width * 0.1418439
	websiteCardOffsetY := p.Height * 0.1604278

	// Calculate the avatar offsets.
	websiteCardAvatarSize := p.Width * 0.1258865 // should be around 284px
	websiteCardAvatarOffsetX := websiteCardOffsetX + 0
	websiteCardAvatarOffsetY := websiteCardOffsetY + 0

	// Calculate the card title and time offsets.
	websiteCardAvatarToTextOffsetX := p.Width * 0.015957
	websiteCardTitleOffsetDiffY := p.Height * 0.08622
	websiteCartTimeOffsetDiffY := p.Height * 0.1504010

	// Calculate the title offset.
	websiteCardTitleSize := titleFontSize * 0.37
	websiteCardTitleOffsetX := websiteCardOffsetX + float64(websiteCardAvatarSize) + websiteCardAvatarToTextOffsetX
	websiteCardTitleOffsetY := websiteCardOffsetY + websiteCardTitleOffsetDiffY

	// Calculate the time card offset.
	websiteCardTimeSize := titleFontSize * 0.25
	websiteCardTimeOffsetX := websiteCardAvatarOffsetX + float64(websiteCardAvatarSize) + websiteCardAvatarToTextOffsetX
	websiteCardTimeOffsetY := websiteCardAvatarOffsetY + websiteCartTimeOffsetDiffY

	// Draw the title.
	rei.Try(dc.LoadFontFace(string(p.TitleFont), titleFontSize))
	dc.DrawStringWrapped(Title, titleOffsetX, titleOffsetY, 0, 0, titleWidth, titleLineSpacing, titleAlign)

	// Draw the website card.
	rei.Try(dc.LoadFontFace(string(p.WebsiteNameFont), websiteCardTitleSize))
	dc.DrawStringAnchored(Name, websiteCardTitleOffsetX, websiteCardTitleOffsetY, 0, 0)

	// Draw the timestamp of the page.
	rei.Try(dc.LoadFontFace(string(p.WebsiteTimeFont), websiteCardTimeSize))
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

// rescaleAvatarAsNeeded rescales the avatar if it's too big.
func (p *PreviewGenerator) rescaleAvatarAsNeeded() {
	avatarOriginal := rei.Must(os.Open(string(p.AvatarFile)))
	defer avatarOriginal.Close()
	originalDecoded := rei.Must(imaging.Decode(avatarOriginal, imaging.AutoOrientation(false)))
	resized := imaging.Resize(originalDecoded, p.calculateAvatarSize(), 0, imaging.Lanczos)
	rei.Try(imaging.Save(resized, p.avatarProperlySizedFile))
}

// calculateAvatarSize calculates the size of the avatar based on the width of the preview.
func (p PreviewGenerator) calculateAvatarSize() int {
	return int(math.Floor(p.Width * 0.120))
}

// Close removes the resized avatar file.
func (p PreviewGenerator) Close() error {
	return os.Remove(p.avatarProperlySizedFile)
}

// saveJpg saves a jpg image from an io.Reader.
func saveJpg(reader io.Reader, filename string) error {
	im, err := imaging.Decode(reader)
	if err != nil {
		return fmt.Errorf("decoding image reader: %v", err)
	}
	if err := imaging.Save(im, filename); err != nil {
		return fmt.Errorf("saving image to %s: %v", filename, err)
	}
	return nil
}
