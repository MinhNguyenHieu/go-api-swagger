package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"external-backend-go/internal/logger"
	"external-backend-go/internal/request"
	"external-backend-go/internal/service"
	"external-backend-go/internal/utility"
)

type ItemHandler struct {
	ItemService *service.ItemService
	Logger      *logger.Logger
}

func NewItemHandler(itemService *service.ItemService, logger *logger.Logger) *ItemHandler {
	return &ItemHandler{ItemService: itemService, Logger: logger}
}

// @Summary Create a new item
// @Description Creates a new item with a name and description. Requires JWT authentication.
// @Tags items
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body request.CreateItemRequest true "Item creation details"
// @Success 201 {object} model.Item "Created item"
// @Failure 400 {object} map[string]string "message: Invalid request data"
// @Failure 401 {object} map[string]string "message: Authentication token required / Invalid token"
// @Failure 500 {object} map[string]string "message: Internal server error"
// @Router /items [post]
func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var req request.CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utility.BadRequestResponse(w, r, fmt.Errorf("Invalid request data"), h.Logger)
		return
	}

	if err := req.Validate(); err != nil {
		utility.BadRequestResponse(w, r, err, h.Logger)
		return
	}

	item, err := h.ItemService.CreateItem(r.Context(), req.Name, req.Description)
	if err != nil {
		utility.InternalServerError(w, r, err, h.Logger)
		return
	}

	utility.JSONResponse(w, http.StatusCreated, item)
}

// @Summary Get an item by ID
// @Description Retrieves a single item by its ID. Requires JWT authentication.
// @Tags items
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Item ID"
// @Success 200 {object} model.Item "Found item"
// @Failure 400 {object} map[string]string "message: Invalid item ID"
// @Failure 401 {object} map[string]string "message: Authentication token required / Invalid token"
// @Failure 404 {object} map[string]string "message: Item not found"
// @Failure 500 {object} map[string]string "message: Internal server error"
// @Router /items/{id} [get]
func (h *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		utility.BadRequestResponse(w, r, fmt.Errorf("Item ID is missing"), h.Logger)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utility.BadRequestResponse(w, r, fmt.Errorf("Invalid item ID: %w", err), h.Logger)
		return
	}

	item, err := h.ItemService.GetItemByID(r.Context(), int32(id))
	if err != nil {
		if err.Error() == "item not found" {
			utility.NotFoundResponse(w, r, h.Logger)
		} else {
			utility.InternalServerError(w, r, err, h.Logger)
		}
		return
	}

	utility.JSONResponse(w, http.StatusOK, item)
}

// @Summary Update an item
// @Description Updates an existing item's name and description by its ID. Requires JWT authentication.
// @Tags items
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Item ID"
// @Param request body request.UpdateItemRequest true "Item update details"
// @Success 200 {object} model.Item "Updated item"
// @Failure 400 {object} map[string]string "message: Invalid request data / Invalid item ID"
// @Failure 401 {object} map[string]string "message: Authentication token required / Invalid token"
// @Failure 404 {object} map[string]string "message: Item not found"
// @Failure 500 {object} map[string]string "message: Internal server error"
// @Router /items/{id} [put]
func (h *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		utility.BadRequestResponse(w, r, fmt.Errorf("Item ID is missing"), h.Logger)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utility.BadRequestResponse(w, r, fmt.Errorf("Invalid item ID: %w", err), h.Logger)
		return
	}

	var req request.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utility.BadRequestResponse(w, r, fmt.Errorf("Invalid request data"), h.Logger)
		return
	}

	if err := req.Validate(); err != nil {
		utility.BadRequestResponse(w, r, err, h.Logger)
		return
	}

	item, err := h.ItemService.UpdateItem(r.Context(), int32(id), req.Name, req.Description)
	if err != nil {
		if err.Error() == "item not found" {
			utility.NotFoundResponse(w, r, h.Logger)
		} else {
			utility.InternalServerError(w, r, err, h.Logger)
		}
		return
	}

	utility.JSONResponse(w, http.StatusOK, item)
}

// @Summary Delete an item
// @Description Deletes an item by its ID. Requires JWT authentication.
// @Tags items
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Item ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "message: Invalid item ID"
// @Failure 401 {object} map[string]string "message: Authentication token required / Invalid token"
// @Failure 404 {object} map[string]string "message: Item not found"
// @Failure 500 {object} map[string]string "message: Internal server error"
// @Router /items/{id} [delete]
func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		utility.BadRequestResponse(w, r, fmt.Errorf("Item ID is missing"), h.Logger)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utility.BadRequestResponse(w, r, fmt.Errorf("Invalid item ID: %w", err), h.Logger)
		return
	}

	err = h.ItemService.DeleteItem(r.Context(), int32(id))
	if err != nil {
		if err.Error() == "item not found" {
			utility.NotFoundResponse(w, r, h.Logger)
		} else {
			utility.InternalServerError(w, r, err, h.Logger)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Get list of items
// @Description Retrieves a paginated list of sample items, requires JWT authentication.
// @Tags items
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number (default 1)"
// @Param pageSize query int false "Number of items per page (default 10)"
// @Success 200 {object} service.PaginatedItems "Paginated list of items"
// @Failure 400 {object} map[string]string "message: Invalid pagination parameters"
// @Failure 401 {object} map[string]string "message: Authentication token required / Invalid token"
// @Failure 500 {object} map[string]string "message: Internal server error"
// @Router /items [get]
func (h *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	paginatedItems, err := h.ItemService.GetItems(r.Context(), page, pageSize)
	if err != nil {
		utility.InternalServerError(w, r, err, h.Logger)
		return
	}

	utility.JSONResponse(w, http.StatusOK, paginatedItems)
}
