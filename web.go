package main

import (
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"bytes"
	"net/http"
	"net/url"
	"strings"
)

// factory to create the scrape.Matcher callback
// used to traverse the retrieved html node finding
// links to other local urls and updating them to
// point to the local links for use inside the markdown
func visitLocalLinks(base url.URL, links *[]url.URL) scrape.Matcher {
	return func(n *html.Node) bool {
		// must check for nil values
		if n.DataAtom == atom.A && n.Parent != nil && n.Parent.Parent != nil {
			for i := range n.Attr {
				if n.Attr[i].Key == "href" {
					link, err := url.Parse(n.Attr[i].Val)
					if err != nil {
						break
					}

					if !link.IsAbs() {
						link = base.ResolveReference(link)
					}

					if link.Host == base.Host && strings.HasPrefix(link.Path, base.Path) {
						// save links in array
						*links = append(*links, *link)
						// set link attribute to markdown link
						name, path := urlToMarkdownUri(base, *link)
						n.Attr[i].Val = path + name
					}
				}
			}
		}
		return false
	}
}

type NetGetter interface {
	Get(string) (*http.Response, error)
}

// abstract away the http.Get
// Method to enable testing
type RealNetGetter struct{}

// uses net/http and
// implements the interface
func (n RealNetGetter) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

// retrieve the local url uri under the url base,
// traverse the html to find links to other local
// pages and update the links to point to the
// local markdown files that will be created later
func getLink(base, uri url.URL, netGetter NetGetter) (*[]url.URL, string, error) {
	if netGetter == nil {
		netGetter = RealNetGetter{}
	}
	resp, err :=
		netGetter.Get(uri.String())

	var links []url.URL

	if err != nil {
		return &links, "", err
	}

	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)

	if err != nil {
		return &links, "", err
	}

	scrape.FindAll(root, visitLocalLinks(base, &links))

	writer := new(bytes.Buffer)
	html.Render(writer, root)

	content := writer.String()
	return &links, content, nil
}
