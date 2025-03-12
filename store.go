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
	"bytes"
	"fmt"
	"log/slog"
	"os"
	tpl "text/template"
	"time"
)

type OutdirData struct {
	Year, Day, Month string
}

func OutDirName(conf *Config) (string, error) {
	tmpl, err := tpl.New("outdir").Parse(conf.Outdir)
	if err != nil {
		return "", fmt.Errorf("failed to parse outdir template: %w", err)
	}

	buf := bytes.Buffer{}

	now := time.Now()
	data := OutdirData{
		Year:  now.Format("2006"),
		Month: now.Format("01"),
		Day:   now.Format("02"),
	}

	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute outdir template: %w", err)
	}

	return buf.String(), nil
}

func WishlistName(conf *Config, list *Wishlist) (string, error) {
	tmpl, err := tpl.New("wishlist").Parse(conf.WishlistFiletemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse wish list template: %w", err)
	}

	buf := bytes.Buffer{}

	err = tmpl.Execute(&buf, list)
	if err != nil {
		return "", fmt.Errorf("failed to execute wish list template: %w", err)
	}

	return buf.String(), nil
}

func WriteWishlist(conf *Config, list *Wishlist, listingfile string) error {
	// write wish list file

	listingfd, err := os.Create(listingfile)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", listingfile, err)
	}
	defer listingfd.Close()

	tmpl, err := tpl.New("wishlist").Parse(Assets[conf.Template])
	if err != nil {
		return fmt.Errorf("failed to parse wishlist template: %w", err)
	}

	err = tmpl.Execute(listingfd, list)
	if err != nil {
		return fmt.Errorf("failed to execute wishlist template: %w", err)
	}

	slog.Info("wrote wish list", "listingfile", listingfile)

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)

	if err != nil {
		// return false on any error
		return false
	}

	return !info.IsDir()
}
