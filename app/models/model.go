package models

import "goblog/pkg/types"

// BaseModel 模型基类
type BaseModel struct {
	ID uint64
}

// GetStringID 获取ID
func (m BaseModel) GetStringID() string {
	return types.Uint64ToString(m.ID)
}
