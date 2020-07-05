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
package namegen

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var (
	// ADJ is a list of adjectives used in the name generator.
	ADJ = []string{"adaptable", "adventurous", "ambitious", "amusing",
		"agreeable", "brave", "bright", "calm", "charming", "considerate",
		"courageous", "creative", "decisive", "diligent", "diplomatic",
		"discreet", "dynamic", "enthusiastic", "exuberant", "faithful",
		"fearless", "friendly", "funny", "generous", "gentle", "gregarious",
		"helpful", "honest", "humorous", "imaginative", "impartial",
		"idependent", "intellectual", "kind", "loving", "loyal", "neat",
		"nice", "passionate", "persistent", "polite", "powerful", "quiet",
		"rational", "reliable", "romantic", "thoughtful", "tidy"}

	// NOUN is a list of nouns used in the name generator.
	NOUN = []string{"aardvark", "albatross", "alligator", "alpaca", "ant",
		"antelope", "badger", "bat", "bear", "bee", "bird", "butterfly",
		"camel", "caibou", "cassowary", "cat", "chicken", "chinchilla",
		"chough", "coati", "cobra", "cod", "crab", "crow", "cuckoo", "deer",
		"dolphin", "dragonfly", "duck", "eagle", "eel", "emu", "falcon",
		"ferret", "finch", "frog", "gecko", "gnu", "bagot", "kiko", "serow",
		"goose", "horse", "hyena", "jellyfish", "kangaroo", "koala", "kudu",
		"lapwing", "lion", "lynx", "mink", "mongoose", "oryx", "otter", "owl",
		"oyster", "unicorn", "pelican", "pony", "turtle", "weasel", "wren",
		"zebra"}
)

var rnd *rand.Rand

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// GetName returns a randomly generated name which is formatted.
func GetName() (name string) {
	name = fmt.Sprintf("%s %s", ADJ[rnd.Intn(len(ADJ))],
		NOUN[rnd.Intn(len(NOUN))])
	name = strings.Title(name)
	return
}
