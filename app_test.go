package main

import (

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("App", func() {
	Describe("When I GET /", func() {
		Context("And I don't have an github oauth2 token stored in session cookie", func() {
			Before
			It("should redirect me to /login", func() {})
		})

		Context("And I do have a github oauth2 token store in session cookie", func() {
			It("should redirect me to /dashboard", func() {})
		})
	})
})
