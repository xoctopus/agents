package installer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	. "github.com/xoctopus/x/testx"
	"github.com/xoctopus/x/testx/bdd"

	"github.com/xoctopus/agents/internal/installer"
)

func chdirAgentsRoot(t testing.TB) string {
	t.Helper()

	root, err := filepath.Abs(filepath.Join("..", ".."))
	Expect(t, err, Succeed())

	orig, err := os.Getwd()
	Expect(t, err, Succeed())
	Expect(t, os.Chdir(root), Succeed())
	t.Cleanup(func() { _ = os.Chdir(orig) })

	return root
}

func skillLinkTargets(t testing.TB, dir string) map[string]string {
	t.Helper()

	entries, err := os.ReadDir(dir)
	Expect(t, err, Succeed())
	Expect(t, len(entries), BeGt(0))

	out := make(map[string]string, len(entries))
	for _, entry := range entries {
		target, err := os.Readlink(filepath.Join(dir, entry.Name()))
		Expect(t, err, Succeed())
		out[entry.Name()] = target
	}
	return out
}

func expectLinkedToSource(t testing.TB, root, dst string) {
	t.Helper()

	src := filepath.Join(root, ".agents", "skills")
	for name := range skillLinkTargets(t, src) {
		got, err := os.Readlink(filepath.Join(dst, name))
		Expect(t, err, Succeed())
		Expect(t, got, Equal(filepath.Join(src, name)))
	}
}

func expectNoSelfSymlink(t testing.TB, src string) {
	t.Helper()

	for name, target := range skillLinkTargets(t, src) {
		Expect(t, target, NotEqual(filepath.Join(src, name)))
	}
}

func TestInstaller_Exec(t *testing.T) {
	bdd.From(t).Given("name is cursor", func(t bdd.T) {
		root := chdirAgentsRoot(t)
		inst := &installer.Installer{Name: installer.CURSOR}

		t.When("Exec", func(t bdd.T) {
			err := inst.Exec(&cobra.Command{})
			t.Then("succeed", bdd.Succeed(err))
			expectLinkedToSource(t, root, filepath.Join(root, ".cursor", "skills"))
		})
	})

	bdd.From(t).Given("name is claude", func(t bdd.T) {
		root := chdirAgentsRoot(t)
		inst := &installer.Installer{Name: installer.CLAUDE}

		t.When("Exec", func(t bdd.T) {
			err := inst.Exec(&cobra.Command{})
			t.Then("succeed", bdd.Succeed(err))
			expectLinkedToSource(t, root, filepath.Join(root, ".claude", "skills"))
		})
	})

	bdd.From(t).Given("name is codex", func(t bdd.T) {
		root := chdirAgentsRoot(t)
		src := filepath.Join(root, ".agents", "skills")
		before := skillLinkTargets(t, src)
		inst := &installer.Installer{Name: installer.CODEX}

		t.When("Exec", func(t bdd.T) {
			err := inst.Exec(&cobra.Command{})
			t.Then("succeed", bdd.Succeed(err))

			after := skillLinkTargets(t, src)
			expectNoSelfSymlink(t, src)
			for name, target := range after {
				if want, ok := before[name]; ok {
					Expect(t, target, Equal(want))
				}
			}
		})
	})

	bdd.From(t).Given("name is empty", func(t bdd.T) {
		root := chdirAgentsRoot(t)
		inst := &installer.Installer{}

		t.When("Exec", func(t bdd.T) {
			err := inst.Exec(&cobra.Command{})
			t.Then("succeed", bdd.Succeed(err))
			expectLinkedToSource(t, root, filepath.Join(root, ".cursor", "skills"))
			expectLinkedToSource(t, root, filepath.Join(root, ".claude", "skills"))
			expectNoSelfSymlink(t, filepath.Join(root, ".agents", "skills"))
		})
	})

	bdd.From(t).Given("path without name", func(t bdd.T) {
		chdirAgentsRoot(t)
		inst := &installer.Installer{Path: "custom/skills"}

		t.When("Exec", func(t bdd.T) {
			err := inst.Exec(&cobra.Command{})
			t.Then("failed", bdd.Failed(err))
			t.Then("error mentions --path", bdd.ErrorContains(err, "--path"))
		})
	})

	bdd.From(t).Given("unsupported name", func(t bdd.T) {
		chdirAgentsRoot(t)
		inst := &installer.Installer{Name: "unknown"}

		t.When("Exec", func(t bdd.T) {
			err := inst.Exec(&cobra.Command{})
			t.Then("failed", bdd.Failed(err))
			t.Then("error mentions unsupported agent", bdd.ErrorContains(err, "不支持的 agent"))
		})
	})
}
