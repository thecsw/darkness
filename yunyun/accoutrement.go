package yunyun

import "strings"

// Accoutrement are settings of a page.
type Accoutrement struct {
	// To prevent unkeyed literars.
	_ struct{}
	// The image for preview
	Preview string
	// PreviewWidth and PreviewHeight are the dimensions of the preview image.
	PreviewWidth  string
	PreviewHeight string
	// PreviewGenerate tells us whether reze needs to create a preview image.
	PreviewGenerate AccoutrementFlip

	// ExcludeHtmlHeadContains is a list of strings that we should match
	// against page's scripts before injecting them into the page.
	// Useful if you want to disable specific scripts on select pages.
	ExcludeHtmlHeadContains ExcludeHtmlHeadContains
	// Draft will prevent rss from showing the page.
	Draft AccoutrementFlip
	// Tomb enables/disables tomb on a page.
	Tomb AccoutrementFlip
	// AuthorImage enables/disable author's header image.
	AuthorImage AccoutrementFlip
	// Math enables/disables math rendering (overrides auto-discovery).
	Math AccoutrementFlip
	// Toc enables/disables table of contents
	Toc AccoutrementFlip
}

// ExcludeHtmlHeadContains is a type to store excluded keywords for html head.
type ExcludeHtmlHeadContains []string

// ShouldExclude returns true if the passed html head element should be excluded.
func (e ExcludeHtmlHeadContains) ShouldExclude(what string) bool {
	for _, excluded := range e {
		if strings.Contains(what, excluded) {
			return true
		}
	}
	return false
}

// ShouldKeep returns true if the passed html head element should be included.
func (e ExcludeHtmlHeadContains) ShouldKeep(what string) bool {
	return !e.ShouldExclude(what)
}

// AccoutrementFlip holds the state of the flag: default, set, unset.
type AccoutrementFlip uint8

const (
	// AccoutrementDefault means that the value is default.
	AccoutrementDefault AccoutrementFlip = iota
	// AccoutrementEnabled means user forced enable.
	AccoutrementEnabled
	// AccoutrementDisabled means user forced disable
	AccoutrementDisabled
)

// IsDefault returns true if the flag was left with no changes to default.
func (a *AccoutrementFlip) IsDefault() bool {
	return *a == AccoutrementDefault
}

// IsEnabled returns true if the flag was manually set.
func (a *AccoutrementFlip) IsEnabled() bool {
	return *a == AccoutrementEnabled
}

// IsDisabled returns true if the flag was manually unset.
func (a *AccoutrementFlip) IsDisabled() bool {
	return *a == AccoutrementDisabled
}

// IsEnabledOrUnchanged returns true if the flag was not set
// or it was enabled.
func (a *AccoutrementFlip) IsEnabledOrUnchanged() bool {
	return a.IsEnabled() || a.IsDefault()
}

// IsDisabledOrDefault returns true if the flag was not set
// or it was disabled.
func (a *AccoutrementFlip) IsDisabledOrDefault() bool {
	return a.IsDisabled() || a.IsDefault()
}

// Enable turns the flag on.
func (a *AccoutrementFlip) Enable() {
	*a = AccoutrementEnabled
}

// Disable turns the flag off.
func (a *AccoutrementFlip) Disable() {
	*a = AccoutrementDisabled
}
