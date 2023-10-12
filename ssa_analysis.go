package main

import (
	"fmt"

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
	// for key, value := range pass.ResultOf {
	// 	fmt.Printf("%s: %d\n", key, value)
	// }

	ssaProg := pass.ResultOf[SSAAnalysis].(*buildssa.SSA)
	for _, fn := range ssaProg.SrcFuncs {
		for _, b := range fn.Blocks {
			for _, i := range b.Instrs {
				if p, ok := i.(*ssa.Phi); ok {
					// fmt.Printf("[%-20T] %s = %s\n", i, p.Name(), i)
					values := make(map[string]bool)
					for _, op := range p.Edges {
						// fmt.Printf("%+v\n", op)
						if binop, ok := op.(*ssa.BinOp); ok {
							if values[binop.String()] {
								fmt.Printf("phi type: [%-15T]\n", binop)
								fmt.Printf("X: %s\n", binop.X.Name())
								fmt.Printf("Y: %s\n", binop.Y.Name())
							}
							values[op.String()] = true

						}
						// fmt.Printf("%v\n", values)
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
