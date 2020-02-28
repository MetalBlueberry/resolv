package tree

import (
	"reflect"
)

type AABBTree struct {
	Root         *aabbTreeNode
	NodeIndexMap map[AABB]*aabbTreeNode
}

func NewAABBTree() *AABBTree {
	return &AABBTree{
		NodeIndexMap: make(map[AABB]*aabbTreeNode),
	}
}

type aabbTreeNodeStack struct {
	data []*aabbTreeNode
}

func newAABBTreeNodeStack() *aabbTreeNodeStack {
	return &aabbTreeNodeStack{
		data: make([]*aabbTreeNode, 0),
	}
}
func (stack *aabbTreeNodeStack) Push(node *aabbTreeNode) {
	stack.data = append(stack.data, node)
}
func (stack *aabbTreeNodeStack) Pop() *aabbTreeNode {
	next := stack.data[len(stack.data)-1]
	stack.data = stack.data[:len(stack.data)-1]
	return next
}
func (stack *aabbTreeNodeStack) Empty() bool {
	return len(stack.data) == 0
}

func (tree *AABBTree) IsEmpty() bool {
	return tree.Root == nil
}

func (tree *AABBTree) Depth() int {
	stack := newAABBTreeNodeStack()
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

func (tree *AABBTree) Insert(object AABB) {
	if tree.NodeIndexMap[object] != nil {
		panic(ErrAlreadyInTree)
	}

	node := newAABBTreeNode(object)
	tree.insertLeaf(node)
	tree.NodeIndexMap[object] = node
}

func (tree *AABBTree) Remove(object AABB) {

	node, ok := tree.NodeIndexMap[object]
	if !ok {
		panic(ErrtNotInTree)
	}
	tree.removeLeaf(node)
	delete(tree.NodeIndexMap, object)
}
func (tree *AABBTree) removeLeaf(node *aabbTreeNode) {
	// // if the leaf is the root then we can just clear the root pointer and return
	// if (leafNodeIndex == _rootNodeIndex)
	if node == tree.Root {
		// {
		// 	_rootNodeIndex = AABB_NULL_NODE;
		tree.Root = nil
		// 	return;
		// }
	}

	// AABBNode& leafNode = _nodes[leafNodeIndex];
	// unsigned parentNodeIndex = leafNode.parentNodeIndex;
	// const AABBNode& parentNode = _nodes[parentNodeIndex];
	parent := node.Parent
	// unsigned grandParentNodeIndex = parentNode.parentNodeIndex;
	grandParent := parent.Parent
	// unsigned siblingNodeIndex = parentNode.leftNodeIndex == leafNodeIndex ? parentNode.rightNodeIndex : parentNode.leftNodeIndex;

	var sibling *aabbTreeNode
	switch node {
	case parent.Left:
		sibling = parent.Right
	case parent.Right:
		sibling = parent.Left
	default:
		panic("parent doesn't contain children, Tree is corrupted.")
	}
	// assert(siblingNodeIndex != AABB_NULL_NODE); // we must have a sibling
	// AABBNode& siblingNode = _nodes[siblingNodeIndex];

	// if (grandParentNodeIndex != AABB_NULL_NODE)
	if grandParent != nil {
		// {
		// 	// if we have a grand parent (i.e. the parent is not the root) then destroy the parent and connect the sibling to the grandparent in its
		// 	// place
		// 	AABBNode& grandParentNode = _nodes[grandParentNodeIndex];
		switch parent {
		case grandParent.Left:
			grandParent.Left = sibling
		case grandParent.Right:
			grandParent.Right = sibling
		}
		// 	if (grandParentNode.leftNodeIndex == parentNodeIndex)
		// 	{
		// 		grandParentNode.leftNodeIndex = siblingNodeIndex;
		// 	}
		// 	else
		// 	{
		// 		grandParentNode.rightNodeIndex = siblingNodeIndex;
		// 	}
		// 	siblingNode.parentNodeIndex = grandParentNodeIndex;
		parent = grandParent
		// 	deallocateNode(parentNodeIndex);
		// 	fixUpwardsTree(grandParentNodeIndex);
		tree.fixUpwardsTree(grandParent)
		// }
	} else {
		// else
		// {
		// 	// if we have no grandparent then the parent is the root and so our sibling becomes the root and has it's parent removed
		tree.Root = sibling
		sibling.Parent = nil
		// 	_rootNodeIndex = siblingNodeIndex;
		// 	siblingNode.parentNodeIndex = AABB_NULL_NODE;
		// 	deallocateNode(parentNodeIndex);
		// }
	}
	node.Parent = nil

	// leafNode.parentNodeIndex = AABB_NULL_NODE;
}

func (tree *AABBTree) insertLeaf(node *aabbTreeNode) {

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

	newParent := newAABBTreeNode(Merge(node, sibling))
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

func (tree *AABBTree) fixUpwardsTree(node *aabbTreeNode) {
	for node != nil {
		node.ObjectAABB = Merge(node.Left, node.Right)
		node = node.Parent
	}
}

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

	if tree.IsEmpty() {
		return overlaps
	}

	stack := newAABBTreeNodeStack()
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
				overlaps = append(overlaps, node.Object)
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
