package service

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
)

type ResourceService interface {
	List(ctx context.Context) ([]communication.ResourceDtoResponse, error)
}

type resourceServiceImpl struct {
	conn db.ConnectionPool

	resourceRepo repositories.ResourceRepository
}

func NewResourceService(conn db.ConnectionPool, repos repositories.Repositories) ResourceService {
	return &resourceServiceImpl{
		conn:         conn,
		resourceRepo: repos.Resource,
	}
}

func (s *resourceServiceImpl) List(ctx context.Context) ([]communication.ResourceDtoResponse, error) {
	tx, err := s.conn.StartTransaction(ctx)
	if err != nil {
		return []communication.ResourceDtoResponse{}, err
	}
	defer tx.Close(ctx)

	resources, err := s.resourceRepo.List(ctx, tx)
	if err != nil {
		return []communication.ResourceDtoResponse{}, err
	}

	var out []communication.ResourceDtoResponse
	for _, resource := range resources {
		dto := communication.ToResourceDtoResponse(resource)
		out = append(out, dto)
	}

	return out, nil
}
