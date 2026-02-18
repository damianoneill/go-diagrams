package main

import (
	"log"

	"github.com/damianoneill/go-diagrams/diagram"
	"github.com/damianoneill/go-diagrams/nodes/k8s"
)

func main() {
	d, err := diagram.New(
		diagram.Label("Kubernetes"),
		diagram.Filename("k8s"),
		diagram.Direction("TB"),
	)

	if err != nil {
		log.Fatal(err)
	}

	ingress := k8s.Network.Ing(diagram.NodeLabel("nginx"))
	svc := k8s.Network.Svc(diagram.NodeLabel("http"))

	d.Connect(ingress, svc)

	g := diagram.NewGroup("pods").Label("Deployment").Connect(svc, k8s.Compute.Pod(diagram.NodeLabel("web server")))

	d.Group(g)

	if err := d.Render(); err != nil {
		log.Fatal(err)
	}
}
