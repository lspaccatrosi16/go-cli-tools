package parser

import (
	"fmt"
	"strings"
)

// takes in a string with formatting directives, and finds them from a map[string]interface{} of values until it reaches the value

func interfaceCheck[T any](v *interface{}) (T, bool) {
	val, ok := (*v).(T)
	return val, ok
}

func getValue(directive string, data Object) (interface{}, []error) {
	// fmt.Printf("parse: %s\n", directive)

	directiveComponents := strings.Split(directive, ".")
	parseErrors := []error{}

	val, exists := data[directiveComponents[0]]

	if !exists {
		parseErrors = append(parseErrors, fmt.Errorf("field %s not found in object", directive))
		return nil, parseErrors
	}

	// fmt.Printf("val: %#v\n", val)

	// has ptr reciever
	asObject, okAsObj := interfaceCheck[Object](&val)

	// fmt.Printf("check: %t %#v\n", okAsObj, asObject)

	if okAsObj {
		newDirective := strings.Join(directiveComponents[1:], ".")

		val, errs := getValue(newDirective, asObject)

		if len(errs) != 0 {
			recievedErr := fmt.Errorf("error recieved whilst parsing %s", directive)
			parseErrors = append(parseErrors, recievedErr)

			parseErrors = append(parseErrors, errs...)

			return nil, parseErrors
		}

		return val, errs
	}

	return val, nil
}

func isLetterOrPoint(ch byte) bool {
	return (ch >= 0b01000001 && ch <= 0b01011010) || (ch >= 0b01100001 && ch <= 0b01111010) || ch == 0b00101110
}

func ParseString(str string, data Object) (string, []error) {
	newString := []byte{}

	parseErrors := []error{}

	for i := 0; i < len(str); {
		ch := str[i]

		if ch != '$' {
			newString = append(newString, ch)
			i++
		} else {
			i++

			directive := []byte{}

			for {
				if i == len(str) {
					break
				}
				ch := str[i]

				if !isLetterOrPoint(ch) {
					break
				}

				directive = append(directive, ch)
				i++
			}

			directiveStr := string(directive)

			valueAtDirective, errs := getValue(directiveStr, data)
			if errs != nil {
				parseErrors = append(parseErrors, errs...)
			} else {
				valAsStr := fmt.Sprintf("%v", valueAtDirective)
				arr := []byte(valAsStr)
				newString = append(newString, arr...)
			}

		}

	}

	return string(newString), parseErrors
}
