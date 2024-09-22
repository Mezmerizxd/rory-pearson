package image_convert

import (
	"testing"
)

func TestImageToIcon(t *testing.T) {
	name, err := Convert("input.png")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("\n\nIcon created: %s", name)

	path, err := GetConvertedZipPath(name, false)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("\n\nZip created: %s", path)

	_, err = GetConvertedZipPath(name, true)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("\n\nZip deleted: %s", path)
}
