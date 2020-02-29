package aabb

import (
	"reflect"
)

type Tree struct {
	Root         *treeNode
	NodeIndexMap map[AABB]*treeNode
}

func NewTree() *Tree {
	return &Tree{
		NodeIndexMap: make(map[AABB]*treeNode),
	}
}

type treeNode struct {
	Object     AABB      `json:"-"`
	ObjectAABB *AABBData `json:"aabb"`

	Parent *treeNode `json:"-"`
	Left   *treeNode `json:"left"`
	Right  *treeNode `json:"right"`

	Depth int `json:"depth"`
}

func newTreeNode(object AABB) *treeNode {
	if reflect.ValueOf(object).Kind() != reflect.Ptr {
		panic(ErrNotAReference)
	}
	return &treeNode{
		Object:     object,
		ObjectAABB: object.AABB(),
	}
}

func (node *treeNode) IsLeaf() bool {
	return node.Left == nil
}

func (node *treeNode) AABB() *AABBData {
	return node.ObjectAABB
}

func (node *treeNode) GetSibling() *treeNode {
	parent := node.Parent
	if parent == nil {
		panic("node doesn't contain a parent")
	}
	switch node {
	case parent.Left:
		return parent.Right
	case parent.Right:
		return parent.Left
	default:
		panic("parent doesn't contain children, Tree is corrupted.")
	}
}

func (node *treeNode) replaceWith(other *treeNode) {
	parent := node.Parent
	if parent == nil {
		panic("node doesn't contain a parent")
	}
	other.Parent = parent

	switch node {
	case parent.Left:
		parent.Left = other
	case parent.Right:
		parent.Right = other
	default:
		panic("parent doesn't contain children, Tree is corrupted.")
	}
}

type TreeNodeStack struct {
	data []*treeNode
}

func newTreeNodeStack() *TreeNodeStack {
	return &TreeNodeStack{
		data: make([]*treeNode, 0),
	}
}
func (stack *TreeNodeStack) Push(node *treeNode) {
	stack.data = append(stack.data, node)
}
func (stack *TreeNodeStack) Pop() *treeNode {
	next := stack.data[len(stack.data)-1]
	stack.data = stack.data[:len(stack.data)-1]
	return next
}
func (stack *TreeNodeStack) Empty() bool {
	return len(stack.data) == 0
}

func (tree *Tree) IsEmpty() bool {
	return tree.Root == nil
}

func (tree *Tree) Depth() int {
	stack := newTreeNodeStack()
	stack.Push(tree.Root)
	var maxDepth int
	for !stack.Empty() {
		next := stack.Pop()
		if next.Left != nil {
			next.Left.Depth = next.Depth + 1
			stack.Push(next.Left)
		}
		if next.Right != nil {
			next.Right.Depth = next.Depth + 1
			stack.Push(next.Right)
		}

		if maxDepth < next.Depth+1 {
			maxDepth = next.Depth + 1
		}
	}
	return maxDepth
}

func (tree *Tree) Insert(object AABB) {
	if tree.NodeIndexMap[object] != nil {
		panic(ErrAlreadyInTree)
	}

	node := newTreeNode(object)
	tree.insertLeaf(node)
	tree.NodeIndexMap[object] = node
}

func (tree *Tree) Remove(object AABB) {

	node, ok := tree.NodeIndexMap[object]
	if !ok {
		panic(ErrtNotInTree)
	}
	tree.removeLeaf(node)
	delete(tree.NodeIndexMap, object)
}
func (tree *Tree) removeLeaf(node *treeNode) {
	if node == tree.Root {
		tree.Root = nil
		return
	}

	parent := node.Parent
	grandParent := parent.Parent
	sibling := node.GetSibling()

	node.Parent = nil

	if grandParent == nil {
		tree.Root = sibling
		sibling.Parent = nil
		return
	}

	parent.replaceWith(sibling)
	tree.fixUpwardsTree(grandParent)
}

func (tree *Tree) insertLeaf(node *treeNode) {

	// if the tree is empty then we make the root the leaf
	if tree.Root == nil {
		tree.Root = node
		return
	}

	// search for the best place to put the new leaf in the tree
	// we use surface area and depth as search heuristics
	currentNode := tree.Root
	for !currentNode.IsLeaf() {

		// because of the test in the while loop above we know we are never a leaf inside it
		leftNode := currentNode.Left
		rightNode := currentNode.Right

		combinedAabb := Merge(currentNode, node)

		newParentNodeCost := 2.0 * combinedAabb.SurfaceArea()
		minimumPushDownCost := 2.0 * (combinedAabb.SurfaceArea() - currentNode.AABB().SurfaceArea())

		costFunc := func(side *treeNode) float64 {
			if side.IsLeaf() {
				return Merge(node, side).SurfaceArea() + minimumPushDownCost
			} else {
				newAABB := Merge(node, side)
				return (newAABB.SurfaceArea() - side.AABB().SurfaceArea()) + minimumPushDownCost
			}
		}
		costLeft := costFunc(leftNode)
		costRight := costFunc(rightNode)

		if newParentNodeCost < costLeft && newParentNodeCost < costRight {
			break
		}

		// 	// otherwise descend in the cheapest direction
		if costLeft < costRight {
			currentNode = leftNode
		} else {
			currentNode = rightNode
		}
	}

	// // the leafs sibling is going to be the node we found above and we are going to create a new
	// // parent node and attach the leaf and this item
	sibling := currentNode
	oldParent := sibling.Parent

	newParent := newTreeNode(Merge(node, sibling))
	newParent.Parent = oldParent
	newParent.Left = sibling
	newParent.Right = node

	node.Parent = newParent
	sibling.Parent = newParent

	switch {
	case oldParent == nil:
		tree.Root = newParent
	case oldParent.Left == sibling:
		oldParent.Left = newParent
	case oldParent.Right == sibling:
		oldParent.Right = newParent
	}

	tree.fixUpwardsTree(node.Parent)
}

func (tree *Tree) fixUpwardsTree(node *treeNode) {
	for node != nil {
		node.ObjectAABB = Merge(node.Left, node.Right)
		node = node.Parent
	}
}

func (tree *Tree) QueryOverlaps(object AABB) []AABB {
	if reflect.ValueOf(object).Kind() != reflect.Ptr {
		panic(ErrNotAReference)
	}
	overlaps := make([]AABB, 0)

	if tree.IsEmpty() {
		return overlaps
	}

	stack := newTreeNodeStack()
	testAABB := object.AABB()
	stack.Push(tree.Root)
	for !stack.Empty() {
		node := stack.Pop()

		if Overlaps(node, testAABB) {
			if node.IsLeaf() && node.Object != object {
				overlaps = append(overlaps, node.Object)
			} else {
				if node.Left != nil {
					stack.Push(node.Left)
				}
				if node.Right != nil {
					stack.Push(node.Right)
				}
			}
		}
	}
	return overlaps
}
