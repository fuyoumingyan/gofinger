// Package cmd /*
package cmd

import (
	"bufio"
	"fmt"
	"github.com/fuyoumingyan/gofinger/pkg/banner"
	"github.com/fuyoumingyan/gofinger/pkg/options"
	"github.com/fuyoumingyan/gofinger/pkg/utils"
	"github.com/fuyoumingyan/gofinger/runner"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"time"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gofinger",
	Short:   "一款指纹识别工具",
	Version: "0.99",
	Long:    banner.Banner,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(banner.Banner)
		gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
		if len(url) == 0 && len(file) == 0 && !stdin {
			return
		}
		gologger.Info().Msg("start fingerprint rule matching ...")
		start := time.Now()
		var urls []string
		if url != "" {
			urls = append(urls, url)
		}
		if file != "" {
			lines, err := utils.ReadLines(file)
			lines = utils.DeduplicateEmptyStrings(lines)
			if err != nil {
				gologger.Fatal().Msg(err.Error())
			}
			for _, line := range lines {
				urls = append(urls, line)
			}
		}
		if stdin {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				urls = append(urls, scanner.Text())
			}
		}
		if len(urls) > 1 {
			err := utils.Mkdir("./result")
			if err != nil {
				gologger.Fatal().Msg(err.Error())
			}
		}
		if screenshot {
			err := utils.Mkdir("./result/screenshots")
			if err != nil {
				gologger.Fatal().Msg(err.Error())
			}
		}
		option := &options.Options{
			Thread:     thread,
			Proxy:      proxy,
			Level:      level,
			Urls:       urls,
			Timeout:    timeout,
			Screenshot: screenshot,
		}
		runner := runner.NewRunner(option)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		runner.RunEnumeration()
		elapsed := time.Since(start)
		gologger.Info().Msgf("this time it takes %v .", elapsed)
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

var (
	url, file, proxy       string
	thread, level, timeout int
	stdin, screenshot      bool
)

func init() {
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "-u https://www.baidu.com")
	rootCmd.Flags().StringVarP(&file, "file", "f", "", "-f targets.txt")
	rootCmd.Flags().IntVarP(&thread, "thread", "t", 50, "-t 25")
	rootCmd.Flags().StringVarP(&proxy, "proxy", "p", "", "-p http://127.0.0.1:8080")
	rootCmd.Flags().IntVarP(&level, "level", "l", 1, "-l 1-2")
	rootCmd.Flags().BoolVarP(&stdin, "stdin", "", false, "--stdin true")
	rootCmd.Flags().IntVarP(&timeout, "timeout", "m", 30, "-m 20")
	rootCmd.Flags().BoolVarP(&screenshot, "screenshot", "s", false, "-s false")
}
