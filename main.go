package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rvolykh/intellij-hcl-schema/intellij"
	"github.com/rvolykh/intellij-hcl-schema/terraform"
)

//go:generate protoc --plugin=./.requirements/bin/protoc-gen-go --plugin=./.requirements/bin/protoc-gen-go-grpc --go_out=paths=source_relative,plugins=grpc:. proto/tfplugin.proto

const (
	help = `Helper to build terraform provider autocompletion for Intellij HCL plugin.
Example:
	intellij-hcl-schema -path ./terraform-provider-hashicups -name hashicups
`
)

var (
	providerPath string
	providerName string
	providerVers string

	showVersion bool
	version     = "custom"
)

func init() {
	flag.StringVar(&providerPath, "path", "", "Path to already build terraform provider")
	flag.StringVar(&providerName, "name", "", "Name to use for provider")
	flag.StringVar(&providerVers, "ver", "0.0.0", "Version to use for provider")
	flag.BoolVar(&showVersion, "version", false, "Show current tool version")

	flag.Usage = func() {
		fmt.Printf("%sFlags:\n", help)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if strings.TrimSpace(providerPath) == "" {
		flag.Usage()
		fmt.Println("Error:\n  -path flag is mandatory")
		os.Exit(2)
	}
	if strings.TrimSpace(providerName) == "" {
		flag.Usage()
		fmt.Println("Error:\n  -name flag is mandatory")
		os.Exit(2)
	}

	tfSchema, err := terraform.GetProviderSchema(context.TODO(), providerPath)
	if err != nil {
		fmt.Printf("Failed to get Provider Schema: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Schema extracted from provider")

	hclSchema := intellij.Export(tfSchema)
	hclSchema.Name = providerName
	hclSchema.Version = providerVers
	fmt.Println("Schema prepared for HCL plugin")

	if err := intellij.LoadToIDE(hclSchema); err != nil {
		fmt.Printf("Failed to save Provider Schema: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Schema is ready to use. Please, restart your IDE")
}
