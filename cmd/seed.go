package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hector3211/night/cmd/flags"
	"github.com/hector3211/night/cmd/program"
	"github.com/hector3211/night/cmd/ui"
	"github.com/hector3211/night/pkg/parse"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var logoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#836FFF")).Bold(true)

const logo = `

 _   _  ___   ____  _   _  _____
| \ | ||_ _| / ___|| | | ||_   _|
|  \| | | | | |  _ | |_| |  | |
| |\  | | | | |_| ||  _  |  | |
|_| \_||___| \____||_| |_|  |_|


`

func init() {
	var flagDBDriver flags.DataBaseDriver
	var flagSqlFilePath flags.File
	var flagConnectionUrl flags.UrlConnection
	rootCmd.AddCommand(seedCmd)

	seedCmd.Flags().VarP(&flagDBDriver, "driver", "d", fmt.Sprintf("Database drivers to use. Allowed values: %s", strings.Join(flags.AllowedDbDrivers, ",")))
	seedCmd.Flags().VarP(&flagSqlFilePath, "path", "p", fmt.Sprintln("Path to SQL seed file example - ./seed.sql"))
	seedCmd.Flags().VarP(&flagConnectionUrl, "url", "u", fmt.Sprintln("Database connection url"))
}

type Options struct {
	DBDriver      *ui.SelectedDriver
	FilePath      *ui.SelectedFilePath
	ConnectionUrl *ui.SelectedUrl
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with data",
	Long:  "This action seeds your database of choice with the sql file provide",
	Run: func(cmd *cobra.Command, args []string) {
		flagDBDriver := flags.DataBaseDriver(cmd.Flag("driver").Value.String())
		flagSqlFilePath := flags.File(cmd.Flag("path").Value.String())
		flagConnectionUrl := flags.UrlConnection(cmd.Flag("url").Value.String())

		log.Printf("called seed with flag: %s %s %s", flagDBDriver.String(), flagSqlFilePath.String(), flagConnectionUrl.String())

		options := Options{
			DBDriver:      &ui.SelectedDriver{},
			FilePath:      &ui.SelectedFilePath{},
			ConnectionUrl: &ui.SelectedUrl{},
		}

		project := &program.Project{
			DBDriver:      flagDBDriver,
			FilePath:      flagSqlFilePath,
			ConnectionUrl: flagConnectionUrl,
		}

		fmt.Printf("%s\n", logoStyle.Render(logo))
		if project.DBDriver == "" {
			tprogram := tea.NewProgram(ui.InitialModel(options.DBDriver, project))
			if _, err := tprogram.Run(); err != nil {
				cobra.CheckErr(fmt.Sprintf("failed runing db driver model: %s", err))
			}
			project.ExitCli(tprogram)

			if options.DBDriver == nil || options.DBDriver.Choice == "" {
				cobra.CheckErr("Database driver is required")
			}

			project.DBDriver = flags.DataBaseDriver(strings.ToLower(options.DBDriver.Choice))
			err := cmd.Flag("driver").Value.Set(project.DBDriver.String())
			if err != nil {
				log.Fatalf("failed to set the driver flag value: %s", err)
			}
		} else {
			options.DBDriver.Choice = project.DBDriver.String()
		}

		if project.FilePath == "" {
			// start sql file picker UI
			tprogram := tea.NewProgram(ui.InitializeFileModel(options.FilePath, project))
			if _, err := tprogram.Run(); err != nil {
				cobra.CheckErr(fmt.Sprintf("failed runing db file model: %s", err))
			}
			project.ExitCli(tprogram)

			if options.FilePath == nil || options.FilePath.Choice == "" {
				cobra.CheckErr("File path is required")
			}
			project.FilePath = flags.File(strings.ToLower(options.FilePath.Choice))

			err := cmd.Flag("path").Value.Set(project.FilePath.String())
			if err != nil {
				log.Fatalf("failed to set the path flag value: %s", err)
			}
		} else {
			options.FilePath.Choice = project.FilePath.String()
		}

		// Can be passed via flag
		var connectionUrl string // sqlite db file
		if project.DBDriver.String() != "sqlite3" {
			if project.ConnectionUrl.String() == "" {
				tprogram := tea.NewProgram(ui.InitializeConnModel(options.ConnectionUrl, project))
				if _, err := tprogram.Run(); err != nil {
					cobra.CheckErr(ui.CreateErrorInputModel(err).Err())
				}
				project.ExitCli(tprogram)

				connValue := options.ConnectionUrl.GetValue()
				if connValue == "" {
					cobra.CheckErr("Connection URL cannot be empty for non-SQLite databases")
				}

				project.ConnectionUrl = flags.UrlConnection(connValue)
				connectionUrl = project.ConnectionUrl.String()
			} else {
				connectionUrl = project.ConnectionUrl.String()
			}
		} else {
			connectionUrl = "mydb.db"
		}

		db, err := OpenDB(
			flags.DataBaseDriver(options.DBDriver.Choice),
			options.FilePath.Choice,
			connectionUrl,
		)
		defer db.DB.Close()

		if err != nil {
			cobra.CheckErr(fmt.Sprintf("failed establishing a database connection with: %s", err))
		}

		fileData, err := os.ReadFile(db.FilePath)
		if err != nil {
			cobra.CheckErr(fmt.Sprintf("failed reading file with: %s", err))
		}

		var query string
		fileType := filepath.Ext(string(project.FilePath))[1:]

		// parse go file
		if fileType == "go" {
			parser := parse.NewParser(project.DBDriver, fileData)
			generatedSql, err := parser.Parse()
			if err != nil {
				cobra.CheckErr(fmt.Sprintf("failed generating SQL from go file: %s", err))
			}

			query = generatedSql

		}
		if fileType == "sql" {
			query = string(fileData)
		}

		// insert
		if _, err := db.DB.Exec(query); err != nil {
			cobra.CheckErr(fmt.Sprintf("failed inserting to DB with: %s", err))
		}

		// run progress bar ui
		tprogram := tea.NewProgram(ui.InitializeProgressModel(project))
		if _, err := tprogram.Run(); err != nil {
			cobra.CheckErr(fmt.Sprintf("failed runing progress model: %s", err))
		}
		project.ExitCli(tprogram)
	},
}
