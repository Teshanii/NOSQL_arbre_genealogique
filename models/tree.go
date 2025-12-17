package models

type TreeNode struct {
	Individual Individual  `json:"individual"`
	Children   []*TreeNode `json:"children"`
}
