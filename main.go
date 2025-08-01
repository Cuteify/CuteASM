package main

import (
	"CuteASM/arch/x86"
	"CuteASM/compiler"
	"CuteASM/lexer"
	"CuteASM/parser"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	path := "./test.asm"
	archType := "x86" // 默认架构
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	if len(os.Args) > 2 {
		archType = os.Args[2] // 从命令行参数获取架构类型
	}
	start := time.Now()
	if archType == "all" {
		for _, arch := range []string{"x86", "mips", "arm", "riscv", "loongarch", "x86_64"} {
			Compile(path, arch)
		}
	} else if strings.Contains(archType, ",") {
		archs := strings.Split(archType, ",")
		for _, arch := range archs {
			Compile(path, strings.TrimSpace(arch))
		}
	} else {
		Compile(path, archType)
	}
	fmt.Println("总耗时", time.Since(start))
}

var btmp = []byte{}

func pr(block *parser.Node, tabnum int) {
	tmp := ""
	for i := 0; i < tabnum; i++ {
		tmp += "\t"
	}
	tmp2 := []byte{}
	if _, ok := block.Value.(*parser.Instruction); ok {
		tmp2 = x86.DoASM(block.Value.(*parser.Instruction), x86.New())
	}
	fmt.Println(tmp, block.Value, tmp2, fmt.Sprintf("%x", tmp2))
	btmp = append(btmp, tmp2...)
	for _, k := range block.Children {
		pr(k, tabnum+1)
	}
	if tabnum == 0 {
		//	fmt.Printf("%b\n", btmp)
	}
}

func Compile(path string, archType string) {
	startTime := time.Now()
	fmt.Println("开始编译:", filepath.Base(path), "架构:", archType)
	arch := x86.New()
	lex := lexer.NewLexer(path)
	p := parser.NewParser(lex, arch)
	p.Parse()
	pr(p.Block, 0)
	// 创建指定架构的编译器
	compiler := compiler.NewCompiler(archType)
	res := compiler.Compile(p.Block)
	// 生成输出文件名
	outPath := path[:len(path)-len(filepath.Ext(path))] + "." + archType + ".asm"
	os.WriteFile(outPath, []byte(res), 0755)
	fmt.Println("编译完成 耗时" + time.Since(startTime).String())
}
