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
	"log/slog"

	"astuart.co/goq"
)

// scrape an ad. uri is the full uri of the ad, dir is the basedir
func ScrapeWishlist(fetch *Fetcher, uri string) error {
	list := &Wishlist{}

	// get the wish list
	slog.Debug("fetching wish list page", "uri", uri)

	body, err := fetch.Get(uri)
	if err != nil {
		return err
	}
	defer body.Close()

	// extract ad contents with goquery/goq
	err = goq.NewDecoder(body).Decode(&list)
	if err != nil {
		return fmt.Errorf("failed to goquery decode HTML wish list body: %w", err)
	}

	// remove empty items
	list = list.Cleanup()

	// prepare wish list dir name
	listfile, err := WishlistName(fetch.Config, list)
	if err != nil {
		return err
	}

	// write list
	err = WriteWishlist(fetch.Config, list, listfile)
	if err != nil {
		return err
	}

	// tell the user
	slog.Debug("extracted wish list", "list", list)

	// stats
	fetch.Config.IncrWishlist()

	return nil
}
