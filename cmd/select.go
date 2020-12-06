package cmd

import (
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"
)

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select task from todolists with either id or name",
	Long: `Select task from todolists with either id or name.
If none is provided, all tasks will be displayed`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetUint("id")
		name, _ := cmd.Flags().GetString("name")
		selectTodolist(id, name)
	},
}

func selectTodolist(id uint, name string) {
	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-todolist.db")
	defer sqliteDatabase.Close()
	selectSQL := "SELECT id, task_name, task_description, status, datecreated, datemodified FROM todo_task"

	whereSQL := " WHERE"
	conditionSQL := ""
	if id > 0 {
		conditionSQL += " id = ?"
	}
	if len(name) > 0 {
		if len(conditionSQL) > 0 {
			conditionSQL += " AND"
		}
		conditionSQL += " task_name LIKE '%" + name + "%'"
	}
	if len(conditionSQL) > 0 {
		selectSQL = selectSQL + whereSQL + conditionSQL
	}
	selectSQL += " ORDER BY datecreated ASC"

	rows, err := sqliteDatabase.Query(selectSQL, id)
	defer rows.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	var rowId, taskName, taskDescription, status, datecreated, datemodified string
	for rows.Next() {
		rows.Scan(&rowId, &taskName, &taskDescription, &status, &datecreated, &datemodified)
		fmt.Println(rowId, taskName, taskDescription, status, datecreated, datemodified)
	}
}

func init() {
	rootCmd.AddCommand(selectCmd)
	selectCmd.Flags().Uint("id", 0, "Give id number if you want to select specific record")
	selectCmd.Flags().StringP("name", "n", "", "Give name if you want to select specific record")
}
