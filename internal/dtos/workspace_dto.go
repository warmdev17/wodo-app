package dtos

import (
	"time"

	"github.com/warmdev17/Wodo-App/internal/repositories"
)

type CreateWorkspaceRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateWorkspaceRequest struct {
	Name *string `json:"name"`
}

type WorkspaceResponse struct {
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToWorkspaceResponse(w repositories.Workspace) WorkspaceResponse {
	return WorkspaceResponse{
		Name:      w.Name,
		Slug:      w.Slug,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}
