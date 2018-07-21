package BLC

import "crypto/sha256"

type Ghw_MerkelTree struct {
	Ghw_RootNode *Ghw_MerkelNode
}

type Ghw_MerkelNode struct {
	Ghw_Left  *Ghw_MerkelNode
	Ghw_Right *Ghw_MerkelNode
	Ghw_Data  []byte
}

func Ghw_NewMerkelTree(data [][]byte) *Ghw_MerkelTree {
	var nodes []Ghw_MerkelNode

	// 如果交易数据不是双数，将最后一个交易复制添加到最后
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	// 生成所有的一级节点，存储到node中
	for _, dataum := range data {
		node := Ghw_NewMerkelNode(nil, nil, dataum)
		nodes = append(nodes, *node)
	}

	// 遍历生成顶层节点
	for i := 0;i<len(data)/2 ;i++{
		var newLevel []Ghw_MerkelNode
		for j:=0 ; j<len(nodes) ;j+=2  {
			node := Ghw_NewMerkelNode(&nodes[j],&nodes[j+1],nil)
			newLevel = append(newLevel,*node)
		}
		nodes = newLevel
	}

	//for ; len(nodes)==1 ;{
	//	var newLevel []Ghw_MerkelNode
	//	for j:=0 ; j<len(nodes) ;j+=2  {
	//		node := Ghw_NewMerkelNode(&nodes[j],&nodes[j+1],nil)
	//		newLevel = append(newLevel,*node)
	//	}
	//	nodes = newLevel
	//}
	mTree := Ghw_MerkelTree{&nodes[0]}
	return &mTree
}

// 新叶节点
func Ghw_NewMerkelNode(left, right *Ghw_MerkelNode, data []byte) *Ghw_MerkelNode {
	mNode := Ghw_MerkelNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Ghw_Data = hash[:]
	} else {
		prevHashes := append(left.Ghw_Data, right.Ghw_Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Ghw_Data = hash[:]
	}

	mNode.Ghw_Left = left
	mNode.Ghw_Right = right

	return &mNode
}
