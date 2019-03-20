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

	"github.com/spf13/viper"

	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
)

var username string
var password string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Get jwt token",
	Long:  `get jwt token from server and save in $HOME/.eagle/token`,
	Run:   login,
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	loginCmd.Flags().StringVarP(&username, "username", "u", "admin", "username of the admin")
	loginCmd.Flags().StringVarP(&password, "password", "p", "admin", "password of the admin")
}

func login(cmd *cobra.Command, args []string) {

	client := graphql.NewClient("http://localhost:9990/api", nil)

	var q struct {
		TokenAuth graphql.String `graphql:"tokenAuth(username: $username, password: $password)"`
	}

	variables := map[string]interface{}{
		"username": graphql.String(username),
		"password": graphql.String(password),
	}

	if err := client.Query(context.Background(), &q, variables); err != nil {
		fmt.Println(err)
		return
	}
	viper.Set("authorization", q.TokenAuth)
	if err := viper.WriteConfig(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Logged in!")
	fmt.Println("---------")
	fmt.Println("Go on and destroy the world...or save it!")
	fmt.Println("WHO CARES?!")
}
