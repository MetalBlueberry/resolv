package resolv

import (
	"encoding/json"
	"math"
	"math/rand"
	"os"
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

func TestAABBTree_Insert(t *testing.T) {
	center := &AABBData{-1, -1, 1, 1}
	left := &AABBData{-2, 0, -1, 1}

	t.Run("Insert two nodes", func(t *testing.T) {
		tree := NewAABBTree()
		assert.Nil(t, tree.Root)

		tree.Insert(center)
		assert.Equal(t, tree.Root.Object, center)

		tree.Insert(left)
		assert.NotEqual(t, tree.Root.Object, center)
		assert.NotEqual(t, tree.Root.Object, left)
	})

	t.Run("Insert the same node twice", func(t *testing.T) {
		tree := NewAABBTree()
		assert.Nil(t, tree.Root)

		tree.Insert(center)
		assert.Equal(t, tree.Root.Object, center)

		tree.Insert(left)
		assert.NotEqual(t, tree.Root.Object, center)
		assert.NotEqual(t, tree.Root.Object, left)

		// Same position doesn't mean equal
		tree.Insert(left.Move(0, 0))

		assert.Panics(t, func() {
			tree.Insert(center)
		})
		assert.Panics(t, func() {
			tree.Insert(left)
		})
	})
}

func TestAABBTree_BigTree(t *testing.T) {
	tree := NewAABBTree()
	base := AABBData{0, 0, 1, 1}
	count := 1000
	rand.Seed(time.Now().Unix())

	brute := make([]AABB, 0, count)
	t.Run("Build big tree", func(t *testing.T) {
		start := time.Now()
		assert.Nil(t, tree.Root)
		for iterations := count; iterations > 0; iterations-- {
			x := rand.Float64()*100 - 50
			y := rand.Float64()*100 - 50

			object := base.Move(float64(x), float64(y))
			tree.Insert(object)
			brute = append(brute, object)
		}
		t.Log(time.Since(start))
		t.Logf("Depth %d, min theorical depth %f", tree.Depth(), math.Log2(float64(count)))
		t.Logf("\nroot %v,\nleft %v,\nright %v", tree.Root.AABB(), tree.Root.Left.AABB(), tree.Root.Right.AABB())

	})
	t.Run("Dump json", func(t *testing.T) {
		file, err := os.Create("dump.json")
		assert.Nil(t, err)
		defer file.Close()
		enc := json.NewEncoder(file)
		enc.SetIndent("", " ")
		err = enc.Encode(tree.Root)
		assert.Nil(t, err)

	})
	t.Run("Performance comparisson", func(t *testing.T) {
		start := time.Now()
		queryOverlaps := tree.QueryOverlaps(&base)
		queryTime := time.Since(start)

		start = time.Now()
		bruteOverlaps := make([]AABB, 0)
		for _, node := range brute {
			if Overlaps(node, &base) {
				bruteOverlaps = append(bruteOverlaps, node)
				assert.Contains(t, queryOverlaps, node)
			}
		}
		bruteTime := time.Since(start)

		assert.Equal(t, len(queryOverlaps), len(bruteOverlaps))
		assert.Less(t, queryTime.Seconds(), bruteTime.Seconds())
		for _, node := range bruteOverlaps {
			assert.Contains(t, queryOverlaps, node)
		}

	})
}

func TestAABBTree_BigTreeNonOverlap(t *testing.T) {
	tree := NewAABBTree()
	base := AABBData{0, 0, 1, 1}
	count := 1000
	rand.Seed(time.Now().Unix())

	brute := make([]AABB, 0, count)
	t.Run("Build big tree", func(t *testing.T) {
		start := time.Now()
		assert.Nil(t, tree.Root)
		rows := int(math.Sqrt(float64(count)))
		positions := make([]int, count, count)
		for iterations := 0; iterations < count; iterations++ {
			positions[iterations] = iterations
		}
		rand.Shuffle(len(positions), func(i, j int) { positions[i], positions[j] = positions[j], positions[i] })

		for _, n := range positions {
			x := n % rows
			y := n / rows
			object := base.Move(float64(x), float64(y))
			tree.Insert(object)
			brute = append(brute, object)

		}
		t.Log(time.Since(start))
		t.Logf("Depth %d, min theorical depth %f", tree.Depth(), math.Log2(float64(count)))
		t.Logf("\nroot %v,\nleft %v,\nright %v", tree.Root.AABB(), tree.Root.Left.AABB(), tree.Root.Right.AABB())

	})
	t.Run("Dump json", func(t *testing.T) {
		file, err := os.Create("dump.json")
		assert.Nil(t, err)
		defer file.Close()
		enc := json.NewEncoder(file)
		enc.SetIndent("", " ")
		err = enc.Encode(tree.Root)
		assert.Nil(t, err)

	})
	t.Run("Performance comparisson", func(t *testing.T) {
		start := time.Now()
		queryOverlaps := tree.QueryOverlaps(&base)
		queryTime := time.Since(start)

		start = time.Now()
		bruteOverlaps := make([]AABB, 0)
		for _, node := range brute {
			if Overlaps(node, &base) {
				bruteOverlaps = append(bruteOverlaps, node)
				assert.Contains(t, queryOverlaps, node)
			}
		}
		bruteTime := time.Since(start)

		assert.Equal(t, len(queryOverlaps), len(bruteOverlaps))
		assert.Less(t, queryTime.Seconds(), bruteTime.Seconds())
		for _, node := range bruteOverlaps {
			assert.Contains(t, queryOverlaps, node)
		}

	})
}

func TestAABBTree_Remove(t *testing.T) {
	tree := NewAABBTree()
	center := NewAABBTreeNode(&AABBData{-1, -1, 1, 1})
	left := NewAABBTreeNode(&AABBData{-2, 0, -1, 1})
	t.Run("Build", func(t *testing.T) {
		assert.Nil(t, tree.Root)

		tree.Insert(center)
		assert.Equal(t, tree.Root.Object, center)

		tree.Insert(left)
		assert.NotEqual(t, tree.Root.Object, center)
		assert.NotEqual(t, tree.Root.Object, left)

		tree.Remove(center)
		assert.Equal(t, tree.Root.Object, left)

	})
}
