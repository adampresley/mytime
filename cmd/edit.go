package cmd

import (
	"errors"
	"fmt"

	"github.com/adampresley/mytime/api/categories"
	"github.com/adampresley/mytime/api/clients"
	"github.com/adampresley/mytime/api/projects"
	"github.com/adampresley/simdb"
	. "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func init() {
	var (
		name     string
		code     string
		rate     float64
		client   string
		category string
	)

	editCmd := &cobra.Command{
		Use:     "edit",
		Aliases: []string{"e", "ed"},
		Short:   `Edit records such as clients, categories, and sessions`,
	}

	editClientCmd := &cobra.Command{
		Use:     "client",
		Aliases: []string{"c"},
		Short:   `Edit a client record`,
		Example: `mytime edit client "test" --name "New Name" --code "New Code"`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("Please provide the code for the client you wish to edit")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err        error
				clientCode string
				client     clients.Client
			)

			if name == "" && code == "" {
				return
			}

			clientCode = args[0]

			if client, err = clientService.GetClientByCode(clientCode); err != nil {
				if errors.Is(err, simdb.ErrZeroRecords) {
					displayError(fmt.Sprintf("Client code %s not found", Green(clientCode)))
				} else {
					displayError(fmt.Sprintf("Cannot load client %s: %s", clientCode, err.Error()))
				}
			}

			if name != "" {
				client.Name = name
			}

			if code != "" {
				client.Code = code
			}

			if err = clientService.UpdateClient(client); err != nil {
				displayError(fmt.Sprintf("Problem updating client record: %s", err.Error()))
			}

			fmt.Printf("Client %s updated!\n", Green(clientCode))
		},
	}

	editCategoryCmd := &cobra.Command{
		Use:     "category",
		Aliases: []string{"cat"},
		Short:   `Edit a category record`,
		Example: `mytime edit category "test" --name "New Name" --code "New Code" --rate 10.00`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("Please provide the code for the category you wish to edit")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err          error
				categoryCode string
				category     categories.Category
			)

			if name == "" && code == "" && rate == -10.00 {
				return
			}

			categoryCode = args[0]

			if category, err = categoryService.GetCategoryByCode(categoryCode); err != nil {
				if errors.Is(err, simdb.ErrZeroRecords) {
					displayError(fmt.Sprintf("Category code %s not found", Green(categoryCode)))
				} else {
					displayError(fmt.Sprintf("Cannot load category %s: %s", categoryCode, err.Error()))
				}
			}

			if name != "" {
				category.Name = name
			}

			if code != "" {
				category.Code = code
			}

			if rate != -10.00 {
				category.Rate = rate
			}

			if err = categoryService.UpdateCategory(category); err != nil {
				displayError(fmt.Sprintf("Problem updating category record: %s", err.Error()))
			}

			fmt.Printf("Client %s updated!\n", Green(categoryCode))
		},
	}

	editProjectCmd := &cobra.Command{
		Use:     "project",
		Aliases: []string{"projects", "proj", "p"},
		Short:   `Edit a project record`,
		Example: `mytime edit project "test" --name "New Name" --code "New Code" --client "New client" --category "New default category"`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("Please provide the code for the project you wish to edit")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err         error
				projectCode string
				project     projects.Project
			)

			if name == "" && code == "" && client == "" && category == "" {
				return
			}

			projectCode = args[0]

			if project, err = projectService.GetProjectByCode(projectCode); err != nil {
				if errors.Is(err, simdb.ErrZeroRecords) {
					displayError(fmt.Sprintf("Project code %s not found", Green(projectCode)))
				} else {
					displayError(fmt.Sprintf("Cannot load project %s: %s", projectCode, err.Error()))
				}
			}

			if name != "" {
				project.Name = name
			}

			if code != "" {
				project.Code = code
			}

			if client != "" {
				var c clients.Client

				if c, err = clientService.GetClientByCode(client); err != nil {
					displayError(fmt.Sprintf("Client %s not found", client))
				}

				project.ClientID = c.ClientID
			}

			if category != "" {
				var c categories.Category

				if c, err = categoryService.GetCategoryByCode(category); err != nil {
					displayError(fmt.Sprintf("Category %s not found", category))
				}

				project.DefaultCategoryID = c.CategoryID
			}

			if err = projectService.UpdateProject(project); err != nil {
				displayError(fmt.Sprintf("Problem updating project record: %s", err.Error()))
			}

			fmt.Printf("Project %s updated!\n", Green(projectCode))
		},
	}

	editClientCmd.Flags().StringVarP(&name, "name", "n", "", "New name for a client")
	editClientCmd.Flags().StringVarP(&code, "code", "c", "", "New code for a client")
	editCategoryCmd.Flags().StringVarP(&name, "name", "n", "", "New name for a category")
	editCategoryCmd.Flags().StringVarP(&code, "code", "c", "", "New code for a category")
	editCategoryCmd.Flags().Float64VarP(&rate, "rate", "r", -10.00, "New rate for a category")
	editProjectCmd.Flags().StringVarP(&name, "name", "n", "", "New name for a project")
	editProjectCmd.Flags().StringVarP(&code, "code", "c", "", "New code for a project")
	editProjectCmd.Flags().StringVarP(&client, "client", "", "", "New client code for a project")
	editProjectCmd.Flags().StringVarP(&category, "category", "", "", "New default category for a project")

	editCmd.AddCommand(editClientCmd, editCategoryCmd, editProjectCmd)
	rootCmd.AddCommand(editCmd)
}
