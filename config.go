/*
Copyright © 2025 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

const (
	VERSION    string = "0.0.1"
	Baseuri    string = "https://www.amazon.de"
	Defaultdir string = "."

	DefaultTemplate string = `template.html` // builtin

	DefaultUserAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	DefaultWishlistFileTemplate string = "{{.Name}}.html"

	DefaultOutdirTemplate string = "."

	// we extract the slug from the uri
	SlugURIPartNum int = 6

	WIN string = "windows"

	// for download throttling
	MinThrottle int = 2
	MaxThrottle int = 20
)

var DirsVisited map[string]int

const Usage string = `This is amazon-my-ass, the wish list downloader.`

type Config struct {
	Verbose              bool   `koanf:"verbose"` // loglevel=info
	Debug                bool   `koanf:"debug"`   // loglevel=debug
	Showversion          bool   `koanf:"version"` // -v
	Showhelp             bool   `koanf:"help"`    // -h
	Outdir               string `koanf:"outdir"`
	Template             string `koanf:"template"`
	WishlistFiletemplate string `koanf:"wishlisttemplate"`
	Loglevel             string `koanf:"loglevel"`
	Limit                int    `koanf:"limit"`
	IgnoreErrors         bool   `koanf:"ignoreerrors"`
	ForceDownload        bool   `koanf:"force"`
	UserAgent            string `koanf:"useragent"` // conf only
	Wishlinks            []string
	StatsCountWishlist   int
}

func (c *Config) IncrWishlist() {
	c.StatsCountWishlist++
}

// load commandline flags and config file
func InitConfig(output io.Writer) (*Config, error) {
	var kloader = koanf.New(".")

	// determine template based on os
	template := DefaultTemplate

	// Load default values using the confmap provider.
	if err := kloader.Load(confmap.Provider(map[string]interface{}{
		"template":         template,
		"outdir":           DefaultOutdirTemplate,
		"loglevel":         "notice",
		"wishlisttemplate": DefaultWishlistFileTemplate,
		"useragent":        DefaultUserAgent,
	}, "."), nil); err != nil {
		return nil, fmt.Errorf("failed to load default values into koanf: %w", err)
	}

	// setup custom usage
	flagset := flag.NewFlagSet("config", flag.ContinueOnError)
	flagset.Usage = func() {
		fmt.Fprintln(output, Usage)
		os.Exit(0)
	}

	// parse commandline flags
	flagset.StringP("config", "c", "", "config file")
	flagset.StringP("outdir", "o", "", "directory where to store ads")
	flagset.IntP("limit", "l", 0, "limit wishlists to be downloaded (default 0, unlimited)")
	flagset.BoolP("verbose", "v", false, "be verbose")
	flagset.BoolP("debug", "d", false, "enable debug log")
	flagset.BoolP("version", "V", false, "show program version")
	flagset.BoolP("help", "h", false, "show usage")
	flagset.BoolP("force", "f", false, "force")
	flagset.BoolP("ignoreerrors", "", false, "ignore image download HTTP errors")

	if err := flagset.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse program arguments: %w", err)
	}

	// generate a  list of config files to try  to load, including the
	// one provided via -c, if any
	var configfiles []string

	configfile, _ := flagset.GetString("config")
	home, _ := os.UserHomeDir()

	if configfile != "" {
		configfiles = []string{configfile}
	} else {
		configfiles = []string{
			"/etc/amazon-my-ass.conf", "/usr/local/etc/amazon-my-ass.conf", // unix variants
			filepath.Join(home, ".config", "amazon-my-ass", "config"),
			filepath.Join(home, ".amazon-my-ass"),
			"amazon-my-ass.conf",
		}
	}

	// Load the config file[s]
	for _, cfgfile := range configfiles {
		path, err := os.Stat(cfgfile)

		if err != nil {
			// ignore non-existent files, but bail out on any other errors
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to stat config file: %w", err)
			}

			continue
		}

		if !path.IsDir() {
			if err := kloader.Load(file.Provider(cfgfile), toml.Parser()); err != nil {
				return nil, fmt.Errorf("error loading config file: %w", err)
			}
		}
	}

	// env overrides config file
	if err := kloader.Load(env.Provider("AMAZONMYASS_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, "AMAZONMYASS_")), "_", ".")
	}), nil); err != nil {
		return nil, fmt.Errorf("error loading environment: %w", err)
	}

	// command line overrides env
	if err := kloader.Load(posflag.Provider(flagset, ".", kloader), nil); err != nil {
		return nil, fmt.Errorf("error loading flags: %w", err)
	}

	// fetch values
	conf := &Config{}
	if err := kloader.Unmarshal("", &conf); err != nil {
		return nil, fmt.Errorf("error unmarshalling: %w", err)
	}

	// adjust loglevel
	switch conf.Loglevel {
	case "verbose":
		conf.Verbose = true
	case "debug":
		conf.Debug = true
	}

	// are there any args left on commandline? if so threat them as adlinks
	conf.Wishlinks = flagset.Args()

	return conf, nil
}
