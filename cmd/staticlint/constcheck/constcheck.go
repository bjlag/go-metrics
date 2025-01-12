// Package constcheck проверяет код на наличие магических чисел/строк, которые можно вынести в константы.
package constcheck

import (
	"go/token"
	"strings"

	"github.com/jgautheron/goconst"
	"golang.org/x/tools/go/analysis"
)

const (
	mockSuffix = "_mock.go"
)

var ConstCheckAnalyzer = &analysis.Analyzer{
	Name: "constcheck",
	Doc:  "find repeated strings that could be replaced by a constant",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	cfg := &goconst.Config{
		IgnoreTests:    true,
		ParseNumbers:   true,
		NumberMin:      1,
		MinOccurrences: 2,
		ExcludeTypes: map[goconst.Type]bool{
			goconst.Call: true,
		},
	}

	issues, err := goconst.Run(pass.Files, pass.Fset, cfg)
	if err != nil {
		return nil, err
	}

	for _, issue := range issues {
		if strings.HasSuffix(issue.Pos.Filename, mockSuffix) {
			continue
		}

		var pos token.Pos
		pass.Fset.Iterate(func(file *token.File) bool {
			if file.Name() == issue.Pos.Filename {
				pos = file.Pos(issue.Pos.Offset)
			}
			return true
		})

		if pos == token.NoPos {
			continue
		}

		pass.Reportf(pos, "could be replaced by a constant '%s', repeated cont %d", issue.Str, issue.OccurrencesCount)
	}

	return nil, nil
}
