package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/warmdev17/Wodo-App/internal/dtos"
	"github.com/warmdev17/Wodo-App/internal/services"
	res "github.com/warmdev17/Wodo-App/pkg/response"
	"github.com/warmdev17/Wodo-App/pkg/utils"
)

type WorkspaceController struct {
	workspaceService *services.WorkspaceService
}

func NewWorkspaceController(workspaceService *services.WorkspaceService) *WorkspaceController {
	return &WorkspaceController{workspaceService: workspaceService}
}

func (c *WorkspaceController) Create(ctx *gin.Context) {
	var req dtos.CreateWorkspaceRequest

	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		if errors.Is(err, services.ErrContextNotFound) {
			res.Unauthorized(ctx, "Unauthorized", "User context not found")
			return
		}
		res.InternalError(ctx, "Error from server side", err.Error())
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	workspace, err := c.workspaceService.CreateWorkspace(ctx, req, userID)
	if err != nil {
		if errors.Is(err, services.ErrInternalServer) {
			res.InternalError(ctx, "Error from server side", err.Error())
			return
		}
	}
	res.Created(ctx, "Workspace created successfully", workspace)
}

func (c *WorkspaceController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")

	workspaceID, err := uuid.Parse(idStr)
	if err != nil {
		res.BadRequest(ctx, "ID is invalid", err.Error())
		return
	}

	workspace, err := c.workspaceService.GetByID(ctx.Request.Context(), workspaceID)
	if err != nil {
		if errors.Is(err, services.ErrWorkspaceNotFound) {
			res.NotFound(ctx, "Workspace not found")
			return
		}
	}
	res.Success(ctx, "Get workspace successfully", workspace)
}
