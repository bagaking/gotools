package procast

import "fmt"

func Recover(handler func(err error), format string, args ...interface{}) {
	str := format
	if len(args) > 0 {
		str = fmt.Sprintf(format, args...)
	}

	if r := recover(); r != nil {
		if v, ok := r.(error); ok {
			if str == "" {
				handler(fmt.Errorf("%w", v))
				return
			}
			handler(fmt.Errorf("%s, %w", str, v))
		} else {
			if str == "" {
				handler(fmt.Errorf("%v", r))
				return
			}
			handler(fmt.Errorf("%s, %v", str, r))
		}
	}
}
