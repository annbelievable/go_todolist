package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	_ "github.com/mattn/go-sqlite3"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todolist",
	Short: "A todo list manager",
	Long: `A todo list manager that has the following abilities:
1) Add a task
2) View a list of tasks
3) Update the status of a task
4) Remove a task`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { fmt.Println("hello todolist user") },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//i will have to init the db/json/csv here
func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//this is an example that tells you anything that needs to be configured on initialization can be done here
	//im curious how they use this configuration file though
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todolist.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	//initialize the database
	_, err := os.Stat("sqlite-todolist.db")

	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create("sqlite-todolist.db")
			if err != nil {
				fmt.Println(err.Error())
			}
			file.Close()
		} else {
			fmt.Println(err)
		}
	}

	sqliteDatabase, _ := sql.Open("sqlite3", "sqlite-todolist.db")
	statement := `SELECT name FROM sqlite_master WHERE type='table' AND name='todo_task';`
	row := sqliteDatabase.QueryRow(statement)

	var name string
	row.Scan(&name)

	if name != "todo_task" {
		statement := `CREATE TABLE IF NOT EXISTS todo_task(
			"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
			"task_name" TEXT NOT NULL,
			"task_description" TEXT,
			"status" TEXT NOT NULL,
			"datecreated" TEXT,
			"datemodified" TEXT
		);`

		query, err := sqliteDatabase.Prepare(statement)
		if err != nil {
			fmt.Println(err.Error())
		}

		query.Exec()
	}

	sqliteDatabase.Close()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".todolist" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".todolist")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
