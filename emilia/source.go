package emilia

import (
	"strings"

	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

// WithSourceCodeTrimmedLeftWhitespace removes leading whitespace from source code blocks
func WithSourceCodeTrimmedLeftWhitespace() yunyun.PageOption {
	return func(page *yunyun.Page) {
		for i, content := range page.Contents.SourceCodeBlocks() {
			lines := strings.Split(content.SourceCode, "\n")
			if len(lines) < 1 {
				continue
			}
			offset := len(lines[0]) - len(strings.TrimLeft(lines[0], " "))
			for i, line := range lines {
				// It's just a whitespace line, then ignore it
				if len(strings.TrimSpace(line)) < 1 {
					continue
				}
				// if the initial offset is bigger, then abort the whole thing
				if len(line)-len(strings.TrimLeft(line, " ")) < offset {
					return
				}
				lines[i] = line[gana.Min(len(lines[i]), offset):]
			}
			page.Contents[i].SourceCode = strings.Join(lines, "\n")
		}
	}
}
