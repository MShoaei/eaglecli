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
	"strconv"

	"github.com/bndr/gotabulate"
	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the information of bot/bots",
	Long:  `list the information of all bots or a specific bot.`,
	Run:   list,
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	flags := listCmd.Flags()
	flags.SortFlags = false
	// flags.String()
	flags.BoolP("all", "a", false, "show all the details of the bot/bots")
	flags.StringArray("fields", []string{}, "The fields that are to be shown in the desired order. accepted values are: id, ip, whoami, os, install-date, admin, av, cpu,gpu, version, last-checkin, last-command, new-command")

	/* flags.Bool("id", false, "show UUID of the bot/bots")
	flags.Bool("ip", false, "show IP of the bot/bots")
	flags.BoolP("whoami", "w", false, "show slaves username")
	flags.Bool("os", false, "show Operating System of the bot/bots")
	flags.Bool("install-date", false, "show Installation Date of the bot/bots")
	flags.Bool("admin", false, "show if the bot has admin permission")
	flags.Bool("av", false, "show Anti-Virus of the bot/bots")
	flags.BoolP("cpu", "c", false, "show CPU of the bot/bots")
	flags.BoolP("gpu", "g", false, "show GPU of the bot/bots")
	flags.BoolP("version", "v", false, "show Version of the bot/bots")
	flags.BoolP("last-checkin", "", false, "show Last Checkin of the bot/bots")
	flags.BoolP("last-command", "l", false, "show the Last completed Command of the bot/bots")
	flags.BoolP("new-command", "n", false, "show to be completed Command of the bot/bots") */
}

func list(cmd *cobra.Command, args []string) {

	if showAll, _ := cmd.Flags().GetBool("all"); showAll {
		src := oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: viper.GetString("authorization"),
			})
		httpclient := oauth2.NewClient(context.Background(), src)
		client := graphql.NewClient("http://localhost:9990/api", httpclient)

		var q struct {
			Bots []struct {
				ID          graphql.String
				IP          graphql.String
				WhoAmI      graphql.String
				OS          graphql.String
				InstallDate graphql.String
				Admin       graphql.Boolean
				AV          graphql.String
				CPU         graphql.String
				GPU         graphql.String
				Version     graphql.String
				LastCheckin graphql.String
				LastCommand graphql.String
				NewCommand  graphql.String
			}
		}
		if err := client.Query(context.Background(), &q, nil); err != nil {
			fmt.Println(err)
			return
		}

		allBots := make([][]string, len(q.Bots))
		for i, item := range q.Bots {

			allBots[i] = []string{
				strconv.Itoa(i + 1),
				string(item.ID),
				string(item.IP),
				string(item.WhoAmI),
				string(item.OS),
				string(item.InstallDate),
				strconv.FormatBool(bool(item.Admin)),
				string(item.AV),
				string(item.CPU),
				string(item.GPU),
				string(item.Version),
				string(item.LastCheckin),
				string(item.LastCommand),
				string(item.NewCommand),
			}
		}

		t := gotabulate.Create(allBots)
		t.SetWrapStrings(true)
		t.SetMaxCellSize(14)
		t.SetAlign("center")
		t.SetHeaders([]string{
			"index",
			"uuid",
			"ip",
			"whoami",
			"os",
			"install date",
			"admin",
			"anti-virus",
			"cpu",
			"gpu",
			"version",
			"last checkin",
			"last command",
			"new command",
		})
		fmt.Println(t.Render("grid"))
	} else {
		fmt.Println("can't do this, yet.")
	}
}
