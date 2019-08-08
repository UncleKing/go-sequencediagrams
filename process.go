package sequence

import "go-sequencediagrams/utils"

type Process struct {
	start    Sequence
	end      Sequence
	parent   *Process
	position utils.Rectangle
}
