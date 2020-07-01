package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adampresley/mytime/api/categories"
	"github.com/adampresley/mytime/api/clients"
	"github.com/adampresley/mytime/api/helpers"
	"github.com/adampresley/mytime/api/projects"
	"github.com/adampresley/mytime/api/sessions"
	"github.com/adampresley/simdb"
	. "github.com/logrusorgru/aurora"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Version       string = "0.0.1"
	DataDirectory string = ".mytime"
)

var (
	configFile string

	rootCmd = &cobra.Command{
		Use:   "mytime",
		Short: "Time tracking, invoicing, and reporting!",
	}

	db              *simdb.Driver
	helperService   helpers.HelperService
	clientService   clients.ClientService
	categoryService categories.CategoryService
	projectService  projects.ProjectService
	sessionService  sessions.SessionService
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	var err error
	var homeDir string
	var fullPath string

	if homeDir, err = os.UserHomeDir(); err != nil {
		displayError(fmt.Sprintf("Unable to find user's home directory: %s", err.Error()))
		os.Exit(1)
	}

	/*
	 * Read config file (if any)
	 */
	fullPath = filepath.Join(homeDir, DataDirectory, "config.yml")

	viper.SetConfigType("yaml")
	viper.SetConfigFile(fullPath)

	_ = viper.ReadInConfig()

	/*
	 * Load database
	 */
	fullPath = filepath.Join(homeDir, DataDirectory)

	if db, err = simdb.New(afero.NewOsFs(), fullPath); err != nil {
		displayError(fmt.Sprintf("Unable to create database: %s", err.Error()))
		os.Exit(1)
	}

	/*
	 * Setup services
	 */
	helperService = helpers.NewHelperService(helpers.HelperServiceConfig{})

	clientService = clients.NewClientService(clients.ClientServiceConfig{
		DB:            db,
		HelperService: helperService,
	})

	categoryService = categories.NewCategoryService(categories.CategoryServiceConfig{
		DB:            db,
		HelperService: helperService,
	})

	projectService = projects.NewProjectService(projects.ProjectServiceConfig{
		ClientService: clientService,
		DB:            db,
		HelperService: helperService,
	})

	sessionService = sessions.NewSessionService(sessions.SessionServiceConfig{
		CategoryService: categoryService,
		ClientService:   clientService,
		DB:              db,
		HelperService:   helperService,
		ProjectService:  projectService,
	})
}

func displayError(msg interface{}) {
	fmt.Printf("%s %v\n", Red("ERROR:"), msg)
	os.Exit(1)
}
