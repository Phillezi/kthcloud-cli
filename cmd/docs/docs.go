package main

import (
	"os"

	"github.com/Phillezi/kthcloud-cli/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra/doc"
)

const (
	docDir = "./docs"
)

func main() {
	if err := os.MkdirAll(docDir, 0755); err != nil {
		logrus.Fatalf("failed to create docs directory: %v", err)
	}

	if err := doc.GenMarkdownTree(cmd.GetRootCMD(), docDir); err != nil {
		logrus.Fatalf("failed to generate markdown docs: %v", err)
	}

	logrus.Println("Docs generated successfully")
}
