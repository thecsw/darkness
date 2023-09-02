package alpha

// Options is used for passing options when initiating emilia.
type Options struct {
	// DarknessConfig is the location of darkness's toml config file.
	DarknessConfig string

	// Url is a custom website url, usually used for serving from localhost.
	Url string

	// OutputExtension overrides whatever is in the file.
	OutputExtension string

	// WorkDir is the working directory of the darkness project.
	WorkDir string

	// Dev enables Url generation through local paths.
	Dev bool

	// Test enables test environment, where darkness config is not needed.
	Test bool

	// VendorGalleries dictates whether we should stub in local gallery images.
	VendorGalleries bool
}
