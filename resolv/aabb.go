package resolv

import "math"

type AABBData struct {
	minX, minY, maxX, maxY float64
}

func (aabb AABBData) AABB() AABBData {
	return aabb
}
func (aabb AABBData) SurfaceArea() float64 {
	a := aabb.maxX - aabb.minX
	b := aabb.maxY - aabb.minY
	return a * b
}

type AABB interface {
	AABB() AABBData
}

type AABBTree struct {
	Root *AABBTreeNode
}

type AABBTreeNode struct {
	Object     AABB
	ObjectAABB AABBData

	Parent, Left, Right *AABBTreeNode
}

func (node *AABBTreeNode) IsLeaf() bool {
	return node.Left == nil
}
func (node *AABBTreeNode) AABB() AABBData {
	return node.ObjectAABB
}

func NewAABBTreeNode(object AABB) *AABBTreeNode {
	return &AABBTreeNode{
		Object:     object,
		ObjectAABB: object.AABB(),
	}
}

func (tree *AABBTree) Add(node *AABBTreeNode) {

	// if the tree is empty then we make the root the leaf
	if tree.Root == nil {
		tree.Root = node
		return
	}

	// search for the best place to put the new leaf in the tree
	// we use surface area and depth as search heuristics
	treeNode := tree.Root
	for !treeNode.IsLeaf() {

		// because of the test in the while loop above we know we are never a leaf inside it
		leftNode := treeNode.Left
		rightNode := treeNode.Right

		combinedAabb := Merge(treeNode, node)

		newParentNodeCost := 2.0 * combinedAabb.SurfaceArea()
		minimumPushDownCost := 2.0 * (combinedAabb.SurfaceArea() - treeNode.AABB().SurfaceArea())

		// use the costs to figure out whether to create a new parent here or descend
		var (
			costLeft  float64
			costRight float64
		)
		if leftNode.IsLeaf() {
			costLeft = Merge(node, leftNode).SurfaceArea() + minimumPushDownCost
		} else {
			newLeftAabb := Merge(node, leftNode)
			costLeft = (newLeftAabb.SurfaceArea() - leftNode.AABB().SurfaceArea()) + minimumPushDownCost
		}
		if rightNode.IsLeaf() {
			costRight = Merge(node, rightNode).SurfaceArea() + minimumPushDownCost
		} else {
			newRightAabb := Merge(node, rightNode)
			costRight = (newRightAabb.SurfaceArea() - rightNode.AABB().SurfaceArea()) + minimumPushDownCost
		}

		// 	// if the cost of creating a new parent node here is less than descending in either direction then
		// 	// we know we need to create a new parent node, errrr, here and attach the leaf to that
		if newParentNodeCost < costLeft && newParentNodeCost < costRight {
			break
		}

		// 	// otherwise descend in the cheapest direction
		if costLeft < costRight {
			treeNode = leftNode
		} else {
			treeNode = rightNode
		}
	}

	// // the leafs sibling is going to be the node we found above and we are going to create a new
	// // parent node and attach the leaf and this item
	sibling := treeNode
	oldParent := sibling.Parent

	newParent := NewAABBTreeNode(Merge(node, sibling))
	newParent.Parent = oldParent
	newParent.Left = sibling
	newParent.Right = node

	node.Parent = newParent
	sibling.Parent = newParent

	if oldParent == nil {
		// the old parent was the root and so this is now the root
		tree.Root = newParent
	} else {
		// the old parent was not the root and so we need to patch the left or right index to
		// point to the new node
		if oldParent.Left == sibling {
			oldParent.Left = newParent
		} else {
			oldParent.Right = newParent
		}
	}

	// // finally we need to walk back up the tree fixing heights and areas
	tree.fixUpwardsTree(node.Parent)
}

func (tree *AABBTree) fixUpwardsTree(node *AABBTreeNode) {
	for node != nil {
		node.ObjectAABB = Merge(node.Left, node.Right)
		node = node.Parent
	}
}

func Merge(a, b AABB) AABBData {
	aData := a.AABB()
	bData := b.AABB()
	return AABBData{
		minX: math.Min(aData.minX, bData.minX),
		minY: math.Min(aData.minY, bData.minY),
		maxX: math.Max(aData.maxX, bData.maxX),
		maxY: math.Max(aData.maxX, bData.maxY),
	}
}
