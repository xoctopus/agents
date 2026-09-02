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
		"github.com/spf13/cobra (",
		"github.com/xoctopus/concx (",
		"github.com/xoctopus/confx (",
		"github.com/xoctopus/genx (",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "skill:") {
		t.Fatalf("output should not contain skill names:\n%s", out)
	}
}
