package ctfd

import (
	"golang.org/x/xerrors"
)

type ChallengesData struct {
	ID         int64  `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Value      int64  `json:"value"`
	Solves     int64  `json:"solves"`
	SolvedByMe bool   `json:"solved_by_me"`
	Category   string `json:"category"`
	// Tags       []string `json:"tags"`
	// Template   string   `json:"template"`
	// Script     string   `json:"script"`
}

type Challenges struct {
	Success bool             `json:"success"`
	Data    []ChallengesData `json:"data"`
}

func (c *Client) ListChallenges() ([]ChallengesData, error) {
	chals := &Challenges{}
	_, err := c.sling.New().Get("challenges").ReceiveSuccess(chals)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if !chals.Success {
		return nil, xerrors.Errorf("success: false")
	}
	return chals.Data, nil
}
