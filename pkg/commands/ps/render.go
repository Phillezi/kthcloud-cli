package ps

import (
	"fmt"
	"os"
	"strings"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kthcloud/go-deploy/dto/v2/body"
)

func renderDeplsTable(depls []body.DeploymentRead, all bool) {

	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Status", "Ping result", "Visibility", "Updated", "Visit"})

	for _, dep := range depls {
		if dep.Status == "resourceDisabled" && !all {
			continue
		}

		dep.Status = strings.Replace(dep.Status, "resource", "", 1)

		pingRes := ""
		if dep.PingResult != nil {
			pingRes = fmt.Sprintf("%d", *dep.PingResult)
		}

		visitStr := ""
		if dep.URL != nil {
			visitStr = fmt.Sprintf("\u001b]8;;%s\u0007%s\u001b]8;;\u0007", *dep.URL, "visit")
		}

		t.AppendRow(table.Row{dep.ID, dep.Name, dep.Status, pingRes, dep.Visibility, util.TimeAgo(dep.UpdatedAt), visitStr})
		t.AppendSeparator()
	}

	t.Render()
}
