package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ani-sche-natalie URL CalenderID",
	Short: "指定したナタリーのアニメまとめページをパースしてGoogleカレンダーに登録します",
	Long: `指定したナタリーのアニメまとめページをパースしてGoogleカレンダーに登録します
	 例えば: https://natalie.mu/comic/anime/season/2024-autumn のようなページから登録します
	 登録先のGoogleカレンダーのCalderIDも必要です。
	 設定>登録するマイカレンダー>カレンダーIDから取得します`,
	Run: func(cmd *cobra.Command, args []string) {
		AnimeSchedule()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
