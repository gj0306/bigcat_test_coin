package tools
import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
)

// 数据实例的接口。模块使用者继承这个接口，外加自己的数据
type Content interface {
	CalculateHash() ([]byte, error)     // 这个是计算数据的hash，由使用者实现
	Equals(other Content) (bool, error) // 这个是判断数据相等的依据，由使用者实现
}

// Merkle Tree 的数据结构
type MerkleTree struct {
	Root         *Node
	merkleRoot   []byte
	Leafs        []*Node
	hashStrategy func() hash.Hash // hash函数。默认sha256
}

// Node 的数据结构
type Node struct {
	Tree   *MerkleTree
	Parent *Node
	Left   *Node
	Right  *Node
	leaf   bool
	dup    bool
	Hash   []byte
	C      Content
}

// 计算节点的hash（一直递归到叶子节点）。使用的前序遍历（根结点 ---> 左子树 ---> 右子树）
func (n *Node) calculateNodeHashFromLeaf() ([]byte, error) {
	if n.leaf {
		return n.C.CalculateHash()
	}

	leftBytes, err := n.Left.calculateNodeHashFromLeaf()
	if err != nil {
		return nil, err
	}

	rightBytes, err := n.Right.calculateNodeHashFromLeaf()
	if err != nil {
		return nil, err
	}

	h := n.Tree.hashStrategy()
	if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// 计算节点的hash（仅仅根据左右子节点的hash计算）
func (n *Node) calculateNodeHash() ([]byte, error) {
	if n.leaf {
		return n.C.CalculateHash()
	}

	h := n.Tree.hashStrategy()
	if _, err := h.Write(append(n.Left.Hash, n.Right.Hash...)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// 创建一棵树
func NewMerkleTree(cs []Content) (*MerkleTree, error) {
	var defaultHashStrategy = sha256.New
	t := &MerkleTree{
		hashStrategy: defaultHashStrategy,
	}
	root, leafs, err := buildWithContent(cs, t)
	if err != nil {
		return nil, err
	}
	t.Root = root
	t.Leafs = leafs
	t.merkleRoot = root.Hash
	return t, nil
}

// 自定义hash函数，创建一棵树
func NewMerkleTreeWithHashStrategy(cs []Content, hashStrategy func() hash.Hash) (*MerkleTree, error) {
	t := &MerkleTree{
		hashStrategy: hashStrategy,
	}
	root, leafs, err := buildWithContent(cs, t)
	if err != nil {
		return nil, err
	}
	t.Root = root
	t.Leafs = leafs
	t.merkleRoot = root.Hash
	return t, nil
}

// 获取某数据在默克尔树中的路径(从查询节点到跟节点的路径的相对路径)。
// 第一个返回结果表示从下往上依次的节点hash值，第二个返回结果表示是左节点还是右节点，0表示左，1表示右
func (m *MerkleTree) GetMerklePath(content Content) ([][]byte, []int64, error) {
	for _, current := range m.Leafs {
		ok, err := current.C.Equals(content)
		if err != nil {
			return nil, nil, err
		}

		if ok {
			currentParent := current.Parent
			var merklePath [][]byte
			var index []int64
			for currentParent != nil {
				if bytes.Equal(currentParent.Left.Hash, current.Hash) {
					merklePath = append(merklePath, currentParent.Right.Hash)
					index = append(index, 1) // right leaf
				} else {
					merklePath = append(merklePath, currentParent.Left.Hash)
					index = append(index, 0) // left leaf
				}
				current = currentParent
				currentParent = currentParent.Parent
			}
			return merklePath, index, nil
		}
	}
	return nil, nil, nil
}

// 构建一棵树。返回 根节点、所有叶子节点
func buildWithContent(cs []Content, t *MerkleTree) (*Node, []*Node, error) {
	if len(cs) == 0 {
		return nil, nil, errors.New("error: cannot construct tree with no content")
	}
	var leafs []*Node
	for _, c := range cs {
		hash, err := c.CalculateHash()
		if err != nil {
			return nil, nil, err
		}

		leafs = append(leafs, &Node{
			Hash: hash,
			C:    c,
			leaf: true,
			Tree: t,
		})
	}
	if len(leafs)%2 == 1 { // 如果叶子节点个数为基数，则最后那个孤立节点复制自己一份
		duplicate := &Node{
			Hash: leafs[len(leafs)-1].Hash,
			C:    leafs[len(leafs)-1].C,
			leaf: true,
			dup:  true,
			Tree: t,
		}
		leafs = append(leafs, duplicate)
	}
	root, err := buildIntermediate(leafs, t)
	if err != nil {
		return nil, nil, err
	}

	return root, leafs, nil
}

// 递归构建所有非叶子节点
// nl 树的每一层的所有节点，第一次调用是最后一层，第二次递归则是倒数第二层
func buildIntermediate(nl []*Node, t *MerkleTree) (*Node, error) {
	var nodes []*Node
	for i := 0; i < len(nl); i += 2 {
		h := t.hashStrategy()
		var left, right int = i, i + 1
		if i+1 == len(nl) { // 如果是基数个，这里会触发。导致最后那个孤立节点复制自己一份，然后他两结合
			right = i
		}
		chash := append(nl[left].Hash, nl[right].Hash...)
		if _, err := h.Write(chash); err != nil {
			return nil, err
		}
		n := &Node{
			Left:  nl[left],
			Right: nl[right],
			Hash:  h.Sum(nil),
			Tree:  t,
		}
		nodes = append(nodes, n)
		nl[left].Parent = n
		nl[right].Parent = n
		if len(nl) == 2 {
			return n, nil
		}
	}
	return buildIntermediate(nodes, t)
}

func (m *MerkleTree) MerkleRoot() []byte {
	return m.merkleRoot
}

// 根据所有叶子节点重建树
func (m *MerkleTree) RebuildTree() error {
	var cs []Content
	for _, c := range m.Leafs {
		cs = append(cs, c.C)
	}
	root, leafs, err := buildWithContent(cs, m)
	if err != nil {
		return err
	}
	m.Root = root
	m.Leafs = leafs
	m.merkleRoot = root.Hash
	return nil
}

// 根据所有数据重建树
func (m *MerkleTree) RebuildTreeWith(cs []Content) error {
	root, leafs, err := buildWithContent(cs, m)
	if err != nil {
		return err
	}
	m.Root = root
	m.Leafs = leafs
	m.merkleRoot = root.Hash
	return nil
}

// 自校验
func (m *MerkleTree) VerifyTree() (bool, error) {
	calculatedMerkleRoot, err := m.Root.calculateNodeHashFromLeaf()
	if err != nil {
		return false, err
	}

	if bytes.Compare(m.merkleRoot, calculatedMerkleRoot) == 0 {
		return true, nil
	}
	return false, nil
}

// 校验这个数据是否在树中
func (m *MerkleTree) VerifyContent(content Content) (bool, error) {
	for _, l := range m.Leafs {
		ok, err := l.C.Equals(content)
		if err != nil {
			return false, err
		}

		if ok {
			currentParent := l.Parent
			for currentParent != nil {
				h := m.hashStrategy()
				rightBytes, err := currentParent.Right.calculateNodeHash()
				if err != nil {
					return false, err
				}

				leftBytes, err := currentParent.Left.calculateNodeHash()
				if err != nil {
					return false, err
				}

				if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
					return false, err
				}
				if bytes.Compare(h.Sum(nil), currentParent.Hash) != 0 {
					return false, nil
				}
				currentParent = currentParent.Parent
			}
			return true, nil
		}
	}
	return false, nil
}

func (n *Node) String() string {
	if n.leaf {
		return fmt.Sprintf("%s, %#v", BufferToHexString(n.Hash, false), n.C)
	}
	return BufferToHexString(n.Hash, false)
}

func (m *MerkleTree) String() string {
	str := ""
	m.Root.CurString(&str)
	return str
}

// 输出节点前序遍历的打印
func (m *Node) CurString(str *string) {
	if m.leaf {
		*str += fmt.Sprintln(m)
		return
	}
	*str += fmt.Sprintln(m)
	m.Left.CurString(str)
	m.Right.CurString(str)
}

func BufferToHexString(data []byte, prefix bool) string {
	hex := fmt.Sprintf(`%x`, data)
	if prefix {
		return `0x` + hex
	}
	return hex
}