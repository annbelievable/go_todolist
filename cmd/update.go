package cmd

import (
	// "database/sql"
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a todolist with new data",
	Long: `Update a todolist filtered through id.
Update the todolist's status or description.`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetUint("id")
		if id <= 0 {
			fmt.Println("id is required")
		} else {
			name, _ := cmd.Flags().GetString("name")
			description, _ := cmd.Flags().GetString("description")
			status, _ := cmd.Flags().GetString("status")
			if len(name) > 0 || len(description) > 0 || len(status) > 0 {
				updateTodolist(id, name, description, status)
			} else {
				fmt.Println("please include the field flag and data that you wish to update")
			}
		}
	},
}

func updateTodolist(id uint, task_name, task_description, status string) {
	var dataMap = make(map[string]interface{})

	if len(task_name) > 0 {
		dataMap["task_name"] = task_name
	}
	if len(task_description) > 0 {
		dataMap["task_description"] = task_description
	}
	if len(status) > 0 {
		dataMap["status"] = status
	}

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-todolist.db")
	defer sqliteDatabase.Close()

	statement := "UPDATE todo_task SET"
	dataMapLen := len(dataMap)
	makeArrLen := dataMapLen
	intCheck := 0
	if id > 0 {
		makeArrLen++
	}
	valueArray := make([]interface{}, makeArrLen)

	for field, value := range dataMap {
		valueArray[intCheck] = value
		statement = statement + " " + field + "=" + "?"
		intCheck++
		if intCheck != dataMapLen {
			statement += ","
		}
	}

	if id > 0 {
		statement += " WHERE id = ?"
		valueArray[dataMapLen] = id
	}

	statement += ";"

	result, err := sqliteDatabase.Exec(statement, valueArray...)
	if err != nil {
		fmt.Println(err.Error())
	}

	affect, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(affect, " record updated")
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().Uint("id", 0, "Give id number to update specific record")
	updateCmd.Flags().StringP("name", "n", "", "Name for the todolist")
	updateCmd.Flags().StringP("description", "d", "", "Description for the todolist")
	updateCmd.Flags().StringP("status", "s", "", "Status for the todolist")
}
