package main

import (
	"os"
	"path/filepath"
)

func ExampleMain() {
	os.Setenv("owner", "cgcgbcbc")
	os.Setenv("repo", "umeng2github-test")
	args := []string{"import", filepath.Join("fixtures", "test.csv")}
	os.Args = append(os.Args, args...)
	main()
}
