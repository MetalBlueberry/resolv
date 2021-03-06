package aabb_test

import (
	"math"
	"math/rand"
	"testing"
	"time"

	. "github.com/SolarLune/resolv/resolv/aabb"
	"github.com/stretchr/testify/assert"
)

func TestAABBTree_Insert(t *testing.T) {
	a := &AABBData{-1, -1, 1, 1}
	b := &AABBData{-2, 0, -1, 1}

	t.Run("Insert two nodes", func(t *testing.T) {
		tree := NewTree()
		assert.Nil(t, tree.Root)

		tree.Insert(a)
		assert.Equal(t, tree.Root.Object, a)

		tree.Insert(b)
		assert.NotEqual(t, tree.Root.Object, a)
		assert.NotEqual(t, tree.Root.Object, b)
	})

	t.Run("Insert the same node twice", func(t *testing.T) {
		tree := NewTree()
		assert.Nil(t, tree.Root)

		tree.Insert(a)
		assert.Equal(t, tree.Root.Object, a)

		tree.Insert(b)
		assert.NotEqual(t, tree.Root.Object, a)
		assert.NotEqual(t, tree.Root.Object, b)

		// Same position doesn't mean equal
		tree.Insert(b.Move(0, 0))

		assert.Panics(t, func() {
			tree.Insert(a)
		})
		assert.Panics(t, func() {
			tree.Insert(b)
		})
	})
}

func TestAABBTree_query(t *testing.T) {
	t.Run("Test performance", func(t *testing.T) {
		tree := NewTree()
		base := AABBData{0, 0, 1, 1}
		count := 10000

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
		})

		t.Run("Performance comparisson", func(t *testing.T) {
			start := time.Now()
			queryOverlaps := tree.QueryOverlaps(&base)
			queryTime := time.Since(start)

			assert.NotEmpty(t, queryOverlaps)

			start = time.Now()
			bruteOverlaps := make([]AABB, 0)
			for _, node := range brute {
				if Overlaps(node, &base) {
					bruteOverlaps = append(bruteOverlaps, node)
					assert.Contains(t, queryOverlaps, node)
				}
			}
			bruteTime := time.Since(start)

			t.Logf("Query Time over Brute Force %%%f", 100.0*queryTime.Seconds()/bruteTime.Seconds())

			assert.Equal(t, len(queryOverlaps), len(bruteOverlaps))
			assert.Less(t, queryTime.Seconds(), bruteTime.Seconds())
			for _, node := range bruteOverlaps {
				assert.Contains(t, queryOverlaps, node)
			}
		})
	})

	t.Run("Query empty tree", func(t *testing.T) {
		a := &AABBData{-1, -1, 1, 1}

		tree := NewTree()

		overlaps := tree.QueryOverlaps(a)
		assert.Empty(t, overlaps)
	})
}

func TestAABBTree_Remove(t *testing.T) {
	t.Run("Simple Insert/Removal", func(t *testing.T) {
		tree := NewTree()
		a := &AABBData{-1, -1, 1, 1}
		b := &AABBData{-2, 0, -1, 1}

		assert.Nil(t, tree.Root)

		t.Run("Insert A", func(t *testing.T) {
			tree.Insert(a)
			assert.Equal(t, tree.Root.Object, a)
			assert.Contains(t, tree.NodeIndexMap, a)
			assert.NotContains(t, tree.NodeIndexMap, b)
		})

		t.Run("Insert B", func(t *testing.T) {
			tree.Insert(b)
			assert.NotEqual(t, tree.Root.Object, a)
			assert.NotEqual(t, tree.Root.Object, b)
			assert.Contains(t, tree.NodeIndexMap, a)
			assert.Contains(t, tree.NodeIndexMap, b)
		})

		t.Run("Remove A", func(t *testing.T) {
			tree.Remove(a)
			assert.Equal(t, tree.Root.Object, b)
			assert.NotContains(t, tree.NodeIndexMap, a)
			assert.Contains(t, tree.NodeIndexMap, b)
		})
	})

	t.Run("Insert/Remove from big tree", func(t *testing.T) {
		a := &AABBData{-1, -1, 1, 1}
		b := &AABBData{-2, 0, -1, 1}

		tree := NewTree()
		base := AABBData{0, 0, 1, 1}
		count := 1000

		t.Run("Insert B", func(t *testing.T) {
			tree.Insert(b)
			assert.NotContains(t, tree.NodeIndexMap, a)
			assert.Contains(t, tree.NodeIndexMap, b)
		})

		brute := make([]AABB, 0, count)
		t.Run("Build big tree", func(t *testing.T) {
			for iterations := count; iterations > 0; iterations-- {
				x := rand.Float64()*100 - 50
				y := rand.Float64()*100 - 50

				object := base.Move(float64(x), float64(y))
				tree.Insert(object)
				brute = append(brute, object)
			}
			assert.NotContains(t, tree.NodeIndexMap, a)
			assert.Contains(t, tree.NodeIndexMap, b)
		})

		t.Run("Insert A", func(t *testing.T) {
			tree.Insert(a)
			assert.Contains(t, tree.NodeIndexMap, a)
			assert.Contains(t, tree.NodeIndexMap, b)
		})

		t.Run("Remove A", func(t *testing.T) {
			tree.Remove(a)
			assert.NotContains(t, tree.NodeIndexMap, a)
			assert.Contains(t, tree.NodeIndexMap, b)
		})
		t.Run("Remove B", func(t *testing.T) {
			tree.Remove(b)
			assert.NotContains(t, tree.NodeIndexMap, a)
			assert.NotContains(t, tree.NodeIndexMap, b)
		})

	})

	t.Run("Remove non existing object", func(t *testing.T) {
		a := &AABBData{-1, -1, 1, 1}
		tree := NewTree()
		assert.PanicsWithValue(t, ErrtNotInTree, func() {
			tree.Remove(a)
		})
	})
}
