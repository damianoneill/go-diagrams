package diagram

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	graphviz "github.com/awalterschulze/gographviz"
)

type Connector interface {
	Connect(start, end *Node, opts ...EdgeOption) Connector
	ConnectByID(start, end string, opts ...EdgeOption) Connector
}

type Diagram struct {
	options Options

	g *graphviz.Escape

	root *Group
}

func New(opts ...Option) (*Diagram, error) {
	options := DefaultOptions(opts...)
	g := graphviz.NewEscape()
	g.SetName("root")
	g.SetDir(true)

	for k, v := range options.attrs() {
		if err := g.AddAttr("root", k, v); err != nil {
			return nil, err
		}
	}

	return newDiagram(g, options), nil
}

func newDiagram(g *graphviz.Escape, options Options) *Diagram {
	return &Diagram{
		g:       g,
		options: options,
		root:    newGroup("root", 0, nil),
	}
}

// SetOutputPath overrides the output directory for rendered files.
// This is useful when you need to control the exact output location,
// for example in tests or when rendering to a temp directory.
func (d *Diagram) SetOutputPath(path string) {
	d.options.Name = path
}

func (d *Diagram) Nodes() []*Node {
	return d.root.Nodes()
}

func (d *Diagram) Edges() []*Edge {
	return d.root.Edges()
}

func (d *Diagram) Groups() []*Group {
	return d.root.Children()
}

func (d *Diagram) Add(ns ...*Node) *Diagram {
	d.root.Add(ns...)
	return d
}

func (d *Diagram) Connect(start, end *Node, opts ...EdgeOption) *Diagram {
	d.Add(start, end)
	return d.ConnectByID(start.ID(), end.ID(), opts...)
}

func (d *Diagram) ConnectByID(start, end string, opts ...EdgeOption) *Diagram {
	d.root.ConnectByID(start, end, opts...)

	return d
}

func (d *Diagram) Group(g *Group) *Diagram {
	d.root.Group(g)
	return d
}

func (d *Diagram) Close() error {
	return nil
}

func (d *Diagram) Render() error {
	return d.render()
}

func (d *Diagram) render() error {
	outdir := d.options.Name
	if err := os.MkdirAll(outdir, os.ModePerm); err != nil {
		return err
	}

	for _, n := range d.root.nodes {
		err := n.render("root", outdir, d.g)
		if err != nil {
			return err
		}
	}

	for _, e := range d.root.edges {
		err := e.render(e.Start(), e.End(), d.g)
		if err != nil {
			return err
		}
	}

	for _, g := range d.root.children {
		if err := g.render(outdir, d.g); err != nil {
			return err
		}
	}

	return d.renderOutput()
}

func (d *Diagram) renderOutput() error {
	// Always save the DOT file first
	if err := d.saveDot(); err != nil {
		return err
	}

	switch d.options.OutFormat {
	case "dot":
		return nil // DOT file already saved
	case "png", "jpg", "svg", "pdf":
		return d.renderImage()
	default:
		return fmt.Errorf("unsupported output format: %s", d.options.OutFormat)
	}
}

// renderImage invokes the Graphviz dot binary to render the DOT file
// to the configured image format (png, svg, jpg, pdf).
func (d *Diagram) renderImage() error {
	dotFile := filepath.Join(d.options.Name, d.options.FileName+".dot")
	outFile := filepath.Join(d.options.Name, d.options.FileName+"."+d.options.OutFormat)

	cmd := exec.Command("dot", "-T"+d.options.OutFormat, "-o", outFile, dotFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("graphviz rendering failed: %w: %s", err, string(output))
	}

	return nil
}

func (d *Diagram) saveDot() error {
	fname := filepath.Join(d.options.Name, d.options.FileName+".dot")

	return os.WriteFile(fname, []byte(d.g.String()), os.ModePerm)
}
