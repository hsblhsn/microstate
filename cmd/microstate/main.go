package main

import (
	"github.com/hsblhsn/microstate/cli"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			l := cli.NewLogger()
			l.Error(r)
		}
	}()
	if err := cli.NewRootCmd().Execute(); err != nil {
		panic(err)
	}
}
