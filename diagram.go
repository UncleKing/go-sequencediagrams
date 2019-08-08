package sequence

import (
	"go-sequencediagrams/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"image/color"
	"image/png"
	"strings"
)

type Diagram struct {
	participants []*Participant
	sequences    []Sequence
	// stores the name of participant and its index -- very much premature optimization but more to save on coding time.
	participantMap map[string]int

	size              utils.Rectangle
	SequenceFont      font.Face
	ParticipantFont   font.Face
	dc                *gg.Context
	participantHeight int
	sequenceEndY      int
	groupList         []*Group
	// will also include other parameters like font and config for sequence colors.
}

const (
	ST_SOLID                = 1
	ST_DOTTED               = 2
	ST_START_PROCESS        = 3
	ST_END_PROCESS          = 4
	ST_START_DOTTED_PROCESS = 5
	ST_END_DOTTED_PROCESS   = 6
	ST_NOTE_OVER            = 7
	ST_NOTE_LEFT            = 8
	ST_NOTE_RIGHT           = 9
	ST_GROUP_MESSAGE        = 10
	ST_ELSE_MESSAGE         = 11
	ST_END_GROUP            = 12
)

const (
	CONFIG_MIN_PADDING_X        = 30
	CONFIG_MIN_PADDING_Y        = 25
	CONFIG_TEXT_PADDING_X       = 15
	CONFIG_TEXT_PADDING_Y       = 5
	CONFIG_ARROW_WIDTH          = 20
	CONFIG_ARROW_HEIGHT         = 14
	CONFIG_SELF_DIAMETER        = 30
	CONFIG_PROCESS_WIDTH        = 10
	CONFIG_MESSAGE_LINE_SPACING = 5.0
	CONFIG_MESSAGE_ALIGN        = gg.AlignLeft
	CONFIG_GROUP_MAX_WIDTH      = 300
	CONFIG_GROUP_BASE_HEIGHT    = 30
)

var CONFIG_SEQUENCE_LINE_COLOR = color.RGBA{0, 0, 0, 255}
var CONFIG_PROCESS_LINE_COLOR = color.RGBA{0, 0, 0, 255}
var CONFIG_GROUP_LINE_COLOR = color.RGBA{0, 0, 0xff, 255}
var CONFIG_GROUP_BG_FILL_COLOR = color.RGBA{0xff, 0xff, 0xff, 255}
var CONFIG_GROUP_TEXT_COLOR = CONFIG_GROUP_LINE_COLOR

func (d *Diagram) GetOrCreateParticipant(name string) *Participant {
	p, ok := d.participantMap[name]
	if !ok {
		// add this participant to the map & the array.
		np := Participant{name: name}
		d.participants = append(d.participants, &np)
		d.participantMap[name] = len(d.participants) - 1
		p = d.participantMap[name]
	}
	return d.participants[p]
}

func (d *Diagram) AddSequence(s Sequence) {
	d.sequences = append(d.sequences, s)
}

func NewDiagram( /*typically we will take config here*/ ) (*Diagram, error) {
	d := Diagram{}
	// Create a temp context for text operations.
	d.dc = gg.NewContext(1, 1)
	d.participantMap = make(map[string]int)
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	d.SequenceFont = truetype.NewFace(font, &truetype.Options{Size: 12})
	d.ParticipantFont = truetype.NewFace(font, &truetype.Options{Size: 14})

	return &d, nil
}

func (d *Diagram) MeasureParticipant(participant *Participant) utils.Rectangle {
	d.dc.Push()
	defer d.dc.Pop()
	d.dc.SetFontFace(d.ParticipantFont)
	w, h := d.dc.MeasureString(participant.name)
	w += CONFIG_TEXT_PADDING_X * 2
	h += CONFIG_TEXT_PADDING_Y * 2

	return utils.Rect(0, 0, int(w), int(h))
}

func (d *Diagram) GetDelta(i1 int, i2 int) (int, error) {
	deltaX := 0
	for i := i1 + 1; i <= i2; i++ {
		p := d.participants[i]
		deltaX += p.delta
	}
	return deltaX, nil

}

func (d *Diagram) AdjustXSpace(p1 *Participant, p2 *Participant, xSpace int) error {

	i1, ok := d.participantMap[p1.name]
	if !ok {
		return fmt.Errorf(" Failed to find participant %s", p1.name)
	}

	i2, ok := d.participantMap[p2.name]
	if !ok {
		return fmt.Errorf(" Failed to find participant %s", p2.name)
	}

	if i1 == i2 {
		// check if there's anyone after us.. if there is then that becomes the i2
		if len(d.participants) > i1+1 {
			i2 = i1 + 1
		} else {
			// we are last we can just return from here
			return nil
		}
	}
	if i2 < i1 {
		i1, i2 = i2, i1
	}

	deltaX, err := d.GetDelta(i1, i2)
	if err != nil {
		return err
	}
	if deltaX > xSpace {
		// we already have enough space and we should go back
		return nil
	}
	avgSpace := int(xSpace / (i2 - i1))

	// check for space between from i1 to i2 for each and ensure that there is atleast
	for i := i1 + 1; i <= i2; i++ {
		if d.participants[i].delta < avgSpace {
			d.participants[i].delta = avgSpace
		}
	}

	return nil
}

func (d *Diagram) RenderParticipant(dc *gg.Context, p *Participant, yOffset int) {

	dc.Push()
	dc.DrawRectangle(float64(p.position.Min.X), float64(p.position.Min.Y+yOffset), float64(p.position.Dx()), float64(p.position.Dy()))

	centerX := float64(p.position.Min.X + p.position.Dx()/2)
	centerY := float64((p.position.Min.Y + p.position.Dy()/2) + yOffset)
	dc.SetFontFace(d.ParticipantFont)
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored(p.name, centerX, centerY, 0.5, 0.5)
	dc.Stroke()
	dc.Pop()
}

func (d *Diagram) RenderParticipantLines(dc *gg.Context, p *Participant) {
	dc.Push()

	dc.SetDash(5, 5)
	rt := p.position
	x1 := float64(rt.Min.X + rt.Dx()/2)
	y1 := float64(rt.Max.Y)
	x2 := float64(rt.Min.X + rt.Dx()/2)
	y2 := float64(d.sequenceEndY)
	dc.SetRGB(0, 0, 0xff)

	dc.DrawLine(x1, y1, x2, y2)
	dc.Stroke()
	dc.Pop()
}

// decide how to measure string without knowing the width

func CreateDiagram(sequence string) ([]byte, error) {

	d, err := NewDiagram()
	if err != nil {
		return []byte{}, err
	}

	// extract participants and store them in a list
	// also we have all the necessary objects created.
	// need to add more types of objects.
	err = d.Parse(sequence)
	if err != nil {
		return []byte{}, err
	}

	//// precompute lengths each participant and place
	d.ComputeParticipantSizeAndPlace()
	//
	//// Precompute length of message and move if necessary
	d.ComputeSequenceMessageAndPlace()

	d.RePlaceParticipants()

	w, h := d.ComputeImageSize()
	//// based on precomputed width and height create the image

	i := d.Render(w, h)

	buf := new(bytes.Buffer)

	err = png.Encode(buf, i)

	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil

}

func (p *Participant) SetDelta(delta int) {
	p.delta = delta
}

func (p *Participant) SetPosition(position utils.Rectangle) {
	p.position = position
}

var methodObjectMap = map[int]func() (Sequence, error){
	ST_SOLID:                func() (Sequence, error) { return new(SolidSequence), nil },
	ST_DOTTED:               func() (Sequence, error) { return new(DottedSequence), nil },
	ST_START_PROCESS:        func() (Sequence, error) { return new(StartProcess), nil },
	ST_END_PROCESS:          func() (Sequence, error) { return new(EndProcess), nil },
	ST_START_DOTTED_PROCESS: func() (Sequence, error) { return new(StartDottedProcess), nil },
	ST_END_DOTTED_PROCESS:   func() (Sequence, error) { return new(EndDottedProcess), nil },
	ST_GROUP_MESSAGE:        func() (Sequence, error) { return new(StartGroupMessage), nil },
	//ST_ELSE_MESSAGE:         func() (Sequence, error) { return new(ElseMessage), nil },
	ST_END_GROUP: func() (Sequence, error) { return new(EndGroupMessage), nil },
}

func (d *Diagram) Parse(sequence string) error {

	groupStack := utils.Stack{}

	if len(sequence) == 0 {
		return fmt.Errorf("Empty sequence")
	}
	lines := strings.Split(sequence, "\n")
	for _, line := range lines {
		seq := make(map[string]interface{})
		jsonstr, typ, err := ParseLine(line)
		if err != nil {
			// figure out a way of sending error line
			return err
		}

		err = json.Unmarshal([]byte(jsonstr), &seq)
		if err != nil {
			// figure out a way of sending error line
			return err
		}

		fun := methodObjectMap[typ]
		if fun == nil {
			return fmt.Errorf("Internal error %d", typ)
		}

		obj, err := fun()
		obj.Init(seq, d, len(d.sequences), typ)
		d.AddSequence(obj)

		if obj.IsStartProcess() {
			p := Process{}
			p.start = obj
			obj.SecondaryParticipant().AddProcess(&p)
			obj.SetStartProcess(&p)
		} else if obj.IsEndProcess() {
			p := obj.PrimaryParticipant().EndProcessAt(obj)
			obj.SetEndProcess(p)
		}

		if typ == ST_GROUP_MESSAGE {
			// add the start to the group stack
			g := Group{}
			g.start = obj.(*StartGroupMessage)
			g.start.group = &g
			groupStack.Push(g)

		}
		if typ == ST_END_GROUP {
			// add the else to the current group stack
			g := groupStack.Pop()
			if g == nil {
				return fmt.Errorf("End without group")
			}
			group := g.(Group)
			group.end = obj.(*EndGroupMessage)
			group.end.group = &group

			d.groupList = append(d.groupList, &group)
		}
		if typ == ST_ELSE_MESSAGE {
			// add the end to the current group stack and compute
			// recompute the group
		}

	}
	// split the string into lines
	// create objects for each line

	return nil

}

func (d *Diagram) ReComputeGroup(g *Group) {
	// Recompute the group when the end is received. At this point all the sequence are measured.

	// compute all the sequences in between and calculate the width required
	partStartIndex := -1
	partEndIndex := -1
	seqStartIndex := g.start.Index()
	startX := -1
	endX := -1
	startY := -1
	endY := g.Position().Max.Y
	// we have initial status of y of group
	for i := seqStartIndex; i < g.end.Index(); i++ {

		s := d.sequences[i]

		partStartIndex, partEndIndex = d.ExpandGroup(s.PrimaryParticipant(), partStartIndex, partEndIndex)
		partStartIndex, partEndIndex = d.ExpandGroup(s.SecondaryParticipant(), partStartIndex, partEndIndex)

		childPos := s.Position()
		if startX > childPos.Min.X || startX == -1 {
			startX = childPos.Min.X
		}
		if endX < childPos.Max.X || endX == -1 {
			endX = childPos.Max.X
		}
		if endY < childPos.Max.Y {
			endY = childPos.Max.Y
		}

		if startY > childPos.Min.Y || startY == -1 {
			startY = childPos.Min.Y
		}
	}

	if partStartIndex == -1 {
		partStartIndex = 0
	}
	if partEndIndex == -1 {
		partEndIndex = len(d.participants)
	}

	endX += CONFIG_MIN_PADDING_X
	endY += CONFIG_MIN_PADDING_Y

	r := utils.Rect(startX, startY, endX, endY)
	g.start.primary = d.participants[partStartIndex]
	g.start.secondary = d.participants[partEndIndex]
	g.end.primary = d.participants[partStartIndex]
	g.end.secondary = d.participants[partEndIndex]
	g.SetPosition(r)

}

func (d *Diagram) ExpandGroup(p *Participant, partStartIndex int, partEndIndex int) (int, int) {
	if p != nil {
		idx, ok := d.participantMap[p.name]
		if ok {
			if idx > partEndIndex || partEndIndex == -1 {
				partEndIndex = idx
			}
			if idx < partStartIndex || partStartIndex == -1 {
				partStartIndex = idx
			}
		}
	}
	return partStartIndex, partEndIndex
}

func (p *Participant) AddProcess(process *Process) {
	if obj := p.processStack.Peek(); obj != nil {
		process.parent = obj.(*Process)
	}
	p.processStack.Push(process)
	p.processes = append(p.processes, process)
}

func (p *Participant) EndProcessAt(s Sequence) *Process {
	// end the process on o top.
	if obj := p.processStack.Pop(); obj != nil {
		process := obj.(*Process)
		process.end = s
		return process
	} // else if its nil we just ignore it ( TODO: To attach warning to the line.
	return nil
}

func (d *Diagram) ComputeParticipantSizeAndPlace() error {
	for idx, p := range d.participants {
		r := d.MeasureParticipant(p)
		if r.Dy() > d.participantHeight {
			d.participantHeight = r.Dy()
		}
		p.SetPosition(r)
		if idx != 0 {
			p.SetDelta(CONFIG_MIN_PADDING_X)
		}
	}
	return nil
}

func (d *Diagram) ComputeSequenceMessageAndPlace() error {
	d.sequenceEndY = CONFIG_MIN_PADDING_Y + d.participantHeight

	for _, s := range d.sequences {

		r := s.MeasureBounds(d, d.SequenceFont)
		r = r.Add(image.Point{X: 0, Y: d.sequenceEndY})
		s.SetPosition(r)
		d.sequenceEndY += r.Dy() + CONFIG_MIN_PADDING_Y

		// ensure that there is enough space between the two participants
		p1 := s.PrimaryParticipant()
		p2 := s.SecondaryParticipant()
		if p2 == nil {
			// this is where we only have a single element ( like notes )
			// only update the y element.
			continue
		}

		/*

			logic for spacing.
			if the two actors are next to each other expand as necessary.
			if there are >= 1  actors in between -- Distribute space between them. if there is already > average skip those two
			also push down the actual add of space to all actors after the last one.
			note that the last one is not same as secondary participant.

		*/

		err := d.AdjustXSpace(p1, p2, r.Dx())
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Diagram) ComputeImageSize() (int, int) {
	// get the last participants and its delta and get the image size
	lastParticipant := d.participants[len(d.participants)-1]
	//lastSequence := d.sequences[len(d.sequences)-1]

	imageWidth := lastParticipant.position.Max.X + CONFIG_MIN_PADDING_X
	imageHeight := d.sequenceEndY + lastParticipant.position.Dy()*2

	return imageWidth, imageHeight

}

func (d *Diagram) Render(width int, height int) image.Image {
	d.dc = gg.NewContext(width, height)

	// draw the participants.
	for _, p := range d.participants {
		// draws the participants and the top and bottom based on config.
		d.RenderParticipant(d.dc, p, 0)
		// draws the dotted lines
		d.RenderParticipantLines(d.dc, p)

		d.RenderParticipant(d.dc, p, d.sequenceEndY)
		d.RenderProcesses(d.dc, p)
	}

	for _, s := range d.sequences {
		s.Render(d, d.dc)
	}

	for _, g := range d.groupList {
		d.RenderGroup(g)
	}

	return d.dc.Image()

}

func (d *Diagram) RePlaceParticipants() {

	x := CONFIG_MIN_PADDING_X

	for idx, p := range d.participants {
		if idx == 0 {
			nr := p.position.Add(image.Point{X: x, Y: 0})
			p.position = nr

		} else {
			// this is minimum width between two center points.

			minWidth := p.delta
			prevPos := d.participants[idx-1].position

			diffPart := prevPos.Dx()/2 + p.position.Dx()/2 + CONFIG_MIN_PADDING_X
			if diffPart > minWidth {
				minWidth = diffPart
			}

			newX := prevPos.Min.X + prevPos.Dx()/2 + minWidth - p.position.Dx()/2
			nr := p.position.Add(image.Point{X: newX, Y: 0})
			p.position = nr
		}
		d.PlaceProcesses(p)
	}
}

func (d *Diagram) PlaceProcesses(p *Participant) {
	for _, process := range p.processes {
		xOffset := p.position.MidX() - CONFIG_PROCESS_WIDTH/2
		if process.parent != nil {
			// check if we have a parent
			xOffset = process.parent.position.Min.X + CONFIG_PROCESS_WIDTH/2
		}
		xStart := xOffset
		xEnd := xStart + CONFIG_PROCESS_WIDTH
		yStart := process.start.Position().MidY()
		yEnd := d.sequenceEndY

		if process.end != nil {
			yEnd = process.end.Position().MidY()
		}
		process.position = utils.Rect(xStart, yStart, xEnd, yEnd)
	}
}

func (d *Diagram) RenderProcesses(dc *gg.Context, p *Participant) {
	// renders all the processes associated with participant
	dc.Push()
	defer dc.Pop()
	dc.SetColor(CONFIG_PROCESS_LINE_COLOR)
	for _, process := range p.processes {
		r := process.position

		dc.DrawRectangle(float64(r.Min.X), float64(r.Min.Y), float64(r.Dx()), float64(r.Dy()))
	}
	dc.Stroke()
}

func (d *Diagram) GetProcessAtSequence(p *Participant, sIndex int) *Process {

	for i := len(p.processes) - 1; i >= 0; i-- {
		process := p.processes[i]
		startIndex := process.start.Index()
		endIndex := len(d.sequences)
		if process.end != nil {
			endIndex = process.end.Index()
		}

		if startIndex <= sIndex && sIndex <= endIndex {
			return process
		}
	}
	return nil
}

func (d *Diagram) RenderGroup(g *Group) {

	d.dc.Push()
	defer d.dc.Pop()
	d.dc.SetFontFace(d.SequenceFont)
	d.dc.SetColor(CONFIG_GROUP_LINE_COLOR)

	x1 := float64(g.start.PrimaryParticipant().position.Min.X)
	w := float64(g.start.SecondaryParticipant().position.Max.X) - x1
	if w < float64(g.position.Dx()) {
		w = float64(g.position.Dx())
	}

	// group's position as rect is already set before
	pos := g.Position()
	// draw the text for group name.

	msgRect := g.MessageRect()

	d.dc.DrawStringWrapped(g.Text(), float64(msgRect.Min.X), float64(msgRect.Min.Y), 1, 1, CONFIG_GROUP_MAX_WIDTH,
		CONFIG_MESSAGE_LINE_SPACING, CONFIG_MESSAGE_ALIGN)

	d.dc.DrawRectangle(x1, float64(g.position.Min.Y), float64(w), float64(g.position.Dy()))
	d.dc.Stroke()

	d.dc.MoveTo(x1, float64(pos.Min.Y))
	// THIS IS YUCK !! please refactor !
	d.dc.LineTo(x1, float64(pos.Min.Y)+CONFIG_MIN_PADDING_Y)
	d.dc.LineTo(x1+CONFIG_MIN_PADDING_X, float64(pos.Min.Y)+CONFIG_MIN_PADDING_Y)
	d.dc.LineTo(x1+CONFIG_MIN_PADDING_X, float64(pos.Min.Y))
	d.dc.SetColor(CONFIG_GROUP_BG_FILL_COLOR)
	d.dc.Fill()

	d.dc.SetColor(CONFIG_GROUP_TEXT_COLOR)
	d.dc.DrawStringAnchored(g.Name(), x1+CONFIG_MIN_PADDING_X/2,
		float64(pos.Min.Y)+CONFIG_MIN_PADDING_Y/2, 0, 0)

	//// draw the bounds of the group.
	//points := g.GetNameBounds()
	//
	//for idx, p := range points {
	//	if idx == 0 {
	//		d.dc.MoveTo(p.X, p.Y)
	//	} else {
	//		d.dc.LineTo(p.X, p.Y)
	//	}
	//
	//}
	//d.dc.Stroke()
	//
}

// parse all lines to fetch the participants and measure them.
// add sequence and add place them between the participants ( ensure space
// for each group
