package handler

import (
	"strings"

	"github.com/go-graphite/carbonapi/pkg/parser"
)

type functionType string
type collectionOfFunctions map[string]functionDescription

type targeter interface {
	getCollection() collectionOfFunctions
	getDescription(funcName string) (functionDescription, bool)
}

type targetFunctionsHandler struct {
	functions collectionOfFunctions
}

func (target targetFunctionsHandler) getDescription(funcName string) (functionDescription, bool) {
	description, ok := target.functions[funcName]
	return description, ok
}

func (target targetFunctionsHandler) getCollection() collectionOfFunctions {
	return target.functions
}

const (
	functionIsWarn functionType = "warn"
	functionIsBad  functionType = "bad"
)

type functionDescription struct {
	Type        functionType `json:"type"`
	Description string       `json:"description"`
}

type FunctionsOfTarget struct {
	SyntaxOk     bool                  `json:"syntax_ok"`
	BadFunctions collectionOfFunctions `json:"bad_functions"`
}

var badFunctions targeter = targetFunctionsHandler{
	functions: collectionOfFunctions{
		"summarize": {Type: functionIsWarn, Description: "Потому что нельзя просто так взять и поставить эту функцию"},
		"bad":       {Type: functionIsBad, Description: "Не используйте плохую функцию"},
	},
}

func targetsHandler(targets []string) *FunctionsOfTarget {
	for _, target := range targets {
		response := targetVerification(target)
		if !response.SyntaxOk {
			return &response
		}

		for _, description := range response.BadFunctions {
			if description.Type == functionIsBad {
				return &response
			}
		}
	}

	return nil
}

func targetVerification(target string) FunctionsOfTarget {
	expr, _, err := parser.ParseExpr(target)
	if err != nil {
		return FunctionsOfTarget{SyntaxOk: false}
	}

	targetResponse := checkExpression(expr)
	targetResponse.SyntaxOk = true

	return targetResponse
}

func checkExpression(expression parser.Expr) FunctionsOfTarget {
	var response FunctionsOfTarget

	if description, ok := badFunctions.getDescription(strings.ToLower(expression.Target())); ok {
		if response.BadFunctions == nil {
			response.BadFunctions = make(collectionOfFunctions)
		}

		response.BadFunctions[expression.Target()] = description
	}

	for _, expr := range expression.Args() {
		targetResponse := checkExpression(expr)
		if targetResponse.BadFunctions != nil {
			for functionName, description := range targetResponse.BadFunctions {
				response.BadFunctions[functionName] = description
			}
		}

		if !targetResponse.SyntaxOk {
			response.SyntaxOk = targetResponse.SyntaxOk
		}
	}

	return response
}
