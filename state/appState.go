package state

import (
	"lazy-queues/util"
)

type AppState struct {
	ProjectID string
}

type Subscription struct {
	chartState string
}

var State AppState = AppState{
	ProjectID: "",
}

func (s *AppState) Validate() bool {
	if s.ProjectID == "" {
		util.Log.Error("Missing projectID")
		return false
	}

	return true
}
