package strings

import (
	"strings"
)

// MqttMatch answers true if the topic matches the pattern according
// to MQTT rules:
// * Pattern is a slash-separated hierarchy.
// Wildcard "+" matches a single level in the hierarchy.
// Wildcard "#" matches all further levels in the hierarchy.
//
// Examples:
// Pattern a/b matches topic a/b
// a/# matches a/b, a/c, a/b/c etc
// a/+/c matches a/b/c, a/c/c, but not a/b, not a/b/c/c, etc.
func MqttMatch(pattern, topic string) bool {
	pattern, nextPattern := mqttChunk(pattern)
	nextTopic := topic
	for len(pattern) > 0 {
		topic, nextTopic = mqttChunk(nextTopic)
		switch pattern {
		case "#":
			return true
		case "+":
		default:
			if pattern != topic {
				return false
			}
		}

		pattern, nextPattern = mqttChunk(nextPattern)
	}
	return nextTopic == ""
}

func mqttChunk(s string) (string, string) {
	idx := strings.Index(s, "/")
	if idx < 0 {
		return s, ""
	} else if idx == 0 {
		if len(s) < 1 {
			return "", ""
		}
		return mqttChunk(s[1:])
	} else {
		next := ""
		if idx < len(s) {
			next = s[idx+1:]
		}
		return s[0:idx], next
	}
}
