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

import "time"

type Wishlist struct {
	Name  string     `goquery:"#profile-list-name,text"`
	Items []Wishitem `goquery:"ul#g-items li .a-list-item ,[html]"`
	Date  time.Time
}

type Wishitem struct {
	Product string `goquery:"h2.a-size-base"`
	Price   string `goquery:".a-price .a-offscreen"`
	Link    string `goquery:"h2.a-size-base a,[href]"`
}

func (list *Wishlist) Cleanup() *Wishlist {
	newlist := &Wishlist{Name: list.Name, Date: time.Now()}

	for _, item := range list.Items {
		if item.Product != "" {
			newlist.Items = append(newlist.Items, item)
		}
	}

	return newlist
}
