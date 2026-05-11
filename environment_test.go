package bundleinstall_test

import (
	"testing"

	bundleinstall "github.com/paketo-buildpacks/bundle-install"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testEnvironment(t *testing.T, context spec.G, it spec.S) {
	var Expect = NewWithT(t).Expect

	context("ParseEnvironment", func() {
		it("parse the environment variables", func() {
			environment, err := bundleinstall.ParseEnvironment([]string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(environment).To(Equal(bundleinstall.Environment{
				KeepGemExtensionBuildFiles: false,
			}))
		})

		context("when BP_KEEP_GEM_EXTENSION_BUILD_FILES is set", func() {
			it("parse the environment variables", func() {
				environment, err := bundleinstall.ParseEnvironment([]string{
					"BP_KEEP_GEM_EXTENSION_BUILD_FILES=true",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(environment).To(Equal(bundleinstall.Environment{
					KeepGemExtensionBuildFiles: true,
				}))
			})
		})

		context("when BP_BUNDLE_WITHOUT is set explicitly", func() {
			it("uses the explicit value", func() {
				environment, err := bundleinstall.ParseEnvironment([]string{
					"BP_BUNDLE_WITHOUT=development:test:assets",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(environment.BundleWithout).To(Equal("development:test:assets"))
			})

			it("wins over RAILS_ENV-derived default", func() {
				environment, err := bundleinstall.ParseEnvironment([]string{
					"BP_BUNDLE_WITHOUT=custom:groups",
					"RAILS_ENV=production",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(environment.BundleWithout).To(Equal("custom:groups"))
			})
		})

		context("when RAILS_ENV is set", func() {
			it("derives 'development:test' for production", func() {
				environment, err := bundleinstall.ParseEnvironment([]string{
					"RAILS_ENV=production",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(environment.BundleWithout).To(Equal("development:test"))
			})

			it("derives 'development:test' for staging", func() {
				environment, err := bundleinstall.ParseEnvironment([]string{
					"RAILS_ENV=staging",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(environment.BundleWithout).To(Equal("development:test"))
			})

			it("derives 'production:test:staging:heroku' for development", func() {
				environment, err := bundleinstall.ParseEnvironment([]string{
					"RAILS_ENV=development",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(environment.BundleWithout).To(Equal("production:test:staging:heroku"))
			})

			it("derives 'production:development:staging:heroku' for test", func() {
				environment, err := bundleinstall.ParseEnvironment([]string{
					"RAILS_ENV=test",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(environment.BundleWithout).To(Equal("production:development:staging:heroku"))
			})

			it("falls back to empty for unrecognized values", func() {
				environment, err := bundleinstall.ParseEnvironment([]string{
					"RAILS_ENV=somethingweird",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(environment.BundleWithout).To(Equal(""))
			})
		})

		context("when only RACK_ENV is set", func() {
			it("derives 'development:test' for production", func() {
				environment, err := bundleinstall.ParseEnvironment([]string{
					"RACK_ENV=production",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(environment.BundleWithout).To(Equal("development:test"))
			})
		})

		context("when both RAILS_ENV and RACK_ENV are set", func() {
			it("RAILS_ENV wins", func() {
				environment, err := bundleinstall.ParseEnvironment([]string{
					"RAILS_ENV=production",
					"RACK_ENV=development",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(environment.BundleWithout).To(Equal("development:test"))
			})
		})

		context("failure cases", func() {
			context("when the BP_KEEP_GEM_EXTENSION_BUILD_FILES env var cannot be parsed", func() {
				it("returns an error", func() {
					_, err := bundleinstall.ParseEnvironment([]string{
						"BP_KEEP_GEM_EXTENSION_BUILD_FILES=banana",
					})
					Expect(err).To(MatchError(ContainSubstring("failed to parse BP_KEEP_GEM_EXTENSION_BUILD_FILES:")))
					Expect(err).To(MatchError(ContainSubstring(`parsing "banana": invalid syntax`)))
				})
			})
		})
	})
}
