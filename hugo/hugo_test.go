package hugo_test

import (
	hugo "github.com/joaodias/hugito-app/hugo"
	"github.com/joaodias/hugito-app/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hugo", func() {
	Describe("When building an HUGO website", func() {
		Context("and there is a problem with the command or the argumet", func() {
			It("should return an error", func() {
				commandExecutor := &mocks.CommandExecutor{
					IsError: true,
				}
				Expect(hugo.BuildSite("stupidCommand", "badDirectory", commandExecutor).Error()).To(Equal("Error building HUGO site."))
			})
		})
		Context("and the website is successfully built", func() {
			It("should return no error", func() {
				commandExecutor := &mocks.CommandExecutor{
					IsError: false,
				}
				Expect(hugo.BuildSite("stupidCommand", "badDirectory", commandExecutor)).To(BeNil())
			})
		})
	})
})
