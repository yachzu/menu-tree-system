package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/xos/Projects/stk-project/backend/internal/dto"
	"github.com/xos/Projects/stk-project/backend/internal/model"
	"github.com/xos/Projects/stk-project/backend/internal/repository"
)

var (
	ErrNotFound      = errors.New("menu not found")
	ErrCircularRef   = errors.New("circular reference: cannot move a menu into its own descendant")
	ErrInvalidParent = errors.New("parent menu not found")
)

type MenuService interface {
	GetTree() ([]*dto.MenuResponse, error)
	GetByID(id uuid.UUID) (*dto.MenuResponse, error)
	Create(req *dto.CreateMenuRequest) (*dto.MenuResponse, error)
	Update(id uuid.UUID, req *dto.UpdateMenuRequest) (*dto.MenuResponse, error)
	Delete(id uuid.UUID) error
	Move(id uuid.UUID, req *dto.MoveMenuRequest) (*dto.MenuResponse, error)
	Reorder(id uuid.UUID, req *dto.ReorderMenuRequest) (*dto.MenuResponse, error)
}

type menuService struct {
	repo repository.MenuRepository
}

func NewMenuService(repo repository.MenuRepository) MenuService {
	return &menuService{repo: repo}
}

func (s *menuService) GetTree() ([]*dto.MenuResponse, error) {
	menus, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch menus: %w", err)
	}

	menuMap := make(map[uuid.UUID]*dto.MenuResponse)
	for i := range menus {
		m := menus[i]
		menuMap[m.ID] = &dto.MenuResponse{
			ID:         m.ID,
			Name:       m.Name,
			ParentID:   m.ParentID,
			Depth:      m.Depth,
			OrderIndex: m.OrderIndex,
			CreatedAt:  m.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  m.UpdatedAt.Format(time.RFC3339),
			Children:   []*dto.MenuResponse{},
		}
	}

	var roots []*dto.MenuResponse
	for _, m := range menus {
		node := menuMap[m.ID]
		if m.ParentID == nil {
			roots = append(roots, node)
		} else {
			if parent, ok := menuMap[*m.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return sortChildren(roots), nil
}

func (s *menuService) GetByID(id uuid.UUID) (*dto.MenuResponse, error) {
	menu, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch menu: %w", err)
	}

	resp := toResponse(menu)
	for _, child := range menu.Children {
		resp.Children = append(resp.Children, toResponse(child))
	}
	resp.Children = sortResponseChildren(resp.Children)
	return resp, nil
}

func (s *menuService) Create(req *dto.CreateMenuRequest) (*dto.MenuResponse, error) {
	depth := 0
	if req.ParentID != nil {
		parent, err := s.repo.FindByID(*req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrInvalidParent
			}
			return nil, fmt.Errorf("failed to fetch parent: %w", err)
		}
		depth = parent.Depth + 1
	}

	maxOrder, err := s.repo.GetMaxOrderIndex(req.ParentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get max order: %w", err)
	}

	menu := &model.Menu{
		Name:       req.Name,
		ParentID:   req.ParentID,
		Depth:      depth,
		OrderIndex: maxOrder + 1,
	}

	if err := s.repo.Create(menu); err != nil {
		return nil, fmt.Errorf("failed to create menu: %w", err)
	}

	return toResponse(menu), nil
}

func (s *menuService) Update(id uuid.UUID, req *dto.UpdateMenuRequest) (*dto.MenuResponse, error) {
	menu, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch menu: %w", err)
	}

	if req.Name != "" {
		menu.Name = req.Name
	}

	if err := s.repo.Update(menu); err != nil {
		return nil, fmt.Errorf("failed to update menu: %w", err)
	}

	return toResponse(menu), nil
}

func (s *menuService) Delete(id uuid.UUID) error {
	if _, err := s.repo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to fetch menu: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete menu: %w", err)
	}
	return nil
}

func (s *menuService) isDescendant(targetID, potentialAncestorID uuid.UUID) (bool, error) {
	if targetID == potentialAncestorID {
		return true, nil
	}
	return s.repo.IsAncestor(potentialAncestorID, targetID)
}

func (s *menuService) Move(id uuid.UUID, req *dto.MoveMenuRequest) (*dto.MenuResponse, error) {
	menu, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch menu: %w", err)
	}

	if req.NewParentID != nil {
		if *req.NewParentID == id {
			return nil, errors.New("cannot move a menu into itself")
		}

		isDesc, err := s.isDescendant(*req.NewParentID, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check circular reference: %w", err)
		}
		if isDesc {
			return nil, ErrCircularRef
		}

		parent, err := s.repo.FindByID(*req.NewParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrInvalidParent
			}
			return nil, fmt.Errorf("failed to fetch parent: %w", err)
		}
		menu.Depth = parent.Depth + 1
	} else {
		menu.Depth = 0
	}

	menu.ParentID = req.NewParentID
	menu.OrderIndex = req.OrderIndex

	if err := s.repo.Update(menu); err != nil {
		return nil, fmt.Errorf("failed to move menu: %w", err)
	}

	return toResponse(menu), nil
}

func (s *menuService) Reorder(id uuid.UUID, req *dto.ReorderMenuRequest) (*dto.MenuResponse, error) {
	menu, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch menu: %w", err)
	}

	menu.OrderIndex = req.OrderIndex

	if err := s.repo.Update(menu); err != nil {
		return nil, fmt.Errorf("failed to reorder menu: %w", err)
	}

	return toResponse(menu), nil
}

func toResponse(m *model.Menu) *dto.MenuResponse {
	return &dto.MenuResponse{
		ID:         m.ID,
		Name:       m.Name,
		ParentID:   m.ParentID,
		Depth:      m.Depth,
		OrderIndex: m.OrderIndex,
		CreatedAt:  m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  m.UpdatedAt.Format(time.RFC3339),
		Children:   []*dto.MenuResponse{},
	}
}

func sortResponseChildren(children []*dto.MenuResponse) []*dto.MenuResponse {
	for i := 0; i < len(children); i++ {
		for j := i + 1; j < len(children); j++ {
			if children[i].OrderIndex > children[j].OrderIndex {
				children[i], children[j] = children[j], children[i]
			}
		}
	}
	return children
}

func sortChildren(roots []*dto.MenuResponse) []*dto.MenuResponse {
	for _, node := range roots {
		node.Children = sortChildren(node.Children)
	}
	return sortResponseChildren(roots)
}
