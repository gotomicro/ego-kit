// Copyright 2021 gotomicro
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tree

import (
	"errors"

	"github.com/gotomicro/ekit"
)

const (
	Red   = false
	Black = true
)

var (
	ErrRBTreeSameRBNode = errors.New("ekit: RBTree不能添加重复节点Key")
	// errRBTreeNotRBNode  = errors.New("ekit: RBTree不存在节点Key")
	// errRBTreeCantRepaceNil = errors.New("ekit: RBTree不能将节点替换为nil")
)

type RBTree[T any, V any] struct {
	root    *RBNode[T, V]
	compare ekit.Comparator[T]
	size    int
}

type RBNode[T any, V any] struct {
	color               bool
	Key                 T
	Value               V
	left, right, parent *RBNode[T, V]
}

// NewRBTree 构建红黑树
func NewRBTree[T any, V any](compare ekit.Comparator[T]) *RBTree[T, V] {
	return &RBTree[T, V]{
		compare: compare,
		root:    nil,
	}
}

func NewRBNode[T any, V any](key T, value V) *RBNode[T, V] {
	return &RBNode[T, V]{
		Key:    key,
		Value:  value,
		color:  Red,
		left:   nil,
		right:  nil,
		parent: nil,
	}
}

func (redBlackTree *RBTree[T, V]) Root() *RBNode[T, V] {
	if redBlackTree == nil {
		return nil
	}
	return redBlackTree.root
}

func (redBlackTree *RBTree[T, V]) Size() int {
	if redBlackTree == nil {
		return 0
	}
	return redBlackTree.size
}

// Add 增加节点
func (redBlackTree *RBTree[T, V]) Add(node *RBNode[T, V]) error {
	return redBlackTree.addNode(node)
}

// Delete 删除节点
func (redBlackTree *RBTree[T, V]) Delete(key T) {
	node := redBlackTree.getRBNode(key)
	if node == nil {
		return
	}
	redBlackTree.deleteNode(node)
}

// Find 查找节点
func (redBlackTree *RBTree[T, V]) Find(key T) *RBNode[T, V] {
	return redBlackTree.getRBNode(key)
}

// addNode 插入新节点
func (redBlackTree *RBTree[T, V]) addNode(node *RBNode[T, V]) error {
	t := redBlackTree.root
	if redBlackTree.root == nil {
		redBlackTree.root = NewRBNode[T, V](node.Key, node.Value)
	}
	cmp := 0
	parent := &RBNode[T, V]{}
	for t != nil {
		parent = t
		cmp = redBlackTree.compare(node.Key, t.Key)
		if cmp < 0 {
			t = t.left
		} else if cmp > 0 {
			t = t.right
		} else if cmp == 0 {
			return ErrRBTreeSameRBNode
		}
	}
	tempNode := &RBNode[T, V]{
		Key:    node.Key,
		parent: parent,
		Value:  node.Value,
		color:  Red,
	}
	if cmp < 0 {
		parent.left = tempNode
	} else {
		parent.right = tempNode
	}
	redBlackTree.size++
	redBlackTree.fixAfterAdd(tempNode)
	return nil
}

// deleteNode 红黑树删除方法
// 删除分两步,第一步取出后继节点,第二部着色旋转
// 取后继节点
// case1:node左右非空子节点,通过getSuccessor获取后继节点
// case2:node左右只有一个非空子节点
// case3:node左右均为空节点
// 着色旋转
// case1:当删除节点非空且为黑色时,会违反红黑树任何路径黑节点个数相同的约束,所以需要重新平衡
// case2:当删除红色节点时,不会破坏任何约束,所以不需要平衡
func (redBlackTree *RBTree[T, V]) deleteNode(node *RBNode[T, V]) {
	// node左右非空,取后继节点
	if node.left != nil && node.right != nil {
		s := redBlackTree.getSuccessor(node)
		node.Key = s.Key
		node.Value = s.Value
		node = s
	}
	var replacement *RBNode[T, V]
	// node节点只有一个非空子节点
	if node.left != nil {
		replacement = node.left
	} else {
		replacement = node.right
	}
	if replacement != nil {
		replacement.parent = node.parent
		if node.parent == nil {
			redBlackTree.root = replacement
		} else if node == node.parent.left {
			node.parent.left = replacement
		} else {
			node.parent.right = replacement
		}
		node.left = nil
		node.right = nil
		node.parent = nil
		if node.color {
			redBlackTree.fixAfterDelete(replacement)
		}
	} else if node.parent == nil {
		// 如果node节点无父节点,说明node为root节点
		redBlackTree.root = nil
	} else {
		// node子节点均为空
		if node.color {
			redBlackTree.fixAfterDelete(node)
		}
		if node.parent != nil {
			if node == node.parent.left {
				node.parent.left = nil
			} else if node == node.parent.right {
				node.parent.right = nil
			}
			node.parent = nil
		}
	}
	redBlackTree.size--
}

// getSuccessor 寻找后继节点
// case1: node节点存在右子节点,则右子树的最小节点是node的后继节点
// case2: node节点不存在右子节点,则其第一个为左节点的祖先的父节点为node的后继节点
func (redBlackTree *RBTree[T, V]) getSuccessor(node *RBNode[T, V]) *RBNode[T, V] {
	if node == nil {
		return nil
	} else if node.right != nil {
		p := node.right
		for p.left != nil {
			p = p.left
		}
		return p
	} else {
		p := node.parent
		ch := node
		for p != nil && ch == p.right {
			ch = p
			p = p.parent
		}
		return p
	}

}

func (redBlackTree *RBTree[T, V]) getRBNode(key T) *RBNode[T, V] {
	node := redBlackTree.root
	for node != nil {
		cmp := redBlackTree.compare(key, node.Key)
		if cmp < 0 {
			node = node.left
		} else if cmp > 0 {
			node = node.right
		} else {
			return node
		}
	}
	return nil
}

// fixAfterAdd 插入时着色旋转
// 如果是空节点、root节点、父节点是黑无需构建
// 可分为3种情况
// fixUncleRed 叔叔节点是红色右节点
// fixAddLeftBlack 叔叔节点是黑色右节点
// fixAddRightBlack 叔叔节点是黑色左节点
func (redBlackTree *RBTree[T, V]) fixAfterAdd(x *RBNode[T, V]) {
	x.color = Red
	for x != nil && x != redBlackTree.root && !x.getParent().getColor() {
		y := x.getUncle()
		if !y.getColor() {
			x = redBlackTree.fixUncleRed(x, y)
			continue
		}
		if x.getParent() == x.getGrandParent().getLeft() {
			x = redBlackTree.fixAddLeftBlack(x)
			continue
		}
		x = redBlackTree.fixAddRightBlack(x)
	}
	redBlackTree.root.setColor(Black)
}

// fixAddLeftRed 叔叔节点是红色右节点，由于不能存在连续红色节点,此时祖父节点x.getParent().getParent()必为黑。另x为红所以叔父节点需要变黑，祖父变红，此时红黑树完成
//
//							  b(b)                    b(r)
//							/		\				/		\
//						  a(r)        y(r)  ->   a(b)        y(b)
//						/   \       /  \         /   \       /  \
//		            x(r)    nil   nil  nil    x (r) nil   nil  nil
//	             	/  \                      /  \
//	            	nil nil                   nil nil
func (redBlackTree *RBTree[T, V]) fixUncleRed(x *RBNode[T, V], y *RBNode[T, V]) *RBNode[T, V] {
	x.getParent().setColor(Black)
	y.setColor(Black)
	x.getGrandParent().setColor(Red)
	x = x.getGrandParent()
	return x
}

// fixAddLeftBlack 叔叔节点是黑色右节点.x节点是父节点左节点,执行左旋，此时x节点变为原x节点的父节点a,也就是左子节点。的接着将x的父节点和爷爷节点的颜色对换。然后对爷爷节点进行右旋转,此时红黑树完成
// 如果x为左节点则跳过左旋操作
//
//							  b(b)                    b(b)                b(r)
//							/		\				/		\            /   \
//						  a(r)        y(b)  ->   a(r)        y(b)  ->  a(b)   y(b)
//						/   \       /  \         /   \       /  \      /  \    /  \
//		               nil   x (r) nil  nil      x(r) nil  nil  nil   x(r) nil nil nil
//	           		 		 /  \               /  \                  / \
//	           		 		nil nil             nil nil              nil nil
func (redBlackTree *RBTree[T, V]) fixAddLeftBlack(x *RBNode[T, V]) *RBNode[T, V] {
	if x == x.getParent().getRight() {
		x = x.getParent()
		redBlackTree.rotateLeft(x)
	}
	x.getParent().setColor(Black)
	x.getGrandParent().setColor(Red)
	redBlackTree.rotateRight(x.getGrandParent())
	return x
}

// fixAddRightBlack 叔叔节点是黑色左节点.x节点是父节点右节点,执行右旋，此时x节点变为原x节点的父节点a,也就是右子节点。接着将x的父节点和爷爷节点的颜色对换。然后对爷爷节点进行右旋转,此时红黑树完成
// 如果x为右节点则跳过右旋操作
//
//							  b(b)                    b(b)                b(r)
//							/		\				/		\            /   \
//						  y(b)       a(r)  ->   y(b)        a(r)  ->  y(b)     a(b)
//						/   \       /  \         /   \       /  \      /  \    /  \
//		               nil   nil x(r)  nil      nil nil  nil  x(r)   nil nil  nil  x(r)
//	           		 		      /  \                         /  \               /  \
//	           		 		      nil nil                    nil nil              nil nil
func (redBlackTree *RBTree[T, V]) fixAddRightBlack(x *RBNode[T, V]) *RBNode[T, V] {
	if x == x.getParent().getLeft() {
		x = x.getParent()
		redBlackTree.rotateRight(x)
	}
	x.getParent().setColor(Black)
	x.getGrandParent().setColor(Red)
	redBlackTree.rotateLeft(x.getGrandParent())
	return x
}

// fixAfterDelete 删除时着色旋转
// 根据x是节点位置分为fixAfterDeleteLeft,fixAfterDeleteRight两种情况
func (redBlackTree *RBTree[T, V]) fixAfterDelete(x *RBNode[T, V]) {
	for x != redBlackTree.root && x.getColor() {
		if x == x.parent.getLeft() {
			x = redBlackTree.fixAfterDeleteLeft(x)
		} else {
			x = redBlackTree.fixAfterDeleteRight(x)
		}
	}
	x.setColor(Black)
}

// fixAfterDeleteLeft 处理x为左子节点时的平衡处理
func (redBlackTree *RBTree[T, V]) fixAfterDeleteLeft(x *RBNode[T, V]) *RBNode[T, V] {
	sib := x.getParent().Right()
	if !sib.getColor() {
		sib.setColor(Black)
		sib.getParent().setColor(Red)
		redBlackTree.rotateLeft(x.getParent())
		sib = x.getParent().getRight()
	}
	if sib.getLeft().getColor() && sib.getRight().getColor() {
		sib.setColor(Red)
		x = x.getParent()
	} else {
		if sib.getRight().getColor() {
			sib.getLeft().setColor(Black)
			sib.setColor(Red)
			redBlackTree.rotateRight(sib)
			sib = x.getParent().getRight()
		}
		sib.setColor(x.getParent().getColor())
		x.getParent().setColor(Black)
		sib.getRight().setColor(Black)
		redBlackTree.rotateLeft(x.getParent())
		x = redBlackTree.root
	}
	return x
}

// fixAfterDeleteRight 处理x为右子节点时的平衡处理
func (redBlackTree *RBTree[T, V]) fixAfterDeleteRight(x *RBNode[T, V]) *RBNode[T, V] {
	sib := x.getParent().Left()
	if !sib.getColor() {
		sib.setColor(Black)
		x.getParent().setColor(Red)
		redBlackTree.rotateRight(x.getParent())
		sib = x.getBrother()
	}
	if sib.getRight().getColor() && sib.getLeft().getColor() {
		sib.setColor(Red)
		x = x.getParent()
	} else {
		if sib.getLeft().getColor() {
			sib.getRight().setColor(Black)
			sib.setColor(Red)
			redBlackTree.rotateLeft(sib)
			sib = x.getParent().getLeft()
		}
		sib.setColor(x.getParent().getColor())
		x.getParent().setColor(Black)
		sib.getLeft().setColor(Black)
		redBlackTree.rotateRight(x.getParent())
		x = redBlackTree.root
	}
	return x
}

// rotateLeft 左旋转
//
//							  b                    a
//							/	\				  /	  \
//						  c       a  ->    		 b     y
//								 / \            /  \
//		                     	x    y     		c	x

func (redBlackTree *RBTree[T, V]) rotateLeft(node *RBNode[T, V]) {
	if node != nil {
		r := node.right
		if r == nil {
			return
		}
		node.right = r.left
		if r.left != nil {
			r.left.parent = node
		}
		r.parent = node.parent
		if node.parent == nil {
			redBlackTree.root = r
		} else if node.parent.left == node {
			node.parent.left = r
		} else {
			node.parent.right = r
		}
		r.left = node
		node.parent = r
	}
}

// rotateRight 右旋转
//
//						  b                    c
//						/	\				  /	  \
//					  c       a  ->    		 x     b
//					 /	\	                       / \
//	                 x  y  	     	 	  	       y  a
func (redBlackTree *RBTree[T, V]) rotateRight(node *RBNode[T, V]) {
	if node != nil {
		l := node.left
		if l == nil {
			return
		}
		node.left = l.right
		if l.right != nil {
			l.right.parent = node
		}
		l.parent = node.parent
		if node.parent == nil {
			redBlackTree.root = l
		} else if node.parent.right == node {
			node.parent.right = l
		} else {
			node.parent.left = l
		}
		l.right = node
		node.parent = l
	}
}

func (node *RBNode[T, V]) getColor() bool {
	if node == nil {
		return Black
	}
	return node.color
}

func (node *RBNode[T, V]) setColor(color bool) {
	if node == nil {
		return
	}
	node.color = color
}

func (node *RBNode[T, V]) getParent() *RBNode[T, V] {
	if node == nil {
		return nil
	}
	return node.parent
}

func (node *RBNode[T, V]) getLeft() *RBNode[T, V] {
	if node == nil {
		return nil
	}
	return node.left
}

func (node *RBNode[T, V]) getRight() *RBNode[T, V] {
	if node == nil {
		return nil
	}
	return node.right
}

func (node *RBNode[T, V]) getUncle() *RBNode[T, V] {
	if node == nil {
		return nil
	}
	return node.getParent().getBrother()
}
func (node *RBNode[T, V]) getGrandParent() *RBNode[T, V] {
	if node == nil {
		return nil
	}
	return node.getParent().getParent()
}
func (node *RBNode[T, V]) getBrother() *RBNode[T, V] {
	if node == nil {
		return nil
	}
	if node == node.getParent().getLeft() {
		return node.getParent().getRight()
	}
	return node.getParent().getLeft()
}
func (node *RBNode[T, V]) SetValue(v V) {
	if node == nil {
		return
	}
	node.Value = v
}
func (node *RBNode[T, V]) Left() *RBNode[T, V] {
	return node.getLeft()
}

func (node *RBNode[T, V]) Right() *RBNode[T, V] {
	return node.getRight()
}

func (node *RBNode[T, V]) Parent() *RBNode[T, V] {
	return node.getParent()
}

// IsRedBlackTree 检测是否满足红黑树
func IsRedBlackTree[T any, V any](root *RBNode[T, V]) bool {
	// 检测节点是否黑色
	if !root.getColor() {
		return false
	}
	// count 取最左树的黑色节点作为对照
	count := 0
	num := 0
	node := root
	for node != nil {
		if node.color {
			count++
		}
		node = node.getLeft()
	}
	return nodeCheck[T](root, count, num)
}

// nodeCheck 节点检测
// 1、是否有连续的红色节点
// 2、每条路径的黑色节点是否一致
func nodeCheck[T any, V any](node *RBNode[T, V], count int, num int) bool {
	if node == nil {
		return true
	}
	if !node.getColor() && !node.parent.getColor() {
		return false
	}
	if node.getColor() {
		num++
	}
	if node.getLeft() == nil && node.getRight() == nil {
		if num != count {
			return false
		}
	}
	return nodeCheck(node.left, count, num) && nodeCheck(node.right, count, num)
}
