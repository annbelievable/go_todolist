package cmd

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a todotask record",
	Long:  `Delete a todotask record with given id`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetUint("id")
		if id <= 0 {
			fmt.Println("id is required")
		} else {
			daleteTodolist(id)
		}
	},
}

func daleteTodolist(id uint) {
	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-todolist.db")
	defer sqliteDatabase.Close()

	statement, err := sqliteDatabase.Prepare(`DELETE FROM todo_task where id = ?`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	statement.Exec(id)
	log.Println("1 row was deleted")
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().Uint("id", 0, "Give id number to delete specific record")
}
