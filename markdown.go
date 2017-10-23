package main

import (
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/lunny/html2md"

)

// convert a web url to the local markdown url
func urlToMarkdownUri(base url.URL, uri url.URL) (name string, path string) {
	relUrl := strings.TrimSuffix(strings.TrimPrefix(uri.Path, base.Path), "/")
  orig_name := "";
	if matches := regexp.MustCompile(`[^\/]([\w\d-_]+(\.[\w]+)?)?$`).FindAllString(relUrl, -1); len(matches) > 0 {
    orig_name = matches[0];

    if len(orig_name) > 0 && orig_name != "/" {
      name = strings.Split(orig_name, ".")[0] + ".md"
    }

  }

  if name == "" {
    name = "index.md"
	}

	path = strings.TrimSuffix(relUrl, orig_name)

  if len(path) == 0 { path = "/" }

	return
}

// create a markdown file for the given url
// by convertign the body to markdown
func saveMarkdownFile(config _config, uri url.URL, body string) {

	name, path := urlToMarkdownUri(config.rooturl, uri)
	directory := config.storagedir + "/" + path
	filename := directory + "/" + name
	filename = strings.Replace(filename, "//", "/", -1)

	config.fs.MkdirAll(directory, os.ModePerm)

	file, err := config.fs.Create(filename)
	defer func() { file.Close() }()

	if err != nil {
		panic(err)
	}

	file.WriteString(regexp.MustCompile(`[\n\s]*\n[\n\s]*`).ReplaceAllString(html2md.Convert(body), "\n"))
}
