package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	tests := []struct {
		slice []string
		value string
		want  bool
		desc  string
	}{
		{[]string{"a", "b", "c"}, "a", true, "contains first element"},
		{[]string{"a", "b", "c"}, "b", true, "contains middle element"},
		{[]string{"a", "b", "c"}, "c", true, "contains last element"},
		{[]string{"a", "b", "c"}, "d", false, "does not contain element"},
		{[]string{}, "a", false, "empty slice"},
		{nil, "a", false, "nil slice"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := Contains(tt.slice, tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceInSlice(t *testing.T) {
	tests := []struct {
		big   []string
		small []string
		want  bool
		desc  string
	}{
		{[]string{"a", "b", "c"}, []string{"a"}, true, "contains single matching element"},
		{[]string{"a", "b", "c"}, []string{"a", "d"}, true, "contains one matching element"},
		{[]string{"a", "b", "c"}, []string{"d", "e"}, false, "no matching elements"},
		{[]string{"a", "b", "c"}, []string{}, false, "small slice is empty"},
		{[]string{}, []string{"a"}, false, "big slice is empty"},
		{nil, []string{"a"}, false, "big slice is nil"},
		{[]string{"a", "b", "c"}, nil, false, "small slice is nil"},
		{nil, nil, false, "both slices are nil"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := SliceInSlice(tt.big, tt.small)
			assert.Equal(t, tt.want, got)
		})
	}
}
