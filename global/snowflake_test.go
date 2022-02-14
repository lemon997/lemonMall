package global

import (
	"context"
	"testing"
)

func TestGenerate(t *testing.T) {
	node, _ := NewNode(1, context.Background())
	var x, y ID
	for i := 0; i < 1e6; i++ {
		y = node.Generate()
		if x == y {
			t.Errorf("x(%d) & y(%d) are the same", x, y)
		}
		x = y
	}
}
