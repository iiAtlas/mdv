package render

import (
	"fmt"
	"regexp"
	"testing"
)

func FuzzGetMaxWidth(f *testing.F) {
	f.Add("narrow")
	f.Add("medium")
	f.Add("wide")
	f.Add("full")
	f.Add("120")
	f.Add("custom")

	numeric := regexp.MustCompile(`^\d+$`)

	f.Fuzz(func(t *testing.T, input string) {
		width := getMaxWidth(input)

		switch input {
		case "narrow":
			if width != "680px" {
				t.Fatalf("expected narrow to map to 680px, got %q", width)
			}
		case "medium":
			if width != "900px" {
				t.Fatalf("expected medium to map to 900px, got %q", width)
			}
		case "wide":
			if width != "1200px" {
				t.Fatalf("expected wide to map to 1200px, got %q", width)
			}
		case "full":
			if width != "" {
				t.Fatalf("expected full to disable width constraint, got %q", width)
			}
		default:
			if numeric.MatchString(input) {
				expected := input + "px"
				if width != expected {
					t.Fatalf("expected numeric input %q to map to %q, got %q", input, expected, width)
				}
			} else if width != "900px" {
				t.Fatalf("expected unknown input %q to fall back to 900px, got %q", input, width)
			}
		}
	})
}

func FuzzExtractBackgroundColor(f *testing.F) {
	f.Add(uint8(0), uint8(0), uint8(0))
	f.Add(uint8(255), uint8(255), uint8(255))

	f.Fuzz(func(t *testing.T, r, g, b uint8) {
		color := fmt.Sprintf("#%02x%02x%02x", r, g, b)
		css := fmt.Sprintf(".markdown-body { background-color: %s; }", color)

		if got := extractBackgroundColor(css); got != color {
			t.Fatalf("expected background color %s to be extracted, got %q", color, got)
		}
	})
}
