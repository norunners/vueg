// Package main is vueg - The Go generator for Vue templates.
// The vueg command expects to be called by go generate.
// Vue template files with the .vue extension are generated into constants in Go source.
// This lets the template constant be available to use in Vue template options.
package main

import (
	"bytes"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/tdewolff/minify"
	mhtml "github.com/tdewolff/minify/html"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	gofile := os.Getenv("GOFILE")
	gopackage := os.Getenv("GOPACKAGE")

	minifier := minify.New()
	minifier.Add("text/html", &mhtml.Minifier{KeepEndTags: true})

	templates, err := filepath.Glob("*.vue")
	must(err)
	for _, template := range templates {
		filebase := strings.TrimSuffix(template, ".vue")
		filename := filebase + ".go"
		// Prevent overwriting the file calling vueg via go generate.
		if filename == gofile {
			panic(fmt.Errorf("file conflict on name: %s", filename))
		}

		// Open the Vue template file.
		file, err := os.Open(template)
		must(err)

		// Read and minify the Vue template.
		buf := bytes.NewBuffer(nil)
		err = minifier.Minify("text/html", buf, file)
		must(err)

		// Parse the Vue template file into html nodes.
		nodes := parse(buf)
		child := firstChild(nodes)

		buf.Reset()
		html.Render(buf, child)

		// Generate the Go source.
		source := jen.NewFile(gopackage)
		comment := fmt.Sprintf("The vueg command generated this source from file: %s, do not edit.", template)
		source.HeaderComment(comment)
		source.Line()
		// Ensure the name of the constant is not exported.
		name := strings.ToLower(filebase[:1]) + filebase[1:]
		source.Const().Id(name).Op("=").Lit(buf.String())

		// Write the source as a Go file.
		file, err = os.Create(filename)
		must(err)
		err = source.Render(file)
		must(err)
		err = file.Close()
		must(err)
	}
}

// parse parses the template into html nodes.
func parse(reader io.Reader) []*html.Node {
	nodes, err := html.ParseFragment(reader, &html.Node{
		Type:     html.ElementNode,
		Data:     "div",
		DataAtom: atom.Div,
	})
	must(err)
	return nodes
}

// firstChild returns the first child element in the template element.
// For example:
// <template>
//     <div>{{ Message }}</div>
// </template>
// Returns the div element.
func firstChild(nodes []*html.Node) *html.Node {
	for _, node := range nodes {
		if node.Type == html.ElementNode && node.Data == "template" {
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				if child.Type == html.ElementNode {
					return child
				}
			}
		}
	}
	must(fmt.Errorf("child element in template element not found in nodes: %v", nodes))
	return nil
}

// must panics on errors.
func must(err error) {
	if err != nil {
		panic(err)
	}
}
