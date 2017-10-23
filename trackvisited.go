package main

import (
	"net/url"
)

// A payload that is used to request a lock for a given url.
// url field designates the url that the routine want's to
// fetch. responseChannel should be provided as a callback
// for the requesting routine. A message of true on the
// responseChannel indicates that no other routine is has
// already handled the url
type visitedUpdateRequest struct {
  url url.URL
  responseChannel chan bool
}

// Dispatch a request on a queue
func (r visitedUpdateRequest) dispatch(queue chan visitedUpdateRequest) (bool) {
  if r.responseChannel == nil {
    r.responseChannel =  make(chan bool)
  }
  queue <- r
  return <- r.responseChannel
}

// Spawns a routine to keep track of already visited urls
// returns a channel that processes visitedUpdateRequest
// strucs
func trackVisitedUrls() (chan visitedUpdateRequest) {
  visited, update_visited :=
    make(map[string]bool), make(chan visitedUpdateRequest)

  go func() {
    for {
      // listen for messages on the
      // channel
      request := <- update_visited
      if visited[request.url.Path] == true {
        request.responseChannel <- false
      } else {
        visited[request.url.Path] = true
        request.responseChannel <- true
      }
    }
  }()

  return update_visited
}
