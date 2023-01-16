package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	Mod "main/Modules/SkyLine_Backend"
	"os"
	"path/filepath"
)

func init() {
	flag.Parse()
}

func main() {
	if *Mod.Bnn {
		Mod.Banner()
	}
	if *Mod.SourceFile != "" {
		if fileext := filepath.Ext(*Mod.SourceFile); fileext == "css" {
			if err := Run(*Mod.SourceFile); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(0)
			}
		} else {
			fmt.Println("Woah there buddy, sorry but that file type is not allowed, files must end in (.csc) not -> ", fileext)
		}
	} else {
		Mod.Banner()
		Mod.Start(os.Stdin, os.Stdout)
		return
	}
}

func Run(filename string) error {
	f, x := os.Open(filename)
	if x != nil {
		fmt.Println("Sorry bro | the SkyLine interperter can not read the file, it may not exist -> ", x)
		os.Exit(0)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var line []string
	for scanner.Scan() {
		line = append(line, scanner.Text())
	}
	if line == nil {
		fmt.Println("Sorry bro | the SkyLine interperter can not run this file through the interpreter, tokens were empty -> ", filename)
	}
	data, x := ioutil.ReadFile(filename)
	if x != nil {
		fmt.Println("Sorry bro | the SkyLine interperter can not read the file, it may not exist -> ", x)
		os.Exit(1)
	}
	parser := Mod.New_Parser(Mod.LexNew(string(data)))
	program := parser.ParseProgram()
	if len(parser.ParserErrors()) > 0 {
		return errors.New(parser.ParserErrors()[0])
	}

	Env := Mod.NewEnvironment()
	result := Mod.Eval(program, Env)
	if _, ok := result.(*Mod.Nil); ok {
		return nil
	}
	_, x = io.WriteString(os.Stdout, result.Inspect()+"\n")
	defer func() {
		if x := recover(); x != nil {
			return nil
			fmt.Println("SkyLine response => null")
		}
	}
	return x
}
