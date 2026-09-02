package installer_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/xoctopus/agents/internal/installer"
)

func TestVersion_Exec(t *testing.T) {
	chdirAgentsRoot(t)

	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)

	if err := (&installer.Version{}).Exec(cmd); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	for _, want := range []string{
		"github.com/xoctopus/concx@v0.2.2",
		"github.com/xoctopus/confx@v0.5.9",
		"github.com/xoctopus/genx@v0.3.8",
		"skills:",
		"\tconcx",
		"\tappx",
		"\tkg",
		"\tgenx",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "github.com/spf13/cobra") {
		t.Fatalf("output should not contain deps without +skill:\n%s", out)
	}
}
