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
	fmt.Println(pass)
	fmt.Println("hello world")

	fmt.Println(len(pass.ResultOf))
	for key, value := range pass.ResultOf {
		fmt.Printf("%s: %d\n", key, value)
	}

	ssaProg := pass.ResultOf[SSAAnalysis].(*buildssa.SSA)
	fmt.Println(ssaProg.SrcFuncs[0].Blocks[0].Instrs[0].String())
	for _, fn := range ssaProg.SrcFuncs {
		for _, b := range fn.Blocks {
			fmt.Printf("%s:\n", b.String())
			for _, i := range b.Instrs {
				if v, ok := i.(ssa.Value); ok {
					fmt.Printf("\t[%-20T] %s = %s\n", i, v.Name(), i)
				} else {
					fmt.Printf("\t[%-20T] %s\n", i, i)
				}
			}
		}
	}
	return nil, nil
}

func main() {
	multichecker.Main(SSAAnalysis, EmbedAnalysis)
}
