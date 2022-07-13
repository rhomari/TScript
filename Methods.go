package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"syscall/js"
)

var Variables = make(map[string]string)

// Methods defined in the tscript language
var Functions = map[string]interface{}{
	"setStyle": func(objectID string, ret string, args ...string) {

		doc := js.Global().Get("document")
		object := doc.Call("getElementById", objectID)
		style := object.Get("style")

		style.Call("setProperty", strings.Trim(args[0], "\""), strings.Trim(args[1], "\""))

	},
	"getStyle": func(objectID string, ret string, args ...string) {

		doc := js.Global().Get("document")
		object := doc.Call("getElementById", objectID)
		style := object.Get("style")

		if args[0] == "" {
			Variables[strings.Trim(ret, " ")] = style.Get("cssText").String()

		} else {

			cssText := style.Get("cssText").String()
			v, err := getProperty(cssText, strings.Trim(args[0], "\""))
			if err != nil {
				fmt.Println(err)
			}
			Variables[strings.Trim(ret, " ")] = v

		}

	},
	"addClass": func(objectID string, ret string, args ...string) {

		doc := js.Global().Get("document")
		object := doc.Call("getElementById", objectID).Get("classList")
		object.Call("add", strings.Trim(args[0], "\""))

	},
	"setHTML": func(objectID string, ret string, args ...string) {

		doc := js.Global().Get("document")
		object :=
			doc.Call("getElementById", objectID)
		if strings.Contains(args[0], "\"") {
			object.Set("innerHTML", strings.Trim(args[0], "\""))
		} else {
			object.Set("innerHTML", Variables[strings.Trim(args[0], " ")])
		}

	},
	"setText": func(objectID string, ret string, args ...string) {

		doc := js.Global().Get("document")
		object := doc.Call("getElementById", objectID)
		doc.Call("getElementById", objectID)
		if strings.Contains(args[0], "\"") {
			object.Set("innerText", strings.Trim(args[0], "\""))
		} else {
			object.Set("innerText", Variables[strings.Trim(args[0], " ")])
		}
	},
	"clickEvent": func(objectID string, ret string, args ...string) {

		doc := js.Global().Get("document")
		object := doc.Call("getElementById", objectID)
		object.Call("click")
	},
	"getHTML": func(objectID string, ret string, args ...string) {

		doc := js.Global().Get("document")
		object := doc.Call("getElementById", objectID)
		Variables[strings.Trim(ret, " ")] = object.Get("innerHTML").String()

	},
	"print": func(objectID string, ret string, args ...string) {

		fmt.Println(Variables[strings.Trim(args[0], " ")])
	},
}

func getProperty(cssText string, name string) (string, error) {

	css := strings.Split(cssText, ";")

	for _, p := range css {
		if strings.Contains(p, ":") {
			kv := strings.Split(p, ":")

			if strings.Trim(kv[0], " ") == name {

				return strings.TrimLeft(kv[1], " "), nil

			}

		}

	}
	return "", fmt.Errorf("CSS Property not found")
}
func eval(expression string) {
	fileset := token.NewFileSet()
	exp, err := parser.ParseExpr("10+6")
	if err == nil {
		ast.Print(fileset, exp)

	}

}
