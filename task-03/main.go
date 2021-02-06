package main

import (
	"regexp"
	"strings"
)

// NginxBlock contains block information
type NginxBlock struct {
	StartLine   string
	EndLine     string
	AllContents string
	// split lines by \n on AllContents,
	// use make to create *[],
	// first create make([]*Type..)
	// then use &var to make it *
	AllLines          *[]*string
	NestedBlocks      []*NginxBlock
	TotalBlocksInside int
}

// IsBlock returns true if line is a block
func (ngBlock *NginxBlock) IsBlock(line string) bool {
	matched, _ := regexp.MatchString(`{|}`, line)
	return matched
}

// IsLine returns true if line is not a block
func (ngBlock *NginxBlock) IsLine(line string) bool {
	return !ngBlock.IsBlock(line)
}

// HasComment return true if the line contains comment
func (ngBlock *NginxBlock) HasComment(line string) bool {
	matched, _ := regexp.MatchString(`#`, line)
	return matched
}

// NginxBlocks contains all the blocks
type NginxBlocks struct {
	blocks      *[]*NginxBlock
	AllContents string
	// split lines by \n on AllContents
	AllLines *[]*string
}

// GetNginxBlock returns an NginxBlock
func GetNginxBlock(lines *[]*string, startIndex, endIndex, recursionMax int) *NginxBlock {
	_, block := getNginxBlockHelper(*lines, startIndex, endIndex, recursionMax)
	return block
}

func getNginxBlockHelper(lines []*string, startIndex, endIndex, recursionMax int) (int, *NginxBlock) {
	b := &NginxBlock{}
	b.StartLine = *lines[startIndex]
	startIndex++
	depth := 0

	for ; startIndex < endIndex; startIndex++ {
		if strings.Contains(*lines[startIndex], "{") {
			if recursionMax > 0 {
				lastIndex, nestedBlock := getNginxBlockHelper(lines, startIndex, endIndex, recursionMax-1)
				startIndex = lastIndex
				b.NestedBlocks = append(b.NestedBlocks, nestedBlock)
			} else {
				depth++
				b.AllContents += *lines[startIndex]
				*b.AllLines = append(*b.AllLines, lines[startIndex])
			}
		} else if strings.Contains(*lines[startIndex], "}") {
			if depth == 0 {
				b.EndLine = *lines[startIndex]
				b.TotalBlocksInside = len(b.NestedBlocks)
				// count blocks
				total := 0
				for _, ch := range b.NestedBlocks {
					if ch != nil {
						total += ch.TotalBlocksInside
					}
				}
				b.TotalBlocksInside += total
				return startIndex, b
			}
			depth--
			b.AllContents += *lines[startIndex]
			*b.AllLines = append(*b.AllLines, lines[startIndex])
		} else {
			b.AllContents += *lines[startIndex]
			*b.AllLines = append(*b.AllLines, lines[startIndex])
		}
	}
	return startIndex, b
}

// GetNginxBlocks return NginxBlocks
func GetNginxBlocks(configContent string) *NginxBlocks {
	lines := strings.Split(configContent, "\n")
	linesPtr := make([]*string, len(lines))
	for i := 0; i < len(lines); i++ {
		linesPtr[i] = &lines[i]
	}
	nbs := &NginxBlocks{}
	nbs.AllLines = &linesPtr
	nbs.AllContents = configContent
	nbs.blocks = &GetNginxBlock(nbs.AllLines, 0, len(*nbs.AllLines), int((^uint(0))>>1)).NestedBlocks

	return nbs
}
func main() {
}
