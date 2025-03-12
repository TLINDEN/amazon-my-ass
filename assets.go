package main

import (
	"bytes"
	"embed"
	"log"
	"path"
)

//go:embed assets/*
var assetfs embed.FS

// put loaded collector assets into global map
var Assets = LoadAssets()

func LoadAssets() map[string]string {
	assets := map[string]string{}

	entries, err := assetfs.ReadDir("assets")
	if err != nil {
		log.Fatalf("failed to read assets dir assets: %s", err)
	}

	for _, file := range entries {
		path := path.Join("assets", file.Name())

		if file.IsDir() {
			continue
		}

		fd, err := assetfs.Open(path)
		if err != nil {
			log.Fatalf("failed to open embedded file %s: %s", path, err)
		}
		defer fd.Close()

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(fd)
		if err != nil {
			log.Fatalf("failed to read embedded file %s: %s", path, err)
		}

		name := file.Name()

		assets[name] = buf.String()
	}

	return assets
}
