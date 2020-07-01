package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/adampresley/mytime/api/categories"
	"github.com/adampresley/mytime/api/clients"
	"github.com/adampresley/mytime/api/projects"

	// . "github.com/logrusorgru/aurora"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	var archived bool
	var name string
	var client string

	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   `Lists clients, projects, or categories`,
	}

	listClientsCmd := &cobra.Command{
		Use:     "clients",
		Aliases: []string{"c", "client"},
		Short:   `Lists all clients`,
		Example: `mt list clients
mt list clients --archived
mt list clients --name "client name"`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var result clients.ClientCollection

			tableData := make([][]string, 5)

			search := clients.ClientSearch{
				Archived: archived,
				Name:     name,
			}

			if result, err = clientService.ListClients(search); err != nil {
				displayError(fmt.Sprintf("Error listing clients: %s", err.Error()))
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Client", "Code"})
			table.SetBorder(false)

			table.SetHeaderColor(
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
				tablewriter.Colors{tablewriter.Bold, tablewriter.FgRedColor},
			)

			for _, c := range result {
				tableData = append(tableData, []string{strconv.Itoa(c.ClientID), c.Name, c.Code})
			}

			table.AppendBulk(tableData)
			table.Render()
		},
	}

	listCategoriesCmd := &cobra.Command{
		Use:     "categories",
		Aliases: []string{"cat", "category"},
		Short:   `Lists all categories.`,
		Example: `mt list categories
mt list categories --archived
mt list categories --name "category name"`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var result categories.CategoryCollection

			tableData := make([][]string, 5)

			search := categories.CategorySearch{
				Archived: archived,
				Name:     name,
			}

			if result, err = categoryService.ListCategories(search); err != nil {
				displayError(fmt.Sprintf("Error listing categories: %s", err.Error()))
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Category", "Code", "Rate"})
			table.SetBorder(false)

			table.SetHeaderColor(
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
				tablewriter.Colors{tablewriter.Bold, tablewriter.FgRedColor},
				tablewriter.Colors{tablewriter.Bold},
			)

			for _, c := range result {
				tableData = append(tableData, []string{strconv.Itoa(c.CategoryID), c.Name, c.Code, fmt.Sprintf("%.2f", c.Rate)})
			}

			table.AppendBulk(tableData)
			table.Render()
		},
	}

	listProjectsCmd := &cobra.Command{
		Use:     "projects",
		Aliases: []string{"p", "project", "proj"},
		Short:   `Lists projects.`,
		Example: `mt list projects 
mt list projects --archived
mt list projects --name "value"
mt list projects --client "client"`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var result projects.ProjectCollection

			tableData := make([][]string, 5)

			search := projects.ProjectSearch{
				Archived: archived,
				Name:     name,
				Client:   client,
			}

			if result, err = projectService.ListProjects(search); err != nil {
				displayError(fmt.Sprintf("Problem listing projects: %s", err.Error()))
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Project", "Code", "Client", "Default Category"})
			table.SetBorder(false)

			table.SetHeaderColor(
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
				tablewriter.Colors{tablewriter.Bold, tablewriter.FgRedColor},
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold},
			)

			for _, p := range result {
				c, _ := clientService.GetClientByID(p.ClientID)
				cat, _ := categoryService.GetCategoryByID(p.DefaultCategoryID)

				tableData = append(tableData, []string{strconv.Itoa(p.ProjectID), p.Name, p.Code, c.Name, cat.Name})
			}

			table.AppendBulk(tableData)
			table.Render()
		},
	}

	listClientsCmd.Flags().BoolVarP(&archived, "archived", "a", false, `Show archived clients. E.g. mt list clients --archived`)
	listClientsCmd.Flags().StringVarP(&name, "name", "n", "", `Filter clients by name or code. E.g. mt list clients --name "client"`)
	listCategoriesCmd.Flags().BoolVarP(&archived, "archived", "a", false, `Show archived categories. E.g. mt list categories --archived`)
	listCategoriesCmd.Flags().StringVarP(&name, "name", "n", "", `Filter categories by name or code. E.g. mt list categories --name "test"`)
	listProjectsCmd.Flags().BoolVarP(&archived, "archived", "a", false, `Show archived projects. E.g. mt list projects --archived`)
	listProjectsCmd.Flags().StringVarP(&name, "name", "n", "", `Filter projects by name or code. E.g. mt list projects --name "test"`)
	listProjectsCmd.Flags().StringVarP(&client, "client", "c", "", `Filter projects by client name or code. E.g. mt list projects --client "client"`)

	listCmd.AddCommand(listClientsCmd, listCategoriesCmd, listProjectsCmd)
	rootCmd.AddCommand(listCmd)
}
