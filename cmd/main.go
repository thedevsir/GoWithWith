package main

import (
	"fmt"
	"os"
	"strings"

	mAdmin "github.com/Gommunity/GoWithWith/app/model/admin"
	"github.com/Gommunity/GoWithWith/config/database"
	"github.com/Gommunity/GoWithWith/services/encrypt"
	"github.com/Gommunity/GoWithWith/services/utility"
	"github.com/spf13/cobra"
	"github.com/zebresel-com/mongodm"
)

func main() {

	utility.LoadEnvironmentVariables(".env")

	var cmdAdminUserInstall = &cobra.Command{
		Use:   "adminRootInstall [CreateSuperUser]",
		Short: "For first-time setup root admin user",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			var username, password = "root", "admin"

			db := &database.Composer{
				Locals:   "resource/locals/locals.json",
				Addrs:    []string{os.Getenv("DBAddrs")},
				Database: os.Getenv("DBName"),
				Username: os.Getenv("DBUsername"),
				Password: os.Getenv("DBPassword"),
				Source:   os.Getenv("DBSource"),
			}
			db.Shoot(map[string]mongodm.IDocumentBase{
				"admin": &mAdmin.Admin{},
			})

			adminModel := database.Connection.Model(mAdmin.AdminCollection)
			admin := &mAdmin.Admin{}
			adminModel.New(admin)
			hash, _ := encrypt.Hash(password)

			admin.Username = strings.ToLower(username)
			admin.Password = hash

			err := admin.Save()

			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println("Administrator SuperUser Successfully Created !")
		},
	}

	var rootCmd = &cobra.Command{Use: "cmd"}
	rootCmd.AddCommand(cmdAdminUserInstall)
	rootCmd.Execute()
}
