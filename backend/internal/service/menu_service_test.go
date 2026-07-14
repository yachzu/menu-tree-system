package service

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/xos/Projects/stk-project/backend/internal/dto"
	"github.com/xos/Projects/stk-project/backend/internal/model"
	"github.com/xos/Projects/stk-project/backend/internal/repository"
)

var (
	rootID       = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	childID      = uuid.MustParse("00000000-0000-0000-0000-000000000002")
	parentID     = uuid.MustParse("00000000-0000-0000-0000-000000000003")
	orphanID     = uuid.MustParse("00000000-0000-0000-0000-000000000004")
	futureID     = uuid.MustParse("00000000-0000-0000-0000-000000000099")
	grandchildID = uuid.MustParse("00000000-0000-0000-0000-000000000005")
)

type mockMenuRepository struct {
	menus map[uuid.UUID]*model.Menu
}

func newMockMenuRepository() *mockMenuRepository {
	return &mockMenuRepository{menus: make(map[uuid.UUID]*model.Menu)}
}

func (m *mockMenuRepository) FindAll() ([]model.Menu, error) {
	var result []model.Menu
	for _, menu := range m.menus {
		result = append(result, *menu)
	}
	return result, nil
}

func (m *mockMenuRepository) FindByID(id uuid.UUID) (*model.Menu, error) {
	menu, ok := m.menus[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	menu.Children = nil
	for _, v := range m.menus {
		if v.ParentID != nil && *v.ParentID == id {
			menu.Children = append(menu.Children, v)
		}
	}
	return menu, nil
}

func (m *mockMenuRepository) FindChildren(parentID *uuid.UUID) ([]model.Menu, error) {
	var result []model.Menu
	for _, menu := range m.menus {
		if menu.ParentID == nil && parentID == nil {
			result = append(result, *menu)
		} else if menu.ParentID != nil && parentID != nil && *menu.ParentID == *parentID {
			result = append(result, *menu)
		}
	}
	return result, nil
}

func (m *mockMenuRepository) Create(menu *model.Menu) error {
	if menu.ID == uuid.Nil {
		menu.ID = uuid.New()
	}
	m.menus[menu.ID] = menu
	return nil
}

func (m *mockMenuRepository) Update(menu *model.Menu) error {
	if _, ok := m.menus[menu.ID]; !ok {
		return gorm.ErrRecordNotFound
	}
	m.menus[menu.ID] = menu
	return nil
}

func (m *mockMenuRepository) Delete(id uuid.UUID) error {
	if _, ok := m.menus[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(m.menus, id)
	return nil
}

func (m *mockMenuRepository) GetMaxOrderIndex(parentID *uuid.UUID) (int, error) {
	max := -1
	for _, menu := range m.menus {
		if menu.ParentID == nil && parentID == nil {
			if menu.OrderIndex > max {
				max = menu.OrderIndex
			}
		} else if menu.ParentID != nil && parentID != nil && *menu.ParentID == *parentID {
			if menu.OrderIndex > max {
				max = menu.OrderIndex
			}
		}
	}
	return max, nil
}

func (m *mockMenuRepository) IsAncestor(ancestorID, descendantID uuid.UUID) (bool, error) {
	current := descendantID
	for {
		menu, ok := m.menus[current]
		if !ok {
			return false, nil
		}
		if menu.ID == ancestorID {
			return true, nil
		}
		if menu.ParentID == nil {
			return false, nil
		}
		current = *menu.ParentID
	}
}

func fixedTime() time.Time {
	return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
}

func testMenu(id uuid.UUID, name string, parentID *uuid.UUID, depth, orderIndex int) *model.Menu {
	t := fixedTime()
	return &model.Menu{
		ID:         id,
		Name:       name,
		ParentID:   parentID,
		Depth:      depth,
		OrderIndex: orderIndex,
		CreatedAt:  t,
		UpdatedAt:  t,
	}
}

func TestNewMenuService(t *testing.T) {
	repo := newMockMenuRepository()
	svc := NewMenuService(repo)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestMenuService_GetTree(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*mockMenuRepository)
		wantIDs  []uuid.UUID
		wantErr  bool
		wantTree func([]*dto.MenuResponse) bool
	}{
		{
			name: "empty tree returns empty slice",
			setup: func(m *mockMenuRepository) {
			},
			wantIDs: nil,
			wantErr: false,
		},
		{
			name: "single root returns one node",
			setup: func(m *mockMenuRepository) {
				m.menus[rootID] = testMenu(rootID, "Root", nil, 0, 1)
			},
			wantIDs: []uuid.UUID{rootID},
			wantErr: false,
		},
		{
			name: "parent-child tree built correctly",
			setup: func(m *mockMenuRepository) {
				m.menus[rootID] = testMenu(rootID, "Root", nil, 0, 1)
				m.menus[childID] = testMenu(childID, "Child", &rootID, 1, 1)
			},
			wantIDs: []uuid.UUID{rootID},
			wantErr: false,
			wantTree: func(resp []*dto.MenuResponse) bool {
				if len(resp) != 1 || len(resp[0].Children) != 1 {
					return false
				}
				return resp[0].Children[0].ID == childID
			},
		},
		{
			name: "multiple roots sorted by order_index",
			setup: func(m *mockMenuRepository) {
				m.menus[childID] = testMenu(childID, "Second", nil, 0, 2)
				m.menus[rootID] = testMenu(rootID, "First", nil, 0, 1)
			},
			wantIDs: []uuid.UUID{rootID, childID},
			wantErr: false,
			wantTree: func(resp []*dto.MenuResponse) bool {
				if len(resp) != 2 {
					return false
				}
				return resp[0].OrderIndex == 1 && resp[1].OrderIndex == 2
			},
		},
		{
			name: "nested three levels",
			setup: func(m *mockMenuRepository) {
				m.menus[rootID] = testMenu(rootID, "Root", nil, 0, 1)
				m.menus[childID] = testMenu(childID, "Child", &rootID, 1, 1)
				m.menus[grandchildID] = testMenu(grandchildID, "Grandchild", &childID, 2, 1)
			},
			wantIDs: []uuid.UUID{rootID},
			wantErr: false,
			wantTree: func(resp []*dto.MenuResponse) bool {
				if len(resp) != 1 || len(resp[0].Children) != 1 {
					return false
				}
				child := resp[0].Children[0]
				if len(child.Children) != 1 {
					return false
				}
				return child.Children[0].ID == grandchildID
			},
		},
		{
			name: "orphan with non-existent parent is dropped from tree",
			setup: func(m *mockMenuRepository) {
				m.menus[rootID] = testMenu(rootID, "Root", nil, 0, 1)
				m.menus[orphanID] = testMenu(orphanID, "Orphan", &futureID, 1, 1)
			},
			wantIDs: []uuid.UUID{rootID},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockMenuRepository()
			if tt.setup != nil {
				tt.setup(repo)
			}
			svc := NewMenuService(repo)
			tree, err := svc.GetTree()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTree() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if tt.wantIDs == nil && len(tree) != 0 {
				t.Errorf("expected empty tree, got %d nodes", len(tree))
				return
			}
			if tt.wantIDs != nil {
				if len(tree) != len(tt.wantIDs) {
					t.Errorf("expected %d roots, got %d", len(tt.wantIDs), len(tree))
					return
				}
				for i, id := range tt.wantIDs {
					if tree[i].ID != id {
						t.Errorf("root[%d].ID = %v, want %v", i, tree[i].ID, id)
					}
				}
			}
			if tt.wantTree != nil && !tt.wantTree(tree) {
				t.Errorf("tree structure check failed")
			}
		})
	}
}

func TestMenuService_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      uuid.UUID
		setup   func(*mockMenuRepository)
		wantErr error
		wantID  uuid.UUID
	}{
		{
			name: "found",
			id:   rootID,
			setup: func(m *mockMenuRepository) {
				m.menus[rootID] = testMenu(rootID, "Root", nil, 0, 1)
			},
			wantErr: nil,
			wantID:  rootID,
		},
		{
			name:    "not found returns ErrNotFound",
			id:      futureID,
			setup:   nil,
			wantErr: ErrNotFound,
			wantID:  uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockMenuRepository()
			if tt.setup != nil {
				tt.setup(repo)
			}
			svc := NewMenuService(repo)
			resp, err := svc.GetByID(tt.id)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetByID() error = %v, wantErr = %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.ID != tt.wantID {
				t.Errorf("resp.ID = %v, want %v", resp.ID, tt.wantID)
			}
		})
	}
}

func TestMenuService_Create(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.CreateMenuRequest
		setup   func(*mockMenuRepository)
		wantErr error
		check   func(*testing.T, *dto.MenuResponse)
	}{
		{
			name: "create without parent sets depth 0",
			req: &dto.CreateMenuRequest{
				Name: "New Root",
			},
			setup:   nil,
			wantErr: nil,
			check: func(t *testing.T, resp *dto.MenuResponse) {
				if resp.Depth != 0 {
					t.Errorf("expected depth 0, got %d", resp.Depth)
				}
				if resp.Name != "New Root" {
					t.Errorf("expected name 'New Root', got %s", resp.Name)
				}
				if resp.ParentID != nil {
					t.Errorf("expected nil ParentID, got %v", resp.ParentID)
				}
				if resp.OrderIndex != 0 {
					t.Errorf("expected OrderIndex 0, got %d", resp.OrderIndex)
				}
			},
		},
		{
			name: "create with valid parent sets depth parent.depth+1",
			req: &dto.CreateMenuRequest{
				Name:     "New Child",
				ParentID: &rootID,
			},
			setup: func(m *mockMenuRepository) {
				m.menus[rootID] = testMenu(rootID, "Root", nil, 0, 1)
			},
			wantErr: nil,
			check: func(t *testing.T, resp *dto.MenuResponse) {
				if resp.Depth != 1 {
					t.Errorf("expected depth 1, got %d", resp.Depth)
				}
				if resp.ParentID == nil || *resp.ParentID != rootID {
					t.Errorf("expected ParentID %v, got %v", rootID, resp.ParentID)
				}
			},
		},
		{
			name: "create with invalid parent returns ErrInvalidParent",
			req: &dto.CreateMenuRequest{
				Name:     "Orphan",
				ParentID: &futureID,
			},
			setup:   nil,
			wantErr: ErrInvalidParent,
		},
		{
			name: "create increments order_index",
			req: &dto.CreateMenuRequest{
				Name: "Second Root",
			},
			setup: func(m *mockMenuRepository) {
				m.menus[rootID] = testMenu(rootID, "First Root", nil, 0, 0)
			},
			wantErr: nil,
			check: func(t *testing.T, resp *dto.MenuResponse) {
				if resp.OrderIndex != 1 {
					t.Errorf("expected OrderIndex 1, got %d", resp.OrderIndex)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockMenuRepository()
			if tt.setup != nil {
				tt.setup(repo)
			}
			svc := NewMenuService(repo)
			resp, err := svc.Create(tt.req)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Create() error = %v, wantErr = %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.ID == uuid.Nil {
				t.Error("expected non-nil ID")
			}
			if tt.check != nil {
				tt.check(t, resp)
			}
		})
	}
}

func TestMenuService_Update(t *testing.T) {
	tests := []struct {
		name    string
		id      uuid.UUID
		req     *dto.UpdateMenuRequest
		setup   func(*mockMenuRepository)
		wantErr error
	}{
		{
			name: "update name successfully",
			id:   rootID,
			req: &dto.UpdateMenuRequest{
				Name: "Updated Name",
			},
			setup: func(m *mockMenuRepository) {
				m.menus[rootID] = testMenu(rootID, "Original", nil, 0, 1)
			},
			wantErr: nil,
		},
		{
			name: "update not found returns ErrNotFound",
			id:   futureID,
			req: &dto.UpdateMenuRequest{
				Name: "Ghost",
			},
			setup:   nil,
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockMenuRepository()
			if tt.setup != nil {
				tt.setup(repo)
			}
			svc := NewMenuService(repo)
			resp, err := svc.Update(tt.id, tt.req)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Update() error = %v, wantErr = %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.Name != tt.req.Name {
				t.Errorf("resp.Name = %s, want %s", resp.Name, tt.req.Name)
			}
		})
	}
}

func TestMenuService_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      uuid.UUID
		setup   func(*mockMenuRepository)
		wantErr error
	}{
		{
			name: "delete existing item",
			id:   rootID,
			setup: func(m *mockMenuRepository) {
				m.menus[rootID] = testMenu(rootID, "To Delete", nil, 0, 1)
			},
			wantErr: nil,
		},
		{
			name:    "delete not found returns ErrNotFound",
			id:      futureID,
			setup:   nil,
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockMenuRepository()
			if tt.setup != nil {
				tt.setup(repo)
			}
			svc := NewMenuService(repo)
			err := svc.Delete(tt.id)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Delete() error = %v, wantErr = %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			_, err = svc.GetByID(tt.id)
			if !errors.Is(err, ErrNotFound) {
				t.Error("item was not deleted")
			}
		})
	}
}

func TestMenuService_Move(t *testing.T) {
	tests := []struct {
		name    string
		id      uuid.UUID
		req     *dto.MoveMenuRequest
		setup   func(*mockMenuRepository)
		wantErr error
		check   func(*testing.T, *dto.MenuResponse)
	}{
		{
			name: "move to different parent",
			id:   childID,
			req: &dto.MoveMenuRequest{
				NewParentID: &parentID,
				OrderIndex:  1,
			},
			setup: func(m *mockMenuRepository) {
				m.menus[childID] = testMenu(childID, "Child", &rootID, 1, 1)
				m.menus[rootID] = testMenu(rootID, "Root", nil, 0, 1)
				m.menus[parentID] = testMenu(parentID, "New Parent", nil, 0, 2)
			},
			wantErr: nil,
			check: func(t *testing.T, resp *dto.MenuResponse) {
				if resp.ParentID == nil || *resp.ParentID != parentID {
					t.Errorf("expected ParentID %v, got %v", parentID, resp.ParentID)
				}
				if resp.Depth != 1 {
					t.Errorf("expected depth 1, got %d", resp.Depth)
				}
				if resp.OrderIndex != 1 {
					t.Errorf("expected OrderIndex 1, got %d", resp.OrderIndex)
				}
			},
		},
		{
			name: "move to root sets ParentID nil and depth 0",
			id:   childID,
			req: &dto.MoveMenuRequest{
				NewParentID: nil,
				OrderIndex:  5,
			},
			setup: func(m *mockMenuRepository) {
				m.menus[childID] = testMenu(childID, "Child", &rootID, 1, 1)
				m.menus[rootID] = testMenu(rootID, "Root", nil, 0, 1)
			},
			wantErr: nil,
			check: func(t *testing.T, resp *dto.MenuResponse) {
				if resp.ParentID != nil {
					t.Errorf("expected nil ParentID, got %v", resp.ParentID)
				}
				if resp.Depth != 0 {
					t.Errorf("expected depth 0, got %d", resp.Depth)
				}
				if resp.OrderIndex != 5 {
					t.Errorf("expected OrderIndex 5, got %d", resp.OrderIndex)
				}
			},
		},
		{
			name: "move to self returns error",
			id:   rootID,
			req: &dto.MoveMenuRequest{
				NewParentID: &rootID,
				OrderIndex:  1,
			},
			setup: func(m *mockMenuRepository) {
				m.menus[rootID] = testMenu(rootID, "Root", nil, 0, 1)
			},
			wantErr: errors.New("cannot move a menu into itself"),
		},
		{
			name: "move to descendant returns ErrCircularRef",
			id:   rootID,
			req: &dto.MoveMenuRequest{
				NewParentID: &childID,
				OrderIndex:  1,
			},
			setup: func(m *mockMenuRepository) {
				m.menus[rootID] = testMenu(rootID, "Root", nil, 0, 1)
				m.menus[childID] = testMenu(childID, "Child", &rootID, 1, 1)
			},
			wantErr: ErrCircularRef,
		},
		{
			name: "move not found returns ErrNotFound",
			id:   futureID,
			req: &dto.MoveMenuRequest{
				NewParentID: &rootID,
				OrderIndex:  1,
			},
			setup: func(m *mockMenuRepository) {
				m.menus[rootID] = testMenu(rootID, "Root", nil, 0, 1)
			},
			wantErr: ErrNotFound,
		},
		{
			name: "move to invalid parent returns ErrInvalidParent",
			id:   childID,
			req: &dto.MoveMenuRequest{
				NewParentID: &futureID,
				OrderIndex:  1,
			},
			setup: func(m *mockMenuRepository) {
				m.menus[childID] = testMenu(childID, "Child", &rootID, 1, 1)
				m.menus[rootID] = testMenu(rootID, "Root", nil, 0, 1)
			},
			wantErr: ErrInvalidParent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockMenuRepository()
			if tt.setup != nil {
				tt.setup(repo)
			}
			svc := NewMenuService(repo)
			resp, err := svc.Move(tt.id, tt.req)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.wantErr.Error() != "" && err.Error() != tt.wantErr.Error() && !errors.Is(err, tt.wantErr) {
					t.Errorf("Move() error = %v, wantErr = %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.ID != tt.id {
				t.Errorf("resp.ID = %v, want %v", resp.ID, tt.id)
			}
			if tt.check != nil {
				tt.check(t, resp)
			}
		})
	}
}

func TestMenuService_Reorder(t *testing.T) {
	tests := []struct {
		name    string
		id      uuid.UUID
		req     *dto.ReorderMenuRequest
		setup   func(*mockMenuRepository)
		wantErr error
		wantIdx int
	}{
		{
			name: "reorder successfully",
			id:   rootID,
			req: &dto.ReorderMenuRequest{
				OrderIndex: 10,
			},
			setup: func(m *mockMenuRepository) {
				m.menus[rootID] = testMenu(rootID, "Root", nil, 0, 1)
			},
			wantErr: nil,
			wantIdx: 10,
		},
		{
			name: "reorder not found returns ErrNotFound",
			id:   futureID,
			req: &dto.ReorderMenuRequest{
				OrderIndex: 5,
			},
			setup:   nil,
			wantErr: ErrNotFound,
			wantIdx: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockMenuRepository()
			if tt.setup != nil {
				tt.setup(repo)
			}
			svc := NewMenuService(repo)
			resp, err := svc.Reorder(tt.id, tt.req)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Reorder() error = %v, wantErr = %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.OrderIndex != tt.wantIdx {
				t.Errorf("resp.OrderIndex = %d, want %d", resp.OrderIndex, tt.wantIdx)
			}
		})
	}
}

func TestMenuService_NestedDepthCalculation(t *testing.T) {
	repo := newMockMenuRepository()

	repo.menus[parentID] = testMenu(parentID, "Level 0", nil, 0, 1)
	repo.menus[childID] = testMenu(childID, "Level 1", &parentID, 1, 1)
	repo.menus[grandchildID] = testMenu(grandchildID, "Level 2", &childID, 2, 1)

	svc := NewMenuService(repo)

	// Create a new child under grandchild - depth should be 3
	newChildID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	req := &dto.CreateMenuRequest{
		Name:     "Level 3",
		ParentID: &grandchildID,
	}
	resp, err := svc.Create(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	repo.menus[newChildID] = &model.Menu{
		ID:         newChildID,
		Name:       resp.Name,
		ParentID:   resp.ParentID,
		Depth:      resp.Depth,
		OrderIndex: resp.OrderIndex,
		CreatedAt:  fixedTime(),
		UpdatedAt:  fixedTime(),
	}

	if resp.Depth != 3 {
		t.Errorf("expected depth 3, got %d", resp.Depth)
	}
}

func TestMenuRepositoryInterface(t *testing.T) {
	var _ repository.MenuRepository = (*mockMenuRepository)(nil)
}
