package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"./common"
	"./plan"
	"./recovery"
	"./sql/lex"
	"./sql/syntax"
)

var lexSwitch = flag.Bool("L", false, "Show the lex parser result.")
var syntaxSwitch = flag.Bool("S", false, "Show the syntax parser result.")
var performanceSwitch = flag.Bool("P", false, "Show the performance informations.")
var inReader = bufio.NewReader(os.Stdin)

func welcome() {
	fmt.Println("MonkeyDB2 @ 2016")
	fmt.Println("Welcome to SQL play ground~")
}

func loop() bool {
	fmt.Print("Monkey>>")
	command := ""

	i := 0
	for !strings.Contains(command, ";") {
		if i > 0 {
			fmt.Print("      ->")
		}
		cmd, _ := inReader.ReadString('\n')
		command += cmd
		i++
	}
	if strings.Contains(command, "quit;") {
		return false
	}

	ts, _ := lex.Parse(*lex.NewByteReader([]byte(command)))
	if *lexSwitch {
		fmt.Println(command)
		fmt.Println(ts)
	}
	stn, err := syntax.Parser(syntax.NewTokenReader(ts))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return true
	}
	if *syntaxSwitch {
		stn.Print(1)
	}
	r, re, err := plan.DirectPlan(stn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return true
	}
	r.Print()
	fmt.Println(re.AffectedRows, " row(s) affected in ", (float64)(re.UsedTime)/1000000000, "s.")
	return true
}

func bye() {
	recovery.SafeExit()
	fmt.Println("Bye!")
}

func main() {
	flag.Parse()
	if *performanceSwitch {
		common.PRINT2 = true
	}
	welcome()
	for loop() {
	}
	bye()
}
