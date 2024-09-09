package cmd

import (
	"fmt"
	"github.com/hector3211/night/cmd/flags"
	"github.com/hector3211/night/cmd/program"
	"github.com/hector3211/night/cmd/ui"
	"github.com/hector3211/night/pkg/parse"
	"log"
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
	var flagSqlFilePath flags.File
	// var flagSeedLanguage flags.SeedLanguage
	rootCmd.AddCommand(seedCmd)

	seedCmd.Flags().VarP(&flagDBDriver, "driver", "d", fmt.Sprintf("Database drivers to use. Allowed values: %s", strings.Join(flags.AllowedDbDrivers, ",")))
	// seedCmd.Flags().VarP(&flagSeedLanguage, "type", "t", fmt.Sprintf("Language type used for seeding. Allowed values %s", strings.Join(flags.AllowedFileTypes, ",")))
	seedCmd.Flags().VarP(&flagSqlFilePath, "path", "p", fmt.Sprintf("Path to SQL seed file example - ./seed.sql"))
}

type Options struct {
	DBDriver      *ui.SelectedDriver
	SeedLanguage  *ui.SelectedLanguage
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
		// flagSeedLanguage := flags.SeedLanguage(cmd.Flag("type").Value.String())

		log.Printf("called seed with flag: %s %s", flagDBDriver.String(), flagSqlFilePath.String())

		options := Options{
			DBDriver:      &ui.SelectedDriver{},
			FilePath:      &ui.SelectedFilePath{},
			ConnectionUrl: &ui.SelectedUrl{},
		}

		project := &program.Project{
			DBDriver: flagDBDriver,
			FilePath: flagSqlFilePath,
			// SeedLanguage:  flagSeedLanguage,
			ConnectionUrl: "",
		}

		fmt.Printf("%s\n", logoStyle.Render(logo))
		if project.DBDriver == "" {
			tprogram := tea.NewProgram(ui.InitialModel(options.DBDriver, project))
			if _, err := tprogram.Run(); err != nil {
				cobra.CheckErr(fmt.Sprintf("failed runing db driver model: %s", err))
			}
			project.ExitCli(tprogram)

			project.DBDriver = flags.DataBaseDriver(strings.ToLower(options.DBDriver.Choice))
			err := cmd.Flag("driver").Value.Set(project.DBDriver.String())
			if err != nil {
				log.Fatalf("failed to set the driver flag value: %s", err)
			}

		}

		// if project.SeedLanguage == flags.UNKNWON {
		// 	tprogram := tea.NewProgram(ui.InitializeChoiceModel(options.SeedLanguage, project))
		// 	if _, err := tprogram.Run(); err != nil {
		// 		cobra.CheckErr(fmt.Sprintf("failed runing seed language model: %s", err))
		// 	}
		// 	project.ExitCli(tprogram)
		//
		// 	if options.SeedLanguage.Choice != "" {
		// 		switch options.SeedLanguage.GetValue() {
		// 		case "Go":
		// 			project.SeedLanguage = flags.SeedLanguage(strings.ToLower(options.SeedLanguage.Choice))
		// 		case "SQL":
		// 			project.SeedLanguage = flags.SeedLanguage(strings.ToLower(options.SeedLanguage.Choice))
		// 		}
		// 	}
		// 	err := cmd.Flag("type").Value.Set(string(project.SeedLanguage))
		// 	if err != nil {
		// 		log.Fatalf("failed to set the seed language type flag value: %s", err)
		// 	}
		// }

		if project.FilePath == "" {
			// start sql file picker UI
			tprogram := tea.NewProgram(ui.InitializeFileModel(options.FilePath, project))
			if _, err := tprogram.Run(); err != nil {
				cobra.CheckErr(fmt.Sprintf("failed runing db file model: %s", err))
			}
			project.ExitCli(tprogram)

			project.FilePath = flags.File(strings.ToLower(options.FilePath.Choice))

			err := cmd.Flag("path").Value.Set(project.FilePath.String())
			if err != nil {
				log.Fatalf("failed to set the path flag value: %s", err)
			}
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
		fileType := strings.Split(project.FilePath.String(), ".")[1]

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
