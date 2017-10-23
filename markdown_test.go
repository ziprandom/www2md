package main

import (
	"testing"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/suite"
	"net/url"
	"github.com/spf13/afero"
	"fmt"
)

type MarkdownTestSuite struct {
    suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *MarkdownTestSuite) SetupTest() {
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMarkdownTestSuite(t *testing.T) {
    suite.Run(t, new(MarkdownTestSuite))
}

func (suite *MarkdownTestSuite) TestUrlToMarkdownUri() {
  base, _ := url.Parse("http://myurl.net/");

  suite.T().Run("turns a link into a file system path", func(t *testing.T) {
    url, _ := url.Parse("http://myurl.net/about/me/index.html");
    name, path := urlToMarkdownUri(*base, *url)
    assert.Equal(t, "index.md", name);
    assert.Equal(t, "about/me/", path);
  })

  suite.T().Run("turns a link into a file system path using the last segment as filename", func(t *testing.T) {
    url, _ := url.Parse("http://myurl.net/about/me/");
    name, path := urlToMarkdownUri(*base, *url)
    assert.Equal(t, "me.md", name);
    assert.Equal(t, "about/", path);
  })

}

func (suite *MarkdownTestSuite) TestSaveToMarkdownFile() {
  base, _ := url.Parse("http://myurl.net/")
  fs := afero.NewMemMapFs()
  config := _config{*base, "test", fs}
  url, _ := url.Parse("http://myurl.net/about/me/");
  body := `
    <html>
      <head></head>
      <body>
        <div>
          <h1 id="welcome">Welcome</h1>
          <p>to this real world example</p>
          <ul>
            <li>of a typical html</li>
            <li>document</li>
          </ul>
          <pre>
            <code>
              and some source code
            </code>
          </pre>
          <blockquote>
            <p>
              and a quote <br>
              that needs to be here  <br>
              and <strong><em>whatnot</em></strong>
            </p>
          </blockquote>
        </div>
      </body>
    </html>
  `
  expected := "\n# Welcome\nto this real world example\n*   of a typical html\n*   document\n`\nand some source code\n`\nand a quote\nthat needs to be here\nand **_whatnot_**\n"

  suite.T().Run("converts an html string to markdown and saves it to the filesystem for a given url", func(t *testing.T) {
    saveMarkdownFile(config, *url, body)
    _, err := fs.Stat("test/about/me.md")
    assert.Nil(t, err, fmt.Sprintf("an error occured %v", err));
    bytes, err := afero.ReadFile(fs, "test/about/me.md")
    assert.Nil(t, err, fmt.Sprintf("an error occured %v", err));
    body := string(bytes)
    assert.Equal(t, expected, body)
  })
}
