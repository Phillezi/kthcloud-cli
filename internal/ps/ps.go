package ps

import (
	"fmt"
	"go-deploy/dto/v2/body"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	log "github.com/sirupsen/logrus"

	"github.com/Phillezi/kthcloud-cli/internal/model"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/spf13/viper"
)

func Ps(all bool) {
	// load the session from the session.json file
	session, err := model.Load(viper.GetString("session-path"))
	if err != nil {
		log.Fatalln("No active session. Please log in")
	}
	if session.AuthSession.IsExpired() {
		log.Fatalln("Session is expired. Please log in again")
	}
	session.SetupClient()

	resp, err := session.Client.Req("/v2/deployments", "GET", nil)
	if err != nil {
		log.Fatal(err)
	}

	depls, err := util.ProcessResponseArr[body.DeploymentRead](resp.String())
	if err != nil {
		log.Fatal(err)
	}

	renderTable(depls, all)

}

func renderTable(depls []body.DeploymentRead, all bool) {

	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Status", "Ping result", "Visibility"})

	for _, dep := range depls {
		if dep.Status == "resourceDisabled" && !all {
			continue
		}

		dep.Status = strings.Replace(dep.Status, "resource", "", 1)

		pingRes := ""
		if dep.PingResult != nil {
			pingRes = fmt.Sprintf("%d", *dep.PingResult)
		}

		t.AppendRow(table.Row{dep.ID, dep.Name, dep.Status, pingRes, dep.Visibility})
		t.AppendSeparator()
	}

	t.Render()
}
