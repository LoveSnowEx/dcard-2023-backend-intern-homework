package db

import (
	"fmt"

	pl "github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db/pagelist"
	"github.com/google/uuid"
)

// Create a new page list
func (db *DB) CreatePageList(list *pl.PageList) error {
	err := db.DB.Create(&list).Error
	if err != nil {
		return fmt.Errorf("create page list: %v", err)
	}
	end := pl.PageNode{
		End:     true,
		ListKey: list.Key,
	}
	err = db.DB.Create(&end).Error
	if err != nil {
		return fmt.Errorf("after create page list: %v", err)
	}
	return nil
}

// Retrieve a page list by key
func (db *DB) GetPageListByKey(key uuid.UUID) (*pl.PageList, error) {
	var list pl.PageList
	err := db.DB.First(&list, "key = ?", key).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// Delete a page list
func (db *DB) DeletePageList(key uuid.UUID) error {
	err := db.DB.Delete(&pl.PageNode{}, "list_key = ?", key).Error
	if err != nil {
		return fmt.Errorf("delete page list: %v", err)
	}
	err = db.DB.Delete(&pl.PageList{}, "key = ?", key).Error
	if err != nil {
		return fmt.Errorf("delete page list: %v", err)
	}
	return nil
}

func (db *DB) ClonePageList(key uuid.UUID) (*pl.PageList, error) {
	list, err := db.GetPageListByKey(key)
	if err != nil {
		return nil, err
	}
	newList := pl.New()
	db.CreatePageList(newList)
	nodes, err := db.GetPageNodesByListKeySorted(list.Key)

	if err != nil {
		return nil, fmt.Errorf("clone page list: %v", err)
	}
	for i := range nodes {
		_, err = db.PushBackPageList(newList.Key, nodes[i].PageID)
		if err != nil {
			return nil, fmt.Errorf("clone page list: %v", err)
		}
	}
	return newList, nil
}

// Retrieve a page node by key
func (db *DB) GetPageNodeByKey(key uuid.UUID) (*pl.PageNode, error) {
	var node pl.PageNode
	err := db.DB.First(&node, "key = ?", key).Error
	if err != nil {
		return nil, fmt.Errorf("get page node by key: %v", err)
	}
	return &node, nil
}

// Retrieve page nodes by list key
func (db *DB) GetPageNodesByListKey(key uuid.UUID) ([]pl.PageNode, error) {
	var nodes []pl.PageNode
	err := db.DB.Find(&nodes, "key = ?", key).Error
	if err != nil {
		return nil, fmt.Errorf("get page nodes by list key: %v", err)
	}
	return nodes, nil
}

// Retrieve sroted page nodes by list key
func (db *DB) GetPageNodesByListKeySorted(key uuid.UUID) ([]pl.PageNode, error) {
	nodes, err := db.GetPageNodesByListKey(key)
	if err != nil {
		return nil, fmt.Errorf("get page nodes by list key sorted: %v", err)
	}
	var end *pl.PageNode
	m := make(map[uuid.UUID]*pl.PageNode)
	for i := range nodes {
		m[nodes[i].Key] = &nodes[i]
		if nodes[i].End {
			end = &nodes[i]
		}
	}
	if end == nil {
		return nil, fmt.Errorf("get page nodes by list key sorted: end node not found")
	}
	var sorted []pl.PageNode
	for cur := m[end.NextKey]; !cur.End; cur = m[cur.NextKey] {
		sorted = append(sorted, *cur)
	}
	sorted = append(sorted, *end)
	return sorted, nil
}

// Retrieve end node by list key
func (db *DB) GetPageListEnd(key uuid.UUID) (*pl.PageNode, error) {
	l, err := db.GetPageListByKey(key)
	if err != nil {
		return nil, fmt.Errorf("get page list end: %v", err)
	}
	var end pl.PageNode
	err = db.DB.Where("list_key = ? AND end = ?", l.Key, true).First(&end).Error
	if err != nil {
		return nil, fmt.Errorf("get page list end node: %v", err)
	}
	return &end, nil
}

// Retrieve begin node by list key
func (db *DB) GetPageListBegin(key uuid.UUID) (*pl.PageNode, error) {
	end, err := db.GetPageListEnd(key)
	if err != nil {
		return nil, fmt.Errorf("get page list begin: %v", err)
	}
	var begin pl.PageNode
	err = db.DB.First(&begin, "key = ?", end.NextKey).Error
	if err != nil {
		return nil, fmt.Errorf("get page list end node: %v", err)
	}
	return &begin, nil
}

// Retrieve next page node
func (db *DB) GetPageNodeNext(key uuid.UUID) (*pl.PageNode, error) {
	node, err := db.GetPageNodeByKey(key)
	if err != nil {
		return nil, fmt.Errorf("get page node next: %v", err)
	}
	var next pl.PageNode
	err = db.DB.First(&next, "key = ?", node.NextKey).Error
	if err != nil {
		return nil, fmt.Errorf("get page node next: %v", err)
	}
	return &next, nil
}

// Retrieve prev page node
func (db *DB) GetPageNodePrev(key uuid.UUID) (*pl.PageNode, error) {
	node, err := db.GetPageNodeByKey(key)
	if err != nil {
		return nil, fmt.Errorf("get page node prev: %v", err)
	}
	var prev pl.PageNode
	err = db.DB.First(&prev, "key = ?", node.PrevKey).Error
	if err != nil {
		return nil, fmt.Errorf("get page node prev: %v", err)
	}
	return &prev, nil
}

// Insert a page node
func (db *DB) InsertPageNode(key uuid.UUID, pageID uint) (*pl.PageNode, error) {
	pos, err := db.GetPageNodeByKey(key)
	if err != nil {
		return nil, fmt.Errorf("insert page: %v", err)
	}
	node := &pl.PageNode{
		ListKey: pos.ListKey,
		PageID:  pageID,
		PrevKey: pos.PrevKey,
		NextKey: pos.Key,
	}
	err = db.DB.Create(&node).Error
	if err != nil {
		return nil, fmt.Errorf("insert page: %v", err)
	}
	err = db.DB.Model(&pl.PageNode{}).Where("key = ?", pos.PrevKey).Update("next_key", node.Key).Error
	if err != nil {
		return nil, fmt.Errorf("insert page: %v", err)
	}
	err = db.DB.Model(&pl.PageNode{}).Where("key = ?", pos.Key).Update("prev_key", node.Key).Error
	if err != nil {
		return nil, fmt.Errorf("insert page: %v", err)
	}
	return node, nil
}

// Erase a page node
func (db *DB) ErasePageNode(key uuid.UUID) (*pl.PageNode, error) {
	pos, err := db.GetPageNodeByKey(key)
	if err != nil {
		return nil, fmt.Errorf("erase page node: %v", err)
	}
	if pos.End {
		return nil, fmt.Errorf("erase page node: cannot erase end node")
	}
	err = db.DB.Delete(&pl.PageNode{}, "key = ?", pos.Key).Error
	if err != nil {
		return nil, fmt.Errorf("erase page node: %v", err)
	}
	err = db.DB.Model(&pl.PageNode{}).Where("key = ?", pos.PrevKey).Update("next_key", pos.NextKey).Error
	if err != nil {
		return nil, fmt.Errorf("erase page node: %v", err)
	}
	err = db.DB.Model(&pl.PageNode{}).Where("key = ?", pos.NextKey).Update("prev_key", pos.PrevKey).Error
	if err != nil {
		return nil, fmt.Errorf("erase page node: %v", err)
	}
	var next pl.PageNode
	err = db.DB.First(&next, "key = ?", pos.NextKey).Error
	if err != nil {
		return nil, fmt.Errorf("erase page node: %v", err)
	}
	return &next, nil
}

// Set page node page id
func (db *DB) SetPageNode(key uuid.UUID, pageID uint) (*pl.PageNode, error) {
	var node pl.PageNode
	err := db.DB.Model(&pl.PageNode{}).Where("key = ?", key).Update("page_id", pageID).First(&node).Error
	if err != nil {
		return nil, fmt.Errorf("set page node: %v", err)
	}
	return &node, nil
}

// Push a page to the back of the list
func (db *DB) PushBackPageList(key uuid.UUID, pageID uint) (*pl.PageNode, error) {
	end, err := db.GetPageListEnd(key)
	if err != nil {
		return nil, fmt.Errorf("push back page list: %v", err)
	}
	node, err := db.InsertPageNode(end.Key, pageID)
	if err != nil {
		return nil, fmt.Errorf("push back page list: %v", err)
	}
	return node, nil
}

// Pop a page from the back of the list
func (db *DB) PopBackPageList(key uuid.UUID) error {
	end, err := db.GetPageListEnd(key)
	if err != nil {
		return fmt.Errorf("pop back page list: %v", err)
	}
	_, err = db.ErasePageNode(end.Key)
	if err != nil {
		return fmt.Errorf("pop back page list: %v", err)
	}
	return nil
}

// Push a page to the front of the list
func (db *DB) PushFrontPageList(key uuid.UUID, pageID uint) (*pl.PageNode, error) {
	begin, err := db.GetPageListBegin(key)
	if err != nil {
		return nil, fmt.Errorf("push front page list: %v", err)
	}
	node, err := db.InsertPageNode(begin.Key, pageID)
	if err != nil {
		return nil, fmt.Errorf("push front page list: %v", err)
	}
	return node, nil
}

// Pop a page from the front of the list
func (db *DB) PopFrontPageList(key uuid.UUID) error {
	begin, err := db.GetPageListBegin(key)
	if err != nil {
		return fmt.Errorf("pop front page list: %v", err)
	}
	_, err = db.ErasePageNode(begin.Key)
	if err != nil {
		return fmt.Errorf("pop front page list: %v", err)
	}
	return nil
}
