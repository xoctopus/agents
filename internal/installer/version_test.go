package installer_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	. "github.com/xoctopus/x/testx"
	"github.com/xoctopus/x/testx/bdd"
	"golang.org/x/mod/modfile"

	"github.com/xoctopus/agents/internal/installer"
)

type skillRequire struct {
	module  string
	version string
	skills  []string
}

func skillRequires(t testing.TB, modpath string) []skillRequire {
	t.Helper()

	data, err := os.ReadFile(modpath)
	Expect(t, err, Succeed())

	f, err := modfile.Parse(modpath, data, nil)
	Expect(t, err, Succeed())

	var out []skillRequire
	for _, r := range f.Require {
		if r.Indirect {
			continue
		}
		var skills []string
		for _, c := range r.Syntax.Comments.Before {
			token := strings.TrimSpace(strings.TrimPrefix(c.Token, "//"))
			if name, found := strings.CutPrefix(token, "+skill:"); found {
				skills = append(skills, name)
			}
		}
		if len(skills) > 0 {
			out = append(out, skillRequire{
				module:  r.Mod.Path,
				version: r.Mod.Version,
				skills:  skills,
			})
		}
	}
	Expect(t, len(out), BeGt(0))
	return out
}

func TestVersion_Exec(t *testing.T) {
	bdd.From(t).Given("go.mod has +skill direct requires", func(t bdd.T) {
		chdirAgentsRoot(t)
		want := skillRequires(t, "go.mod")

		var buf bytes.Buffer
		cmd := &cobra.Command{}
		cmd.SetOut(&buf)

		t.When("Exec", func(t bdd.T) {
			err := (&installer.Version{}).Exec(cmd)
			t.Then("succeed", bdd.Succeed(err))

			out := buf.String()
			for _, info := range want {
				modLine := info.module + "@" + info.version
				Expect(t, out, ContainsSubString(modLine))
				for _, skill := range info.skills {
					Expect(t, out, ContainsSubString("\t"+skill+"\n"))
				}
			}
			Expect(t, strings.Contains(out, "github.com/spf13/cobra"), BeFalse())
		})
	})
}
