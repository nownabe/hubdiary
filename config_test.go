package main

import (
	"os"
	"testing"

	"github.com/adrg/xdg"
	"github.com/motemen/go-gitconfig"
)

func assertEqualConfig(t *testing.T, actual, expect *config) {
	t.Helper()

	if actual.Repo != expect.Repo {
		t.Errorf("Repo expects %s, but got %s", expect.Repo, actual.Repo)
	}

	if actual.User != expect.User {
		t.Errorf("User expects %s, but got %s", expect.User, actual.User)
	}

	if actual.Email != expect.Email {
		t.Errorf("Email expects %s, but got %s", expect.Email, actual.Email)
	}

	if actual.PAT != expect.PAT {
		t.Errorf("PAT expects %s, but got %s", expect.PAT, actual.PAT)
	}
}

func Test_newConfig(t *testing.T) {
	type args struct {
		configPath string
		repo       string
		user       string
		email      string
		pat        string
	}

	cases := map[string]struct {
		input     args
		expect    *config
		expectErr bool
	}{
		"home_config": {
			args{"", "", "", "", ""},
			&config{"home/home", "home", "home@example.com", "homepat"},
			false,
		},
		"set_all": {
			args{"test/set_all.json", "", "", "", ""},
			&config{"testowner/testrepo", "testuser", "testuser@example.com", "testpat"},
			false,
		},
		"overwrite": {
			args{"test/set_all.json", "over/over", "over", "over@example.com", "overpat"},
			&config{"over/over", "over", "over@example.com", "overpat"},
			false,
		},
		"no_config": {
			args{"notexist", "", "", "", ""},
			&config{"gitcfg/diary", "gitcfg", "gitcfg@example.com", "envpat"},
			false,
		},
		"invalid_json": {
			args{"test/invalid.json", "", "", "", ""},
			&config{},
			true,
		},
	}

	xdg.ConfigHome = "./test"

	orgPAT := os.Getenv("GITHUB_PAT")
	os.Setenv(patEnvName, "envpat")
	defer func() {
		os.Setenv(patEnvName, orgPAT)
	}()

	gitConfig := gitconfig.File("test/gitconfig")

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			// t.Parallel()
			cfg, err := newConfig(gitConfig, c.input.configPath,
				c.input.repo, c.input.user, c.input.email, c.input.pat)

			if c.expectErr && err == nil {
				t.Fatal("error is expected")
			}

			if !c.expectErr {
				if err != nil {
					t.Fatalf("error is not expected: %v", err)
				}
				assertEqualConfig(t, cfg, c.expect)
			}
		})
	}
}
