// Copyright © 2019 Mohammad Shoaei <mohammad.shoaei@outlook.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"context"
	"fmt"

	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set command for bot/bots",
	Long:  `set a new command for selected bot/bots`,
	Run:   set,
}

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	setCmd.Flags().StringP("command", "c", "", "the new command to be set for the bot/bots")
}

func set(cmd *cobra.Command, args []string) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: viper.GetString("authorization"),
		})
	httpclient := oauth2.NewClient(context.Background(), src)
	client := graphql.NewClient("http://localhost:3002/api", httpclient)

	var m struct {
		Done graphql.Boolean `graphql:"setCommand(ids: $bots, command: $command)"`
	}
	bots := make([]graphql.ID, 0, len(args))
	// bots := [1]graphql.ID{}

	for _, id := range args {
		bots = append(bots, id)
		// bots[0] = graphql.ID(id)
	}

	command, _ := cmd.Flags().GetString("command")
	variables := map[string]interface{}{
		"bots":    bots,
		"command": graphql.String(command),
	}
	if err := client.Mutate(context.Background(), &m, variables); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Done: %t", m.Done)
}
