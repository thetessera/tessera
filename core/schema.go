package core

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Schema describes which blocks and props a page may contain.
type Schema struct {
	Tessera int              `json:"tessera"`
	Name    string           `json:"name"`
	Version int              `json:"version"`
	Limits  Limits           `json:"limits"`
	Root    Block            `json:"root"`
	Blocks  map[string]Block `json:"blocks"`
}

// Limits caps the size of a whole page.
type Limits struct {
	MaxBlocks int `json:"maxBlocks"`
	MaxDepth  int `json:"maxDepth"`
}

// Block describes one block type: its props and the slots that hold other
// blocks. The page root has the same shape.
type Block struct {
	Label       string          `json:"label"`
	Description string          `json:"description"`
	Props       map[string]Prop `json:"props"`
	Slots       map[string]Slot `json:"slots"`
}

// Slot describes a place inside a block where other blocks can be put.
type Slot struct {
	Min   int                  `json:"min"`
	Max   *int                 `json:"max"`
	Allow map[string]AllowRule `json:"allow"`
}

// AllowRule limits how one block type may be used inside a slot.
type AllowRule struct {
	Min      int    `json:"min"`
	Max      *int   `json:"max"`
	Position string `json:"position"`
}

// Prop describes one field of a block. Which constraints apply depends on Type.
type Prop struct {
	// all types
	Type     string `json:"type"`
	Label    string `json:"label"`
	Required bool   `json:"required"`

	// string, richtext
	MaxLength *int `json:"maxLength"`

	// url
	Protocols []string `json:"protocols"`
	Hosts     []string `json:"hosts"`

	// enum
	Options []string `json:"options"`
	Default string   `json:"default"`

	// object
	Fields map[string]Prop `json:"fields"`

	// list
	Of       *Prop `json:"of"`
	MinItems int   `json:"minItems"`
	MaxItems *int  `json:"maxItems"`

	// richtext
	Nodes []string  `json:"nodes"`
	Marks []string  `json:"marks"`
	Links *LinkRule `json:"links"`
}

// LinkRule says which links a richtext prop may contain.
type LinkRule struct {
	Protocols []string `json:"protocols"`
}

// ParseSchema reads a schema from JSON and checks that it is well formed.
// Unknown keys are rejected, so a misspelled constraint is an error and not
// a silently missing limit.
func ParseSchema(data []byte) (*Schema, error) {
	var schema Schema
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&schema); err != nil {
		return nil, fmt.Errorf("parse schema: %w", err)
	}

	if err := schema.check(); err != nil {
		return nil, err
	}
	return &schema, nil
}

func (s *Schema) check() error {
	if s.Version < 1 {
		return fmt.Errorf("version: must be 1 or greater, got %d", s.Version)
	}
	if len(s.Name) == 0 {
		return fmt.Errorf("name: is required")
	}
	if s.Tessera != 1 {
		return fmt.Errorf("tessera: unsupported format version %d", s.Tessera)
	}
	if err := s.checkBlock("root", s.Root); err != nil {
		return err
	}
	for name, block := range s.Blocks {
		if err := s.checkBlock("blocks."+name, block); err != nil {
			return err
		}
	}
	return nil
}

func (s *Schema) checkBlock(path string, b Block) error {
	for slotName, slot := range b.Slots {
		for blockType := range slot.Allow {
			if _, ok := s.Blocks[blockType]; !ok {
				return fmt.Errorf("%s.slots.%s.allow: unknown block %q", path, slotName, blockType)
			}
		}
	}
	for propName, prop := range b.Props {
		if err := s.checkProp(path+".props."+propName, prop); err != nil {
			return err
		}
	}
	return nil
}

func (s *Schema) checkProp(path string, p Prop) error {
	switch p.Type {
	case "string", "richtext":
	case "url":
		if len(p.Protocols) == 0 {
			return fmt.Errorf("%s: url requires protocols", path)
		}
	case "enum":
		if len(p.Options) == 0 {
			return fmt.Errorf("%s: enum requires options", path)
		}
	case "object":
		if len(p.Fields) == 0 {
			return fmt.Errorf("%s: object requires fields", path)
		}
		for name, field := range p.Fields {
			if err := s.checkProp(path+".fields."+name, field); err != nil {
				return err
			}
		}
	case "list":
		if p.Of == nil {
			return fmt.Errorf("%s: list requires of", path)
		}
		err := s.checkProp(path+".of", *p.Of)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("%s: unknown property type %q", path, p.Type)
	}
	return nil
}
