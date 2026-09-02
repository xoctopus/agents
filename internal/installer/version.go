package installer

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

// Version 打印skill版本
type Version struct{}

func (v *Version) Exec(cmd *cobra.Command, _ ...string) error {
	imports, err := directs("go.mod")
	if err != nil {
		return fmt.Errorf("读取版本信息失败: %w", err)
	}

	for _, i := range imports {
		cmd.Println(i.module + "@" + i.version)
		cmd.Println("skills:")
		for _, name := range i.skills {
			cmd.Printf("\t%s\n", name)
		}
	}

	return nil
}

type Info struct {
	module  string
	version string
	skills  []string
}

func directs(modpath string) ([]*Info, error) {
	data, err := os.ReadFile(modpath)
	if err != nil {
		return nil, fmt.Errorf("读取 go.mod 失败: %w", err)
	}

	f, err := modfile.Parse(modpath, data, nil)
	if err != nil {
		return nil, fmt.Errorf("解析 go.mod 失败: %w", err)
	}

	directives := make([]*Info, 0)
	for _, r := range f.Require {
		if r.Indirect {
			continue
		}
		skills := make([]string, 0)
		for _, c := range r.Syntax.Comments.Before {
			if len(c.Token) > 0 {
				token := strings.TrimSpace(strings.TrimPrefix(c.Token, "//"))
				if name, found := strings.CutPrefix(token, "+skill:"); found {
					skills = append(skills, name)
				}
			}
		}
		if len(skills) > 0 {
			directives = append(directives, &Info{
				module:  r.Mod.Path,
				version: r.Mod.Version,
				skills:  skills,
			})
		}
	}

	sort.Slice(directives, func(i, j int) bool {
		return directives[i].module < directives[j].module
	})

	return directives, nil
}
