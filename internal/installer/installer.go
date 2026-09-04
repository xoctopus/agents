// Package installer xoctopus skill 安装器
// +genx:doc
package installer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xoctopus/genx/pkg/agent"
	"github.com/xoctopus/x/misc/must"
)

type Agent string

const (
	CLAUDE Agent = "CLAUDE"
	CURSOR Agent = "CURSOR"
	CODEX  Agent = "CODEX"
)

const source = ".agents/skills"

var agents = []Agent{CLAUDE, CURSOR, CODEX}

var gDefaultPath = map[Agent]string{
	CLAUDE: ".claude/skills",
	CURSOR: ".cursor/skills",
	CODEX:  ".agents/skills",
}

// Installer 安装依赖skill
type Installer struct {
	// Name 名称 如 cursor claude codex. 不指定则安装全部
	Name Agent `cmd:"name"`
	// Path 安装路径. 不指定则使用 Agent 默认路径; 仅可与 Name 同时使用
	Path string `cmd:"path"`
	// Mode 安装模式. 0: 项目级(当前路径) 1: 当前用户
	Mode int `cmd:"mode,default=0"`
	// UpgradeToLatest 是否更新最新go module依赖
	UpgradeToLatest bool `cmd:"upgrade,short=u,noopdef=true"`
}

func (i *Installer) Exec(cmd *cobra.Command, _ ...string) error {
	targets := agents
	if len(i.Name) > 0 {
		targets = []Agent{Agent(strings.ToUpper(string(i.Name)))}
	} else if len(i.Path) > 0 {
		return fmt.Errorf("--path 仅可与 --name 同时使用")
	}

	if err := update(cmd, i.UpgradeToLatest); err != nil {
		return err
	}

	if err := (&agent.Installer{}).Install(context.Background()); err != nil {
		return fmt.Errorf("安装失败: %w", err)
	}

	src := must.NoErrorV(filepath.Abs(source))
	_, err := os.Stat(src)
	must.NoErrorF(err, "读取源目录错误")

	for _, name := range targets {
		target, ok := gDefaultPath[name]
		if !ok {
			cmd.Println("支持的 agent: cursor, claude, codex")
			return fmt.Errorf("不支持的 agent: %s", name)
		}

		cmd.Printf("安装 %s...\n", name)

		if len(i.Path) > 0 {
			target = i.Path
		}

		if i.Mode == 1 {
			target = filepath.Join(must.NoErrorV(os.UserHomeDir()), target)
		}

		dst := must.NoErrorV(filepath.Abs(target))

		if src == dst {
			cmd.Printf("安装成功: ==> %s\n", target)
			continue
		}

		if err = os.MkdirAll(dst, 0o755); err != nil {
			return fmt.Errorf("目标目录创建失败: %s [%w]", dst, err)
		}

		if err = sync(cmd, src, dst); err != nil {
			return err
		}
	}
	return nil
}

func update(cmd *cobra.Command, upgrade bool) error {
	if !upgrade {
		return nil
	}

	cmd.Println("开始更新依赖...")
	info := must.NoErrorV(os.Stat("go.mod"))
	must.BeTrue(!info.IsDir())

	for _, args := range [][]string{
		{"get", "-u", "./..."},
		{"mod", "tidy"},
	} {
		c := exec.Command("go", args...)
		c.Env = append(os.Environ(), "GOWORK=off")

		c.Stdout = cmd.OutOrStdout()
		c.Stderr = cmd.OutOrStderr()

		if err := c.Run(); err != nil {
			return fmt.Errorf("依赖更新失败: %w", err)
		}
	}

	return nil
}

func sync(cmd *cobra.Command, src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("读取源目录失败: %s [%w]", src, err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	sort.Strings(names)

	maxName := 0
	for _, name := range names {
		if len(name) > maxName {
			maxName = len(name)
		}
	}

	for _, name := range names {
		srcPath := filepath.Join(src, name)
		dstPath := filepath.Join(dst, name)

		_ = os.RemoveAll(dstPath)
		if err := os.Symlink(srcPath, dstPath); err != nil {
			return fmt.Errorf("安装失败 %s -> %s [%w]", dstPath, srcPath, err)
		}
		pad := strings.Repeat(" ", maxName-len(name))
		cmd.Printf("安装成功 [%s]%s ==> %s\n", name, pad, dstPath)
	}

	return nil
}
