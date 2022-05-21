package main

import (
	"os"
	"testing"

	"github.com/adrg/xdg"
	"github.com/motemen/go-gitconfig"
)

func Test_newConfig(t *testing.T) {
	expected := &config{
		Owner:          "gitcfg",
		Repo:           "diary",
		Branch:         "main",
		CommitterName:  "gitcfg",
		CommitterEmail: "gitcfg@example.com",
		PAT:            "envpat",
		Editor:         "enveditor",
	}

	orgPAT := os.Getenv(patEnvName)
	orgEditor := os.Getenv(editorEnvName)
	os.Setenv(patEnvName, "envpat")
	os.Setenv(editorEnvName, "enveditor")
	defer func() {
		os.Setenv(patEnvName, orgPAT)
		os.Setenv(editorEnvName, orgEditor)
	}()

	gitConfig := gitconfig.File("test/gitconfig")
	actual, err := newConfig(gitConfig)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	assertEqualConfig(t, actual, expected)
}

func Test_getConfigPath(t *testing.T) {
	orgConfigHome := xdg.ConfigHome
	xdg.ConfigHome = "./test"
	defer func() { xdg.ConfigHome = orgConfigHome }()

	cases := map[string]struct {
		input  string
		expect string
	}{
		"without_input": {"", "test/hubdiary/config.json"},
		"with_input":    {"/path/to/myconfig.json", "/path/to/myconfig.json"},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			cfg := &config{}
			actual := cfg.getConfigPath(c.input)

			if actual != c.expect {
				t.Errorf("Exptected %s, but got %s", c.expect, actual)
			}
		})
	}
}

func Test_loadFile(t *testing.T) {
	orgConfigHome := xdg.ConfigHome
	xdg.ConfigHome = "./test"
	defer func() { xdg.ConfigHome = orgConfigHome }()

	cases := map[string]struct {
		input     string
		expect    *config
		expectErr bool
	}{
		"set_all": {
			"test/set_all.json",
			&config{"cfgowner", "cfgrepo", "cfgbranch", "cfgname", "cfgemail", "cfgpat", "cfgeditor"},
			false,
		},
		"set_owner": {
			"test/set_owner.json",
			&config{Owner: "cfgowner", Repo: "diary"},
			false,
		},
		"home_config": {
			"",
			&config{Owner: "home", Repo: "homerepo"},
			false,
		},
		"invalid":   {"test/invalid.json", &config{}, true},
		"not_exist": {"notexist.json", &config{}, true},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			cfg := &config{Owner: "gituser", Repo: "diary"}
			err := cfg.loadFile(c.input)

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

func assertEqualConfig(t *testing.T, actual, expect *config) {
	t.Helper()

	if actual.Owner != expect.Owner {
		t.Errorf("Owner expects %s, but got %s", expect.Owner, actual.Owner)
	}

	if actual.Repo != expect.Repo {
		t.Errorf("Repo expects %s, but got %s", expect.Repo, actual.Repo)
	}

	if actual.Branch != expect.Branch {
		t.Errorf("Branch expects %s, but got %s", expect.Branch, actual.Branch)
	}

	if actual.CommitterName != expect.CommitterName {
		t.Errorf("CommitterName expects %s, but got %s", expect.CommitterName, actual.CommitterName)
	}

	if actual.CommitterEmail != expect.CommitterEmail {
		t.Errorf("CommitterEmail expects %s, but got %s", expect.CommitterEmail, actual.CommitterEmail)
	}

	if actual.PAT != expect.PAT {
		t.Errorf("PAT expects %s, but got %s", expect.PAT, actual.PAT)
	}

	if actual.Editor != expect.Editor {
		t.Errorf("Editor expects %s, but got %s", expect.Editor, actual.Editor)
	}
}
