package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gen0cide/h3"
	"github.com/urfave/cli"
)

var (
	discovercmd = cli.Command{
		Name:      "discover",
		Usage:     "Creates an empty JSON database and editable CSV file for competencies",
		UsageText: "h3 discover INPUT_FILE JSON_OUTPUT CSV_OUTPUT TREE_OUTPUT",
		Action:    discoverCommand,
	}
)

func discoverCommand(c *cli.Context) error {
	if len(c.Args()) != 4 {
		return errors.New("must provide three arguments: INPUT_FILE JSON_OUTPUT CSV_OUTPUT TREE_OUTPUT")
	}
	data, err := ioutil.ReadFile(c.Args().Get(0))
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")

	u := h3.NewUniverse()
	var e *h3.Element
	var tope *h3.Element
	curOffset := 0

	for idx, x := range lines {
		offset := strings.Count(x, "\t")
		name := strings.Replace(x, "\t", "", -1)
		switch {
		case offset == 0:
			e = u.NewTopLevelElement(name)
			tope = e
			continue
		case offset == 1:
			e = tope.NewChild(name)
			curOffset = 1
			continue
		case offset == curOffset:
			e = e.Parent.NewChild(name)
			continue
		case offset == (curOffset + 1):
			e = e.NewChild(name)
			curOffset++
			continue
		case offset < curOffset:
			delta := curOffset - offset
			for z := 0; z < delta; z++ {
				e = e.Parent
			}
			curOffset -= delta
			e = e.Parent.NewChild(name)
			continue
		case offset > (curOffset + 1):
			spew.Dump(offset)
			spew.Dump(curOffset)
			spew.Dump(x)
			spew.Dump(lines[idx-1])
		default:
			l.Errorf("FML")
			spew.Dump(offset)
			spew.Dump(curOffset)
			spew.Dump(x)
			spew.Dump(lines[idx-1])
		}
	}

	err = u.Normalize()
	if err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(u.TopLevel, "", "  ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(c.Args().Get(1), jsonData, 0644)

	fw, err := os.Create(c.Args().Get(2))
	if err != nil {
		return err
	}
	csvWriter := csv.NewWriter(fw)
	err = csvWriter.WriteAll(u.ToTable())

	if err = csvWriter.Error(); err != nil {
		return err
	}

	err = ioutil.WriteFile(c.Args().Get(3), []byte(u.Tree.Print()), 0644)
	return err
}
