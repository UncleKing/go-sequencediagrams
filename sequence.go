package sequence

import (
	"go-sequencediagrams/utils"
	"github.com/fogleman/gg"
	"golang.org/x/image/font"

	"math"
)

type Sequence interface {
	PrimaryParticipant() *Participant
	Text() string
	SecondaryParticipant() *Participant
	SetPosition(rectangle utils.Rectangle)
	Position() utils.Rectangle
	Render(d *Diagram, dc *gg.Context)
	MeasureBounds(d *Diagram, sequenceFont font.Face) utils.Rectangle

	IsStartProcess() bool
	IsEndProcess() bool

	SetStartProcess(*Process)
	SetEndProcess(*Process)
	Index() int

	Type() int

	Init(data map[string]interface{}, d *Diagram, index int, seqType int) error
}

type BaseSequence struct {
	primary   *Participant
	secondary *Participant
	message   string
	// stores the y position of the sequence.
	position utils.Rectangle

	startProcess *Process
	endProcess   *Process
	index        int

	seqType int
}

func (s *BaseSequence) Type() int {
	return s.seqType
}

func (s *BaseSequence) SetStartProcess(sp *Process) {
	s.startProcess = sp
}

func (s *BaseSequence) SetEndProcess(ep *Process) {
	s.endProcess = ep
}

func (s *BaseSequence) PrimaryParticipant() *Participant {
	return s.primary
}
func (s *BaseSequence) Text() string {
	return s.message
}
func (s *BaseSequence) SecondaryParticipant() *Participant {
	return s.secondary
}

func (s *BaseSequence) Index() int {
	return s.index
}

func (s *BaseSequence) Init(data map[string]interface{}, d *Diagram, index int, seqType int) error {
	s.primary = d.GetOrCreateParticipant(data["src"].(string))
	s.secondary = d.GetOrCreateParticipant(data["dest"].(string))
	s.message = data["text"].(string)
	s.seqType = seqType
	s.index = index
	return nil
}

func (s *BaseSequence) MeasureBounds(d *Diagram, sequenceFont font.Face) utils.Rectangle {
	dc := d.dc
	dc.Push()
	defer dc.Pop()
	dc.SetFontFace(sequenceFont)
	w, h := dc.MeasureString(s.Text())

	w += CONFIG_TEXT_PADDING_X*2 + CONFIG_ARROW_WIDTH
	h += CONFIG_TEXT_PADDING_Y * 2

	// special condition check if primary and secondary are same
	if s.primary == s.secondary {
		h += CONFIG_SELF_DIAMETER
	}
	// measures the text width and adds padding towards the end.
	return utils.Rect(0, 0, int(w), int(h))
}

func (s *BaseSequence) SetPosition(rectangle utils.Rectangle) {
	s.position = rectangle
}

func (s *BaseSequence) Position() utils.Rectangle {
	return s.position
}

func (s *BaseSequence) IsStartProcess() bool {
	return false
}

func (s *BaseSequence) IsEndProcess() bool {
	return false
}

type SolidSequence struct {
	BaseSequence
}

func (s *SolidSequence) Render(d *Diagram, dc *gg.Context) {
	s.RenderSequence(d, dc, false)
}

func (s *StartProcess) Render(d *Diagram, dc *gg.Context) {
	s.RenderSequence(d, dc, false)
}

func (s *EndProcess) Render(d *Diagram, dc *gg.Context) {
	s.RenderSequence(d, dc, false)
}

func (s *StartDottedProcess) Render(d *Diagram, dc *gg.Context) {
	s.RenderSequence(d, dc, true)
}

func (s *EndDottedProcess) Render(d *Diagram, dc *gg.Context) {
	s.RenderSequence(d, dc, true)
}

type DottedSequence struct {
	BaseSequence
}

func (s *DottedSequence) Render(d *Diagram, dc *gg.Context) {
	s.RenderSequence(d, dc, true)
}

type StartProcess struct {
	BaseSequence
}

func (s *StartProcess) IsStartProcess() bool {
	return true
}

type EndProcess struct {
	BaseSequence
}

func (s *EndProcess) IsEndProcess() bool {
	return true
}

type StartDottedProcess struct {
	BaseSequence
}

func (s *StartDottedProcess) IsStartProcess() bool {
	return true
}

type EndDottedProcess struct {
	BaseSequence
}

func (s *EndDottedProcess) IsEndProcess() bool {
	return true
}

// zero angle is >
func (b *BaseSequence) DrawArrow(dc *gg.Context, width float64, height float64, x int, y int, angle float64) {
	adc := gg.NewContext(int(width), int(height))
	// set transparent background
	adc.SetHexColor("ffffff00")
	adc.Clear()

	adc.MoveTo(0, 0)
	adc.LineTo(width, height/2)
	adc.LineTo(0, height)
	adc.ClosePath()
	adc.SetColor(CONFIG_SEQUENCE_LINE_COLOR)
	adc.FillPreserve()
	adc.Stroke()
	dc.Push()
	dc.RotateAbout(gg.Radians(angle), float64(x), float64(y))
	dc.DrawImage(adc.Image(), x, y-int(height/2))
	dc.Pop()
}

func (b BaseSequence) RenderSequence(d *Diagram, dc *gg.Context, isDotted bool) {

	dc.Push()
	defer dc.Pop()

	//special condition if the start and end are same
	if b.primary == b.secondary {
		// we draw a line and a semi circle and we come back and end at the same participant
		// we then draw the text over it.

		p1 := b.primary.position
		x1 := float64(p1.Min.X + p1.Dx()/2)

		// first we check if we are starting a process
		process := d.GetProcessAtSequence(b.primary, b.index)
		if process != nil {
			x1 = float64(process.position.Max.X)
		}

		x2 := x1 + float64(b.position.Dx())/2
		if isDotted {
			dc.SetDash(5, 5)
		}
		dc.DrawLine(x1, float64(b.position.Min.Y), x2, float64(b.position.Min.Y))
		dc.Stroke()
		dc.DrawEllipticalArc(x2, float64(b.position.Min.Y)+CONFIG_SELF_DIAMETER/2, CONFIG_SELF_DIAMETER/2, CONFIG_SELF_DIAMETER/2, gg.Radians(90), gg.Radians(-90))
		dc.Stroke()
		dc.DrawLine(x1, float64(b.position.Min.Y)+CONFIG_SELF_DIAMETER, x2, float64(b.position.Min.Y)+CONFIG_SELF_DIAMETER)
		dc.Stroke()

		arrowStartX := x1 + CONFIG_ARROW_WIDTH
		arrowAngle := 180.0
		b.DrawArrow(dc, CONFIG_ARROW_WIDTH, CONFIG_ARROW_HEIGHT, int(arrowStartX), b.position.Min.Y+CONFIG_SELF_DIAMETER, arrowAngle)

		dc.SetFontFace(d.SequenceFont)
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(b.message, x2, float64(b.position.Min.Y), 0.5, -0.2)

		return
	} else {

		if isDotted {
			dc.SetDash(5, 5)
		}

		p1 := b.primary.position
		p2 := b.secondary.position
		x1 := float64(p1.Min.X + p1.Dx()/2)
		y := float64(b.position.Min.Y + b.position.Dy()/2)
		x2 := float64(p2.Min.X + p2.Dx()/2)

		isReverse := false
		if x1 > x2 {
			isReverse = true
		}

		process := d.GetProcessAtSequence(b.primary, b.index)
		if process != nil {
			if isReverse {
				x1 = float64(process.position.Min.X)
			} else {
				x1 = float64(process.position.Max.X)
			}
		}

		process = d.GetProcessAtSequence(b.secondary, b.index)
		if process != nil {
			if isReverse {
				x2 = float64(process.position.Max.X)
			} else {
				x2 = float64(process.position.Min.X)
			}
		}

		arrowStartX := x2 - CONFIG_ARROW_WIDTH

		arrowAngle := 0.0
		if isReverse {
			// your primary is after the secondary.
			// arrow direction changes
			x1, x2 = x2, x1
			arrowStartX = x1 + CONFIG_ARROW_WIDTH
			arrowAngle = 180
		}
		centerX := x1 + math.Abs(x2-x1)/2

		dc.SetRGB(0, 0, 0xff)
		dc.DrawLine(x1, y, x2, y)
		dc.Stroke()
		dc.SetFontFace(d.SequenceFont)
		b.DrawArrow(dc, CONFIG_ARROW_WIDTH, CONFIG_ARROW_HEIGHT, int(arrowStartX), int(y), arrowAngle)
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(b.message, centerX, y, 0.5, -0.2)

	}
}
