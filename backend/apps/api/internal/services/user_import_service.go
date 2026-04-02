package services

import (
	"context"
	"errors"
	"strings"

	"api/internal/models"
	"github.com/google/uuid"
)

type UserImportService struct{}

func NewUserImportService() *UserImportService {
	return &UserImportService{}
}

func (s *UserImportService) StartJob(_ context.Context, tenantID string, req models.CSVImportJobRequest) (models.CSVImportJob, error) {
	if strings.TrimSpace(req.FileAssetID) == "" || req.DefaultRole == "" {
		return models.CSVImportJob{}, errors.New("invalid csv import request")
	}

	return models.CSVImportJob{
		JobID:        uuid.NewString(),
		TenantID:     tenantID,
		Status:       "queued",
		TotalRows:    0,
		SuccessCount: 0,
		FailedCount:  0,
	}, nil
}

func (s *UserImportService) ValidateRows(rows []map[string]string) []models.CSVImportJobError {
	errorsOut := make([]models.CSVImportJobError, 0)
	for index, row := range rows {
		email := strings.TrimSpace(strings.ToLower(row["email"]))
		fullName := strings.TrimSpace(row["full_name"])
		if email == "" || fullName == "" {
			errorsOut = append(errorsOut, models.CSVImportJobError{
				Row:     index + 1,
				Message: "email and full_name are required",
			})
		}
	}
	return errorsOut
}
