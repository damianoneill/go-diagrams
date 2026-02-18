package diagram_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/damianoneill/go-diagrams/diagram"
	"github.com/damianoneill/go-diagrams/nodes/generic"
)

func TestRenderDot(t *testing.T) {
	dir := t.TempDir()
	outDir := filepath.Join(dir, "output")

	d, err := diagram.New(
		diagram.Label("Test"),
		diagram.Filename("test-dot"),
	)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	d.SetOutputPath(outDir)

	fw := generic.Network.Firewall(diagram.NodeLabel("fw"))
	sw := generic.Network.Switch(diagram.NodeLabel("sw"))
	d.Connect(fw, sw, diagram.Forward())

	if err := d.Render(); err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	dotFile := filepath.Join(outDir, "test-dot.dot")
	data, err := os.ReadFile(dotFile)
	if err != nil {
		t.Fatalf("failed to read DOT file: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("DOT file is empty")
	}
	t.Logf("DOT file size: %d bytes", len(data))
}

func TestRenderPNG(t *testing.T) {
	if _, err := exec.LookPath("dot"); err != nil {
		t.Skip("graphviz not installed, skipping PNG render test")
	}

	dir := t.TempDir()
	outDir := filepath.Join(dir, "output")

	d, err := diagram.New(
		diagram.Label("PNG Test"),
		diagram.Filename("test-png"),
		diagram.OutputFormat("png"),
	)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	d.SetOutputPath(outDir)

	fw := generic.Network.Firewall(diagram.NodeLabel("Firewall"))
	sw := generic.Network.Switch(diagram.NodeLabel("Switch"))
	rt := generic.Network.Router(diagram.NodeLabel("Router"))

	d.Connect(fw, rt, diagram.Forward())
	d.Connect(rt, sw, diagram.Forward())

	if err := d.Render(); err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	pngFile := filepath.Join(outDir, "test-png.png")
	info, err := os.Stat(pngFile)
	if err != nil {
		t.Fatalf("PNG file not found: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("PNG file is empty")
	}
	t.Logf("PNG file size: %d bytes", info.Size())
}
