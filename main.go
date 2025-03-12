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
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/inconshreveable/mousetrap"
	"github.com/tlinden/yadu"
)

const LevelNotice = slog.Level(2)

func main() {
	os.Exit(Main(os.Stdout))
}

func init() {
	// if we're running on Windows  AND if the user double clicked the
	// exe  file from explorer, we  tell them and then  wait until any
	// key has been hit, which  will make the cmd window disappear and
	// thus give the user time to read it.
	if runtime.GOOS == "windows" {
		if mousetrap.StartedByExplorer() {
			fmt.Println("Do no double click kleingebaeck.exe!")
			fmt.Println("Please open a command shell and run it from there.")
			fmt.Println()
			fmt.Print("Press any key to quit: ")
			_, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				panic(err)
			}
		}
	}
}

func Main(output io.Writer) int {
	logLevel := &slog.LevelVar{}

	conf, err := InitConfig(output)
	if err != nil {
		return Die(err)
	}

	if conf.Showversion {
		fmt.Fprintf(output, "This is kleingebaeck version %s\n", VERSION)

		return 0
	}

	if conf.Showhelp {
		fmt.Fprintln(output, Usage)

		return 0
	}

	if conf.Verbose {
		logLevel.Set(slog.LevelInfo)
	}

	if conf.Debug {
		buildInfo, _ := debug.ReadBuildInfo()
		opts := &yadu.Options{
			Level:     logLevel,
			AddSource: true,
			//NoColor:   IsNoTty(),
		}

		logLevel.Set(slog.LevelDebug)

		handler := yadu.NewHandler(output, opts)
		debuglogger := slog.New(handler).With(
			slog.Group("program_info",
				slog.Int("pid", os.Getpid()),
				slog.String("go_version", buildInfo.GoVersion),
			),
		)
		slog.SetDefault(debuglogger)
	}

	slog.Debug("config", "conf", conf)

	// prepare output dir
	outdir, err := OutDirName(conf)
	if err != nil {
		return Die(err)
	}
	conf.Outdir = outdir

	// used for all HTTP requests
	fetch, err := NewFetcher(conf)
	if err != nil {
		return Die(err)
	}

	switch {
	case len(conf.Wishlinks) >= 1:
		// directly backup wish list[s]
		for _, uri := range conf.Wishlinks {
			err := ScrapeWishlist(fetch, uri)
			if err != nil {
				return Die(err)
			}
		}
	default:
		return Die(errors.New("invalid or no wishlist link specified"))
	}

	if conf.StatsCountWishlist > 0 {
		fmt.Fprintf(output, "Successfully downloaded %d wishlists to %s.\n",
			conf.StatsCountWishlist, conf.Outdir)
	} else {
		fmt.Fprintf(output, "No wish lists found.")
	}

	return 0
}

func Die(err error) int {
	slog.Error("Failure", "error", err.Error())

	return 1
}
