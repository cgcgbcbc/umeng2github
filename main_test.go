package main

import (
	"os"
	"path/filepath"
)

func ExampleMain() {
	os.Setenv("OWNER", "cgcgbcbc")
	os.Setenv("REPO", "umeng2github-test")
	os.Setenv("SHORT", "true")
	args := []string{"umeng2github", "import", filepath.Join("fixtures", "test.csv")}
	os.Args = args
	main()
	// Output:
	// success
}
