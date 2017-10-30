package main

import (
	_ "fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/url"
	"os"
	"testing"
)

type ConfigTestSuite struct {
	suite.Suite
	fs afero.Fs
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *ConfigTestSuite) SetupTest() {
	suite.fs = afero.NewMemMapFs()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func removeDirectoryAndReportExistence(fs afero.Fs, path string) bool {
	dirStat, err := fs.Stat(path)

	if err == nil {
		defer func() {
			fs.Remove(path)
		}()
		return dirStat.IsDir()
	} else {
		return false
	}

}

func (suite *ConfigTestSuite) TestConfigArgsProvidedAndDefault() {

	suite.T().Run("empty options initialize with default values", func(t *testing.T) {
		default_url, _ := url.Parse(defaultRootURL)
		config, err := readFromArgs([]string{"command_name"}, suite.fs)
		assert.Nil(suite.T(), err, "should not throw an error")
		assert.Equal(suite.T(), config, _config{*default_url, "./.tmp", suite.fs}, "didn't initialize with default values")

		assert.True(
			suite.T(),
			removeDirectoryAndReportExistence(suite.fs, "./.tmp"),
			"didn't create the directory",
		)
	})

	suite.T().Run("partially empty options initialize with default values", func(t *testing.T) {
		expected_url, _ := url.Parse("http://golang.org")
		config, err := readFromArgs([]string{"command_name", "http://golang.org"}, suite.fs)
		assert.Nil(suite.T(), err, "should not throw an error")
		assert.Equal(suite.T(), config, _config{*expected_url, "./.tmp", suite.fs}, "didn't initialize with default values")

		assert.True(
			suite.T(),
			removeDirectoryAndReportExistence(suite.fs, "./.tmp"),
			"didn't create the directory",
		)
	})

	suite.T().Run("options initialize when all provided", func(t *testing.T) {
		expected_url, _ := url.Parse("http://golang.org")
		config, err := readFromArgs([]string{"command_name", "http://golang.org", "/output_dir"}, suite.fs)

		assert.Nil(suite.T(), err, "should not throw an error")
		assert.Equal(suite.T(), config, _config{*expected_url, "/output_dir", suite.fs}, "didn't initialize with default values")

		assert.True(suite.T(),
			removeDirectoryAndReportExistence(suite.fs, "/output_dir"),
			"didn't create the directory",
		)
	})

}

func (suite *ConfigTestSuite) TestConfigFixUrl() {

	suite.T().Run("adds a slash to the end of the url if it doesn't end in a file extension", func(t *testing.T) {
		defer func() {
			suite.fs.Remove("/output_dir")
		}()
		config, err := readFromArgs([]string{"command_name", "http://golang.org/public", "/output_dir"}, suite.fs)
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), config.rooturl.String(), "http://golang.org/public/", "didn't add the trailing slash")
	})

	suite.T().Run("doesn't add a slash to the end of the url if it ends in a file extension", func(t *testing.T) {
		defer func() {
			suite.fs.Remove("/output_dir")
		}()
		config, err := readFromArgs([]string{"command_name", "http://golang.org/index.php", "/output_dir"}, suite.fs)
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), config.rooturl.String(), "http://golang.org/index.php", "didn't add the trailing slash")
	})

}

func (suite *ConfigTestSuite) TestConfigErrors() {

	suite.T().Run("fails if the url provided is invalid", func(t *testing.T) {
		defer func() {
			suite.fs.Remove("/output_dir")
		}()
		_, err := readFromArgs([]string{"command_name", "::hallo", "/output_dir"}, suite.fs)
		assert.EqualError(suite.T(), err, "parse ::hallo: missing protocol scheme")
	})

	suite.T().Run("fails if directory provided already exists", func(t *testing.T) {
		defer func() {
			suite.fs.Remove("/output_dir")
		}()
		suite.fs.Mkdir("/output_dir", os.ModePerm)
		_, err := readFromArgs([]string{"command_name", "http://www.golang.org", "/output_dir"}, suite.fs)
		assert.EqualError(suite.T(), err, "directory does already exist!")
	})

}
