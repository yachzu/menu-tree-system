package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/xos/Projects/stk-project/backend/internal/dto"
	"github.com/xos/Projects/stk-project/backend/internal/service"
)

type MenuHandler struct {
	svc service.MenuService
}

func NewMenuHandler(svc service.MenuService) *MenuHandler {
	return &MenuHandler{svc: svc}
}

func (h *MenuHandler) GetTree(c *gin.Context) {
	tree, err := h.svc.GetTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "an unexpected error occurred",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tree})
}

// GetByID godoc
// @Summary Get single menu item
// @Description Get a single menu item by ID including its children
// @Tags menus
// @Accept json
// @Produce json
// @Param id path string true "Menu ID"
// @Success 200 {object} dto.MenuResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /menus/{id} [get]
func (h *MenuHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "invalid menu ID format",
		})
		return
	}

	menu, err := h.svc.GetByID(id)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "menu not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "an unexpected error occurred",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": menu})
}

// Create godoc
// @Summary Create new menu item
// @Description Create a new menu item with optional parent
// @Tags menus
// @Accept json
// @Produce json
// @Param request body dto.CreateMenuRequest true "Menu data"
// @Success 201 {object} dto.MenuResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /menus [post]
func (h *MenuHandler) Create(c *gin.Context) {
	var req dto.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	menu, err := h.svc.Create(&req)
	if err != nil {
		if err == service.ErrInvalidParent {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "invalid_parent",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "an unexpected error occurred",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": menu})
}

// Update godoc
// @Summary Update menu item
// @Description Update a menu item's name
// @Tags menus
// @Accept json
// @Produce json
// @Param id path string true "Menu ID"
// @Param request body dto.UpdateMenuRequest true "Update data"
// @Success 200 {object} dto.MenuResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /menus/{id} [put]
func (h *MenuHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "invalid menu ID format",
		})
		return
	}

	var req dto.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	menu, err := h.svc.Update(id, &req)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "menu not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "an unexpected error occurred",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": menu})
}

// Delete godoc
// @Summary Delete menu item
// @Description Delete a menu item. Children are cascade deleted at database level.
// @Tags menus
// @Accept json
// @Produce json
// @Param id path string true "Menu ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /menus/{id} [delete]
func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "invalid menu ID format",
		})
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "menu not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "an unexpected error occurred",
		})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// Move godoc
// @Summary Move menu item to different parent
// @Description Move a menu item to a different parent node, including circular reference detection
// @Tags menus
// @Accept json
// @Produce json
// @Param id path string true "Menu ID"
// @Param request body dto.MoveMenuRequest true "Move data"
// @Success 200 {object} dto.MenuResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /menus/{id}/move [patch]
func (h *MenuHandler) Move(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "invalid menu ID format",
		})
		return
	}

	var req dto.MoveMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	menu, err := h.svc.Move(id, &req)
	if err != nil {
		if err == service.ErrNotFound || err == service.ErrInvalidParent {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "invalid_request",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "an unexpected error occurred",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": menu})
}

// Reorder godoc
// @Summary Reorder menu item within same level
// @Description Change the order index of a menu item among its siblings
// @Tags menus
// @Accept json
// @Produce json
// @Param id path string true "Menu ID"
// @Param request body dto.ReorderMenuRequest true "Reorder data"
// @Success 200 {object} dto.MenuResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /menus/{id}/reorder [patch]
func (h *MenuHandler) Reorder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "invalid menu ID format",
		})
		return
	}

	var req dto.ReorderMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	menu, err := h.svc.Reorder(id, &req)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "menu not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "an unexpected error occurred",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": menu})
}
