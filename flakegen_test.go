package flakegen

import (
	"testing"
)

func TestNewNode(t *testing.T) {
	t.Run("测试场景", func(t *testing.T) {
		t.Log(NewNode(1, 1))
	})
}

func TestNode_GetID(t *testing.T) {
	t.Run("测试场景", func(t *testing.T) {
		node, _ := NewNode(1, 1)
		t.Log(node.GetID())
	})
}
