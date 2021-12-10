// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/10
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "bytes"
    _ "embed"
    "entgo.io/ent/entc/gen"
    "fmt"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/spf13/cobra"
    "log"
    "os"
    "path/filepath"
    "strings"
    "text/template"
)

const (
    defaultSchema = "internal/ent/schema"
)

// schema template for the "init" command.
var tmpl = template.Must(template.New("schema").
    Parse(`package schema

import (
    "entgo.io/ent"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// {{ . }} holds the schema definition for the {{ . }} entity.
type {{ . }} struct {
    ent.Schema
}

// Fields of the {{ . }}.
func ({{ . }}) Fields() []ent.Field {
    return []ent.Field{}
}

// Edges of the {{ . }}.
func ({{ . }}) Edges() []ent.Edge {
    return nil
}

func ({{ . }}) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
    }
}
`))

func InitCmd() *cobra.Command {
    var target, customtmpl string
    cmd := &cobra.Command{
        Use:   "init [flags] [schemas]",
        Short: "initialize an environment with zero or more schemas",
        Example: examples(
            "ent init Example",
            "ent init --target entv1/schema User Group",
            "ent init --target entv1/schema --template tmpl/default.tmpl User Group",
        ),
        Run: func(cmd *cobra.Command, names []string) {
            if err := initEnv(target, names); err != nil {
                log.Fatalln(fmt.Errorf("ent/init: %w", err))
            }
        },
    }
    cmd.Flags().StringVar(&target, "target", defaultSchema, "target directory for schemas")
    cmd.Flags().StringVar(&customtmpl, "template", "", "target template for schemas")
    return cmd
}

func initEnv(target string, names []string) error {
    if err := createDir(target); err != nil {
        return fmt.Errorf("create dir %s: %w", target, err)
    }

    for _, name := range names {
        name = utils.StrToFirstUpper(name)
        if err := gen.ValidSchemaName(name); err != nil {
            return fmt.Errorf("init schema %s: %w", name, err)
        }
        b := bytes.NewBuffer(nil)
        if err := tmpl.Execute(b, name); err != nil {
            return fmt.Errorf("executing template %s: %w", name, err)
        }
        newFileTarget := filepath.Join(target, strings.ToLower(name+".go"))
        if err := os.WriteFile(newFileTarget, b.Bytes(), 0644); err != nil {
            return fmt.Errorf("writing file %s: %w", newFileTarget, err)
        }
    }
    return nil
}
