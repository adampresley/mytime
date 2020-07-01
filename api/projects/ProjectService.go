package projects

import (
	"fmt"
	"strings"

	"github.com/adampresley/mytime/api/clients"
	"github.com/adampresley/mytime/api/helpers"
	"github.com/adampresley/simdb"
)

type ProjectServicer interface {
	ArchiveProject(code string) error
	CreateProject(project Project) (int, error)
	GetProjectByCode(code string) (Project, error)
	GetProjectByID(id int) (Project, error)
	ListProjects(search ProjectSearch) (ProjectCollection, error)
	UpdateProject(project Project) error
}

type ProjectServiceConfig struct {
	ClientService clients.ClientServicer
	DB            *simdb.Driver
	HelperService helpers.HelperServicer
}

type ProjectService struct {
	ClientService clients.ClientServicer
	DB            *simdb.Driver
	HelperService helpers.HelperServicer
}

func NewProjectService(config ProjectServiceConfig) ProjectService {
	return ProjectService{
		ClientService: config.ClientService,
		DB:            config.DB,
		HelperService: config.HelperService,
	}
}

func (s ProjectService) ArchiveProject(code string) error {
	var err error
	var project Project

	if project, err = s.GetProjectByCode(code); err != nil {
		return err
	}

	project.Archived = true

	if err = s.DB.Open(Project{}).Update(project); err != nil {
		return fmt.Errorf("Error updating s.DB: %w", err)
	}

	return nil
}

func (s ProjectService) CreateProject(project Project) (int, error) {
	var err error

	project.ProjectID = s.DB.Open(Project{}).GetNextNumericID()

	if err = s.DB.Insert(project); err != nil {
		return 0, fmt.Errorf("Error creating new project: %w", err)
	}

	return project.ProjectID, nil
}

func (s ProjectService) GetProjectByCode(code string) (Project, error) {
	var err error
	var project Project

	if err = s.DB.Open(Project{}).Where("code", "=", code).First().AsEntity(&project); err != nil {
		return project, fmt.Errorf("Error querying for project: %w", err)
	}

	return project, nil
}

func (s ProjectService) GetProjectByID(id int) (Project, error) {
	var err error
	var project Project

	if err = s.DB.Open(Project{}).Where("projectID", "=", id).First().AsEntity(&project); err != nil {
		return project, fmt.Errorf("Error querying for project: %w", err)
	}

	return project, nil
}

func (s ProjectService) ListProjects(search ProjectSearch) (ProjectCollection, error) {
	var err error
	var d *simdb.Driver

	result := make(ProjectCollection, 0, 10)

	d = s.DB.Open(Project{})

	if search.Archived {
		d = d.Where("archived", "=", true)
	} else {
		d = d.Where("archived", "=", false)
	}

	if search.Name != "" {
		d = d.Where("name", "contains", search.Name)
	}

	if err = d.Get().AsEntity(&result); err != nil {
		return result, fmt.Errorf("Error querying for Projects: %w", err)
	}

	if search.Client != "" {
		result = s.filterProjectsByClient(result, search.Client)
	}

	return result, nil
}

func (s ProjectService) UpdateProject(project Project) error {
	return s.DB.Update(project)
}

func (s ProjectService) filterProjectsByClient(projects ProjectCollection, client string) ProjectCollection {
	result := make(ProjectCollection, 0, 50)

	for _, p := range projects {
		c, _ := s.ClientService.GetClientByID(p.ClientID)
		lowerClient := strings.ToLower(client)

		if strings.Contains(strings.ToLower(c.Name), lowerClient) || strings.Contains(strings.ToLower(c.Code), lowerClient) {
			result = append(result, p)
		}
	}

	return result
}
