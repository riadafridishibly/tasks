package main

import (
	"strings"
	"testing"
)

func TestNginxBlock_IsBlock(t *testing.T) {
	type fields struct {
		StartLine         string
		EndLine           string
		AllContents       string
		AllLines          *[]*string
		NestedBlocks      []*NginxBlock
		TotalBlocksInside int
	}
	type args struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Block Opening 1",
			args: args{
				"server {",
			},
			want: true,
		},
		{
			name: "Block Opeing 2",
			args: args{
				"localhost { # some comment",
			},
			want: true,
		},
		{
			name: "Block Closing 1",
			args: args{
				"} ",
			},
			want: true,
		},
		{
			name: "Block Closing 2",
			args: args{
				"} # some comment",
			},
			want: true,
		},
		{
			name: "Not a block 1",
			args: args{
				"# some comment",
			},
			want: false,
		},
		{
			name: "Not a block 2",
			args: args{
				"port 8080;",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ngBlock := &NginxBlock{
				StartLine:         tt.fields.StartLine,
				EndLine:           tt.fields.EndLine,
				AllContents:       tt.fields.AllContents,
				AllLines:          tt.fields.AllLines,
				NestedBlocks:      tt.fields.NestedBlocks,
				TotalBlocksInside: tt.fields.TotalBlocksInside,
			}
			if got := ngBlock.IsBlock(tt.args.line); got != tt.want {
				t.Errorf("NginxBlock.IsBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNginxBlock_IsLine(t *testing.T) {
	type fields struct {
		StartLine         string
		EndLine           string
		AllContents       string
		AllLines          *[]*string
		NestedBlocks      []*NginxBlock
		TotalBlocksInside int
	}
	type args struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Is a line",
			args: args{
				line: "",
			},
			want: true,
		},
		{
			name: "Is a block",
			args: args{
				line: "hello {",
			},
			want: false,
		},
		{
			name: "Is a normal line",
			args: args{
				line: "port 5432;",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ngBlock := &NginxBlock{
				StartLine:         tt.fields.StartLine,
				EndLine:           tt.fields.EndLine,
				AllContents:       tt.fields.AllContents,
				AllLines:          tt.fields.AllLines,
				NestedBlocks:      tt.fields.NestedBlocks,
				TotalBlocksInside: tt.fields.TotalBlocksInside,
			}
			if got := ngBlock.IsLine(tt.args.line); got != tt.want {
				t.Errorf("NginxBlock.IsLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNginxBlock_HasComment(t *testing.T) {
	type fields struct {
		StartLine         string
		EndLine           string
		AllContents       string
		AllLines          *[]*string
		NestedBlocks      []*NginxBlock
		TotalBlocksInside int
	}
	type args struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "No comment",
			args: args{
				line: "server {",
			},
			want: false,
		},
		{
			name: "Comment",
			args: args{
				line: "server { # some comment",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ngBlock := &NginxBlock{
				StartLine:         tt.fields.StartLine,
				EndLine:           tt.fields.EndLine,
				AllContents:       tt.fields.AllContents,
				AllLines:          tt.fields.AllLines,
				NestedBlocks:      tt.fields.NestedBlocks,
				TotalBlocksInside: tt.fields.TotalBlocksInside,
			}
			if got := ngBlock.HasComment(tt.args.line); got != tt.want {
				t.Errorf("NginxBlock.HasComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func splitToString(data string) *[]*string {
	lines := strings.Split(data, "\n")
	lptr := make([]*string, len(lines))
	for i := 0; i < len(lines); i++ {
		lptr[i] = &lines[i]
	}
	return &lptr
}

// Need to add some generated test!
func TestGetNginxBlock(t *testing.T) {
	type args struct {
		lines        *[]*string
		startIndex   int
		endIndex     int
		recursionMax int
	}
	tests := []struct {
		name string
		args args
		want *NginxBlock
	}{
		{
			name: "Simple Nginx Config Start Line & End Line",
			args: args{
				lines:        splitToString("ab {\n}"),
				startIndex:   0,
				endIndex:     2,
				recursionMax: 4,
			},
			want: &NginxBlock{
				StartLine:         "ab {",
				EndLine:           "}",
				AllContents:       "ab {\n}",
				AllLines:          nil,
				NestedBlocks:      nil,
				TotalBlocksInside: 0,
			},
		},
		{
			name: "Simple Nginx Config Start Line & End Line",
			args: args{
				lines:        splitToString("ab {\n}"),
				startIndex:   0,
				endIndex:     2,
				recursionMax: 4,
			},
			want: &NginxBlock{
				StartLine:         "ab {",
				EndLine:           "}",
				AllContents:       "ab {\n}",
				AllLines:          nil,
				NestedBlocks:      nil,
				TotalBlocksInside: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNginxBlock(tt.args.lines, tt.args.startIndex, tt.args.endIndex, tt.args.recursionMax); got.StartLine != tt.want.StartLine && got.EndLine != tt.want.EndLine {
				t.Errorf("GetNginxBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNginxBlocks(t *testing.T) {
	type args struct {
		configContent string
	}
	tests := []struct {
		name string
		args args
		want *NginxBlocks
	}{
		{
			name: "All content test",
			args: args{
				configContent: "abc {\ndef{\n}\n}",
			},
			want: &NginxBlocks{
				AllContents: "abc {\ndef{\n}\n}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNginxBlocks(tt.args.configContent); got.AllContents != tt.want.AllContents {
				t.Errorf("GetNginxBlocks() = %v, want %v", got, tt.want)
			}
		})
	}
}
