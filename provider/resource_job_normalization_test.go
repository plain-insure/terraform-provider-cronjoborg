// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestNormalizeScheduleSlice(t *testing.T) {
    cases := []struct {
        in   []int
        want []int
    }{
        {[]int{-1}, []int{}},
        {[]int{0, 5}, []int{0, 5}},
        {[]int{}, []int{}},
    }
    for _, c := range cases {
        got := normalizeScheduleSlice(c.in)
        if len(got) != len(c.want) {
            t.Fatalf("length mismatch: in=%v got=%v want=%v", c.in, got, c.want)
        }
        for i := range got {
            if got[i] != c.want[i] {
                t.Fatalf("value mismatch: in=%v got=%v want=%v", c.in, got, c.want)
            }
        }
    }
}
