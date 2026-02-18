package diagram_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/damianoneill/go-diagrams/diagram"
	"github.com/damianoneill/go-diagrams/nodes/generic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEdgeLabelOption(t *testing.T) {
	opts := diagram.DefaultEdgeOptions(diagram.EdgeLabel("10Gbps"))
	assert.Equal(t, "10Gbps", opts.Label)
}

func TestEdgeColorOption(t *testing.T) {
	opts := diagram.DefaultEdgeOptions(diagram.EdgeColor("#FF0000"))
	assert.Equal(t, "#FF0000", opts.Color)
}

func TestEdgeStyleOption(t *testing.T) {
	opts := diagram.DefaultEdgeOptions(diagram.EdgeStyle("dashed"))
	assert.Equal(t, "dashed", opts.Style)
}

func TestEdgeFontColorOption(t *testing.T) {
	opts := diagram.DefaultEdgeOptions(diagram.EdgeFontColor("#00FF00"))
	assert.Equal(t, "#00FF00", opts.Font.Color)
}

func TestEdgeAttributeOption(t *testing.T) {
	opts := diagram.DefaultEdgeOptions(diagram.EdgeAttribute("constraint", "false"))
	assert.Equal(t, "false", opts.Attributes["constraint"])
}

func TestEdgeOptionsCombined(t *testing.T) {
	opts := diagram.DefaultEdgeOptions(
		diagram.EdgeLabel("PoE"),
		diagram.EdgeColor("#333333"),
		diagram.EdgeStyle("dotted"),
		diagram.Forward(),
	)
	assert.Equal(t, "PoE", opts.Label)
	assert.Equal(t, "#333333", opts.Color)
	assert.Equal(t, "dotted", opts.Style)
	assert.True(t, opts.Forward)
	assert.False(t, opts.Reverse)
}

func TestEdgeOptionsOverrideDefaults(t *testing.T) {
	// Default color is #7B8894; override it.
	defaults := diagram.DefaultEdgeOptions()
	assert.Equal(t, "#7B8894", defaults.Color)

	overridden := diagram.DefaultEdgeOptions(diagram.EdgeColor("#000000"))
	assert.Equal(t, "#000000", overridden.Color)
}

func TestEdgeDirectionOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    []diagram.EdgeOption
		forward bool
		reverse bool
	}{
		{"Forward", []diagram.EdgeOption{diagram.Forward()}, true, false},
		{"Reverse", []diagram.EdgeOption{diagram.Reverse()}, false, true},
		{"Bidirectional", []diagram.EdgeOption{diagram.Bidirectional()}, true, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opts := diagram.DefaultEdgeOptions(tc.opts...)
			assert.Equal(t, tc.forward, opts.Forward)
			assert.Equal(t, tc.reverse, opts.Reverse)
		})
	}
}

func TestEdgeMultipleAttributes(t *testing.T) {
	opts := diagram.DefaultEdgeOptions(
		diagram.EdgeAttribute("constraint", "false"),
		diagram.EdgeAttribute("weight", "2"),
		diagram.EdgeAttribute("penwidth", "3.0"),
	)
	assert.Equal(t, "false", opts.Attributes["constraint"])
	assert.Equal(t, "2", opts.Attributes["weight"])
	assert.Equal(t, "3.0", opts.Attributes["penwidth"])
}

// renderDOTWithEdges is a helper that creates a simple diagram with two nodes
// connected by an edge using the given EdgeOptions, renders it to DOT, and
// returns the DOT content as a string.
func renderDOTWithEdges(t *testing.T, edgeOpts ...diagram.EdgeOption) string {
	t.Helper()

	dir := t.TempDir()
	outDir := filepath.Join(dir, "output")

	d, err := diagram.New(
		diagram.Label("Edge Test"),
		diagram.Filename("edge-test"),
	)
	require.NoError(t, err)
	d.SetOutputPath(outDir)

	fw := generic.Network.Firewall(diagram.NodeLabel("Firewall"))
	sw := generic.Network.Switch(diagram.NodeLabel("Switch"))
	d.Connect(fw, sw, edgeOpts...)

	require.NoError(t, d.Render())

	dotFile := filepath.Join(outDir, "edge-test.dot")
	data, err := os.ReadFile(dotFile)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	return string(data)
}

func TestEdgeLabelRenderedInDOT(t *testing.T) {
	dot := renderDOTWithEdges(t, diagram.Forward(), diagram.EdgeLabel("10Gbps"))
	assert.Contains(t, dot, "10Gbps", "DOT output should contain edge label")
}

func TestEdgeColorRenderedInDOT(t *testing.T) {
	dot := renderDOTWithEdges(t, diagram.Forward(), diagram.EdgeColor("#FF0000"))
	assert.Contains(t, dot, "#FF0000", "DOT output should contain edge color")
}

func TestEdgeStyleRenderedInDOT(t *testing.T) {
	dot := renderDOTWithEdges(t, diagram.Forward(), diagram.EdgeStyle("dashed"))
	assert.Contains(t, dot, "dashed", "DOT output should contain edge style")
}

func TestEdgeFontColorRenderedInDOT(t *testing.T) {
	dot := renderDOTWithEdges(t, diagram.Forward(), diagram.EdgeFontColor("#00FF00"))
	assert.Contains(t, dot, "#00FF00", "DOT output should contain edge font color")
}

func TestEdgeAttributeRenderedInDOT(t *testing.T) {
	dot := renderDOTWithEdges(t, diagram.Forward(), diagram.EdgeAttribute("penwidth", "3.0"))
	assert.Contains(t, dot, "penwidth", "DOT output should contain custom attribute name")
	assert.Contains(t, dot, "3.0", "DOT output should contain custom attribute value")
}

func TestEdgeAllOptionsRenderedInDOT(t *testing.T) {
	dot := renderDOTWithEdges(t,
		diagram.Forward(),
		diagram.EdgeLabel("PoE"),
		diagram.EdgeColor("#333333"),
		diagram.EdgeStyle("dotted"),
		diagram.EdgeFontColor("#111111"),
	)
	assert.Contains(t, dot, "PoE")
	assert.Contains(t, dot, "#333333")
	assert.Contains(t, dot, "dotted")
	assert.Contains(t, dot, "#111111")
}

func TestEdgeDefaultOmitsEmptyFields(t *testing.T) {
	// When no label or style is set, those fields should not appear in the DOT
	// output because trimAttrs strips empty strings.
	dot := renderDOTWithEdges(t, diagram.Forward())

	// Find the edge line (contains "->")
	lines := strings.Split(dot, "\n")
	var edgeLine string
	for _, line := range lines {
		if strings.Contains(line, "->") {
			edgeLine = line
			break
		}
	}
	require.NotEmpty(t, edgeLine, "should find an edge line in DOT output")

	// label and style should be absent (trimmed) when not explicitly set
	assert.NotContains(t, edgeLine, "label=", "empty label should be trimmed")
	assert.NotContains(t, edgeLine, "style=", "empty style should be trimmed")
}

func TestEdgeBidirectionalRenderedInDOT(t *testing.T) {
	dot := renderDOTWithEdges(t, diagram.Bidirectional(), diagram.EdgeLabel("sync"))
	assert.Contains(t, dot, "both", "DOT output should contain dir=both for bidirectional")
	assert.Contains(t, dot, "sync", "DOT output should contain the edge label")
}

func TestEdgeReverseRenderedInDOT(t *testing.T) {
	dot := renderDOTWithEdges(t, diagram.Reverse())
	assert.Contains(t, dot, "back", "DOT output should contain dir=back for reverse")
}
