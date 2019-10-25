package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mapleque/vpsagent/client"
	"github.com/mapleque/vpsagent/server"
)

func main() {
	var (
		c            string
		body         string
		bodyFilePath string

		port        int
		token       string
		tlsKeyPath  string
		tlsCertPath string
		ipAllows    string
		verbose     bool
	)
	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.StringVar(&token, "sign-token", "", "token for signature")

	flag.StringVar(&c, "c", "", "run in client mode, pass host of server, such as: http://127.0.0.1:8888/")
	flag.StringVar(&body, "script", "", "for client mode, script")
	flag.StringVar(&bodyFilePath, "script-file", "", "for client mode, script file path")

	flag.IntVar(&port, "p", 0, "run in server mode, listening port")
	flag.StringVar(&tlsKeyPath, "tls-key-file", "", "for server mode, tls key file path")
	flag.StringVar(&tlsCertPath, "tls-cert-file", "", "for server mode, tls cert file path")
	flag.StringVar(&ipAllows, "ip-allows", "127.0.0.1", "for server mode, client ip white list, multiple splited with ','")

	flag.Parse()

	if c != "" {
		cli := client.New(c, token, verbose)
		if body != "" {
			cli.Do([]byte(body))
		}
		if bodyFilePath != "" {
			f, err := os.Open(bodyFilePath)
			if err != nil {
				panic(err)
			}

			fbody, err := ioutil.ReadAll(f)
			if err != nil {
				panic(err)
			}
			cli.Do(fbody)
		}
		fmt.Println("script or script-file is required in client mode, run 'vpsagent -h' for more information")
		os.Exit(1)
	}

	if port == 0 {
		flag.Usage()
		os.Exit(1)
	}

	agent := server.New(
		port,
		token,
		tlsKeyPath,
		tlsCertPath,
		strings.Split(ipAllows, ","),
		verbose,
	)

	agent.Run()
}
