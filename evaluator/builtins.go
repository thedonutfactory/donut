package evaluator

import (
	"github.com/thedonutfactory/donutbox/object"
)

var builtinFunctions = map[string]*object.Builtin{
	"len":   object.GetBuiltinByName("len"),
	"print": object.GetBuiltinByName("print"),
	"first": object.GetBuiltinByName("first"),
	"last":  object.GetBuiltinByName("last"),
	"rest":  object.GetBuiltinByName("rest"),
	"push":  object.GetBuiltinByName("push"),
	"pop":   object.GetBuiltinByName("pop"),
	"split": object.GetBuiltinByName("split"),
	"join":  object.GetBuiltinByName("join"),
}
