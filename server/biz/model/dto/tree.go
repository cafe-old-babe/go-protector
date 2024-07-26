package dto

import "reflect"

type TreeNode struct {
	ID       uint64     `json:"id"`
	PID      uint64     `json:"pid"`
	Name     string     `json:"name"`
	Selected bool       `json:"selected"`
	Children []TreeNode `json:"children"`
}

// 4-18	【实战】部门与岗位管理-部门接口开发-掌握GO语言反射基本操作
func GenerateTree(slice any, rootId uint64, idField, pidField, nameField string,
	selectedSet map[uint64]any) (node *TreeNode) {
	// 判断 slice 是否为 slice 类型
	if selectedSet == nil {
		selectedSet = make(map[uint64]any)
	}
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return
	}
	// 判断 slice 是否为空
	if reflect.ValueOf(slice).Len() == 0 {
		return
	}
	// 创建 TreeNode 对象
	if rootId == 0 {
		node = &TreeNode{
			ID:       rootId,
			Name:     "根节点",
			Children: make([]TreeNode, 0),
		}
	} else {
		// 遍历 slice
		for i := 0; i < reflect.ValueOf(slice).Len(); i++ {
			// 获取当前元素的值
			item := reflect.ValueOf(slice).Index(i)
			// 获取当前元素的 id 字段的值
			id := item.FieldByName(idField).Uint()
			name := item.FieldByName(nameField).String()
			if id == rootId {
				node = &TreeNode{
					ID:       rootId,
					Name:     name,
					Children: make([]TreeNode, 0),
				}
				_, node.Selected = selectedSet[node.ID]
				break
			}
		}
	}
	if node == nil {
		return
	}
	// 遍历 slice
	for i := 0; i < reflect.ValueOf(slice).Len(); i++ {
		// 获取当前元素的值
		item := reflect.ValueOf(slice).Index(i)
		// 获取当前元素的 id 和 pid 字段的值
		id := item.FieldByName(idField).Uint()
		pid := item.FieldByName(pidField).Uint()
		// 判断当前元素的 pid 是否等于 rootId
		if pid == rootId {
			// 创建子节点
			childNode := GenerateTree(slice, id, idField, pidField, nameField, selectedSet)
			childNode.PID = pid
			if childNode != nil {
				// 将子节点添加到当前节点的 children 字段中
				node.Children = append(node.Children, *childNode)
			}

		}
	}
	return
}
