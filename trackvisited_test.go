package main;

import (
  "testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"time"
	"net/url"
)

type TrackvisitedTestSuite struct {
    suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *TrackvisitedTestSuite) SetupTest() {}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTrackvisitedTestSuite(t *testing.T) {
    suite.Run(t, new(TrackvisitedTestSuite))
}

func (suite *TrackvisitedTestSuite) TestUrlToMarkdownUri() {

  suite.T().Run("keeps thread safe track of visited urls", func(t *testing.T) {
    tracker := trackVisitedUrls();
    uri, _ := url.Parse("http://myurl.net/my-test-url");
    numer_of_runs := 100;
    swarm := make([]int, numer_of_runs);

    results := make(chan bool)
    granted_locks := 0
    rejected_locks := 0
    // firing 100 requests in
    // parallel into the service
    // tracking the results in
    // the results channel
    for range swarm {
      go func() {
        results <- (visitedUpdateRequest{url: *uri}).dispatch(tracker)
      }()
    }

    terminate := time.After(time.Second / 200)

    for true {
      select {
      case res := <-results:
        if res == true {
          granted_locks += 1
        } else {
          rejected_locks += 1
        }
        continue
      case <-terminate:
      }
      break
    }

    assert.Equal(t, 1, granted_locks, "the service allowed more than one client to reserver the resource!")
    assert.Equal(t, numer_of_runs - 1, rejected_locks, "the test didn't finish all the runs before the timeout")

  })
}
