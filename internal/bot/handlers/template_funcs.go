package handlers

import (
	"reflect"
	"strconv"
	"text/template"

	"github.com/pkg/errors"
)

var templateFuncs = template.FuncMap{
	"sum": sum,
}

func sum(arg0 reflect.Value, args ...reflect.Value) (reflect.Value, error) {
	switch arg0.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sum := arg0.Int()
		for _, arg := range args {
			if err := addInt64(&sum, arg); err != nil {
				return reflect.Value{}, err
			}
		}
		return reflect.ValueOf(sum), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		sum := int64(arg0.Uint())
		if u, us := arg0.Uint(), uint64(sum); u < us {
			return reflect.Value{}, errors.New("overflow")
		}

		for _, arg := range args {
			if err := addInt64(&sum, arg); err != nil {
				return reflect.Value{}, err
			}
		}
		return reflect.ValueOf(sum), nil
	case reflect.String:
		arg0str := arg0.String()
		sum, err := strconv.ParseInt(arg0str, 10, 64)
		if err != nil {
			return reflect.Value{}, errors.Wrapf(err, "parse %q as int", arg0str)
		}
		for _, arg := range args {
			argInt, err := strconv.ParseInt(arg.String(), 10, 64)
			if err != nil {
				return reflect.Value{}, errors.Wrapf(err, "parse %q as int", arg.String())
			}
			if err := addInt64(&sum, reflect.ValueOf(argInt)); err != nil {
				return reflect.Value{}, err
			}
		}
		return reflect.ValueOf(sum), nil
	default:
		return reflect.Value{}, errors.Errorf("unsupported type: %s", arg0.Kind())
	}
}

func mul(arg0 reflect.Value, args ...reflect.Value) (reflect.Value, error) {
	switch arg0.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sum := arg0.Int()
		for _, arg := range args {
			if err := mulInt64(&sum, arg); err != nil {
				return reflect.Value{}, err
			}
		}
		return reflect.ValueOf(sum), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		sum := int64(arg0.Uint())
		if u, us := arg0.Uint(), uint64(sum); u < us {
			return reflect.Value{}, errors.New("overflow")
		}

		for _, arg := range args {
			if err := mulInt64(&sum, arg); err != nil {
				return reflect.Value{}, err
			}
		}
		return reflect.ValueOf(sum), nil
	case reflect.String:
		arg0str := arg0.String()
		sum, err := strconv.ParseInt(arg0str, 10, 64)
		if err != nil {
			return reflect.Value{}, errors.Wrapf(err, "parse %q as int", arg0str)
		}
		for _, arg := range args {
			argInt, err := strconv.ParseInt(arg.String(), 10, 64)
			if err != nil {
				return reflect.Value{}, errors.Wrapf(err, "parse %q as int", arg.String())
			}
			if err := mulInt64(&sum, reflect.ValueOf(argInt)); err != nil {
				return reflect.Value{}, err
			}
		}
		return reflect.ValueOf(sum), nil
	default:
		return reflect.Value{}, errors.Errorf("unsupported type: %s", arg0.Kind())
	}
}

func addInt64(acc *int64, val reflect.Value) error {
	var intVal int64
	switch val.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		intVal = int64(val.Uint())
		if u, us := val.Uint(), uint64(intVal); u < us {
			return errors.New("overflow")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal = val.Int()

	default:
		return errors.Errorf("unsupported type: %s", val.Kind())
	}

	before := *acc
	*acc += intVal
	if intVal > 0 && *acc < before {
		return errors.New("overflow")
	} else if intVal < 0 && *acc > before {
		return errors.New("underflow")
	}
	return nil
}

func mulInt64(acc *int64, val reflect.Value) error {
	var intVal int64
	switch val.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		intVal = int64(val.Uint())
		if u, us := val.Uint(), uint64(intVal); u < us {
			return errors.New("overflow")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal = val.Int()
	default:
		return errors.Errorf("unsupported type: %s", val.Kind())
	}
	before := *acc
	*acc *= intVal
	if intVal > 0 && *acc < before {
		return errors.New("overflow")
	} else if intVal < 0 && *acc > before {
		return errors.New("underflow")
	}
	return nil
}
