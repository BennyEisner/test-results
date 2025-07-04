package application

import (
	"context"
	"errors"
	"testing"

	"github.com/BennyEisner/test-results/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectRepository is a mock implementation of domain.ProjectRepository
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id int64) (*domain.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectRepository) GetAll(ctx context.Context) ([]*domain.Project, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Project), args.Error(1)
}

func (m *MockProjectRepository) GetByName(ctx context.Context, name string) (*domain.Project, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectRepository) Create(ctx context.Context, p *domain.Project) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProjectRepository) Update(ctx context.Context, id int64, name string) (*domain.Project, error) {
	args := m.Called(ctx, id, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectRepository) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func TestProjectService_GetProjectByID(t *testing.T) {
	tests := []struct {
		name          string
		id            int64
		setupMock     func(*MockProjectRepository)
		expectedError error
		expectedID    int64
	}{
		{
			name: "success",
			id:   1,
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Project{ID: 1, Name: "Test Project"}, nil)
			},
			expectedError: nil,
			expectedID:    1,
		},
		{
			name: "invalid id",
			id:   0,
			setupMock: func(repo *MockProjectRepository) {
				// No mock setup needed for invalid input
			},
			expectedError: domain.ErrInvalidInput,
			expectedID:    0,
		},
		{
			name: "not found",
			id:   999,
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, int64(999)).Return(nil, nil)
			},
			expectedError: domain.ErrProjectNotFound,
			expectedID:    0,
		},
		{
			name: "database error",
			id:   1,
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, int64(1)).Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("failed to get project by ID 1: database error"),
			expectedID:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProjectRepository)
			tt.setupMock(mockRepo)

			service := NewProjectService(mockRepo)
			ctx := context.Background()

			project, err := service.GetProjectByID(ctx, tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, project.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_GetAllProjects(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*MockProjectRepository)
		expectedError error
		expectedCount int
	}{
		{
			name: "success",
			setupMock: func(repo *MockProjectRepository) {
				projects := []*domain.Project{
					{ID: 1, Name: "Project 1"},
					{ID: 2, Name: "Project 2"},
				}
				repo.On("GetAll", mock.Anything).Return(projects, nil)
			},
			expectedError: nil,
			expectedCount: 2,
		},
		{
			name: "empty list",
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetAll", mock.Anything).Return([]*domain.Project{}, nil)
			},
			expectedError: nil,
			expectedCount: 0,
		},
		{
			name: "database error",
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetAll", mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("failed to get all projects: database error"),
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProjectRepository)
			tt.setupMock(mockRepo)

			service := NewProjectService(mockRepo)
			ctx := context.Background()

			projects, err := service.GetAllProjects(ctx)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, projects, tt.expectedCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_CreateProject(t *testing.T) {
	tests := []struct {
		name          string
		projectName   string
		setupMock     func(*MockProjectRepository)
		expectedError error
		expectedID    int64
	}{
		{
			name:        "success",
			projectName: "New Project",
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByName", mock.Anything, "New Project").Return(nil, nil)
				repo.On("Create", mock.Anything, mock.MatchedBy(func(p *domain.Project) bool {
					return p.Name == "New Project"
				})).Run(func(args mock.Arguments) {
					p := args.Get(1).(*domain.Project)
					p.ID = 1
				}).Return(nil)
			},
			expectedError: nil,
			expectedID:    1,
		},
		{
			name:        "empty name",
			projectName: "",
			setupMock: func(repo *MockProjectRepository) {
				// No mock setup needed for invalid input
			},
			expectedError: domain.ErrInvalidInput,
			expectedID:    0,
		},
		{
			name:        "duplicate project",
			projectName: "Existing Project",
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByName", mock.Anything, "Existing Project").Return(&domain.Project{ID: 1, Name: "Existing Project"}, nil)
			},
			expectedError: domain.ErrDuplicateProject,
			expectedID:    0,
		},
		{
			name:        "database error on create",
			projectName: "New Project",
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByName", mock.Anything, "New Project").Return(nil, nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedError: errors.New("failed to create project: database error"),
			expectedID:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProjectRepository)
			tt.setupMock(mockRepo)

			service := NewProjectService(mockRepo)
			ctx := context.Background()

			project, err := service.CreateProject(ctx, tt.projectName)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, project.ID)
				assert.Equal(t, tt.projectName, project.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_UpdateProject(t *testing.T) {
	tests := []struct {
		name          string
		id            int64
		newName       string
		setupMock     func(*MockProjectRepository)
		expectedError error
		expectedID    int64
	}{
		{
			name:    "success",
			id:      1,
			newName: "Updated Project",
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Project{ID: 1, Name: "Old Name"}, nil)
				repo.On("GetByName", mock.Anything, "Updated Project").Return(nil, nil)
				repo.On("Update", mock.Anything, int64(1), "Updated Project").Return(&domain.Project{ID: 1, Name: "Updated Project"}, nil)
			},
			expectedError: nil,
			expectedID:    1,
		},
		{
			name:    "invalid id",
			id:      0,
			newName: "Updated Project",
			setupMock: func(repo *MockProjectRepository) {
				// No mock setup needed for invalid input
			},
			expectedError: domain.ErrInvalidInput,
			expectedID:    0,
		},
		{
			name:    "empty name",
			id:      1,
			newName: "",
			setupMock: func(repo *MockProjectRepository) {
				// No mock setup needed for invalid input
			},
			expectedError: domain.ErrInvalidInput,
			expectedID:    0,
		},
		{
			name:    "project not found",
			id:      999,
			newName: "Updated Project",
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))
			},
			expectedError: errors.New("failed to check project existence: not found"),
			expectedID:    0,
		},
		{
			name:    "duplicate name",
			id:      1,
			newName: "Existing Project",
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Project{ID: 1, Name: "Old Name"}, nil)
				repo.On("GetByName", mock.Anything, "Existing Project").Return(&domain.Project{ID: 2, Name: "Existing Project"}, nil)
			},
			expectedError: domain.ErrDuplicateProject,
			expectedID:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProjectRepository)
			tt.setupMock(mockRepo)

			service := NewProjectService(mockRepo)
			ctx := context.Background()

			project, err := service.UpdateProject(ctx, tt.id, tt.newName)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, project.ID)
				assert.Equal(t, tt.newName, project.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_DeleteProject(t *testing.T) {
	tests := []struct {
		name          string
		id            int64
		setupMock     func(*MockProjectRepository)
		expectedError error
	}{
		{
			name: "success",
			id:   1,
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Project{ID: 1, Name: "Test Project"}, nil)
				repo.On("Delete", mock.Anything, int64(1)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "invalid id",
			id:   0,
			setupMock: func(repo *MockProjectRepository) {
				// No mock setup needed for invalid input
			},
			expectedError: domain.ErrInvalidInput,
		},
		{
			name: "project not found",
			id:   999,
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, int64(999)).Return(nil, nil)
			},
			expectedError: domain.ErrProjectNotFound,
		},
		{
			name: "database error on delete",
			id:   1,
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Project{ID: 1, Name: "Test Project"}, nil)
				repo.On("Delete", mock.Anything, int64(1)).Return(errors.New("database error"))
			},
			expectedError: errors.New("failed to delete project: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProjectRepository)
			tt.setupMock(mockRepo)

			service := NewProjectService(mockRepo)
			ctx := context.Background()

			err := service.DeleteProject(ctx, tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_GetProjectByName(t *testing.T) {
	tests := []struct {
		name          string
		projectName   string
		setupMock     func(*MockProjectRepository)
		expectedError error
		expectedID    int64
	}{
		{
			name:        "success",
			projectName: "Test Project",
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByName", mock.Anything, "Test Project").Return(&domain.Project{ID: 1, Name: "Test Project"}, nil)
			},
			expectedError: nil,
			expectedID:    1,
		},
		{
			name:        "empty name",
			projectName: "",
			setupMock: func(repo *MockProjectRepository) {
				// No mock setup needed for invalid input
			},
			expectedError: domain.ErrInvalidInput,
			expectedID:    0,
		},
		{
			name:        "project not found",
			projectName: "Non-existent Project",
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByName", mock.Anything, "Non-existent Project").Return(nil, nil)
			},
			expectedError: domain.ErrProjectNotFound,
			expectedID:    0,
		},
		{
			name:        "database error",
			projectName: "Test Project",
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByName", mock.Anything, "Test Project").Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("failed to get project by name: database error"),
			expectedID:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProjectRepository)
			tt.setupMock(mockRepo)

			service := NewProjectService(mockRepo)
			ctx := context.Background()

			project, err := service.GetProjectByName(ctx, tt.projectName)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, project.ID)
				assert.Equal(t, tt.projectName, project.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
