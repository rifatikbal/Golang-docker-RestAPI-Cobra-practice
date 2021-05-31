
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mytest",
	Short: "This is just a server ",
	Long: `This is just a Rest API server`,

	Run: func(cmd *cobra.Command, args []string) { fmt.Println("this is inside rootcmd")},
}


func Execute() {

	fmt.Println("root is executing ")
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	fmt.Println("cobra root initialize is working")
}


