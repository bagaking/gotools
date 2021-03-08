package strs

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"text/scanner"

	"github.com/bagaking/gotools/reflectool"
)

var (
	ErrKVStrInvalidToken = errors.New("kv-str invalid token")
	ErrKVStrReflectError = errors.New("kv-str reflect error")
)

// KVStr is a string containing a set of key value pairs
//
// By convention, KV strings are optionally concatenated with
// space-separated key=value pairs. Each key=value pair is
// separated by a comma (U+002C ',') or semicolon (U+003B ';')
//
// Each key is a non-empty string that satisfies the syntax of
// go identity. Value can have four types:
//
// numeric: 	Any integer, floating point, or complex number
//			  	that satisfies the go syntax
//				e.g. `a=1,b=2.,c=.3,d=01 ,e=10e2`
//				=> {a:1, b:2, c:.3, d:01, e:10e2}
//
// character: 	Any character enclosed by single-quotes `'`
//				e.g. `char='x'` => {char:x}
//
// boolean: 	Only the value of `false` means that the value
//				is false, in other cases you can just write key,
//				the `=value` part can be omitted
//				e.g. `switch,window` => {switch:true, window:true}
//
// string: 		Any character enclosed by back-quotes ```, or
//				any character concatenation to the next separator
//				that is not one of the above cases, where
//				consecutive whitespace characters are treated as
//				a single whitespace character
//				e.g. `str=this   is a string,str2=`quoted string``
//				=> {str:"this is a string", str2="`quoted string`"}
type KVStr string

// ToMap will return the result of converting kv str to map
func (kv KVStr) ToMap() (map[string]string, error) {
	ret := make(map[string]string)
	if err := kv.ForEach(func(key, val string) {
		ret[key] = val
	}); err != nil {
		return nil, err
	}
	return ret, nil
}

// ForEach iterates over all key-value pairs
func (kv KVStr) ForEach(fn func(key, val string)) error {
	reader := strings.NewReader(string(kv))
	scanN := scanner.Scanner{}
	// scanN.Mode = scanner.ScanStrings | scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats
	scanN.Init(reader)

	for {
		tokenType := scanN.Scan()
		identKey := scanN.TokenText()
		if tokenType != scanner.Ident {
			return fmt.Errorf("%w, should be a key identity, token=`%s`(type=%v) @%s#%d", ErrKVStrInvalidToken,
				identKey, scanner.TokenString(tokenType), kv, scanN.Position.Column)
		}

		tokenType = scanN.Scan()
		switch tokenType {
		case '=':
			switch token := scanN.Scan(); token {
			case scanner.Ident, scanner.Char, scanner.String, scanner.Int, scanner.Float, scanner.RawString:
				value := scanN.TokenText()
				prevWord, prevPos := value, scanN.Position

				for next := scanN.Scan(); next != ',' && next != ';'; next = scanN.Scan() {
					if next == scanner.EOF {
						fn(identKey, value)
						return nil
					}
					if scanN.Offset-prevPos.Offset > len(prevWord) {
						prevWord = scanN.TokenText()
						value += " " + prevWord
						continue
					} else {
						prevWord = scanN.TokenText()
						value += prevWord
					}
					prevPos = scanN.Position
				}
				fn(identKey, value)
			default:
				return fmt.Errorf(
					"%w, should be a value identity, got token=`%s`(type=%v) @ %s #%d", ErrKVStrInvalidToken,
					scanN.TokenText(), scanner.TokenString(token), kv, scanN.Position.Column)
			}
		case ',', ';':
			fn(scanN.TokenText(), "")
			scanN.Next()
		case scanner.EOF:
			fn(scanN.TokenText(), "")
			scanN.Next()
			return nil
		default:
			return fmt.Errorf(
				"%w, should be a assign symbol,  token=`%s`(type=%v) @%s #%d", ErrKVStrInvalidToken,
				scanN.TokenText(), scanner.TokenString(tokenType), kv, scanN.Position.Column)
		}
	}
}

// ReflectTo copies kv by name to a struct
//
// The key used to assign a value should be equal to
// the name of the field corresponding to the target
// struct, or the snake form of its name.
//
// e.g. The field `FieldName` can be assigned by the
//		key `FieldName` or `field_name`.
func (kv KVStr) ReflectTo(target interface{}) (extra map[string]string, err error) {
	r := reflect.ValueOf(target)
	if r.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("%w, target must be an interface{}", ErrKVStrReflectError)
	}

	kvMap, err := kv.ToMap()
	if err != nil {
		return nil, fmt.Errorf("format error, %w", err)
	}

	var reflector reflectool.FieldHandler = func(field *reflect.Value, fieldType reflect.StructField) error {
		v, ok := kvMap[fieldType.Name]
		if !ok {
			snake := Conv2Snake(fieldType.Name)
			v, ok = kvMap[snake]
			if !ok {
				return nil
			}
		}

		converted, err := Conv2PlainType(v, field.Type())
		if err != nil {
			return err
		}

		val := reflect.ValueOf(converted)
		field.Set(val)
		delete(kvMap, fieldType.Name) // todo: check this
		return nil
	}

	if err := reflectool.ForEachField(target, reflector); err != nil {
		return nil, fmt.Errorf("reflect error, %w", err)
	}

	return kvMap, nil
}
