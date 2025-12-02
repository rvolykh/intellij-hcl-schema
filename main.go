package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rvolykh/intellij-hcl-schema/intellij"
	"github.com/rvolykh/intellij-hcl-schema/terraform"

)

//go:generate protoc --plugin=./.requirements/bin/protoc-gen-go --plugin=./.requirements/bin/protoc-gen-go-grpc --go_out=paths=source_relative,plugins=grpc:. proto/tfplugin5/tfplugin5.proto
//go:generate protoc --plugin=./.requirements/bin/protoc-gen-go --plugin=./.requirements/bin/protoc-gen-go-grpc --go_out=paths=source_relative,plugins=grpc:. proto/tfplugin6/tfplugin6.proto

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
	timeout      time.Duration

	showVersion bool
	version     = "custom"
)

func init() {
	flag.StringVar(&providerPath, "path", "", "Path to already build terraform provider")
	flag.StringVar(&providerName, "name", "", "Name to use for provider")
	flag.StringVar(&providerVers, "ver", "0.0.0", "Version to use for provider")
	flag.DurationVar(&timeout, "timeout", 60*time.Second, "Timeout for provider schema extraction")
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
		fmt.Fprintln(os.Stderr, "Error:\n  -path flag is mandatory")
		os.Exit(2)
	}
	if strings.TrimSpace(providerName) == "" {
		flag.Usage()
		fmt.Fprintln(os.Stderr, "Error:\n  -name flag is mandatory")
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	tfSchema, err := terraform.GetProviderSchema(ctx, providerPath)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Fprintf(os.Stderr, "Failed to get Provider Schema: operation timed out after %v\n", timeout)
		} else {
			fmt.Fprintf(os.Stderr, "Failed to get Provider Schema: %s\n", err)
		}
		os.Exit(1)
	}
	fmt.Println("Schema extracted from provider")

	tfSchema.Name = providerName
	tfSchema.Version = providerVers
	fmt.Println("Schema prepared for HCL plugin")

	if err := intellij.LoadToIDE(tfSchema); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save Provider Schema: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Schema is ready to use. Please, restart your IDE")
}
