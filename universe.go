package h3

import (
	"fmt"

	gotree "github.com/DiSiqueira/GoTree"
	"github.com/iancoleman/strcase"
	"github.com/kennygrant/sanitize"
)

var (
	csvColumns = []string{
		"type",
		"path",
		"common_name",
		"weight",
		"title",
		"description",
		"rubric_unfamiliar",
		"rubric_beginner",
		"rubric_moderate",
		"rubric_advanced",
		"rubric_expert",
	}
)

// Universe manages the top level skill tree
type Universe struct {
	Elms     []*Element          `json:"-"`
	ByPath   map[string]*Element `json:"-"`
	TopLevel map[string]*Element `json:"universe"`
	Tree     gotree.Tree         `json:"-"`
}

// NewUniverse makes a new universe.
func NewUniverse() *Universe {
	return &Universe{
		Elms:     []*Element{},
		ByPath:   map[string]*Element{},
		TopLevel: map[string]*Element{},
		Tree:     gotree.New("."),
	}
}

// NewTopLevelElement creates a new top level element in the universe
func (u *Universe) NewTopLevelElement(name string) *Element {
	normalized := strcase.ToSnake(sanitize.Path(name))
	if c, ok := u.ByPath[normalized]; ok {
		return c
	}
	c := &Element{
		Name:     name,
		CN:       normalized,
		Children: map[string]*Element{},
		Universe: u,
		Tree:     u.Tree.Add(normalized),
	}
	u.Elms = append(u.Elms, c)
	u.TopLevel[normalized] = c
	return c
}

// Normalize attempts to normalize the fields within the universe
func (u *Universe) Normalize() error {
	for _, x := range u.Elms {
		if x.Parent != nil {
			x.Path = x.AbsPath()
		}
		if len(x.Children) > 0 {
			x.Type = "CATEGORY"
			x.Weight = 1
			continue
		}
		x.Type = "COMPETENCY"
		x.Title = "__DEFINE_COMPETENCY_TITLE__"
		x.Description = "__DEFINE_COMPETENCY_DESCRIPTION__"
		x.Rubric = EmptyRubric()
	}
	return nil
}

// ToTable returns a two dimensional table
func (u *Universe) ToTable() [][]string {
	t := [][]string{csvColumns}
	n := "N/A"
	for _, x := range u.Elms {
		var weight, title, desc, rub0, rub1, rub2, rub3, rub4 string
		switch x.Type {
		case "CATEGORY":
			weight = fmt.Sprintf("%d", x.Weight)
			title, desc, rub0, rub1, rub2, rub3, rub4 = n, n, n, n, n, n, n
		case "COMPETENCY":
			weight = n
			title = x.Title
			desc = x.Description
			rub0 = x.Rubric[Level(0)]
			rub1 = x.Rubric[Level(1)]
			rub2 = x.Rubric[Level(2)]
			rub3 = x.Rubric[Level(3)]
			rub4 = x.Rubric[Level(4)]
		}
		r := []string{
			x.Type,
			x.Path,
			x.CN,
			weight,
			title,
			desc,
			rub0,
			rub1,
			rub2,
			rub3,
			rub4,
		}
		t = append(t, r)
	}
	return t
}
