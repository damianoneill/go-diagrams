# Go-Diagrams

Create beautiful system diagrams with Go.

This is a modernised fork of [blushft/go-diagrams](https://github.com/blushft/go-diagrams), which is a loose port of the Python [diagrams](https://github.com/mingrammer/diagrams) library.

## What changed in this fork

| Area | Before (blushft) | After (this fork) |
|------|-------------------|-------------------|
| Go version | 1.14 | 1.22+ |
| Asset embedding | fileb0x (generated code, 1138 files) | embed.FS (standard library) |
| Deprecated APIs | ioutil.ReadFile, ioutil.WriteFile, manual rand seeding | os.ReadFile, os.WriteFile, math/rand/v2 |
| Output formats | DOT only (manual dot command for images) | DOT, PNG, SVG, JPG, PDF (direct rendering via Graphviz) |
| Dependencies | fileb0x, golang.org/x/net, golang.org/x/exp, pkg/errors | Removed - only gographviz, uuid, jennifer, strcase, testify remain |

The API is backward compatible. Existing code that imports blushft/go-diagrams only needs an import path change.

## Prerequisites

- **Go 1.22** or later
- **Graphviz** (required for image rendering; DOT output works without it)

macOS:

    brew install graphviz

Debian / Ubuntu:

    apt-get install graphviz

Alpine (Docker):

    apk add graphviz

## Installation

    go get github.com/damianoneill/go-diagrams

## Quick start

    package main

    import (
        "log"

        "github.com/damianoneill/go-diagrams/diagram"
        "github.com/damianoneill/go-diagrams/nodes/generic"
    )

    func main() {
        d, err := diagram.New(
            diagram.Label("Simple Network"),
            diagram.Filename("network"),
            diagram.Direction("LR"),
            diagram.OutputFormat("png"),
        )
        if err != nil {
            log.Fatal(err)
        }

        fw := generic.Network.Firewall(diagram.NodeLabel("Firewall"))
        rt := generic.Network.Router(diagram.NodeLabel("Router"))
        sw := generic.Network.Switch(diagram.NodeLabel("Switch"))

        d.Connect(fw, rt, diagram.Forward())
        d.Connect(rt, sw, diagram.Forward())

        if err := d.Render(); err != nil {
            log.Fatal(err)
        }
        // Output: go-diagrams/network.dot + go-diagrams/network.png
    }

## Features

### Output formats

Set the output format with diagram.OutputFormat():

    // Render directly to PNG (requires Graphviz)
    d, _ := diagram.New(diagram.OutputFormat("png"))

    // Or SVG, JPG, PDF
    d, _ := diagram.New(diagram.OutputFormat("svg"))

    // DOT only (no Graphviz required) - this is the default
    d, _ := diagram.New(diagram.OutputFormat("dot"))

When using png, svg, jpg, or pdf, the library invokes the dot binary from Graphviz automatically. The DOT source file is always written alongside the rendered image.

### Custom output directory

By default, output goes into a go-diagrams/ directory. Override with SetOutputPath():

    d, _ := diagram.New(diagram.Filename("my-diagram"), diagram.OutputFormat("png"))
    d.SetOutputPath("/tmp/output")
    d.Render() // writes to /tmp/output/my-diagram.dot and /tmp/output/my-diagram.png

### Clusters (groups)

Group nodes into visual clusters:

    d, _ := diagram.New(diagram.Label("App"), diagram.Filename("app"), diagram.OutputFormat("png"))

    dns := gcp.Network.Dns(diagram.NodeLabel("DNS"))
    lb := gcp.Network.LoadBalancing(diagram.NodeLabel("NLB"))
    cache := gcp.Database.Memorystore(diagram.NodeLabel("Cache"))
    db := gcp.Database.Sql(diagram.NodeLabel("Database"))

    dc := diagram.NewGroup("GCP")
    dc.NewGroup("services").
        Label("Service Layer").
        Add(
            gcp.Compute.ComputeEngine(diagram.NodeLabel("Server 1")),
            gcp.Compute.ComputeEngine(diagram.NodeLabel("Server 2")),
            gcp.Compute.ComputeEngine(diagram.NodeLabel("Server 3")),
        ).
        ConnectAllFrom(lb.ID(), diagram.Forward()).
        ConnectAllTo(cache.ID(), diagram.Forward())

    dc.NewGroup("data").Label("Data Layer").Add(cache, db).Connect(cache, db)

    d.Connect(dns, lb, diagram.Forward()).Group(dc)
    d.Render()

Produces:

![app-diagram](images/app-diagram.png)

### Available node providers

| Provider | Import path | Examples |
|----------|-------------|---------|
| Generic | nodes/generic | Network (Router, Switch, Firewall, VPN), Compute, Database, Storage |
| AWS | nodes/aws | Compute, Database, Network, Storage, ML, and more |
| GCP | nodes/gcp | Compute, Database, Network, Storage, ML, and more |
| Azure | nodes/azure | Compute, Database, Network, Storage, and more |
| Kubernetes | nodes/k8s | Compute, Network, Storage, Control Plane, RBAC |
| Apps | nodes/apps | Docker, Kafka, Redis, Nginx, Jenkins, Git, and more |
| Alibaba Cloud | nodes/alibabacloud | Compute, Database, Network, Storage |
| Oracle Cloud | nodes/oci | Compute, Database, Network, Storage |
| OpenStack | nodes/openstack | Compute, Networking, Storage, and more |
| Firebase | nodes/firebase | Develop, Grow, Quality |
| Elastic | nodes/elastic | Elasticsearch, Observability, Security |
| SaaS | nodes/saas | Chat, Logging, Monitoring, CDN, Identity |
| Programming | nodes/programming | Languages and Frameworks |

### Edge options

    // Directional edges
    d.Connect(a, b, diagram.Forward())       // a -> b
    d.Connect(a, b, diagram.Reverse())       // a <- b
    d.Connect(a, b, diagram.Bidirectional()) // a <-> b

## More examples

See the [_examples](_examples) directory for complete working examples:

- **[app](_examples/app)** - GCP application architecture
- **[webservice](_examples/webservice)** - Web service with load balancing
- **[workers](_examples/workers)** - Worker pool pattern
- **[kubernetes](_examples/kubernetes)** - Kubernetes deployment

## Migrating from blushft/go-diagrams

1. Update your import paths from github.com/blushft/go-diagrams to github.com/damianoneill/go-diagrams.

2. Add diagram.OutputFormat("png") to render images directly instead of running the dot command manually.

Everything else is API-compatible.

## License

[MIT](LICENSE)
