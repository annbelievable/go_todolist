package cmd

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task into todolists",
	Long:  `Add a task in the order of: name, description and status`,
	Run: func(cmd *cobra.Command, args []string) {
		//if a command has any flags, get them here like the example below
		//the flag can be processed here or passed to the next function
		//you can run mutilple funtions
		// fstatus, _ := cmd.Flags().GetBool("float")

		//use this to get the full map/struct of the args
		// cmdflags := cmd.Flags()
		// fmt.Println(cmdflags)
		// cmdArgs := cmd.Flags().Args()
		// fmt.Println(cmdArgs)

		//using getstring
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")

		message := ""
		if len(name) == 0 {
			message += "name is required;"
		}
		if len(description) == 0 {
			message += "description is required;"
		}
		if len(message) > 0 {
			fmt.Println(message)
		} else {
			addTodoList(name, description, status)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	//make the name flag required
	addCmd.Flags().StringP("name", "n", "", "Name for the todolist")
	addCmd.Flags().StringP("description", "d", "", "Description for the todolist")
	addCmd.Flags().StringP("status", "s", "", "Status for the todolist")
}

//using this way i do not need the flags anymore
//i still need flag in future
//need to check for empty flag
func addTodoList(name, description, status string) {
	//this is the way when using args to check
	/*
		if len(args) == 0 {
			fmt.Println("arguments are needed for this subcommand")
			return
		}

		name := args[0]
		description := args[1]
		status := args[2]
	*/

	//open the db
	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-todolist.db")
	defer sqliteDatabase.Close()

	insertSQL := `INSERT INTO todo_task( task_name, task_description, status, datecreated, datemodified ) VALUES ( ?, ?, ?, ?, ? )`
	statement, err := sqliteDatabase.Prepare(insertSQL)
	if err != nil {
		fmt.Println(err.Error())
	}
	timeNow := time.Now()
	results, err := statement.Exec(name, description, status, timeNow, timeNow)
	if err != nil {
		fmt.Println(err.Error())
	}

	//get the inserted record id
	id, err := results.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("new record id:", id)
}
