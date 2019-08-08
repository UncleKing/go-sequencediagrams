package sequence

import (
	"github.com/fogleman/gg"
	"go-sequencediagrams/utils"
	"golang.org/x/image/font"
)

type Group struct {
	start    *StartGroupMessage
	end      *EndGroupMessage
	position utils.Rectangle
}

type BaseGroupMessage struct {
	message   string
	name      string
	position  utils.Rectangle
	index     int
	primary   *Participant
	secondary *Participant
	seqType   int
	group     *Group
}

type StartGroupMessage struct {
	BaseGroupMessage
}

type EndGroupMessage struct {
	BaseGroupMessage
}

func (bg *BaseGroupMessage) Type() int {
	return bg.seqType
}

func (g *Group) SetPosition(rectangle utils.Rectangle) {
	g.position = rectangle
}

func (eg *EndGroupMessage) MeasureBounds(d *Diagram, sequenceFont font.Face) utils.Rectangle {

	d.ReComputeGroup(eg.group)
	// at this point the other bounds are ready..
	eg.position = utils.Rect(0, 0, eg.group.position.Dx(), CONFIG_GROUP_BASE_HEIGHT)
	return eg.position
}

func (bg *BaseGroupMessage) PrimaryParticipant() *Participant {
	return bg.primary
}

func (bg *BaseGroupMessage) Text() string {
	return bg.message
}

func (bg *BaseGroupMessage) SecondaryParticipant() *Participant {
	return bg.secondary
}

func (bg *BaseGroupMessage) SetPosition(rectangle utils.Rectangle) {
	bg.position = rectangle
}

func (bg *BaseGroupMessage) Position() utils.Rectangle {
	return bg.position
}

func (bg *BaseGroupMessage) Render(d *Diagram, dc *gg.Context) {
	// we don't render it here.. only done by the main group
}

func (bg *BaseGroupMessage) MeasureBounds(d *Diagram, font font.Face) utils.Rectangle {

	dc := d.dc
	dc.Push()
	defer dc.Pop()
	dc.SetFontFace(font)

	textWidth, textHeight := dc.MeasureString(bg.Text())

	if textWidth > CONFIG_GROUP_MAX_WIDTH {
		lines := dc.WordWrap(bg.Text(), CONFIG_GROUP_MAX_WIDTH)
		textHeight = float64(len(lines)) * float64(textHeight+CONFIG_MESSAGE_LINE_SPACING)
		textWidth = CONFIG_GROUP_MAX_WIDTH
	}
	totalHeight := int(textHeight + CONFIG_MESSAGE_LINE_SPACING*2)

	return utils.Rect(0, 0, int(textWidth), int(totalHeight))

}

func (bg *BaseGroupMessage) IsStartProcess() bool {
	return false
}

func (bg *BaseGroupMessage) IsEndProcess() bool {
	return false
}

func (bg *BaseGroupMessage) SetStartProcess(*Process) {
	return
}

func (bg *BaseGroupMessage) SetEndProcess(*Process) {
	return
}

func (bg *BaseGroupMessage) Index() int {
	return bg.index
}

func (bg *BaseGroupMessage) Init(data map[string]interface{}, d *Diagram, index int, seqType int) error {
	bg.index = index
	bg.message = data["text"].(string)
	bg.name = data["name"].(string)
	bg.seqType = seqType
	return nil
}

func (eg *EndGroupMessage) Init(data map[string]interface{}, d *Diagram, index int, seqType int) error {
	eg.index = index
	eg.seqType = seqType
	return nil
}

func (g *Group) Position() utils.Rectangle {
	return g.position
}
func (g *Group) Name() string {
	return g.start.name
}

func (g *Group) MessageRect() utils.Rectangle {
	return utils.Rect(0, 0, 0, 0)
}

func (g *Group) Text() string {
	return g.start.Text()
}

func (g *Group) GetNameBounds() []gg.Point {
	return make([]gg.Point, 5)
}
