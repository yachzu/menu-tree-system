package dto

import "github.com/google/uuid"

type CreateMenuRequest struct {
	Name     string     `json:"name" binding:"required,min=1,max=255"`
	ParentID *uuid.UUID `json:"parent_id"`
}

type UpdateMenuRequest struct {
	Name string `json:"name" binding:"omitempty,min=1,max=255"`
}

type MoveMenuRequest struct {
	NewParentID *uuid.UUID `json:"new_parent_id"`
	OrderIndex  int        `json:"order_index" binding:"min=0"`
}

type ReorderMenuRequest struct {
	OrderIndex int `json:"order_index" binding:"required,min=0"`
}

type MenuResponse struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	ParentID   *uuid.UUID      `json:"parent_id"`
	Depth      int             `json:"depth"`
	OrderIndex int             `json:"order_index"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
	Children   []*MenuResponse `json:"children,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
