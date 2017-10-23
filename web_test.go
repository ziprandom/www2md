package main

import (
  "net/http"
	"testing"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/suite"
	"strings"
	"errors"
	"io"
	"net/url"
)

type WebTestSuite struct {
    suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *WebTestSuite) SetupTest() {
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestWebTestSuite(t *testing.T) {
    suite.Run(t, new(WebTestSuite))
}

type FakeNetGetter struct {
  urls map[string]string
}

type FakeReadCloser struct {
  io.Reader
}

func (rc FakeReadCloser) Close() (error) {return nil }

func (g FakeNetGetter) Get(url string) (reader *http.Response, err error) {
  body, ok := g.urls[url]
  if ok  {
    reader = &http.Response{Body: FakeReadCloser{strings.NewReader(body)} }
  } else {
    err = errors.New("doesn't exist!")
  }
  return
}

func makeHtmlWithLinks(base string, links []string) (body string) {
  body = "<html><body>";
  for _, link := range links {
    body += "<a href=\"" + base + "/" + link + "\">a link</a>";
  }
  body +="</body></html>"
  return
}

func (suite *WebTestSuite) TestWeb() {

  urlGetter := FakeNetGetter{
    urls: map[string]string{
      "http://myurl.net/":        makeHtmlWithLinks("http://myurl.net",    []string{"/", "/home", "/about", "/posts"}),
      "http://myurl.net/home":    makeHtmlWithLinks("http://myurl.net",    []string{"/", "/home", "/about", "/posts"}),
      "http://myurl.net/about":   makeHtmlWithLinks("http://myurl.net",    []string{"/", "/home", "/about", "/posts"}),
      "http://myurl.net/posts":   makeHtmlWithLinks("http://myurl.net",    []string{"/", "/posts/1", "/posts/2"}),
      "http://myurl.net/posts/1": makeHtmlWithLinks("http://myurl.net",    []string{"/", "/posts/1", "/posts/2"}),
      "http://myurl.net/posts/2": makeHtmlWithLinks("http://external.net", []string{"/", "/link-a/", "/link-b"}),
    },
  }

  base, _ := url.Parse("http://myurl.net/")

  suite.T().Run("returns an array of all the urls encountered on the page", func(t *testing.T) {
    new_urls, _, _ := getLink(*base, *base, urlGetter)
    expected_urls := []url.URL{}
    for _, link := range []string{"/", "/home", "/about", "/posts"} {
      url, _ := url.Parse("http://myurl.net/" + link)
      expected_urls = append(expected_urls, *url)
    }
    assert.Equal(suite.T(), expected_urls, *new_urls)
  });

  suite.T().Run("returns the body with the links replaced with the corresponding markdown links", func(t *testing.T) {
    _, body, _ := getLink(*base, *base, urlGetter)
    assert.Equal(suite.T(), "<html><head></head><body><a href=\"/index.md\">a link</a><a href=\"/home.md\">a link</a><a href=\"/about.md\">a link</a><a href=\"/posts.md\">a link</a></body></html>", body )
  });


}
