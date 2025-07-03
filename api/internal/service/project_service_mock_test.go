package service

import (
	"context"
	"errors"
	"testing"

	"github.com/BennyEisner/test-results/internal/models"
	mock_models "github.com/BennyEisner/test-results/internal/models/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProjectService_GetProjectByID_WithMock(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock repository
	mockRepo := mock_models.NewMockProjectRepository(ctrl)

	// Create the service with the mock
	service := NewProjectService(mockRepo)

	t.Run("success", func(t *testing.T) {
		expectedProject := &models.Project{ID: 1, Name: "Test Project"}

		// Set up mock expectations
		mockRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(expectedProject, nil)

		// Call the service method
		result, err := service.GetProjectByID(1)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedProject, result)
	})

	t.Run("not found", func(t *testing.T) {
		// Set up mock expectations for not found
		mockRepo.EXPECT().GetByID(gomock.Any(), int64(999)).Return(nil, errors.New("not found"))

		// Call the service method
		result, err := service.GetProjectByID(999)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProjectService_CreateProject_WithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_models.NewMockProjectRepository(ctrl)
	service := NewProjectService(mockRepo)

	t.Run("success", func(t *testing.T) {
		projectName := "New Project"
		expectedProject := &models.Project{ID: 1, Name: projectName}

		// Set up mock expectations
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, p *models.Project) error {
				p.ID = 1
				return nil
			})

		// Call the service method
		result, err := service.CreateProject(projectName)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedProject.ID, result.ID)
		assert.Equal(t, expectedProject.Name, result.Name)
	})

	t.Run("database error", func(t *testing.T) {
		projectName := "Fail Project"
		dbError := errors.New("database connection failed")

		// Set up mock expectations
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbError)

		// Call the service method
		result, err := service.CreateProject(projectName)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "repository error creating project")
	})
}

func TestProjectService_GetAllProjects_WithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_models.NewMockProjectRepository(ctrl)
	service := NewProjectService(mockRepo)

	t.Run("success", func(t *testing.T) {
		expectedProjects := []*models.Project{
			{ID: 1, Name: "Project 1"},
			{ID: 2, Name: "Project 2"},
		}

		// Set up mock expectations
		mockRepo.EXPECT().GetAll(gomock.Any()).Return(expectedProjects, nil)

		// Call the service method
		result, err := service.GetAllProjects()

		// Assertions
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Project 1", result[0].Name)
		assert.Equal(t, "Project 2", result[1].Name)
	})

	t.Run("empty result", func(t *testing.T) {
		// Set up mock expectations for empty result
		mockRepo.EXPECT().GetAll(gomock.Any()).Return([]*models.Project{}, nil)

		// Call the service method
		result, err := service.GetAllProjects()

		// Assertions
		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})
}
