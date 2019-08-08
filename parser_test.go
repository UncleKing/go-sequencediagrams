package sequence

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadSolid(t *testing.T) {
	actualOutput := make(map[string]interface{})
	outputJsonObj := make(map[string]interface{})
	inputStr := "a->n: Mes sage"
	outputJson := `{"src": "a","dest":"n", "type":"solid", "text": "Mes sage"}`
	err := json.Unmarshal([]byte(outputJson), &outputJsonObj)

	assert.NoError(t, err)

	output, typ, err := ParseLine(inputStr)
	assert.Equal(t, typ, ST_SOLID)
	err = json.Unmarshal([]byte(output), &actualOutput)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, outputJsonObj, actualOutput)

	inputStr = "a -> n : Mes sage"
	outputJson = `{"src": "A","dest":"B", "type":"solid", "text": "Mes sage"}`
	json.Unmarshal([]byte(outputJson), outputJsonObj)

	output, typ, err = ParseLine(inputStr)
	assert.Equal(t, typ, ST_SOLID)
	err = json.Unmarshal([]byte(output), &actualOutput)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, outputJsonObj, actualOutput)

}

func TestReadDotted(t *testing.T) {
	actualOutput := make(map[string]interface{})
	outputJsonObj := make(map[string]interface{})
	inputStr := "A-->B:Mess age"
	outputJson := `{"src": "A","dest":"B", "type":"dotted", "text": "Mess age"}`
	err := json.Unmarshal([]byte(outputJson), &outputJsonObj)

	assert.NoError(t, err)

	output, typ, err := ParseLine(inputStr)
	assert.Equal(t, typ, ST_DOTTED)
	err = json.Unmarshal([]byte(output), &actualOutput)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, outputJsonObj, actualOutput)

	inputStr = "A --> B : Mess age"
	outputJson = `{"src": "A","dest":"B", "type":"dotted", "text": "Mess age"}`
	json.Unmarshal([]byte(outputJson), outputJsonObj)

	output, typ, err = ParseLine(inputStr)
	err = json.Unmarshal([]byte(output), &actualOutput)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, outputJsonObj, actualOutput)

}

func TestReadDottedStart(t *testing.T) {
	actualOutput := make(map[string]interface{})
	outputJsonObj := make(map[string]interface{})
	inputStr := "A-->+B:Mess age"
	outputJson := `{"src": "A","dest":"B", "type":"start_dotted_process", "text": "Mess age"}`
	err := json.Unmarshal([]byte(outputJson), &outputJsonObj)

	assert.NoError(t, err)

	output, typ, err := ParseLine(inputStr)
	assert.Equal(t, typ, ST_START_DOTTED_PROCESS)
	err = json.Unmarshal([]byte(output), &actualOutput)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, outputJsonObj, actualOutput)

	inputStr = "A -->+ B : Mess age"
	outputJson = `{"src": "A","dest":"B", "type":"start_dotted_process", "text": "Mess age"}`
	json.Unmarshal([]byte(outputJson), outputJsonObj)

	output, typ, err = ParseLine(inputStr)
	assert.Equal(t, typ, ST_START_DOTTED_PROCESS)
	err = json.Unmarshal([]byte(output), &actualOutput)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, outputJsonObj, actualOutput)

}

func TestReadDottedEnd(t *testing.T) {
	actualOutput := make(map[string]interface{})
	outputJsonObj := make(map[string]interface{})
	inputStr := "A-->-B:Mess age"
	outputJson := `{"src": "A","dest":"B", "type":"end_dotted_process", "text": "Mess age"}`
	err := json.Unmarshal([]byte(outputJson), &outputJsonObj)

	assert.NoError(t, err)

	output, typ, err := ParseLine(inputStr)
	assert.Equal(t, typ, ST_END_DOTTED_PROCESS)
	err = json.Unmarshal([]byte(output), &actualOutput)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, outputJsonObj, actualOutput)

	inputStr = "A -->- B : Mess age"
	outputJson = `{"src": "A","dest":"B", "type":"end_dotted_process", "text": "Mess age"}`
	json.Unmarshal([]byte(outputJson), outputJsonObj)

	output, typ, err = ParseLine(inputStr)
	assert.Equal(t, typ, ST_END_DOTTED_PROCESS)
	err = json.Unmarshal([]byte(output), &actualOutput)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, outputJsonObj, actualOutput)

}

func TestNotesLeft(t *testing.T) {
	actualOutput := make(map[string]interface{})
	outputJsonObj := make(map[string]interface{})
	inputStr := "note left of A:Message B"
	outputJson := `{"src": ["A"], "type":"notes", "side": "left", "text": "Message B"}`
	err := json.Unmarshal([]byte(outputJson), &outputJsonObj)

	assert.NoError(t, err)

	output, typ, err := ParseLine(inputStr)
	assert.Equal(t, typ, ST_NOTE_LEFT)
	err = json.Unmarshal([]byte(output), &actualOutput)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, outputJsonObj, actualOutput)

	inputStr = "note left of A : Message B"
	outputJson = `{"src": ["A"], "type":"notes", "side": "left", "text": "Message B"}`
	json.Unmarshal([]byte(outputJson), outputJsonObj)

	output, typ, err = ParseLine(inputStr)
	assert.Equal(t, typ, ST_NOTE_LEFT)
	err = json.Unmarshal([]byte(output), &actualOutput)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, outputJsonObj, actualOutput)
}

func TestNotesRight(t *testing.T) {
	actualOutput := make(map[string]interface{})
	outputJsonObj := make(map[string]interface{})
	inputStr := "note right of A:Message B"
	outputJson := `{"src": ["A"], "type":"notes", "side": "right", "text": "Message B"}`
	err := json.Unmarshal([]byte(outputJson), &outputJsonObj)

	assert.NoError(t, err)

	output, typ, err := ParseLine(inputStr)
	assert.Equal(t, typ, ST_NOTE_RIGHT)
	err = json.Unmarshal([]byte(output), &actualOutput)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, outputJsonObj, actualOutput)

	inputStr = "note right of A : Message B"
	outputJson = `{"src": ["A"], "type":"notes", "side": "right", "text": "Message B"}`
	json.Unmarshal([]byte(outputJson), outputJsonObj)

	output, typ, err = ParseLine(inputStr)
	assert.Equal(t, typ, ST_NOTE_RIGHT)
	err = json.Unmarshal([]byte(output), &actualOutput)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, outputJsonObj, actualOutput)
}

func TestNotesLeftMultiLine(t *testing.T) {

}
func TestNotesRightMultiLine(t *testing.T) {

}

func TestNotesOver(t *testing.T) {
	actualOutput := make(map[string]interface{})
	outputJsonObj := make(map[string]interface{})
	inputStr := "note over A:Message B"
	outputJson := `{"src": ["A"], "type":"notes", "side": "top", "text": "Message B"}`
	err := json.Unmarshal([]byte(outputJson), &outputJsonObj)

	assert.NoError(t, err)

	output, typ, err := ParseLine(inputStr)
	assert.Equal(t, typ, ST_NOTE_OVER)
	err = json.Unmarshal([]byte(output), &actualOutput)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, outputJsonObj, actualOutput)

	outputJsonObj = make(map[string]interface{})
	inputStr = "note over A : Message B"
	outputJson = `{"src": ["A"], "type":"notes", "side": "top", "text": "Message B"}`
	err = json.Unmarshal([]byte(outputJson), &outputJsonObj)
	assert.NoError(t, err)
	output, typ, err = ParseLine(inputStr)
	assert.Equal(t, typ, ST_NOTE_OVER)
	err = json.Unmarshal([]byte(output), &actualOutput)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, outputJsonObj, actualOutput)

	outputJsonObj = make(map[string]interface{})
	inputStr = "note over A,C : Message B"
	outputJson = `{"src": ["A", "C"], "type":"notes", "side": "top", "text": "Message B"}`
	err = json.Unmarshal([]byte(outputJson), &outputJsonObj)
	assert.NoError(t, err)
	output, typ, err = ParseLine(inputStr)
	assert.Equal(t, typ, ST_NOTE_OVER)
	err = json.Unmarshal([]byte(output), &actualOutput)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, outputJsonObj, actualOutput)

	outputJsonObj = make(map[string]interface{})
	inputStr = "note over A, C, D : Message B"
	outputJson = `{"src": ["A", "C", "D"], "type":"notes", "side": "top", "text": "Message B"}`
	err = json.Unmarshal([]byte(outputJson), &outputJsonObj)
	assert.NoError(t, err)
	output, typ, err = ParseLine(inputStr)
	assert.Equal(t, typ, ST_NOTE_OVER)
	err = json.Unmarshal([]byte(output), &actualOutput)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, outputJsonObj, actualOutput)

}

func TestNotesOverMultiLine(t *testing.T) {

}

func TestNotesOverMultipleMultiLine(t *testing.T) {

}

//
//func TestImagesNormal(t *testing.T) {
//	infile1, err := os.Open("test/blank.gif")
//	assert.NoError(t, err)
//	defer infile1.Close()
//
//	infile2, err := os.Open("test/image2.gif")
//	assert.NoError(t, err)
//	defer infile2.Close()
//
//	//blank, err := os.Open("test/blank.gif")
//	//assert.NoError(t, err)
//	//defer blank.Close()
//
//	same, err := compareImages(infile1, infile2)
//	assert.NoError(t, err)
//	assert.Equal(t, true, same)
//
//	//same, err = compareImages(infile1, blank)
//	//assert.NoError(t, err)
//	//assert.Equal(t, true, same)
//}
//
//
//
//func compareImages(img1, img2 io.Reader) (bool, error) {
//
//	image1, format1, err1 := image.Decode(img1)
//	image2, format2, err2 := image.Decode(img2)
//
//	if err1 != nil {
//		return false, fmt.Errorf("Image failed to decode err1: %s", err1.Error())
//	}
//	if err2 != nil {
//		return false, fmt.Errorf("Image failed to decode err2: %s", err2.Error())
//	}
//	fmt.Printf("image1 format: %s, image2 format: %s", format1, format2)
//	if(image1.Bounds().Size() != image2.Bounds().Size()) {
//		return false, fmt.Errorf("Image size is different image1: %s, image2: %s", image1.Bounds().Size(), image2.Bounds().Size())
//	}
//	bounds := image1.Bounds()
//	w, h := bounds.Max.X, bounds.Max.Y
//	for x := 0; x < w; x++ {
//		for y := 0; y < h; y++ {
//			r1, g1, b1, a1 := image1.At(x, y).RGBA()
//			r2, g2, b2, a2 := image2.At(x, y).RGBA()
//			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
//				return false, fmt.Errorf("Images are different")
//			}
//		}
//	}
//
//	return true, nil
//}
