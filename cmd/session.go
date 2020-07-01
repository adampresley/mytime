package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/adampresley/mytime/api/categories"
	"github.com/adampresley/mytime/api/clients"
	"github.com/adampresley/mytime/api/projects"
	"github.com/adampresley/mytime/api/sessions"
	"github.com/adampresley/simdb"
	"github.com/eiannone/keyboard"
	. "github.com/logrusorgru/aurora"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	var (
		interactive  bool
		categoryCode string
		clientCode   string
		projectCode  string
		paid         bool
		invoiced     bool
		sessionID    int
		sessionIDs   []int
	)

	sessionCmd := &cobra.Command{
		Use:     "session",
		Aliases: []string{"s", "sessions"},
		Short:   `Tools to work with sessions (timing)`,
	}

	startSessionCmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"s", "time"},
		Short:   `Starts a session timing against a project`,
		Example: `mytime session start "projectCode" "notes" - Starts timing using the default category code
mytime session start "projectCode" "notes" --category "categoryCode" - Starts timing using a specific category code
mytime session start "projectCode" "notes" --interactive - Starts timing and displays a time, waiting for you to press Q to stop`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("Please provide the project code to start timing for, and a small note describing this session")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err         error
				projectCode string
				notes       string

				project  projects.Project
				client   clients.Client
				category categories.Category

				hasActiveSession bool
				wait             *sync.WaitGroup
				activeSession    sessions.ActiveSession
				startTime        time.Time
			)

			projectCode = args[0]
			notes = args[1]

			if project, err = projectService.GetProjectByCode(projectCode); err != nil {
				if errors.Is(err, simdb.ErrZeroRecords) {
					displayError(fmt.Sprintf("Project code %s not found", Green(projectCode)))
				} else {
					displayError(fmt.Sprintf("Problem getting project: %s", err.Error()))
				}
			}

			if client, err = clientService.GetClientByID(project.ClientID); err != nil {
				if errors.Is(err, simdb.ErrZeroRecords) {
					displayError(fmt.Sprintf("Client ID %d not found", Green(project.ClientID)))
				} else {
					displayError(fmt.Sprintf("Problem getting client information: %s", err.Error()))
				}
			}

			if categoryCode != "" {
				if category, err = categoryService.GetCategoryByCode(categoryCode); err != nil {
					if errors.Is(err, simdb.ErrZeroRecords) {
						displayError(fmt.Sprintf("Category code %s not found", Green(categoryCode)))
					} else {
						displayError(fmt.Sprintf("Problem getting category: %s", err.Error()))
					}
				}
			} else {
				if category, err = categoryService.GetCategoryByID(project.DefaultCategoryID); err != nil {
					if errors.Is(err, simdb.ErrZeroRecords) {
						displayError(fmt.Sprintf("Category id %d not found", Green(project.DefaultCategoryID)))
					} else {
						displayError(fmt.Sprintf("Problem getting category: %s", err.Error()))
					}
				}
			}

			/*
			 * Don't allow the user to continue if there is an active session. That
			 * means something went wrong! Otherwise, start a session.
			 */
			if hasActiveSession, err = sessionService.HasActiveSession(); err != nil {
				displayError(fmt.Errorf("Problem determining if there is an active session in progress: %s", err.Error()))
			}

			if hasActiveSession {
				displayError("You already have an active session in progress!")
			}

			if activeSession, startTime, err = sessionService.StartActiveSession(project.ProjectID, category.CategoryID, client.ClientID, notes); err != nil {
				displayError(fmt.Errorf("Problem starting session: %s", err.Error()))
			}

			fmt.Printf("Timing for %s\nProject: %s\nCategory %s\nStart Time: %s\n", Green(client.Name), Green(project.Name), Cyan(category.Name), startTime.Format("3:04 PM"))

			if interactive {
				/*
				 * Start timing
				 */
				c, cancel := context.WithCancel(context.Background())
				wait = &sync.WaitGroup{}
				wait.Add(1)

				fmt.Printf("\nPress '%s' to stop timing.\n\n", BrightRed("q"))

				go func() {
					ticker := time.NewTicker(1 * time.Second)
					done := false

					for {
						select {
						case <-ticker.C:
							diff := time.Now().Sub(startTime)

							fmt.Printf("\rTime: %s", time.Time{}.Add(diff).Format("15:04:05"))

						case <-c.Done():
							activeSession.EndTime = time.Now()

							if !done {
								done = true
								wait.Done()
							}

							break
						}
					}
				}()

				/*
				 * Wait for the letter 'q' to be pressed to stop timing
				 */
				if err = keyboard.Open(); err != nil {
					displayError(fmt.Sprintf("Problem capturing keyboard input: %s", err.Error()))
				}

				for {
					char, _, err := keyboard.GetKey()

					if err != nil {
						displayError(fmt.Sprintf("Problem reading key from keyboard: %s", err.Error()))
						return
					}

					if char == 'q' {
						break
					}
				}

				keyboard.Close()
				cancel()
				wait.Wait()

				fmt.Printf("\n")

				/*
				 * Store the session
				 */
				session := sessions.Session{
					ClientID:      project.ClientID,
					ProjectID:     project.ProjectID,
					CategoryID:    category.CategoryID,
					StartDateTime: activeSession.StartTime,
					EndDateTime:   activeSession.EndTime,
					Notes:         notes,
					Invoiced:      false,
					Paid:          false,
				}

				if _, err = sessionService.CreateSession(session); err != nil {
					displayError(fmt.Sprintf("Problem recording session to database: %s", err.Error()))
				}

				fmt.Printf("\nSession recorded!\n")
				sessionService.DeleteActiveSessions()
			}
		},
	}

	stopSessionCmd := &cobra.Command{
		Use:     "stop",
		Aliases: []string{"st"},
		Short:   "Stops an active timing session",
		Example: `mytime session stop`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err           error
				activeSession sessions.ActiveSession
			)

			if activeSession, err = sessionService.GetActiveSession(); err != nil {
				displayError(err.Error())
			}

			activeSession.EndTime = time.Now()

			session := sessions.Session{
				ClientID:      activeSession.ClientID,
				ProjectID:     activeSession.ProjectID,
				CategoryID:    activeSession.CategoryID,
				StartDateTime: activeSession.StartTime,
				EndDateTime:   activeSession.EndTime,
				Notes:         activeSession.Notes,
				Invoiced:      false,
				Paid:          false,
			}

			if _, err = sessionService.CreateSession(session); err != nil {
				displayError(fmt.Sprintf("Problem recording session to database: %s", err.Error()))
			}

			diff := activeSession.EndTime.Sub(activeSession.StartTime)

			fmt.Printf("Start Time: %s\n", activeSession.StartTime.Format("3:04:05 PM"))
			fmt.Printf("End Time: %s\n", activeSession.EndTime.Format("3:04:05 PM"))
			fmt.Printf("Total time: %s\n", Green(time.Time{}.Add(diff).Format("15:04:05")))
			fmt.Printf("\nSession recorded!\n")
			sessionService.DeleteActiveSessions()

		},
	}

	sessionStatusCmd := &cobra.Command{
		Use:     "status",
		Short:   `Display the status of a current session`,
		Example: `mytime session status`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err      error
				client   clients.Client
				project  projects.Project
				category categories.Category

				activeSession sessions.ActiveSession
			)

			if activeSession, err = sessionService.GetActiveSession(); err != nil {
				displayError(err.Error())
			}

			if client, err = clientService.GetClientByID(activeSession.ClientID); err != nil {
				displayError(fmt.Sprintf("Problem getting client: %s", err.Error()))
			}

			if project, err = projectService.GetProjectByID(activeSession.ProjectID); err != nil {
				displayError(fmt.Sprintf("Problem getting project: %s", err.Error()))
			}

			if category, err = categoryService.GetCategoryByID(activeSession.CategoryID); err != nil {
				displayError(fmt.Sprintf("Problem getting category: %s", err.Error()))
			}

			diff := time.Now().Sub(activeSession.StartTime)

			fmt.Printf("Timing for %s\n", Green(client.Name))
			fmt.Printf("Project: %s\n", Green(project.Name))
			fmt.Printf("Category: %s\n", Cyan(category.Name))
			fmt.Printf("Start Time: %s\n", activeSession.StartTime.Format("3:04 PM"))
			fmt.Printf("Current Duration: %s\n", time.Time{}.Add(diff).Format("15:04:05"))
		},
	}

	sessionCloseCmd := &cobra.Command{
		Use:     "close",
		Aliases: []string{"c"},
		Short:   `Closes a session, marking it as invoiced and paid today (useful for time entries that really aren't invoiced)`,
		Example: `mytime session close "projectcode" 1. In this example the '1' is the ID of the session to remove.`,
		Args: func(cmd *cobra.Command, args []string) error {
			var err error

			if len(args) < 2 {
				return fmt.Errorf("Please provide the project code and session ID")
			}

			if _, err = strconv.Atoi(args[1]); err != nil {
				return fmt.Errorf("Please provide a numeric ID for the session ID")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	sessionReportCmd := &cobra.Command{
		Use:     "report",
		Aliases: []string{"r", "reports"},
		Short:   `Report on sessions`,
		Example: `mytime session report
mytime session report --category "categoryCode"
mytime session report --client "clientCode"
mytime session report --project "projectCode"
mytime session report --paid
mytime session report --invoiced
mytime session report --id 2
mytime session report --ids 2,54,3`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var result sessions.SessionCollection
			tableData := make([][]string, 9)

			search := sessions.SessionSearch{
				CategoryCode: categoryCode,
				ClientCode:   clientCode,
				Invoiced:     invoiced,
				Paid:         paid,
				ProjectCode:  projectCode,
				SessionID:    sessionID,
				SessionIDs:   sessionIDs,
			}

			if result, err = sessionService.ListSessions(search); err != nil {
				displayError(err.Error())
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Client", "Project", "Category", "Date", "Time", "Duration", "Invoiced", "Paid"})
			table.SetBorder(false)

			table.SetHeaderColor(
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold},
			)

			for _, s := range result {
				c, _ := clientService.GetClientByID(s.ClientID)
				p, _ := projectService.GetProjectByID(s.ProjectID)
				cat, _ := categoryService.GetCategoryByID(s.CategoryID)

				dateFormatted := s.StartDateTime.Format("Mon Jan _2 2006")
				startTime := s.StartDateTime.Format("3:04:05PM")
				endTime := s.EndDateTime.Format("3:04:05PM")

				diff := s.EndDateTime.Sub(s.StartDateTime)
				t := fmt.Sprintf("%s - %s", startTime, endTime)
				duration := time.Time{}.Add(diff).Format("15:04:05")

				invoiced := ""
				paid := ""

				if s.Invoiced {
					invoiced = s.InvoiceDate.Format("Mon Jan _2 2006")
				}

				if s.Paid {
					paid = s.PaidDate.Format("Mon Jan _2 2006")
				}

				tableData = append(tableData, []string{strconv.Itoa(s.SessionID), c.Name, p.Name, cat.Name, dateFormatted, t, duration, invoiced, paid})
			}

			table.AppendBulk(tableData)
			table.Render()

		},
	}

	startSessionCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Starts a timing session in interactive mode")
	startSessionCmd.Flags().StringVarP(&categoryCode, "category", "c", "", "Category to use in this timing session")
	sessionReportCmd.Flags().StringVarP(&categoryCode, "category", "a", "", "Filter sessions by category code")
	sessionReportCmd.Flags().StringVarP(&clientCode, "client", "c", "", "Filter sessions by client code")
	sessionReportCmd.Flags().StringVarP(&projectCode, "project", "p", "", "Filter sessions by project code")
	sessionReportCmd.Flags().BoolVarP(&paid, "paid", "m", false, "Filter sessions for those that are paid")
	sessionReportCmd.Flags().BoolVarP(&invoiced, "invoiced", "i", false, "Filter sessions for those that are invoiced")
	sessionReportCmd.Flags().IntVarP(&sessionID, "id", "", 0, "Filter sessions by ID")
	sessionReportCmd.Flags().IntSliceVarP(&sessionIDs, "ids", "", []int{}, "Filter sessions by a list of IDs")

	sessionCmd.AddCommand(startSessionCmd, stopSessionCmd, sessionStatusCmd, sessionCloseCmd, sessionReportCmd)
	rootCmd.AddCommand(sessionCmd)
}
