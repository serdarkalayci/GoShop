package data

import (
	"testing"
)

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name:  "Teavana",
		Price: 1.00,
		SKU:   "hda-as-fg",
	}
	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
