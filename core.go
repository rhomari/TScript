package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"syscall/js"
)

func registerCallbacks() {
	js.Global().Set("Execute", js.FuncOf(Execute))
}
func main() {

	registerCallbacks()

	select {}
}

func Execute(this js.Value, args []js.Value) interface{} {

	document := js.Global().Get("document")

	tscripttag := document.Call("querySelector", "script[language=tscript]")

	if fmt.Sprintf("%s", tscripttag) != "<null>" {
		if tscripttag.Get("src").String() != "" {
			go func() {
				resp, err := http.Get(tscripttag.Get("src").String())
				if err != nil {
					fmt.Println(err)
				}
				script, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
				}
				processCode(string(script))

			}()

		} else {
			processCode(tscripttag.Get("innerText").String())
		}
	} else {
		fmt.Println("tscript tag not found")
	}

	return nil

}
func processCode(script string) {
	lines := strings.Split(script, "\n")
	for _, line := range lines {
		insts := getParts(line)
		for _, inst := range insts {
			for _, method := range inst.methods {

				_method, ok := Functions[method.name]
				if ok {
					_method.(func(string, string, ...string))(inst.object, method.ret, method.args...)
				} else {
					fmt.Println("function not found", method.name)
				}
			}

		}

	}
}

type Method struct {
	name string
	args []string
	ret  string
}
type Instruction struct {
	object  string
	methods []Method
}

func getParts(line string) []Instruction {
	if line == "" {
		return []Instruction{}
	}
	operands := strings.Split(line, "=")
	var leftop, rightop string
	if len(operands) > 1 {
		leftop = operands[0]
		rightop = operands[1]
	} else {
		rightop = line

	}

	rightop = strings.TrimSpace(rightop)

	instructions := []Instruction{}

	lineparts := strings.Split(rightop, ".")
	object := lineparts[0]
	var methods []Method
	for i := 1; i < len(lineparts); i++ {
		methodparts := strings.Split(lineparts[i], "(")
		args := strings.Split(strings.TrimRight(methodparts[1], ")"), ",")
		if leftop != "" {
			methods = append(methods, Method{name: methodparts[0], args: args, ret: leftop})
		} else {
			methods = append(methods, Method{name: methodparts[0], args: args, ret: ""})
		}

	}
	instructions = append(instructions, Instruction{object, methods})

	return instructions

}
