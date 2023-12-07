package parse

import (
	"io"

	d2 "github.com/FurqanSoftware/goldmark-d2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"

	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func Markdown(contents []byte, w io.Writer) error {
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			extension.GFM,
			// &pikchr.Extender{},
			highlighting.NewHighlighting(
				highlighting.WithStyle("hrdark"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
			&d2.Extender{
				Layout:  d2dagrelayout.DefaultLayout,
				ThemeID: &d2themescatalog.NeutralDefault.ID,
				// Sketch:  true,
			},
		),
	)

	return md.Convert(contents, w)
}
