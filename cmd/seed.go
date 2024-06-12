package cmd

import (
	"fmt"
	"log"
	"night/cmd/flags"
	"night/cmd/program"
	"night/cmd/ui"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	logoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#836FFF")).Bold(true)
)

const logo = `

 _   _  ___   ____  _   _  _____
| \ | ||_ _| / ___|| | | ||_   _|
|  \| | | | | |  _ | |_| |  | |
| |\  | | | | |_| ||  _  |  | |
|_| \_||___| \____||_| |_|  |_|


`

func init() {
	var flagDBDriver flags.DataBaseDriver
	var flagSqlFilePath flags.SqlFile
	rootCmd.AddCommand(seedCmd)

	seedCmd.Flags().VarP(&flagDBDriver, "driver", "d", fmt.Sprintf("Database drivers to use. Allowed values: %s", strings.Join(flags.AllowedDbDrivers, ",")))
	seedCmd.Flags().VarP(&flagSqlFilePath, "path", "p", fmt.Sprintf("Path to SQL seed file example - ./seed.sql"))
}

type Options struct {
	DBDriver      *ui.SelectedDriver
	SqlFilePath   *ui.SelectedFilePath
	ConnectionUrl *ui.SelectedUrl
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with data",
	Long:  "This action seeds your database of choice with the sql file provide",
	Run: func(cmd *cobra.Command, args []string) {
		flagDBDriver := flags.DataBaseDriver(cmd.Flag("driver").Value.String())
		flagSqlFilePath := flags.SqlFile(cmd.Flag("path").Value.String())

		log.Printf("called seed with flag: %s %s", flagDBDriver.String(), flagSqlFilePath.String())

		options := Options{
			DBDriver:      &ui.SelectedDriver{},
			SqlFilePath:   &ui.SelectedFilePath{},
			ConnectionUrl: &ui.SelectedUrl{},
		}

		project := &program.Project{
			DBDriver:      flagDBDriver,
			SqlFilePath:   flagSqlFilePath,
			ConnectionUrl: "",
		}

		fmt.Printf("%s\n", logoStyle.Render(logo))
		if project.DBDriver == "" {
			// spin up driver picker UI
			tprogram := tea.NewProgram(ui.InitialModel(options.DBDriver, project))
			if _, err := tprogram.Run(); err != nil {
				cobra.CheckErr(fmt.Sprintf("failed runing db driver model: %s", err))
			}
			project.ExitCli(tprogram)

			project.DBDriver = flags.DataBaseDriver(strings.ToLower(options.DBDriver.Choice))
			// err := cmd.Flag("driver").Value.Set(project.DBDriver.Type())
			// if err != nil {
			// 	log.Fatalf("failed to set the driver flag value: %s", err)
			// }

		}

		if project.SqlFilePath == "" {
			// start sql file picker UI
			tprogram := tea.NewProgram(ui.InitializeFileModel(options.SqlFilePath, project))
			if _, err := tprogram.Run(); err != nil {
				cobra.CheckErr(fmt.Sprintf("failed runing db file model: %s", err))
			}
			project.ExitCli(tprogram)

			project.SqlFilePath = flags.SqlFile(strings.ToLower(options.SqlFilePath.Choice))
			// err := cmd.Flag("path").Value.Set(string(project.SqlFilePath))
			// if err != nil {
			// 	log.Fatalf("failed to set the path flag value: %s", err)
			// }
		}

		connectionUrl := "mydb.db" // sqlite db file
		if project.DBDriver.String() != "sqlite3" {
			if project.ConnectionUrl == "" {
				tprogram := tea.NewProgram(ui.InitializeConnModel(options.ConnectionUrl, project))
				if _, err := tprogram.Run(); err != nil {
					cobra.CheckErr(ui.CreateErrorInputModel(err).Err())
				}
				project.ExitCli(tprogram)

				project.ConnectionUrl = options.ConnectionUrl.GetValue()
				connectionUrl = project.ConnectionUrl // Rewrite connectionUrl
			}
		}

		db, err := OpenDB(
			flags.DataBaseDriver(options.DBDriver.Choice),
			options.SqlFilePath.Choice,
			connectionUrl,
		)
		defer db.DB.Close()
		if err != nil {
			cobra.CheckErr(fmt.Sprintf("failed establishing a DB connection with: %s", err))
			// log.Fatalf("failed establishing a DB connection: %s", err)
			// panic("failed establishing a DB connection")
		}

		// read from file
		fileData, err := os.ReadFile(db.SqlFilePath)
		if err != nil {
			cobra.CheckErr(fmt.Sprintf("failed reading sql file with: %s", err))
		}
		query := string(fileData)

		// insert
		if _, err := db.DB.Exec(query); err != nil {
			cobra.CheckErr(fmt.Sprintf("failed inserting to DB with: %s", err))
		}

	},
}
