package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUmeng2github(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Umeng2github Suite")
}

var _ = BeforeSuite()