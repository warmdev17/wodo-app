package dtos

import "time"

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
