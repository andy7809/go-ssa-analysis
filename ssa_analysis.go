package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

var EmbedAnalysis = &analysis.Analyzer{
	Name:     "embedanalysis",
	Doc:      "reports embeddings",
	Run:      run,
	Requires: []*analysis.Analyzer{buildssa.Analyzer},
}

var SSAAnalysis = buildssa.Analyzer

func run(pass *analysis.Pass) (interface{}, error) {
	ssaProg := pass.ResultOf[SSAAnalysis].(*buildssa.SSA)
	for _, fn := range ssaProg.SrcFuncs {
		for _, b := range fn.Blocks {
			for _, i := range b.Instrs {
				if p, ok := i.(*ssa.Phi); ok {
					values := make(map[string]bool)
					for _, op := range p.Edges {
						if binop, ok := op.(*ssa.BinOp); ok {
							if values[binop.String()] {
								pass.Reportf(binop.Pos(), "found equivalent SSA values")
							}
							values[op.String()] = true
						}
					}
				}
			}
		}
	}
	return nil, nil
}

func main() {
	multichecker.Main(SSAAnalysis, EmbedAnalysis)
}
