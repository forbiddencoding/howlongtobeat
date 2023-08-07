/*
 * BSD 3-Clause License
 *
 * Copyright (c) 2023. Edgar Schmidt
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted provided that the
 * following conditions are met:
 *
 * Redistributions of source code must retain the above copyright notice, this list of conditions and the following
 * disclaimer.
 *
 * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following
 * disclaimer in the documentation and/or other materials provided with the distribution.
 *
 * Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products
 * derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
 * INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
 * WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
 * THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package howlongtobeat

import "testing"

func Test_calculateJaccardSimilarity(t *testing.T) {
	tables := []struct {
		gameTitle1 string
		gameTitle2 string
		expected   float64
	}{
		{"Super Mario", "Super Mario", 1},
		{"Super Mario", "Mario Super", 1},
		{"Super Mario Bros", "Super Luigi Bros", 0.5},
		{"Super Mario Bros", "Super Luigi BROS", 0.5}, // case insensitive
		{"Halo", "Call Of Duty", 0},
		{"", "", 0},
		{"Final Fantasy", "Final Fast", 0.33},
		{"Multi word title here", "Multi multi word word title title here here", 1.0}, // repeating words
	}

	for _, table := range tables {
		result := calculateJaccardSimilarity(table.gameTitle1, table.gameTitle2)
		if result != table.expected {
			t.Errorf("calculateJaccardSimilarity(%s, %s) was incorrect, got: %f, want: %f.", table.gameTitle1, table.gameTitle2, result, table.expected)
		}
	}
}
