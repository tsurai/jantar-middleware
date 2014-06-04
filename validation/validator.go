package validation

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"time"
)

/* TODO: allow custom cookie name */

// Validator is a helper for validating user supplied data and returning error messages.
// It offers various validation functions and can save errors in a Cookie.
type Validator struct {
	rw        http.ResponseWriter
	hasErrors bool
	errors    map[string][]string
}

type validatorError struct {
	validation *Validator
	name       string
	index      int
}

func NewValidator(rw http.ResponseWriter) *Validator {
	return &Validator{rw, false, make(map[string][]string)}
}

// SaveErrors saves current validation error in a http.Cookie
func (v *Validator) SaveErrors() {
	if v.hasErrors {
		values := url.Values{}
		for key, array := range v.errors {
			for _, val := range array {
				values.Add(key, val)
			}
		}

		http.SetCookie(v.rw, &http.Cookie{Name: "JANTAR_ERRORS", Value: values.Encode(), Secure: false, HttpOnly: true, Path: "/"})
	}
}

// HasErrors returns true of an validation error occured. Otherwise false is returned
func (v *Validator) HasErrors() bool {
	return v.hasErrors
}

func (v *Validator) addError(name string, message string) *validatorError {
	result := &validatorError{v, name, -1}

	v.hasErrors = true
	v.errors[name] = append(v.errors[name], message)
	result.index = len(v.errors[name]) - 1

	return result
}

// Required checks the existance of given obj. How exactly this check is being performed depends on the type of obj. Valid types are: int, string, time.Time and slice.
// The given name determines the association of this error in the resulting validation error map.
func (v *Validator) Required(name string, obj interface{}) *validatorError {
	valid := false
	defaultMessage := "Required"

	switch value := obj.(type) {
	case nil:
		valid = false
	case int:
		valid = value != 0
	case string:
		valid = len(value) > 0
	case time.Time:
		valid = value.IsZero()
	default:
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Slice {
			valid = v.Len() > 0
		}
	}

	if !valid {
		return v.addError(name, defaultMessage)
	}

	return nil
}

// TODO: add time.Time to Min

// Min checks if given obj is smaller or equal to min. How exactly this check is being performed depends on the type of obj. Valid types are: int, string and slice.
// The given name determines the association of this error in the resulting validation error map.
func (v *Validator) Min(name string, obj interface{}, min int) *validatorError {
	valid := false
	defaultMessage := fmt.Sprintf("Must be larger than %d", min)

	switch value := obj.(type) {
	case nil:
		valid = false
	case int:
		valid = value >= min
	case string:
		valid = len(value) >= min
	default:
		v := reflect.ValueOf(obj)
		if v.Kind() == reflect.Slice {
			valid = v.Len() >= min
		}
	}

	if !valid {
		return v.addError(name, defaultMessage)
	}

	return nil
}

// TODO: add time.Time to Max

// Max checks if given obj is smaller or equal max. How exactly this check is being performed depends on the type of obj. Valid types are: int, string and slice.
// The given name determines the association of this error in the resulting validation error map.
func (v *Validator) Max(name string, obj interface{}, max int) *validatorError {
	valid := false
	defaultMessage := fmt.Sprintf("Must be smaller than %d", max)

	switch value := obj.(type) {
	case nil:
		valid = false
	case int:
		valid = value <= max
	case string:
		valid = len(value) <= max
	default:
		v := reflect.ValueOf(obj)
		if v.Kind() == reflect.Slice {
			valid = v.Len() <= max
		}
	}

	if !valid {
		return v.addError(name, defaultMessage)
	}

	return nil
}

// TODO: add time.Time to MinMax

// MinMax compiles Min and Max in one call.
func (v *Validator) MinMax(name string, obj interface{}, min int, max int) *validatorError {
	valid := false
	defaultMessage := fmt.Sprintf("Must be larger %d and smaller %d", min, max)

	switch value := obj.(type) {
	case nil:
		valid = false
	case int:
		valid = value >= min && value <= max
	case string:
		valid = len(value) >= min && len(value) <= max
	default:
		v := reflect.ValueOf(obj)
		if v.Kind() == reflect.Slice {
			valid = v.Len() >= min && v.Len() <= max
		}
	}

	if !valid {
		return v.addError(name, defaultMessage)
	}

	return nil
}

// Length checks the exact length of obj. How exactly this check is being performed depends on the type of obj. Valid types are: int, string and slice.
// The given name determines the association of this error in the resulting validation error map.
func (v *Validator) Length(name string, obj interface{}, length int) *validatorError {
	valid := false
	defaultMessage := fmt.Sprintf("Must be %d symbols long", length)

	switch value := obj.(type) {
	case nil:
		valid = false
	case int:
		valid = int(math.Ceil(math.Log10(float64(value)))) == length
	case string:
		valid = len(value) == length
	default:
		v := reflect.ValueOf(obj)
		if v.Kind() == reflect.Slice {
			valid = v.Len() == length
		}
	}

	if !valid {
		return v.addError(name, defaultMessage)
	}

	return nil
}

// Equals tests for deep equality of two given objects.
// The given name determines the association of this error in the resulting validation error map.
func (v *Validator) Equals(name string, obj interface{}, obj2 interface{}) *validatorError {
	defaultMessage := fmt.Sprintf("%v does not equal %v", obj, obj2)

	if obj == nil || obj2 == nil || !reflect.DeepEqual(obj, obj2) {
		return v.addError(name, defaultMessage)
	}

	return nil
}

func (v *Validator) MatchRegex(name string, obj interface{}, pattern string) *validatorError {
	valid := true
	defaultMessage := fmt.Sprintf("Must match regex %s", pattern)

	if obj == nil {
		valid = false
	} else {
		match, err := regexp.MatchString(pattern, reflect.ValueOf(obj).String())
		if err != nil || !match {
			valid = false
		}
	}

	if !valid {
		return v.addError(name, defaultMessage)
	}

	return nil
}

func (v *Validator) Custom(name string, match bool, message string) *validatorError {
	if match {
		return v.addError(name, message)
	}

	return nil
}

func (vr *validatorError) Message(msg string) *validatorError {
	if vr != nil && vr.index != -1 {
		vr.validation.errors[vr.name][vr.index] = msg
	}

	return vr
}
