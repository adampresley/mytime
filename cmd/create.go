package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/adampresley/mytime/api/categories"
	"github.com/adampresley/mytime/api/clients"
	"github.com/adampresley/mytime/api/projects"
	"github.com/adampresley/simdb"
	. "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func init() {
	createCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   `Creates a new client, project, or category`,
	}

	createClientCmd := &cobra.Command{
		Use:     "client",
		Aliases: []string{"c"},
		Short:   `Creates a new client`,
		Example: `mytime create client "Client A" "clientcode"`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("Please provide a name and code for this client!")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err        error
				clientName string
				clientCode string
			)

			clientName = args[0]
			clientCode = args[1]

			client := clients.Client{
				Name:     clientName,
				Code:     clientCode,
				Archived: false,
			}

			if _, err = clientService.CreateClient(client); err != nil {
				displayError(fmt.Sprintf("Error creating client: %s", err.Error()))
			}

			fmt.Printf("New client %s created!\n", Green(clientName))
		},
	}

	createCategoryCmd := &cobra.Command{
		Use:     "category",
		Aliases: []string{"cat"},
		Short:   `Creates a new category. `,
		Example: `mytime create category "Development" "dev" 50.00`,
		Args: func(cmd *cobra.Command, args []string) error {
			var err error

			if len(args) < 3 {
				return fmt.Errorf("Please provide the name, code, and rate for your new category!")
			}

			if _, err = strconv.ParseFloat(args[2], 64); err != nil {
				return fmt.Errorf("Invalid rate. Must be a decimal number!")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err          error
				categoryName string
				code         string
				rate         float64
			)

			categoryName = args[0]
			code = args[1]
			rate, _ = strconv.ParseFloat(args[2], 64)

			category := categories.Category{
				Name:     categoryName,
				Code:     code,
				Rate:     rate,
				Archived: false,
			}

			if _, err = categoryService.CreateCategory(category); err != nil {
				displayError(fmt.Sprintf("Error creating category: %s", err.Error()))
			}

			fmt.Printf("New category '%s' created!\n", Green(categoryName))
		},
	}

	createProjectCmd := &cobra.Command{
		Use:     "project",
		Aliases: []string{"p", "proj"},
		Short:   `Create a new project.`,
		Long:    `Creates a new project. Projects are tied to clients, and are what time is tracked against.`,
		Example: `mytime create project "Name" "code" "clientCode" "defaultCategoryCode"`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 4 {
				return fmt.Errorf("Please provide a name, code, client code, and default category code for your new project!")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err                 error
				name                string
				code                string
				clientCode          string
				defaultCategoryCode string

				client   clients.Client
				category categories.Category

				newProjectID int
			)

			name = args[0]
			code = args[1]
			clientCode = args[2]
			defaultCategoryCode = args[3]

			if client, err = clientService.GetClientByCode(clientCode); err != nil {
				if errors.Is(err, simdb.ErrZeroRecords) {
					displayError(fmt.Sprintf("Client code %s not found", Green(clientCode)))
				} else {
					displayError(err.Error())
				}
			}

			if category, err = categoryService.GetCategoryByCode(defaultCategoryCode); err != nil {
				if errors.Is(err, simdb.ErrZeroRecords) {
					displayError(fmt.Sprintf("Category code %s not found", Green(defaultCategoryCode)))
				} else {
					displayError(err.Error())
				}
			}

			newProject := projects.Project{
				Name:              name,
				Code:              code,
				ClientID:          client.ClientID,
				DefaultCategoryID: category.CategoryID,
				Archived:          false,
			}

			if newProjectID, err = projectService.CreateProject(newProject); err != nil {
				displayError(fmt.Sprintf("Problem creating project: %s", err.Error()))
			}

			fmt.Printf("New project created!\n\nID: %d\nName: %s\nCode: %s\n", newProjectID, name, Green(code))
		},
	}

	createCmd.AddCommand(createClientCmd, createCategoryCmd, createProjectCmd)
	rootCmd.AddCommand(createCmd)
}
