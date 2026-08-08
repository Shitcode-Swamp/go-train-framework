package parsers

import (
	"go/ast"
	"go/token"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/ctxgen/injection"
	utils_strings "github.com/Shitcode-Swamp/go-train-framework/src/utils/strings"
)

type LexicalParser struct {}

func NewLexicalParser() *LexicalParser {
	return &LexicalParser{}
}

func (p *LexicalParser) ParseFuncParams(v *ast.FuncType, file *ast.File) []injection.Injectee {
	params := v.Params
	if params == nil {
		return nil
	}
	injectees := make([]injection.Injectee, 0, len(params.List))
	for i, param := range params.List {
		start := param.Pos()
		var end token.Pos
		if i+1 < len(params.List) {
			end = params.List[i+1].Pos()
		} else {
			end = params.Closing
		}

		commentText := p.findComment(file.Comments, start, end)
		if commentText == "" && param.Comment != nil {
			commentText = param.Comment.Text()
		}

		t := utils_strings.Extract(commentText, `@Inject\(([^)]+)\)`)
		if t == nil {
			continue
		}
		injectees = append(injectees, injection.NewInjectee(*t))
	}
	return injectees
}

func (p *LexicalParser) findComment(groups []*ast.CommentGroup, start, end token.Pos) string {
	for _, cg := range groups {
		if cg.Pos() >= start && cg.Pos() < end {
			return cg.Text()
		}
	}
	return ""
}
