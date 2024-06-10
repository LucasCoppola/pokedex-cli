package main

import "os"

func commandExit(args ...string) error {
	os.Exit(0)
	return nil
}
