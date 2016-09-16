package vimlparser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/haya14busa/go-vimlparser/ast"
	internal "github.com/haya14busa/go-vimlparser/go"
	"github.com/haya14busa/go-vimlparser/internal/exporter"
)

// ParseOption is option for Parse().
type ParseOption struct {
	Neovim bool
}

// ParseFile parses Vim script.
func ParseFile(r io.Reader, opt *ParseOption) (node *ast.File, err error) {
	defer func() {
		if r := recover(); r != nil {
			node = nil
			err = fmt.Errorf("go-vimlparser:Parse: %v", r)
			// log.Printf("%s", debug.Stack())
		}
	}()
	lines := readlines(r)
	reader := internal.NewStringReader(lines)
	neovim := false
	if opt != nil {
		neovim = opt.Neovim
	}
	node = exporter.NewNode(internal.NewVimLParser(neovim).Parse(reader)).(*ast.File)
	return
}

// Parse parses Vim script. TODO: depricated
func Parse(r io.Reader, opt *ParseOption) (node *BaseNode, err error) {
	defer func() {
		if r := recover(); r != nil {
			node = nil
			err = fmt.Errorf("go-vimlparser:Parse: %v", r)
			// log.Printf("%s", debug.Stack())
		}
	}()
	lines := readlines(r)
	reader := internal.NewStringReader(lines)
	neovim := false
	if opt != nil {
		neovim = opt.Neovim
	}
	node = newNode(internal.NewVimLParser(neovim).Parse(reader))
	return
}

// ParseExpr parses Vim expression. TODO: depricated
func ParseExpr(r io.Reader) (node *BaseNode, err error) {
	defer func() {
		if r := recover(); r != nil {
			node = nil
			err = fmt.Errorf("go-vimlparser:Parse: %v", r)
			// log.Printf("%s", debug.Stack())
		}
	}()
	lines := readlines(r)
	reader := internal.NewStringReader(lines)
	p := internal.NewExprParser(reader)
	node = newNode(p.Parse())
	return
}

// Compile compiles Vim script AST into S-expression like format. TODO: depricated
func Compile(w io.Writer, node *BaseNode) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("go-vimlparser:Compile: %v", r)
			// log.Printf("%s", debug.Stack())
		}
	}()
	c := internal.NewCompiler()
	out := c.Compile(newExportNode(node))
	_, err = w.Write([]byte(strings.Join(out, "\n")))
	return nil
}

func readlines(r io.Reader) []string {
	lines := []string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
