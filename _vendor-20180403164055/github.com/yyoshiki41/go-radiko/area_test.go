package radiko

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestAreaID(t *testing.T) {
	areaID, err := AreaID()
	if err != nil {
		t.Errorf("Failed to get area id: %s", err)
	}
	if !(strings.HasPrefix(areaID, "JP") || areaID == "OUT") {
		t.Errorf("Invalid area id.\nAreaID: %s", areaID)
	}
}

func TestProcessSpanNode(t *testing.T) {
	const expected = "JP13"
	s := `document.write('<span class="` + expected + `">TOKYO JAPAN</span>');`

	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		t.Errorf("Parse HTML: %s", err)
	}

	areaID := processSpanNode(doc)
	if areaID != expected {
		t.Errorf(
			"Failed to process span node.\nAreaID: %s", areaID)
	}
}
