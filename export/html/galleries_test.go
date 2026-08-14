package html

import (
	"testing"

	"github.com/thecsw/darkness/v3/emilia/rem"
)

func TestExtractCustomVFlex(t *testing.T) {
	tests := []struct {
		line string
		want uint
	}{
		{line: "image.jpg :v-flex 40", want: 40},
		{line: "image.jpg", want: 0},
		{line: "image.jpg :v-flex 0", want: 0},
		{line: "image.jpg :v-flex 101", want: 0},
		{line: "image.jpg :v-flex nope", want: 0},
	}

	for _, test := range tests {
		if got := extractCustomVFlex(test.line); got != test.want {
			t.Errorf("extractCustomVFlex(%q) = %d, want %d", test.line, got, test.want)
		}
	}
}

func TestHasFlexBreak(t *testing.T) {
	if !hasFlexBreak("image.jpg :flex-break") {
		t.Error("expected :flex-break to be found")
	}
	if hasFlexBreak("image.jpg :flex-breakpoint") {
		t.Error("unexpected :flex-break match")
	}
}

func TestMakeGalleryColumnsVerticalFlexConsumesRemainingSpace(t *testing.T) {
	items := []rem.GalleryItem{
		{OriginalLine: "one.jpg :v-flex 40"},
		{OriginalLine: "two.jpg :v-flex 50"},
		{OriginalLine: "three.jpg"},
		{OriginalLine: "four.jpg"},
	}

	columns := makeGalleryColumns(items, 3)
	if len(columns) != 2 {
		t.Fatalf("got %d columns, want 2", len(columns))
	}
	if columns[0].width != 3 {
		t.Errorf("first column width = %d, want 3", columns[0].width)
	}
	if len(columns[0].items) != 3 {
		t.Fatalf("first column has %d rows, want 3", len(columns[0].items))
	}

	wantShares := []float64{40, 30, 30}
	for index, want := range wantShares {
		if got := columns[0].items[index].verticalShare; got != want {
			t.Errorf("first column row %d share = %g, want %g", index, got, want)
		}
	}

	if len(columns[1].items) != 1 || columns[1].items[0].verticalShare != 100 {
		t.Errorf("unannotated item should be its own full-height column: %#v", columns[1])
	}
}

func TestMakeGalleryColumnsVerticalFlexUsesPercentageOfRemainingSpace(t *testing.T) {
	items := []rem.GalleryItem{
		{OriginalLine: "one.jpg :v-flex 50"},
		{OriginalLine: "two.jpg :v-flex 50"},
		{OriginalLine: "three.jpg"},
	}

	columns := makeGalleryColumns(items, 3)
	if len(columns) != 1 || len(columns[0].items) != 3 {
		t.Fatalf("got columns %#v, want one three-row column", columns)
	}

	wantShares := []float64{50, 25, 25}
	for index, want := range wantShares {
		if got := columns[0].items[index].verticalShare; got != want {
			t.Errorf("column row %d share = %g, want %g", index, got, want)
		}
	}
}

func TestMakeGalleryColumnsFlexBreakStartsNewColumn(t *testing.T) {
	items := []rem.GalleryItem{
		{OriginalLine: "one.jpg :v-flex 50"},
		{OriginalLine: "two.jpg :v-flex 50 :flex-break"},
		{OriginalLine: "three.jpg"},
	}

	columns := makeGalleryColumns(items, 3)
	if len(columns) != 2 {
		t.Fatalf("got %d columns, want 2", len(columns))
	}
	if got := columns[0].items[0].verticalShare; got != 50 {
		t.Errorf("first column share = %g, want 50", got)
	}

	wantSecondColumn := []float64{50, 50}
	if len(columns[1].items) != len(wantSecondColumn) {
		t.Fatalf("second column has %d rows, want %d", len(columns[1].items), len(wantSecondColumn))
	}
	for index, want := range wantSecondColumn {
		if got := columns[1].items[index].verticalShare; got != want {
			t.Errorf("second column row %d share = %g, want %g", index, got, want)
		}
	}
}

func TestMakeGalleryColumnsUsesFirstItemFlexWidth(t *testing.T) {
	items := []rem.GalleryItem{
		{OriginalLine: "one.jpg :flex 2 :v-flex 40"},
		{OriginalLine: "two.jpg"},
	}

	columns := makeGalleryColumns(items, 3)
	if len(columns) != 1 {
		t.Fatalf("got %d columns, want 1", len(columns))
	}
	if columns[0].width != 2 {
		t.Errorf("column width = %d, want 2", columns[0].width)
	}
}
