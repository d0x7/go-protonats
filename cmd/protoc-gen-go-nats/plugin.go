package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "print the version and exit")
	flag.BoolVar(&showVersion, "v", false, "print the version and exit")
	flag.Parse()
	if showVersion {
		fmt.Printf("%s %s\n", filepath.Base(os.Args[0]), version)
		return
	}

	var (
		flags flag.FlagSet
	)
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL |
			pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
		gen.SupportedEditionsMinimum = descriptorpb.Edition_EDITION_PROTO3
		gen.SupportedEditionsMaximum = descriptorpb.Edition_EDITION_2023

		for _, f := range gen.Files {
			if f.Generate {
				generateFile(gen, f)
			}
		}
		return nil
	})
}
