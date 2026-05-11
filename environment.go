package bundleinstall

import (
	"fmt"
	"strconv"
	"strings"
)

type Environment struct {
	KeepGemExtensionBuildFiles bool
	// BundleWithout is the colon-separated list of Bundler groups to exclude
	// from `bundle install` in BOTH the build and launch layers. When empty,
	// the build layer installs all groups (legacy behaviour) and the launch
	// layer falls back to its hardcoded "development:test" default.
	BundleWithout string
}

// ParseEnvironment reads the relevant BP_* and runtime env variables and
// returns the resolved Environment.
//
// BundleWithout precedence:
//  1. Explicit BP_BUNDLE_WITHOUT (highest priority).
//  2. Derived from RAILS_ENV / RACK_ENV (RAILS_ENV wins if both set):
//     - production | staging   -> "development:test"
//     - development            -> "production:test:staging:heroku"
//     - test                   -> "production:development:staging:heroku"
//     - any other / unset      -> "" (legacy behaviour, install everything in build)
func ParseEnvironment(environ []string) (Environment, error) {
	var environment Environment
	var (
		explicitWithout string
		railsEnv        string
		rackEnv         string
	)

	for _, variable := range environ {
		if value, found := strings.CutPrefix(variable, "BP_KEEP_GEM_EXTENSION_BUILD_FILES="); found {
			var err error
			environment.KeepGemExtensionBuildFiles, err = strconv.ParseBool(value)
			if err != nil {
				return Environment{}, fmt.Errorf("failed to parse BP_KEEP_GEM_EXTENSION_BUILD_FILES: %w", err)
			}
		}
		if value, found := strings.CutPrefix(variable, "BP_BUNDLE_WITHOUT="); found {
			explicitWithout = value
		}
		if value, found := strings.CutPrefix(variable, "RAILS_ENV="); found {
			railsEnv = value
		}
		if value, found := strings.CutPrefix(variable, "RACK_ENV="); found {
			rackEnv = value
		}
	}

	switch {
	case explicitWithout != "":
		environment.BundleWithout = explicitWithout
	default:
		appEnv := railsEnv
		if appEnv == "" {
			appEnv = rackEnv
		}
		environment.BundleWithout = defaultBundleWithoutFor(appEnv)
	}

	return environment, nil
}

func defaultBundleWithoutFor(appEnv string) string {
	switch appEnv {
	case "production", "staging":
		return "development:test"
	case "development":
		return "production:test:staging:heroku"
	case "test":
		return "production:development:staging:heroku"
	default:
		return ""
	}
}
