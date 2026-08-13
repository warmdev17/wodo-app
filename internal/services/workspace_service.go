package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/warmdev17/Wodo-App/internal/dtos"
	"github.com/warmdev17/Wodo-App/internal/repositories"
)

type WorkspaceService struct {
	repo *repositories.Queries
}

func NewWorkspaceService(repo *repositories.Queries) *WorkspaceService {
	return &WorkspaceService{repo: repo}
}

func (s *WorkspaceService) CreateWorkspace(ctx context.Context, req dtos.CreateWorkspaceRequest, ownerID uuid.UUID) (dtos.WorkspaceResponse, error) {
	baseSlug := slug.Make(req.Name)
	slugPattern := fmt.Sprintf("%s%%", baseSlug)

	slugs, err := s.repo.ListExistingSlug(ctx, slugPattern)
	if err != nil {
		return dtos.WorkspaceResponse{}, ErrInternalServer
	}

	if len(slugs) > 0 {
		maxSuffix := 0
		hasBaseSlug := false

		prefix := baseSlug + "-" // VD: "my-workspace-"

		for _, s := range slugs {
			if s == baseSlug {
				hasBaseSlug = true
				continue
			}

			if suffixStr, hasPrefix := strings.CutPrefix(s, prefix); hasPrefix {
				if suffix, err := strconv.Atoi(suffixStr); err == nil && suffix > maxSuffix {
					maxSuffix = suffix
				}
			}
		}

		if hasBaseSlug || maxSuffix > 0 {
			baseSlug = fmt.Sprintf("%s-%d", baseSlug, maxSuffix+1)
		}
	}

	args := repositories.CreateWorkspaceParams{
		Name:    req.Name,
		Slug:    baseSlug,
		OwnerID: ownerID,
	}

	workspace, err := s.repo.CreateWorkspace(ctx, args)
	if err != nil {
		return dtos.WorkspaceResponse{}, ErrInternalServer
	}

	return dtos.WorkspaceResponse{
		Name:      workspace.Name,
		Slug:      workspace.Slug,
		CreatedAt: workspace.CreatedAt,
		UpdatedAt: workspace.UpdatedAt,
	}, nil
}

func (s *WorkspaceService) GetWorkspace(ctx context.Context, id uuid.UUID) (dtos.WorkspaceResponse, error) {
	workspace, err := s.repo.GetWorkspaceByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dtos.WorkspaceResponse{}, ErrWorkspaceNotFound
		}
		return dtos.WorkspaceResponse{}, ErrInternalServer
	}

	return dtos.ToWorkspaceResponse(workspace), nil
}
