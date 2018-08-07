package tfschema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/config/configschema"
	"github.com/olekukonko/tablewriter"
)

// Block is wrapper for configschema.Block.
// This ia a layer for customization not enough for Terraform's core.
// Most of the structure is the smae as the core, but some are different.
type Block struct {
	// Attributes is a map of any attributes.
	Attributes []*Attribute `json:"attributes"`
	// BlockTypes is a map of any nested block types.
	BlockTypes []*NestedBlock `json:"block_types"`
}

// NewBlock creates a new Block instance.
func NewBlock(b *configschema.Block) *Block {
	return &Block{
		Attributes: NewAttributes(b.Attributes),
		BlockTypes: NewBlockTypes(b.BlockTypes),
	}
}

// NewAttributes creates a new map of Attributes.
func NewAttributes(as map[string]*configschema.Attribute) []*Attribute {
	m := make([]*Attribute, len(as))

	i := 0
	for k, v := range as {
		m[i] = NewAttribute(v, k)
		i++
	}

	return m
}

// NewBlockTypes creates a new map of NestedBlocks.
func NewBlockTypes(bs map[string]*configschema.NestedBlock) []*NestedBlock {
	m := make([]*NestedBlock, len(bs))

	i := 0
	for k, v := range bs {
		m[i] = NewNestedBlock(v, k)
	}

	return m
}

// FormatJSON returns a formatted string in JSON format.
func (b *Block) FormatJSON() (string, error) {
	bytes, err := json.MarshalIndent(b, "", "    ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// FormatTable returns a formatted string in table format.
func (b *Block) FormatTable() (string, error) {
	return b.renderBlock()
}

// renderBlock returns a formatted string in table format for Block.
func (b *Block) renderBlock() (string, error) {
	buf := new(bytes.Buffer)
	attributes, err := b.renderAttributes()
	if err != nil {
		return "", err
	}
	buf.WriteString(attributes)

	blockTypes, err := b.renderBlockTypes()
	if err != nil {
		return "", err
	}
	buf.WriteString(blockTypes)

	return buf.String(), nil
}

// renderAttributes returns a formatted string in table format for Attributes.
func (b *Block) renderAttributes() (string, error) {
	buf := new(bytes.Buffer)
	table := tablewriter.NewWriter(buf)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	table.SetHeader([]string{"attribute", "type", "required", "optional", "computed", "sensitive"})

	for _, v := range b.Attributes {
		typeName, err := v.Type.Name()
		if err != nil {
			return "", err
		}

		row := []string{
			v.Name,
			typeName,
			strconv.FormatBool(v.Required),
			strconv.FormatBool(v.Optional),
			strconv.FormatBool(v.Computed),
			strconv.FormatBool(v.Sensitive),
		}
		table.Append(row)
	}

	table.Render()

	return buf.String(), nil
}

// renderBlockTypes returns a formatted string in table format for BlockTypes.
func (b *Block) renderBlockTypes() (string, error) {
	if len(b.BlockTypes) == 0 {
		return "", nil
	}

	buf := new(bytes.Buffer)

	for _, v := range b.BlockTypes {
		blockType := fmt.Sprintf("\nblock_type: %s, nesting: %s, min_items: %d, max_items: %d\n", v.Name, v.Nesting, v.MinItems, v.MaxItems)
		buf.WriteString(blockType)

		block, err := v.renderBlock()
		if err != nil {
			return "", err
		}

		buf.WriteString(block)
	}

	return buf.String(), nil
}
