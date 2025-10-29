package widgets

import (
	"image/color"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// XMLViewer displays XML with syntax highlighting
type XMLViewer struct {
	widget.BaseWidget
	richText *widget.RichText
	scroll   *container.Scroll
}

// NewXMLViewer creates a new XML viewer with syntax highlighting
func NewXMLViewer() *XMLViewer {
	viewer := &XMLViewer{
		richText: widget.NewRichText(),
	}
	viewer.richText.Wrapping = fyne.TextWrapOff
	viewer.scroll = container.NewScroll(viewer.richText)
	viewer.ExtendBaseWidget(viewer)
	return viewer
}

// SetXML sets the XML content with syntax highlighting
func (x *XMLViewer) SetXML(xml string) {
	segments := x.highlightXML(xml)
	x.richText.Segments = segments
	x.richText.Refresh()
}

// highlightXML applies syntax highlighting to XML
func (x *XMLViewer) highlightXML(xml string) []widget.RichTextSegment {
	var segments []widget.RichTextSegment

	// Define colors for syntax highlighting
	tagColor := color.NRGBA{R: 63, G: 127, B: 95, A: 255}         // Green for tags
	attrNameColor := color.NRGBA{R: 152, G: 118, B: 170, A: 255}  // Purple for attributes
	attrValueColor := color.NRGBA{R: 206, G: 145, B: 120, A: 255} // Orange for values
	commentColor := color.NRGBA{R: 106, G: 153, B: 85, A: 255}    // Green for comments
	textColor := color.NRGBA{R: 220, G: 220, B: 220, A: 255}      // Light gray for text

	lines := strings.Split(xml, "\n")

	for i, line := range lines {
		// Replace leading spaces with non-breaking spaces to preserve indentation
		leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
		if leadingSpaces > 0 {
			line = strings.Repeat("\u00A0", leadingSpaces) + strings.TrimLeft(line, " ")
		}

		// Process line for XML highlighting
		segments = append(segments, x.highlightLine(line, tagColor, attrNameColor, attrValueColor, commentColor, textColor)...)

		// Add newline except for last line
		if i < len(lines)-1 {
			segments = append(segments, &widget.TextSegment{
				Text:  "\n",
				Style: widget.RichTextStyle{},
			})
		}
	}

	return segments
}

// highlightLine highlights a single line of XML
func (x *XMLViewer) highlightLine(line string, tagColor, attrNameColor, attrValueColor, commentColor, textColor color.Color) []widget.RichTextSegment {
	var segments []widget.RichTextSegment

	// Check for comment
	if strings.Contains(line, "<!--") {
		segments = append(segments, &widget.TextSegment{
			Text: line,
			Style: widget.RichTextStyle{
				ColorName: "",
				Inline:    true,
			},
		})
		return segments
	}

	// Pattern to match XML tags and content
	tagPattern := regexp.MustCompile(`(<[^>]+>)|([^<>]+)`)
	matches := tagPattern.FindAllStringSubmatch(line, -1)

	for _, match := range matches {
		if match[1] != "" {
			// This is a tag
			tag := match[1]
			segments = append(segments, x.highlightTag(tag, tagColor, attrNameColor, attrValueColor)...)
		} else if match[2] != "" {
			// This is text content (including whitespace)
			text := match[2]
			segments = append(segments, &widget.TextSegment{
				Text: text,
				Style: widget.RichTextStyle{
					ColorName: "foreground",
					Inline:    true,
				},
			})
		}
	}

	return segments
}

// highlightTag highlights an XML tag with attributes
func (x *XMLViewer) highlightTag(tag string, tagColor, attrNameColor, attrValueColor color.Color) []widget.RichTextSegment {
	var segments []widget.RichTextSegment

	// Check if it's a closing tag, opening tag, or self-closing tag
	if strings.HasPrefix(tag, "</") {
		// Closing tag - all in tag color
		segments = append(segments, &widget.TextSegment{
			Text: tag,
			Style: widget.RichTextStyle{
				ColorName: "success",
				Inline:    true,
			},
		})
	} else if strings.HasPrefix(tag, "<?") {
		// Processing instruction - special color
		segments = append(segments, &widget.TextSegment{
			Text: tag,
			Style: widget.RichTextStyle{
				ColorName: "warning",
				Inline:    true,
			},
		})
	} else {
		// Opening or self-closing tag - highlight tag name and attributes
		// Extract tag name and attributes
		attrPattern := regexp.MustCompile(`^(<[a-zA-Z0-9_:-]+)(.*?)(\/?>)$`)
		match := attrPattern.FindStringSubmatch(tag)

		if match != nil {
			// Tag opening bracket and name
			segments = append(segments, &widget.TextSegment{
				Text: match[1],
				Style: widget.RichTextStyle{
					ColorName: "success",
					Inline:    true,
				},
			})

			// Attributes
			if match[2] != "" {
				attrSegs := x.highlightAttributes(match[2], attrNameColor, attrValueColor)
				segments = append(segments, attrSegs...)
			}

			// Closing bracket
			segments = append(segments, &widget.TextSegment{
				Text: match[3],
				Style: widget.RichTextStyle{
					ColorName: "success",
					Inline:    true,
				},
			})
		} else {
			// Fallback - entire tag in tag color
			segments = append(segments, &widget.TextSegment{
				Text: tag,
				Style: widget.RichTextStyle{
					ColorName: "success",
					Inline:    true,
				},
			})
		}
	}

	return segments
}

// highlightAttributes highlights XML attributes
func (x *XMLViewer) highlightAttributes(attrs string, attrNameColor, attrValueColor color.Color) []widget.RichTextSegment {
	var segments []widget.RichTextSegment

	// Pattern to match attribute="value"
	attrPattern := regexp.MustCompile(`(\s+)([a-zA-Z0-9_:-]+)(=)("([^"]*)"|'([^']*)')`)
	lastIndex := 0

	matches := attrPattern.FindAllStringSubmatchIndex(attrs, -1)

	for _, match := range matches {
		// Add any text before this match
		if match[0] > lastIndex {
			segments = append(segments, &widget.TextSegment{
				Text: attrs[lastIndex:match[0]],
				Style: widget.RichTextStyle{
					Inline: true,
				},
			})
		}

		// Whitespace
		segments = append(segments, &widget.TextSegment{
			Text: attrs[match[2]:match[3]],
			Style: widget.RichTextStyle{
				Inline: true,
			},
		})

		// Attribute name
		segments = append(segments, &widget.TextSegment{
			Text: attrs[match[4]:match[5]],
			Style: widget.RichTextStyle{
				ColorName: "primary",
				Inline:    true,
			},
		})

		// Equals sign
		segments = append(segments, &widget.TextSegment{
			Text: attrs[match[6]:match[7]],
			Style: widget.RichTextStyle{
				Inline: true,
			},
		})

		// Attribute value (with quotes)
		segments = append(segments, &widget.TextSegment{
			Text: attrs[match[8]:match[9]],
			Style: widget.RichTextStyle{
				ColorName: "warning",
				Inline:    true,
			},
		})

		lastIndex = match[1]
	}

	// Add any remaining text
	if lastIndex < len(attrs) {
		segments = append(segments, &widget.TextSegment{
			Text: attrs[lastIndex:],
			Style: widget.RichTextStyle{
				Inline: true,
			},
		})
	}

	return segments
}

// CreateRenderer implements fyne.Widget
func (x *XMLViewer) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(x.scroll)
}

// GetText returns the current text content
func (x *XMLViewer) GetText() string {
	var text strings.Builder
	for _, seg := range x.richText.Segments {
		if textSeg, ok := seg.(*widget.TextSegment); ok {
			text.WriteString(textSeg.Text)
		}
	}
	return text.String()
}
