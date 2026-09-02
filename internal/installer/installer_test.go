package installer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	"github.com/xoctopus/agents/internal/installer"
)

func chdirAgentsRoot(t *testing.T) string {
	t.Helper()

	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	return root
}

func TestInstaller_Exec_cursor(t *testing.T) {
	root := chdirAgentsRoot(t)

	cmd := &cobra.Command{}
	inst := &installer.Installer{Name: installer.CURSOR}
	if err := inst.Exec(cmd); err != nil {
		t.Fatal(err)
	}

	src := filepath.Join(root, ".agents", "skills")
	dst := filepath.Join(root, ".cursor", "skills")

	entries, err := os.ReadDir(src)
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range entries {
		name := entry.Name()
		want := filepath.Join(src, name)
		got, err := os.Readlink(filepath.Join(dst, name))
		if err != nil {
			t.Fatalf("readlink %s: %v", name, err)
		}
		if got != want {
			t.Fatalf("%s link = %q, want %q", name, got, want)
		}
	}
}

func TestInstaller_Exec_codex(t *testing.T) {
	chdirAgentsRoot(t)

	cmd := &cobra.Command{}
	inst := &installer.Installer{Name: installer.CODEX}
	if err := inst.Exec(cmd); err != nil {
		t.Fatal(err)
	}
}
