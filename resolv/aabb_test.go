package resolv

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
				a: &AABBData{0, 0, 3, 3},
				b: &AABBData{2, 0, 5, 2},
			},
			want: &AABBData{0, 0, 5, 3},
		},
		{
			name: "Merge next to each other",
			args: args{
				a: &AABBData{0, 0, 1, 1},
				b: &AABBData{2, 0, 3, 1},
			},
			want: &AABBData{0, 0, 3, 1},
		},
		{
			name: "Merge inside",
			args: args{
				a: &AABBData{-2, -2, 2, 2},
				b: &AABBData{-1, -1, 1, 1},
			},
			want: &AABBData{-2, -2, 2, 2},
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
	center := NewAABBTreeNode(&AABBData{-1, -1, 1, 1})
	left := NewAABBTreeNode(&AABBData{-2, 0, -1, 1})
	t.Run("Build", func(t *testing.T) {
		assert.Nil(t, tree.Root)

		tree.insertLeaf(center)
		assert.Equal(t, tree.Root, center)

		tree.insertLeaf(left)
		assert.NotEqual(t, tree.Root, center)
		assert.NotEqual(t, tree.Root, left)
	})
}

func TestAABBTree_Add_Random(t *testing.T) {
	move := func(data AABBData, x, y float64) *AABBData {
		data.MinX += x
		data.MaxX += x
		data.MinY += y
		data.MaxY += y
		return &data
	}
	tree := &AABBTree{}
	base := AABBData{0, 0, 1, 1}
	count := 100000
	rand.Seed(time.Now().Unix())

	brute := make([]AABB, 0, count)
	t.Run("Build", func(t *testing.T) {
		start := time.Now()
		assert.Nil(t, tree.Root)
		for count > 0 {
			count--
			x := rand.Float64()*100 - 50
			y := rand.Float64()*100 - 50

			object := move(base, float64(x), float64(y))
			tree.Insert(object)
			brute = append(brute, object)
		}
		t.Log(time.Since(start))
		t.Logf("Depth %d", tree.Depth())
		t.Logf("\nroot %v,\nleft %v,\nright %v", tree.Root.AABB(), tree.Root.Left.AABB(), tree.Root.Right.AABB())
		start = time.Now()
		overlaps := tree.QueryOverlaps(&base)
		t.Log(time.Since(start))
		t.Logf("overlap count %d", len(overlaps))

		start = time.Now()
		bruteOverlaps := make([]AABB, 0)
		for _, node := range brute {
			if Overlaps(node, &base) {
				bruteOverlaps = append(bruteOverlaps, node)
				assert.Contains(t, overlaps, node)
			}
		}
		t.Logf("brute overlap count %d", len(bruteOverlaps))
		t.Log(time.Since(start))

	})
}
