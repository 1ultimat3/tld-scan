package main

import (
	"fmt"
	"net/http"
	"io"
	"os"
	"strings"
	"sync"
	"net"
	"github.com/spf13/cobra"
)


func worker(domains chan string, wg *sync.WaitGroup) {
	//worker code
	for fqdm := range domains {
		fqdm_with_port := fmt.Sprintf("%s:%d", fqdm, 80)
		conn, err := net.Dial("tcp", fqdm_with_port)
		if err == nil {
			fmt.Println(strings.ToLower(fqdm))
			conn.Close()
		}
		wg.Done()
	}
}

func domainProducer(base string, domains chan string, wg *sync.WaitGroup) {
	res, err := http.Get("https://data.iana.org/TLD/tlds-alpha-by-domain.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	domainsRaw, _ := io.ReadAll(res.Body)
	for _, domain := range strings.Split(string(domainsRaw), "\n") {
		if strings.HasPrefix(domain, "#") || domain  == "" {
			continue
		}
		wg.Add(1)
		domains <- fmt.Sprintf("%s.%s", base, domain)
	}
	//finished producer work
	wg.Done()
}

func TldScan(base string, workers int) {
	domains := make(chan string, workers)
	var wg sync.WaitGroup
	//signalize producer needs to be initialized
	wg.Add(1)
	go domainProducer(base, domains, &wg)
	for i := 0; i < cap(domains); i++ {
		go worker(domains, &wg)
	}
	wg.Wait()
}


func main() {
	var workers int

	var cmdMain =  &cobra.Command {
    		Use:   "scan [basename]",
    		Short: "Determine all reserved <domain>.<top-level-domain>",
    		Long: `This program determines all used domains by checking if there is any web application running on <fqdm>:80.`,
    		Args: cobra.MinimumNArgs(1),
    		Run: func(cmd *cobra.Command, args []string) {
					base := args[0]
					TldScan(base, workers)
    		},
	}
	cmdMain.Flags().IntVarP(&workers, "workers", "w", 10, "amount of workers running in parallel")
	var rootCmd = &cobra.Command{Use: "tld"}
	rootCmd.AddCommand(cmdMain)
	rootCmd.Execute()
}
