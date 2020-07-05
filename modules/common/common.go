// Neat Note. A notes sharing platform for university students.
// Copyright (C) 2020 Humaid AlQassimi
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package common

// ContainsInt64 returns whether an array if it contains a specific elements.
func ContainsInt64(a []int64, i int64) bool {
	for _, v := range a {
		if v == i {
			return true
		}
	}
	return false
}

func RemoveInt64(a []int64, v int64) []int64 {
	for i, k := range a {
		if k == v {
			return append(a[:i], a[i+1:]...)
		}
	}
	return a
}
