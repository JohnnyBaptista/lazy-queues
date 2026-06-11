package state

import (
	"lazy-queues/util"
)

type AppState struct {
	SubscriptionID string
	ProjectID      string
}

var State AppState = AppState{
	SubscriptionID: "",
	ProjectID:      "",
}

func (s *AppState) Validate() bool {
	if s.ProjectID == "" {
		util.Log.Error("Missing projectID")
		return false
	}

	if s.SubscriptionID == "" {
		util.Log.Error("Missing SubscriptionID")
		return false
	}

	return true
}
