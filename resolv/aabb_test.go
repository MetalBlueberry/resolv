package resolv

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	type args struct {
		a AABB
		b AABB
	}
	tests := []struct {
		name string
		args args
		want AABB
	}{
		{
			name: "Merge overlapping",
			args: args{
				a: AABBData{0, 0, 3, 3},
				b: AABBData{2, 0, 5, 2},
			},
			want: AABBData{0, 0, 5, 3},
		},
		{
			name: "Merge next to each other",
			args: args{
				a: AABBData{0, 0, 1, 1},
				b: AABBData{2, 0, 3, 1},
			},
			want: AABBData{0, 0, 3, 1},
		},
		{
			name: "Merge inside",
			args: args{
				a: AABBData{-2, -2, 2, 2},
				b: AABBData{-1, -1, 1, 1},
			},
			want: AABBData{-2, -2, 2, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Merge(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAABBTree_Add(t *testing.T) {
	tree := &AABBTree{}
	center := NewAABBTreeNode(AABBData{-1, -1, 1, 1})
	left := NewAABBTreeNode(AABBData{-2, 0, -1, 1})
	right := NewAABBTreeNode(AABBData{1, 2, -1, 1})
	t.Run("Build", func(t *testing.T) {
		assert.Nil(t, tree.Root)

		tree.Add(center)
		assert.Equal(t, tree.Root, center)

		tree.Add(left)
		assert.NotEqual(t, tree.Root, center)
		assert.NotEqual(t, tree.Root, left)

		tree.Add(right)
		tree.Add(left)
		tree.Add(center)
		assert.Equal(t, tree.Root, left)
	})
}
