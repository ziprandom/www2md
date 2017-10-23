package main

import (
	"fmt"
  "net"
	"net/url"
	"os"
	"errors"
	"sort"
	"github.com/spf13/afero"
)


func appendUnique(links *[]url.URL, link url.URL) []url.URL {
	for _, _link := range *links {
		if _link.Path == link.Path {
			return *links
		}
	}
	return append(*links, link)
}


func retrieve(control_channels controlChannels, config _config, uri url.URL, visited chan visitedUpdateRequest) {

	if !( visitedUpdateRequest{url: uri} ).dispatch(visited) {
    control_channels.skipped <- errors.New("url already visited")
    return
  }

	links, body, error :=
    getLink(config.rooturl, uri, nil)

  if error != nil {
    control_channels.skipped <- error
  } else {
		saveMarkdownFile(config, uri, body)
  }

	additional_fetches := 0
	for _, link := range *links {
    additional_fetches++
    go retrieve(control_channels, config, link, visited)
	}

	control_channels.fetched <- url_subsequent_visits{uri, additional_fetches}

}

type url_subsequent_visits struct {
  url url.URL
  subsequent_requests int
}

type controlChannels struct {
  fetched chan url_subsequent_visits
  skipped chan error
}

func main() {

  config, err :=
    readFromArgs(os.Args, afero.NewOsFs());

  if  err != nil {
    fmt.Printf("\nusage: %v <root_url> <storage_directory>\n\n", os.Args[0])
    fmt.Printf("you didn't provide valid arguments:\n\n  %s\n\n", err)
    os.Exit(1)
  }

	fmt.Printf("using the following arguments: \n %#v\n", config)

  track_visited := trackVisitedUrls()

  control_channels := controlChannels{
    make(chan url_subsequent_visits),
    make(chan error),
  }

	go retrieve(control_channels, config, config.rooturl, track_visited)

	tocollect := 1

  urls_fetched := []url.URL{}

	for n := 0; n < tocollect; n++ {
    select {
    case url_subsequent := <-control_channels.fetched:
      urls_fetched = append(urls_fetched, url_subsequent.url)
      tocollect += url_subsequent.subsequent_requests
    case err := <-control_channels.skipped:
      if oe, ok := err.(*net.OpError); ok {
        fmt.Println(oe)
      } else {
        // url has been skipped because it was already requested
        // we dont treat this as an error, just to track the
        // requests for channel syncronisation
      }
    }
  }

  fmt.Printf("Finished scraping `%s`\n", config.rooturl.String())
  fmt.Println("locations processed:")

  sort.Slice(urls_fetched, func(i, j int) (bool) {
    return urls_fetched[i].Path < urls_fetched[j].Path
  })

  for i, link := range  urls_fetched {
    fmt.Printf("%3d %s\n", i+1, link.String())
  }

  fmt.Println("bye bye...")

}
