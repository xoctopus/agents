// Package installer xoctopus skill 安装器
// +genx:doc
package installer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

var gDefaultPath = map[Agent]string{
	CLAUDE: ".claude/skills",
	CURSOR: ".cursor/skills",
	CODEX:  ".agents/skills",
}

// Installer 安装依赖skill
type Installer struct {
	// AgentName 名称 如 cursor claude codex
	Name Agent `cmd:"name,require"`
	// Path 安装路径. 不指定则使用 Agent 默认路径
	Path string `cmd:"path"`
	// Mode 安装模式. 0: 项目级(当前路径) 1: 当前用户
	Mode int `cmd:"mode,default=0"`
}

func (i *Installer) Exec(cmd *cobra.Command, _ ...string) error {
	if err := update(cmd); err != nil {
		return err
	}

	if err := (&agent.Installer{}).Install(context.Background()); err != nil {
		return fmt.Errorf("安装失败: %w", err)
	}

	i.Name = Agent(strings.ToUpper(string(i.Name)))
	target, ok := gDefaultPath[i.Name]
	if !ok {
		cmd.Println("支持的 agent: cursor, claude, codex")
		return fmt.Errorf("不支持的 agent: %s", i.Name)
	}
	if len(i.Path) > 0 {
		target = i.Path
	}

	if i.Mode == 1 {
		target = filepath.Join(must.NoErrorV(os.UserHomeDir()), target)
	}

	src := must.NoErrorV(filepath.Abs(source))
	dst := must.NoErrorV(filepath.Abs(target))

	_, err := os.Stat(src)
	must.NoErrorF(err, "读取源目录错误")

	if src == dst {
		cmd.Printf("安装成功: ==> %s\n", target)
		return nil
	}

	if err = os.MkdirAll(dst, 0o755); err != nil {
		return fmt.Errorf("目标目录创建失败: %s [%w]", dst, err)
	}

	return sync(cmd, src, dst)
}

func update(cmd *cobra.Command) error {
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

	for _, entry := range entries {
		name := entry.Name()
		srcPath := filepath.Join(src, name)
		dstPath := filepath.Join(dst, name)

		_ = os.RemoveAll(dstPath)
		if err := os.Symlink(srcPath, dstPath); err != nil {
			return fmt.Errorf("安装失败 %s -> %s [%w]", dstPath, srcPath, err)
		}
		cmd.Printf("安装成功 [%s] ==> %s\n", name, dstPath)
	}

	return nil
}
