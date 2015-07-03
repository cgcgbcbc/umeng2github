package main

import (
	"fmt"
	"time"
)

func ExampleFormatRecord() {
	testDate, err := time.Parse("2006-01-02 03:04:05", "2015-06-10 09:21:43")
	if err != nil {
		fmt.Print(err)
		return
	}
	record := Record{
		Title:      "test",
		Count:      2,
		FirstDate:  testDate,
		Version:    "0.1.0",
		StackTrace: "java.lang.NullPointerException",
	}

	result, err := FormatRecord(record)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(result)

	// Output:
	// 发生次数 | 首次发生时间 | 版本
	// -----|-----|-----
	// 2 | 2015-06-10 09:21:43 | 0.1.0
	//
	// StackTrace:
	// ```
	// java.lang.NullPointerException
	// ```
}
