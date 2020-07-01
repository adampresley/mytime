package categories

import (
	"fmt"
	"strings"

	"github.com/adampresley/mytime/api/helpers"
	"github.com/adampresley/simdb"
)

type CategoryServicer interface {
	CreateCategory(category Category) (int, error)
	ListCategories(search CategorySearch) (CategoryCollection, error)
	GetCategoryByCode(code string) (Category, error)
	GetCategoryByID(id int) (Category, error)
	UpdateCategory(category Category) error
}

type CategoryServiceConfig struct {
	DB            *simdb.Driver
	HelperService helpers.HelperServicer
}

type CategoryService struct {
	DB            *simdb.Driver
	HelperService helpers.HelperServicer
}

func NewCategoryService(config CategoryServiceConfig) CategoryService {
	return CategoryService{
		DB:            config.DB,
		HelperService: config.HelperService,
	}
}

func (s CategoryService) CreateCategory(category Category) (int, error) {
	var err error

	category.CategoryID = s.DB.Open(Category{}).GetNextNumericID()

	if category.Code == "" {
		category.Code = s.HelperService.CreateAutoCode(category.Name)
	}

	if err = s.DB.Insert(category); err != nil {
		return 0, fmt.Errorf("error inserting new category: %w", err)
	}

	return category.CategoryID, nil
}

func (s CategoryService) ListCategories(search CategorySearch) (CategoryCollection, error) {
	var err error

	result := make(CategoryCollection, 0, 10)
	d := s.DB.Open(Category{})

	if search.Archived {
		d = d.Where("archived", "=", true)
	} else {
		d = d.Where("archived", "=", false)
	}

	if err = d.Get().AsEntity(&result); err != nil {
		return result, fmt.Errorf("Error querying for Categories: %w", err)
	}

	if search.Name != "" {
		result = s.filterCategoriesByName(result, search.Name)
	}

	return result, nil
}

func (s CategoryService) GetCategoryByCode(code string) (Category, error) {
	var err error
	var category Category

	if err = s.DB.Open(Category{}).Where("code", "=", code).First().AsEntity(&category); err != nil {
		return category, fmt.Errorf("Error finding category with a code '%s': %w", code, err)
	}

	return category, nil
}

func (s CategoryService) GetCategoryByID(id int) (Category, error) {
	var err error
	var category Category

	if err = s.DB.Open(Category{}).Where("categoryID", "=", id).First().AsEntity(&category); err != nil {
		return category, fmt.Errorf("Error finding category with an ID '%s': %w", id, err)
	}

	return category, nil
}

func (s CategoryService) UpdateCategory(category Category) error {
	return s.DB.Update(category)
}

func (s CategoryService) filterCategoriesByName(categories CategoryCollection, name string) CategoryCollection {
	result := make(CategoryCollection, 0, 20)

	for _, c := range categories {
		lowerName := strings.ToLower(name)

		if strings.Contains(strings.ToLower(c.Name), lowerName) || strings.Contains(strings.ToLower(c.Code), lowerName) {
			result = append(result, c)
		}
	}

	return result
}
