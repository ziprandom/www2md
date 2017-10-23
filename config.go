package main

import (
	"net/url"
  "os"
	"errors"
	"fmt"
	"regexp"
	"github.com/spf13/afero"
)

// a struct holding the configuration
// information for the scraping run
type _config struct {
  rooturl url.URL
  storagedir string
  fs afero.Fs
}

// default values
const defaultRootURL = "http://www.catb.org/esr/writings/taoup/html/"
const defaultStorageDir = "./.tmp"

func (c _config) GoString() (string) {
  return fmt.Sprintf(
    `base url: "%s", storage directory: "%s"`,
    c.rooturl.String(), c.storagedir);
}

// create the _config struct from the given
// slice of strings optionally filling in
// missing values with the default values
func readFromArgs(args []string, fs afero.Fs) (_config, error) {

  // we don't need the
  // command name
  args = args[1:]

  var root *url.URL
  var storageDir string
  var err error

  switch len(args) {
  case 0:
    root, err = url.Parse(defaultRootURL)
    storageDir = defaultStorageDir
  case 1:
    root, err = url.Parse(args[0])
    storageDir = defaultStorageDir
  default:
    root, err = url.Parse(args[0])
    storageDir = args[1]
  }

  if err != nil {
    return _config{}, err
  }

  // throw an error if directory already exists or
  // if creation failes.
  if _, err = fs.Stat(storageDir); os.IsNotExist(err) {
    err = fs.MkdirAll(storageDir, os.ModePerm)
  } else if err == nil {
    err = errors.New("directory does already exist!")
  }

  // add a trailing / unless the path ends in a file name
	if regexp.MustCompile(`\/[\w\d_-]+$`).MatchString(root.Path) {
		root.Path += "/"
	}

  if err != nil {
    return _config{}, err
  } else {
    return _config{*root, storageDir, fs}, err
  }
}
