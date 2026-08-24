package main

import (
	"context"
	"fmt"
	"os"

	"github.com/xoctopus/genx/pkg/agent"
)

import (
	_ "github.com/xoctopus/concx/pkg/chanx"
	_ "github.com/xoctopus/confx/pkg/types"
	_ "github.com/xoctopus/logx"
	_ "github.com/xoctopus/sqlx/pkg/builder"
	_ "github.com/xoctopus/x/testx"
)

func main() {
	if err := (&agent.Installer{}).Install(context.Background()); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}
}
