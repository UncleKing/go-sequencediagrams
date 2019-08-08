package sequence

import "go-sequencediagrams/utils"

type Participant struct {
	name     string
	position utils.Rectangle
	// This is the distance from the previous participant
	delta        int
	processStack utils.Stack
	processes    []*Process
}
