package h3

import (
	"strings"

	"github.com/kennygrant/sanitize"

	"github.com/DiSiqueira/GoTree"

	"github.com/iancoleman/strcase"

	"github.com/dariubs/percent"
)

// Element describes a single leaf in the skill tree
type Element struct {
	Name        string              `json:"name"`
	CN          string              `json:"cn"`
	Path        string              `json:"path"`
	Type        string              `json:"type"`
	Title       string              `json:"title,omitempty"`
	Description string              `json:"description,omitempty"`
	Weight      int                 `json:"weight,omitempty"`
	Parent      *Element            `json:"-"`
	Universe    *Universe           `json:"-"`
	Children    map[string]*Element `json:"children,omitempty"`
	UserScore   Level               `json:"user_score,omitempty"`
	Rubric      Rubric              `json:"rubric,omitempty"`
	Tree        gotree.Tree         `json:"-"`
}

// NewChild creates a new element leaf from an element
func (e *Element) NewChild(name string) *Element {
	normalized := strcase.ToSnake(sanitize.Path(name))
	if c, ok := e.Children[normalized]; ok {
		return c
	}
	c := &Element{
		Name:     name,
		Parent:   e,
		CN:       normalized,
		Children: map[string]*Element{},
		Universe: e.Universe,
		Tree:     e.Tree.Add(normalized),
	}
	e.Universe.Elms = append(e.Universe.Elms, c)
	e.Children[normalized] = c
	return c
}

// Count returns the number of elements in this leaf
func (e *Element) Count() int {
	if len(e.Children) == 0 {
		return 1
	}
	count := 0
	for _, x := range e.Children {
		count += x.Count()
	}
	return count
}

// Percent calculates the delta between max score of this element and the user percent
func (e *Element) Percent() float64 {
	count := e.Count()
	max := count * int(Expert)
	score := e.Score()
	return percent.PercentOf(score, max)
}

// Score calculates the user supplied score of this element (or it's tree)
func (e *Element) Score() int {
	if len(e.Children) == 0 {
		return int(e.UserScore)
	}

	score := 0
	for _, x := range e.Children {
		score += x.Score()
	}

	return score
}

// AbsPath returns an elements qualified path from root
func (e *Element) AbsPath() string {
	path := []string{e.CN}
	parent := e.Parent
	for {
		if parent == nil {
			break
		}
		path = append([]string{parent.CN}, path...)
		parent = parent.Parent
	}
	return strings.Join(path, ".")
}
