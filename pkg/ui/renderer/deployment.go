package renderer

import (
	"fmt"
	"io"
	"text/tabwriter"
)

type DeploymentLike interface {
	GetID() string
	GetName() string
	GetOwner() string
	GetStatus() string
}

func renderDeployments(w io.Writer, deployments []DeploymentLike) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tName\tOwner\tStatus")

	for _, d := range deployments {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", d.GetID(), d.GetName(), d.GetOwner(), d.GetStatus())
	}

	return tw.Flush()
}

type DeploymentAdapter struct {
	ID, Name, Owner, Status string
}

func (d DeploymentAdapter) GetID() string     { return d.ID }
func (d DeploymentAdapter) GetName() string   { return d.Name }
func (d DeploymentAdapter) GetOwner() string  { return d.Owner }
func (d DeploymentAdapter) GetStatus() string { return d.Status }
