package services

import (
	"context"
)

type ProductService struct {
}

func (this *ProductService) GetProdStock(ctx context.Context, in *ProdRequest) (*ProdResponse, error) {
	return &ProdResponse{ProdStock: 20}, nil
}
