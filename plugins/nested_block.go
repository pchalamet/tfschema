package tfschema

import (
	"github.com/hashicorp/terraform/config/configschema"
)

// NestedBlock is wrapper for configschema.NestedBlock
type NestedBlock struct {
	Name string `json:"name"`
	// Block is a nested child block.
	Block
	// Nesting is a nesting mode.
	Nesting configschema.NestingMode `json:"nesting"`
	// MinItems is a lower limit on number of nested child blocks.
	MinItems int `json:"min_items"`
	// MaxItems is a upper limit on number of nested child blocks.
	MaxItems int `json:"max_items"`
}

// NewNestedBlock creates a new NestedBlock instance.
func NewNestedBlock(b *configschema.NestedBlock, name string) *NestedBlock {
	block := NewBlock(&b.Block)
	return &NestedBlock{
		Name:     name,
		Block:    *block,
		Nesting:  b.Nesting,
		MinItems: b.MinItems,
		MaxItems: b.MaxItems,
	}
}
