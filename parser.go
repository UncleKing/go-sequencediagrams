package sequence

import (
	"fmt"
	"regexp"
	"strings"
)

// ParseLine validated at
// https://regex-golang.appspot.com/assets/html/index.html
func ParseLine(str string) (string, int, error) {

	// A -> B: Message
	solid := regexp.MustCompile(`^\s*([a-zA-Z0-9]+)\s*->\s*([a-zA-Z0-9]+)\s*:\s*([a-zA-Z0-9\s]+)\s*$`)

	match := solid.FindAllStringSubmatch(str, -1)
	if len(match) > 0 {
		return fmt.Sprintf(`{"src": "%s","dest":"%s", "type":"solid", "text": "%s"}`, match[0][1], match[0][2], match[0][3]), ST_SOLID, nil

	}

	// A --> B : Message
	dotted := regexp.MustCompile(`^\s*([a-zA-Z0-9]+)\s*-->\s*([a-zA-Z0-9]+)\s*:\s*([a-zA-Z0-9\s]+)\s*$`)
	match = dotted.FindAllStringSubmatch(str, -1)
	if len(match) > 0 {
		return fmt.Sprintf(`{"src": "%s","dest":"%s", "type":"dotted", "text": "%s"}`, match[0][1], match[0][2], match[0][3]), ST_DOTTED, nil
	}

	//A ->+ B: Message
	startProcess := regexp.MustCompile(`^\s*([a-zA-Z0-9]+)\s*->\+\s*([a-zA-Z0-9]+)\s*:\s*([a-zA-Z0-9\s]+)\s*$`)
	match = startProcess.FindAllStringSubmatch(str, -1)
	if len(match) > 0 {
		return fmt.Sprintf(`{"src": "%s","dest":"%s", "type":"start_process", "text": "%s"}`, match[0][1], match[0][2], match[0][3]), ST_START_PROCESS, nil
	}

	// A ->- B: Message
	endProcess := regexp.MustCompile(`^\s*([a-zA-Z0-9]+)\s*->\-\s*([a-zA-Z0-9]+)\s*:\s*([a-zA-Z0-9\s]+)\s*$`)

	match = endProcess.FindAllStringSubmatch(str, -1)
	if len(match) > 0 {
		return fmt.Sprintf(`{"src": "%s","dest":"%s", "type":"end_process", "text": "%s"}`, match[0][1], match[0][2], match[0][3]), ST_END_PROCESS, nil
	}

	// A -->+ B: Message
	startDottedProcess := regexp.MustCompile(`^\s*([a-zA-Z0-9]+)\s*-->\+\s*([a-zA-Z0-9]+)\s*:\s*([a-zA-Z0-9\s]+)\s*$`)
	match = startDottedProcess.FindAllStringSubmatch(str, -1)
	if len(match) > 0 {
		return fmt.Sprintf(`{"src": "%s","dest":"%s", "type":"start_dotted_process", "text": "%s"}`, match[0][1], match[0][2], match[0][3]), ST_START_DOTTED_PROCESS, nil
	}

	// A -->- B: Message
	endDottedProcess := regexp.MustCompile(`^\s*([a-zA-Z0-9]+)\s*-->\-\s*([a-zA-Z0-9]+)\s*:\s*([a-zA-Z0-9\s]+)\s*$`)
	match = endDottedProcess.FindAllStringSubmatch(str, -1)
	if len(match) > 0 {
		return fmt.Sprintf(`{"src": "%s","dest":"%s", "type":"end_dotted_process", "text": "%s"}`, match[0][1], match[0][2], match[0][3]), ST_END_DOTTED_PROCESS, nil
	}

	// note over A, B: Message
	notesOver := regexp.MustCompile(`^\s*note over\s*([,\w\s]*)\s*:\s*([\w\s]+)$`)

	match = notesOver.FindAllStringSubmatch(str, -1)
	if len(match) > 0 {
		csvKeys := strings.Split(match[0][1], ",")
		// i wish i paid more attention in the regex classes.
		var keys []string
		for _, n := range csvKeys {
			keys = append(keys, fmt.Sprintf(`"%s"`, strings.Trim(n, " ")))
		}
		allKeys := strings.Join(keys, ",")
		return fmt.Sprintf(`{"src": [%s],"type":"notes", "side": "top", "text": "%s"}`, allKeys, match[0][2]), ST_NOTE_OVER, nil
	}

	// note right of  A : message
	notesRightOf := regexp.MustCompile(`^\s*note right of\s*([\w]*)\s*:\s*([\w\s]+)$`)

	match = notesRightOf.FindAllStringSubmatch(str, -1)
	if len(match) > 0 {
		return fmt.Sprintf(`{"src": ["%s"],"type":"notes", "side": "right", "text": "%s"}`, match[0][1], match[0][2]), ST_NOTE_RIGHT, nil
	}

	// note left of A : Message
	notesLeftOf := regexp.MustCompile(`^\s*note left of\s*([\w]*)\s*:\s*([\w\s]+)$`)
	match = notesLeftOf.FindAllStringSubmatch(str, -1)
	if len(match) > 0 {
		return fmt.Sprintf(`{"src": ["%s"],"type":"notes", "side": "left", "text": "%s"}`, match[0][1], match[0][2]), ST_NOTE_LEFT, nil
	}

	alt := regexp.MustCompile(`^\s*alt \s*(.*)$`)
	match = alt.FindAllStringSubmatch(str, -1)
	if len(match) > 0 {
		return fmt.Sprintf(`{"type":"group","text": "%s", "name": "alt"}`, match[0][1]), ST_GROUP_MESSAGE, nil
	}

	loop := regexp.MustCompile(`^\s*loop \s*(.*)$`)
	match = loop.FindAllStringSubmatch(str, -1)
	if len(match) > 0 {
		return fmt.Sprintf(`{"type":"group","text": "%s", "name": "loop"}`, match[0][1]), ST_GROUP_MESSAGE, nil
	}

	stelse := regexp.MustCompile(`^\s*else \s*(.*)$`)
	match = stelse.FindAllStringSubmatch(str, -1)
	if len(match) > 0 {
		return fmt.Sprintf(`{"type":"else","text": "%s"}`, match[0][1]), ST_ELSE_MESSAGE, nil
	}

	end := regexp.MustCompile(`^\s*end\s*$`)
	match = end.FindAllStringSubmatch(str, -1)
	if len(match) > 0 {
		return `{"type":"end"}`, ST_END_GROUP, nil
	}

	return "", 0, fmt.Errorf("No matching format supported")
}
