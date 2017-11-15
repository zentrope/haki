//
// Copyright © 2017-present Keith Irwin
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published
// by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package interpreter

import (
	"errors"
	"fmt"
	"strings"
)

var builtins = map[string]primitiveFunc{
	"+":    _add,
	"-":    _minus,
	"prn":  _prn,
	"=":    _equals,
	"head": _head,
	"tail": _tail,
}

type primitiveFunc func(args []Expression) (Expression, error)

func isIntegral(val float64) bool {
	return val == float64(int64(val))
}

func _head(args []Expression) (Expression, error) {
	if len(args) == 0 {
		return NilExpression, errors.New("head requires a parameter")
	}

	if !args[0].IsList() {
		return NilExpression, errors.New("head requires a list parameter")
	}

	list := args[0].list
	if len(list) == 0 {
		return NilExpression, nil
	}
	return list[0], nil
}

func _tail(args []Expression) (Expression, error) {

	if len(args) == 0 {
		return NilExpression, errors.New("tail requires a parameter")
	}

	if !args[0].IsList() {
		return NilExpression, errors.New("tail requires a list parameter")
	}

	list := args[0].list

	return NewExpr(ExpList, list[1:]), nil
}

func _equals(args []Expression) (Expression, error) {
	// Return true if all the arguments are equal to each other in value
	// and type.

	if len(args) < 1 {
		return FalseExpression, errors.New("wrong number of args for '=', must be at least one")
	}

	sentinel := args[0]

	for _, a := range args[1:] {
		if !a.Equals(sentinel) {
			return FalseExpression, nil
		}
	}
	return TrueExpression, nil
}

func _add(args []Expression) (Expression, error) {
	var result float64
	for _, arg := range args {
		switch arg.tag {
		case ExpFloat:
			result = result + float64(arg.float)
		case ExpInteger:
			result = result + float64(arg.integer)
		default:
			return NilExpression, fmt.Errorf("unknown argument type, int/float expected, got [%#v]", arg)
		}
	}

	if isIntegral(result) {
		return NewExpr(ExpInteger, int64(result)), nil
	}
	return NewExpr(ExpFloat, result), nil
}

func _minus(args []Expression) (Expression, error) {
	if len(args) < 1 {
		return NilExpression, errors.New("'-' requires 1 or more args")
	}

	var result float64

	switch args[0].tag {
	case ExpFloat:
		result = float64(args[0].float)
	case ExpInteger:
		result = float64(args[0].integer)
	default:
		return NilExpression, fmt.Errorf("unknown argument type, int/float expected, got [%#v]", args[0])
	}

	if len(args) == 1 {
		return NewExpr(ExpFloat, -1.0*result), nil
	}

	for _, arg := range args[1:] {
		switch arg.tag {
		case ExpFloat:
			result = result - float64(arg.float)
		case ExpInteger:
			result = result - float64(arg.integer)
		default:
			return NilExpression, fmt.Errorf("unknown argument type, int/float expected, got [%#v]", arg)
		}
	}
	return NewExpr(ExpFloat, result), nil
}

func _prn(args []Expression) (Expression, error) {
	values := make([]string, 0)
	for _, a := range args {
		values = append(values, fmt.Sprintf("%v", a))
	}
	fmt.Println(strings.Join(values, " "))
	return NilExpression, nil
}
