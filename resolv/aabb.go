package resolv

// Inspired by https://www.azurefromthetrenches.com/introductory-guide-to-aabb-tree-collision-detection/

import (
	"math"
	"reflect"
)

type AABBData struct {
	MinX, MinY, MaxX, MaxY float64
}

func (aabb *AABBData) AABB() *AABBData {
	return aabb
}
func (aabb *AABBData) SurfaceArea() float64 {
	a := aabb.MaxX - aabb.MinX
	b := aabb.MaxY - aabb.MinY
	return a * b
}

func (aabb *AABBData) Merge(other *AABBData) *AABBData {
	return Merge(aabb, other)
}

func Merge(a, b AABB) *AABBData {
	aData := a.AABB()
	bData := b.AABB()
	return &AABBData{
		MinX: math.Min(aData.MinX, bData.MinX),
		MinY: math.Min(aData.MinY, bData.MinY),
		MaxX: math.Max(aData.MaxX, bData.MaxX),
		MaxY: math.Max(aData.MaxX, bData.MaxY),
	}
}

func (aabb *AABBData) Overlaps(other *AABBData) bool {
	return Overlaps(aabb, other)
}

func Overlaps(a, b AABB) bool {
	aData := a.AABB()
	bData := b.AABB()
	return aData.MaxX > bData.MinX &&
		aData.MinX < bData.MaxX &&
		aData.MaxY > bData.MinY &&
		aData.MinY < bData.MaxY
}

type AABB interface {
	AABB() *AABBData
}

type AABBTree struct {
	Root *AABBTreeNode
}

type AABBTreeNode struct {
	Object AABB

	Parent, Left, Right *AABBTreeNode

	objectAABB *AABBData
	depth      int
}

func NewAABBTreeNode(object AABB) *AABBTreeNode {
	if reflect.ValueOf(object).Kind() != reflect.Ptr {
		panic("provided object must be a pointer")
	}
	return &AABBTreeNode{
		Object:     object,
		objectAABB: object.AABB(),
	}
}

func (node *AABBTreeNode) IsLeaf() bool {
	return node.Left == nil
}
func (node *AABBTreeNode) AABB() *AABBData {
	return node.objectAABB
}

type AABBTreeNodeStack struct {
	data []*AABBTreeNode
}

func NewAABBTreeNodeStack() *AABBTreeNodeStack {
	return &AABBTreeNodeStack{
		data: make([]*AABBTreeNode, 0),
	}
}
func (stack *AABBTreeNodeStack) Push(node *AABBTreeNode) {
	stack.data = append(stack.data, node)
}
func (stack *AABBTreeNodeStack) Pop() *AABBTreeNode {
	next := stack.data[len(stack.data)-1]
	stack.data = stack.data[:len(stack.data)-1]
	return next
}
func (stack *AABBTreeNodeStack) Empty() bool {
	return len(stack.data) == 0
}

func (tree *AABBTree) Depth() int {
	stack := NewAABBTreeNodeStack()
	stack.Push(tree.Root)
	var maxDepth int
	for !stack.Empty() {
		next := stack.Pop()
		if next.Left != nil {
			next.Left.depth = next.depth + 1
			stack.Push(next.Left)
		}
		if next.Right != nil {
			next.Right.depth = next.depth + 1
			stack.Push(next.Right)
		}

		if maxDepth < next.depth+1 {
			maxDepth = next.depth + 1
		}
	}
	return maxDepth
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
		node.objectAABB = Merge(node.Left, node.Right)
		node = node.Parent
	}
}

// QueryOverlaps checks for overlapping AABB objects. Be aware that if two
func (tree *AABBTree) QueryOverlaps(object AABB) []AABB {
	if reflect.ValueOf(object).Kind() != reflect.Ptr {
		panic("provided object must be a pointer")
	}
	// std::forward_list<std::shared_ptr<IAABB>> AABBTree::queryOverlaps(const std::shared_ptr<IAABB>& object) const
	// {
	// 	std::forward_list<std::shared_ptr<IAABB>> overlaps;
	// 	std::stack<unsigned> stack;
	// 	AABB testAabb = object->getAABB();
	overlaps := make([]AABB, 0)
	stack := NewAABBTreeNodeStack()
	testAABB := object.AABB()
	// 	stack.push(_rootNodeIndex);
	// 	while(!stack.empty())
	// 	{
	stack.Push(tree.Root)
	for !stack.Empty() {
		// 		unsigned nodeIndex = stack.top();
		// 		stack.pop();
		node := stack.Pop()

		// 		if (nodeIndex == AABB_NULL_NODE) continue;

		// 		const AABBNode& node = _nodes[nodeIndex];
		if Overlaps(node, testAABB) {
			// 		if (node.aabb.overlaps(testAabb))
			// 		{
			if node.IsLeaf() && node.Object != object {
				// 			if (node.isLeaf() && node.object != object)
				// 			{
				overlaps = append(overlaps, node)
				// 				overlaps.push_front(node.object);
				// 			}
				// 			else
			} else {
				if node.Left != nil {
					stack.Push(node.Left)
				}
				if node.Right != nil {
					stack.Push(node.Right)
				}
				// 			{
				// 				stack.push(node.leftNodeIndex);
				// 				stack.push(node.rightNodeIndex);
				// 			}
				// 		}
				// 	}
			}
		}
	}
	return overlaps
	// 	return overlaps;
	// }
}
