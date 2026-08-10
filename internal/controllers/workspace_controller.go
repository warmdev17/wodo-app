package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/warmdev17/Wodo-App/internal/dtos"
	"github.com/warmdev17/Wodo-App/internal/services"
	"github.com/warmdev17/Wodo-App/pkg/response"
)

type WorkspaceController struct {
	workspaceService *services.WorkspaceService
}

func NewWorkspaceController(workspaceService *services.WorkspaceService) *WorkspaceController {
	return &WorkspaceController{workspaceService: workspaceService}
}

func (c *WorkspaceController) Create(ctx *gin.Context) {
	var req dtos.CreateWorkspaceRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}
}
