package ctfd

import (
	"path"
	"strconv"

	"golang.org/x/xerrors"
)

type Hint struct {
	ID      int64  `json:"id"`
	Cost    int64  `json:"cost"`
	Content string `json:"content"`
}

type Challenge struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Value       int64  `json:"value"`
	Description string `json:"description"`
	// ConnectionInfo
	Category    string `json:"category"`
	State       string `json:"state"`
	MaxAttempts int64  `json:"max_attempts"`
	Type        string `json:"type"`
	// TypeData
	Solves     int64    `json:"solves"`
	SolvedByMe bool     `json:"solved_by_me"`
	Attempts   int64    `json:"attempts"`
	Files      []string `json:"files"`
	// Tags
	Hints []Hint `json:"hints"`
	// View
}

func (c *Client) GetChallenge(id int64) (*Challenge, error) {
	chal := new(struct {
		Success bool      `json:"success"`
		Data    Challenge `json:"data"`
	})
	_, err := c.sling.New().Get(path.Join("challenges", strconv.FormatInt(id, 10))).ReceiveSuccess(chal)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if !chal.Success {
		return nil, xerrors.Errorf("success: false")
	}
	return &chal.Data, nil
}
