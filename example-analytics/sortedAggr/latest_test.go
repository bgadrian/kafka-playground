package sortedAggr

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewLatest(t *testing.T) {
	lat := NewLatest(10)
	lat.Add("1")
	lat.Add("2")
	lat.Add("3")

	should := []string{"3", "2", "1"}
	got := lat.Get()
	if cmp.Equal(should, got) == false {
		t.Errorf("get is wrong, got %s should %s", got, should)
	}
}
