package tools

import (
	"crypto/md5"
	"fmt"
	"runtime"
	"testing"
)

type Data struct {
	s string
}

func (data *Data) CalculateHash() ([]byte, error) {
	h := md5.New()
	if _, err := h.Write([]byte(data.s)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func (data *Data) Equals(other Content) (bool, error) {
	return data.s == other.(*Data).s, nil
}

func TestExampleGetMerklePath(t *testing.T) {
	data1 := Data{
		s: `111`,
	}
	data2 := Data{
		s: `222`,
	}
	data3 := Data{
		s: `333`,
	}
	tree, err := NewMerkleTree([]Content{&data1, &data2, &data3})
	if err != nil {
		panic(err)
	}

	fmt.Println(tree)
	a, b, err := tree.GetMerklePath(&data1)
	if err != nil {
		panic(err)
	}
	fmt.Println(a, b)
	// Output:
	// 99b04e1afa591732e6e779c4042e5be416bbdc79ff83baf018a6e1bb2a89aa61
	// 5adde34a358fb17ffbc93d842b48c82110f134664b95e8aa897be88d416fbf53
	// 698d51a19d8a121ce581499d7b701668, &go_blkchain_merkle_tree.Data{s:"111"}
	// bcbe3365e6ac95ea2c0343a2395834dd, &go_blkchain_merkle_tree.Data{s:"222"}
	// 5d755d4a14c4ddf4274088e6aaf578402164c73f3ea94c9e379f110ddf092a47
	// 310dcbbf4cce62f762a2aaa148d556bd, &go_blkchain_merkle_tree.Data{s:"333"}
	// 310dcbbf4cce62f762a2aaa148d556bd, &go_blkchain_merkle_tree.Data{s:"333"}
	//
	// [[188 190 51 101 230 172 149 234 44 3 67 162 57 88 52 221] [93 117 93 74 20 196 221 244 39 64 136 230 170 245 120 64 33 100 199 63 62 169 76 158 55 159 17 13 223 9 42 71]] [1 1]
}

func TestExampleMerkleTree_VerifyTree(t *testing.T) {
	data1 := Data{
		s: `111`,
	}
	data2 := Data{
		s: `222`,
	}
	tree, err := NewMerkleTree([]Content{&data1, &data2})
	if err != nil {
		panic(err)
	}
	bool_, err := tree.VerifyTree()
	if err != nil {
		panic(err)
	}
	fmt.Println(bool_)

	// Output:
	// true
}

func TestExampleMerkleTree_VerifyContent(t *testing.T) {
	data1 := Data{
		s: `111`,
	}
	data2 := Data{
		s: `222`,
	}
	tree, err := NewMerkleTree([]Content{&data1, &data2})
	if err != nil {
		panic(err)
	}
	data3 := Data{
		s: `333`,
	}
	bool_, err := tree.VerifyContent(&data3)
	if err != nil {
		panic(err)
	}
	fmt.Println(bool_)

	// Output:
	// false
	runtime.GC()
}
