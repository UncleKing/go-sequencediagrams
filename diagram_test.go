package sequence

import (
	"github.com/fogleman/gg"
	"github.com/stretchr/testify/assert"
	"go-sequencediagrams/utils"
	"testing"
)

func TestReadParse(t *testing.T) {
	d, err := NewDiagram()
	assert.NoError(t, err, "NewDiagram gave error !")
	seq := "a->n:message"
	err = d.Parse(seq)
	assert.NoError(t, err, "Parse gave error !")

	assert.Len(t, d.sequences, 1, "len should be 1")
	assert.Len(t, d.participants, 2, "participants should be 2")

}

func TestDiagram_ComputeParticipantSizeAndPlace(t *testing.T) {
	d, err := NewDiagram()
	dc := gg.NewContext(10, 10)
	dc.SetFontFace(d.ParticipantFont)

	assert.NoError(t, err, "NewDiagram gave error !")
	seq := "W->i:message"
	err = d.Parse(seq)
	assert.NoError(t, err, "Parse gave error !")

	assert.Len(t, d.sequences, 1, "len should be 1")
	assert.Len(t, d.participants, 2, "participants should be 2")

	d.ComputeParticipantSizeAndPlace()
	w1, h1 := dc.MeasureString("W")
	w2, h2 := dc.MeasureString("i")

	r1 := utils.Rect(0, 0, int(w1+CONFIG_TEXT_PADDING_X*2), int(h1+CONFIG_TEXT_PADDING_Y*2))
	r2 := utils.Rect(0, 0, int(w2+CONFIG_TEXT_PADDING_X*2), int(h2+CONFIG_TEXT_PADDING_Y*2))

	p1 := d.GetOrCreateParticipant("W")
	assert.Equal(t, r1, p1.position, "Rect should be of same size")
	assert.Equal(t, 0, p1.delta)

	p2 := d.GetOrCreateParticipant("i")
	assert.Equal(t, r2, p2.position, "Rect should be of same size")
	assert.Equal(t, CONFIG_MIN_PADDING_X, p2.delta)

}

func TestDiagram_MeasureSequence(t *testing.T) {
	d, err := NewDiagram()
	dc := gg.NewContext(10, 10)
	dc.SetFontFace(d.SequenceFont)
	assert.NoError(t, err, "NewDiagram gave error !")

	seq := "W->i:message"
	err = d.Parse(seq)
	assert.NoError(t, err, "Parse gave error !")

	assert.Len(t, d.sequences, 1, "len should be 1")
	assert.Len(t, d.participants, 2, "participants should be 2")
	//ractual := d.MeasureSequence(d.sequences[0])

	//w, h := dc.MeasureString("message")

	//rexpected := image.Rect(0, 0, int(w+CONFIG_TEXT_PADDING_X*2), int(h+CONFIG_TEXT_PADDING_Y*2))
	//assert.Equal(t, rexpected, ractual, "Sequence width should be of same size.")

}

func TestDiagram_GetDelta(t *testing.T) {
	d, err := NewDiagram()
	dc := gg.NewContext(10, 10)
	dc.SetFontFace(d.ParticipantFont)

	assert.NoError(t, err, "NewDiagram gave error !")
	seq := "W->i:message"
	err = d.Parse(seq)
	assert.NoError(t, err, "Parse gave error !")

	assert.Len(t, d.sequences, 1, "len should be 1")
	assert.Len(t, d.participants, 2, "participants should be 2")

	d.ComputeParticipantSizeAndPlace()
	w1, h1 := dc.MeasureString("W")
	w2, h2 := dc.MeasureString("i")

	r1 := utils.Rect(0, 0, int(w1+CONFIG_TEXT_PADDING_X*2), int(h1+CONFIG_TEXT_PADDING_Y*2))
	r2 := utils.Rect(0, 0, int(w2+CONFIG_TEXT_PADDING_X*2), int(h2+CONFIG_TEXT_PADDING_Y*2))

	p1 := d.GetOrCreateParticipant("W")
	assert.Equal(t, r1, p1.position, "Rect should be of same size")

	p2 := d.GetOrCreateParticipant("i")
	assert.Equal(t, r2, p2.position, "Rect should be of same size")

	delta, err := d.GetDelta(0, 1)
	assert.NoError(t, err, "Get Delta gave error")
	assert.Equal(t, CONFIG_MIN_PADDING_X, delta)

}

func TestDiagram_AdjustXSpace(t *testing.T) {
	d, err := NewDiagram()
	dc := gg.NewContext(10, 10)
	dc.SetFontFace(d.ParticipantFont)

	assert.NoError(t, err, "NewDiagram gave error !")
	seq := "W->i:message"
	err = d.Parse(seq)
	assert.NoError(t, err, "Parse gave error !")

	assert.Len(t, d.sequences, 1, "len should be 1")
	assert.Len(t, d.participants, 2, "participants should be 2")

	d.ComputeParticipantSizeAndPlace()
	w1, h1 := dc.MeasureString("W")
	w2, h2 := dc.MeasureString("i")

	r1 := utils.Rect(0, 0, int(w1+CONFIG_TEXT_PADDING_X*2), int(h1+CONFIG_TEXT_PADDING_Y*2))
	r2 := utils.Rect(0, 0, int(w2+CONFIG_TEXT_PADDING_X*2), int(h2+CONFIG_TEXT_PADDING_Y*2))

	p1 := d.GetOrCreateParticipant("W")
	assert.Equal(t, r1, p1.position, "Rect should be of same size")

	p2 := d.GetOrCreateParticipant("i")
	assert.Equal(t, r2, p2.position, "Rect should be of same size")

	d.AdjustXSpace(p1, p2, 200)
	assert.Equal(t, 200, p2.delta)

}

//
//func TestDiagram_ComputeSequenceMessageAndPlace(t *testing.T) {
//
//	d, err := NewDiagram()
//	dc := gg.NewContext(10, 10)
//	dc.SetFontFace(d.ParticipantFont)
//
//	assert.NoError(t, err, "NewDiagram gave error !")
//	seq := "W->i:message"
//	err = d.Parse(seq)
//	assert.NoError(t, err, "Parse gave error !")
//
//	assert.Len(t, d.sequences, 1, "len should be 1")
//	assert.Len(t, d.participants, 2, "participants should be 2")
//
//	d.ComputeParticipantSizeAndPlace()
//	w1, h1 := dc.MeasureString("W")
//	w2, h2 := dc.MeasureString("i")
//
//	r1 := image.Rect(0, 0, int(w1+CONFIG_TEXT_PADDING_X*2), int(h1+CONFIG_TEXT_PADDING_Y*2))
//	r2 := image.Rect(0, 0, int(w2+CONFIG_TEXT_PADDING_X*2), int(h2+CONFIG_TEXT_PADDING_Y*2))
//
//	d.ComputeSequenceMessageAndPlace()
//
//}

func BenchmarkCreateDiagramSlow(b *testing.B) {
	str := `A ->+ B: Start
A ->+ B: Start
B -> C : Hello
C -> D : Hello D
D -> E : Hello E
B->- A : Done
B->- A : Done`

	for i := 0; i < b.N; i++ {
		CreateDiagram(str)
	}
	//		s := NewSheet(sizeX, sizeY)
	//		_, _ = s.FindPath(fromx , fromy , tox , toy)
	//	}
}

func BenchmarkCreateDiagram(b *testing.B) {
	str := `A ->+ B: Start
A ->+ B: Start
B -> E : Hello E
B->- A : Done`

	for i := 0; i < b.N; i++ {
		CreateDiagram(str)
	}
	//		s := NewSheet(sizeX, sizeY)
	//		_, _ = s.FindPath(fromx , fromy , tox , toy)
	//	}
}
