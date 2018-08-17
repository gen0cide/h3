package h3

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// Unfamiliar means the user had zero knowledge of the topic
	Unfamiliar Level = iota

	// Beginner means the user had a beginner level understanding of the topic
	Beginner

	// Moderate means the user had a moderate level understanding of the topic
	Moderate

	// Advanced means the user had an advanced level understanding of the topic
	Advanced

	// Expert means the user had an expert level understanding of the topic
	Expert
)

var (
	// ErrInvalidLevel is thrown when a level could not be parsed from it's string
	ErrInvalidLevel = errors.New("invalid level supplied")

	// ValidLevels is a map of valid levels and their descriptions
	ValidLevels = map[Level]string{
		Unfamiliar: "Subject was unfamiliar with the subject",
		Beginner:   "Subject had a beginner level understanding of the subject",
		Moderate:   "Subject had a moderate level understanding of the subject",
		Advanced:   "Subject had an advanced level understanding of the subject",
		Expert:     "Subject had an expert level understanding of the subject",
	}
)

// Level represents the skill level that can be assigned to a skill in the tree
type Level int

// String implements the Stringer interface for Level
func (l Level) String() string {
	switch l {
	case Unfamiliar:
		return "unfamiliar"
	case Beginner:
		return "beginner"
	case Moderate:
		return "moderate"
	case Advanced:
		return "advanced"
	case Expert:
		return "expert"
	default:
		return "unknown"
	}
}

// LevelFromWord attempts to parse a Level from it's string representation
func LevelFromWord(w string) (Level, error) {
	lw := strings.ToLower(strings.TrimSpace(w))
	switch lw {
	case "unfamiliar":
		return Unfamiliar, nil
	case "beginner":
		return Beginner, nil
	case "moderate":
		return Moderate, nil
	case "advanced":
		return Advanced, nil
	case "expert":
		return Expert, nil
	default:
		return -1, ErrInvalidLevel
	}
}

// MarshalJSON implements the JSONMarshaler for the Level type
func (l Level) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", strings.ToUpper(l.String()))), nil
}

// UnmarshalJSON implements the JSONUnmarshaler for the Level type
func (l Level) UnmarshalJSON(b []byte) error {
	lev, err := LevelFromWord(strings.Replace(string(b), "\"", "", -1))
	if err != nil {
		return err
	}
	l = lev
	return nil
}
