package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type UserRole string

var (
	ErrValidationLengthString      = errors.New("validation error length string")
	ErrValidationLengthSliceString = errors.New("validation error length slice string")
	ErrValidationRegexp            = errors.New("validation error regexp")
	ErrValidationMin               = errors.New("validation error min")
	ErrValidationMax               = errors.New("validation error max")
	ErrValidationIn                = errors.New("validation error in")
)

type Validator interface {
	Validate(fieldName string, v any) (*ValidationError, error)
}

func validatorFactory(name string, constraint string) (Validator, error) {
	switch name {
	case "len":
		i, err := strconv.Atoi(constraint)
		if err != nil {
			return nil, err
		}
		return NewLength(i), nil
	case "min":
		i, err := strconv.Atoi(constraint)
		if err != nil {
			return nil, err
		}
		return NewMin(i), nil
	case "max":
		i, err := strconv.Atoi(constraint)
		if err != nil {
			return nil, err
		}
		return NewMax(i), nil
	case "regexp":
		return NewRegexp(strings.ReplaceAll(constraint, "\\\\", "\\")), nil
	case "in":
		return NewIn(strings.Split(constraint, ",")), nil

	default:
		return nil, fmt.Errorf("validator for %v doesn't exists", name)
	}
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	result := ""
	for _, err := range v {
		result += fmt.Sprintf("Field: %v, has error: %v. ", err.Field, err.Err.Error())
	}
	return strings.Trim(result, " ")
}

func Validate(v interface{}) (ValidationErrors, error) {
	var resultValidation ValidationErrors
	reflectionType := reflect.TypeOf(v)
	reflectionValue := reflect.ValueOf(v)

	if reflectionValue.Kind() != reflect.Struct {
		return nil, errors.New("input value myst be structure")
	}

	for i := 0; i < reflectionValue.NumField(); i++ {
		fieldV := reflectionValue.Field(i)
		fieldT := reflectionType.Field(i)
		tag := string(fieldT.Tag)
		validators := getValidators(tag)

		for _, v := range validators {
			res, err := v.Validate(fieldT.Name, fieldV.Interface())
			if err != nil {
				return nil, err
			}

			if res != nil {
				resultValidation = append(resultValidation, *res)
			}
		}
	}
	return resultValidation, nil
}

func getValidators(tag string) []Validator {
	tags := strings.Split(tag, " ")
	var result []Validator

	for _, tag := range tags {
		re := regexp.MustCompile(`validate:(.*)`)
		validatorsReResult := re.FindAllStringSubmatch(tag, -1)
		if validatorsReResult == nil || len(validatorsReResult[0]) < 2 {
			continue
		}

		validatorsStr := validatorsReResult[0][1]
		validatorsStr = strings.ReplaceAll(validatorsStr, "\"", "")

		validators := strings.Split(validatorsStr, "|")
		for _, validator := range validators {
			validatorParts := strings.Split(validator, ":")
			validatorStruct, err := validatorFactory(validatorParts[0], validatorParts[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
			result = append(result, validatorStruct)
		}
	}
	return result
}

type Length struct {
	Constraint int
}

func (l Length) Validate(fieldName string, v any) (*ValidationError, error) {
	switch vType := v.(type) {
	case string:
		if len(vType) != l.Constraint {
			return &ValidationError{
				Field: fieldName,
				Err:   errors.Wrapf(ErrValidationLengthString, "The length value must be %v. Given: %v", l.Constraint, len(vType)),
			}, nil
		}
	case []string:
		incorrectValues := make([]string, 0)
		for _, v := range v.([]string) {
			if len(v) != l.Constraint {
				incorrectValues = append(incorrectValues, v)
			}
		}

		if len(incorrectValues) > 0 {
			return &ValidationError{
				Field: fieldName,
				Err: errors.Wrapf(
					ErrValidationLengthSliceString,
					"The length values: %v must be %v", incorrectValues, l.Constraint,
				),
			}, nil
		}
	default:
		return nil, errors.New("unsupported type for validator Length")
	}
	return nil, nil
}

func NewLength(constraint int) Length {
	return Length{
		Constraint: constraint,
	}
}

type Min struct {
	Constraint int
}

func (m Min) Validate(fieldName string, v any) (*ValidationError, error) {
	switch vType := v.(type) {
	case int:
		if vType < m.Constraint {
			return &ValidationError{
				Field: fieldName,
				Err:   errors.Wrapf(ErrValidationMin, "Value: %v must be greater than: %v", vType, m.Constraint),
			}, nil
		}
	default:
		return nil, errors.New("unsupported type for validator Min")
	}
	return nil, nil
}

func NewMin(constraint int) Min {
	return Min{
		Constraint: constraint,
	}
}

type Max struct {
	Constraint int
}

func (m Max) Validate(fieldName string, v any) (*ValidationError, error) {
	switch vType := v.(type) {
	case int:
		if vType > m.Constraint {
			return &ValidationError{
				Field: fieldName,
				Err:   errors.Wrapf(ErrValidationMax, "Value: %v must be greater than: %v", vType, m.Constraint),
			}, nil
		}
	default:
		return nil, errors.New("unsupported type for validator Max")
	}
	return nil, nil
}

func NewMax(constraint int) Max {
	return Max{
		Constraint: constraint,
	}
}

type Regexp struct {
	Constraint string
}

func (r Regexp) Validate(fieldName string, v any) (*ValidationError, error) {
	switch vType := v.(type) {
	case string:
		matched, err := regexp.MatchString(r.Constraint, vType)
		if err != nil {
			return nil, err
		}
		if !matched {
			return &ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("%w Value: %v does not match the pattern: %v", ErrValidationRegexp, vType, r.Constraint),
			}, nil
		}
	default:
		return nil, errors.New("unsupported type for validator Regexp")
	}
	return nil, nil
}

func NewRegexp(constraint string) Regexp {
	return Regexp{
		Constraint: constraint,
	}
}

type In struct {
	Constraint []string
}

func (i In) Validate(fieldName string, v any) (*ValidationError, error) {
	if !slices.Contains(i.Constraint, fmt.Sprint(v)) {
		return &ValidationError{
			Field: fieldName,
			Err:   fmt.Errorf("%w Value: %v not in: %v", ErrValidationIn, fmt.Sprint(v), i.Constraint),
		}, nil
	}
	return nil, nil
}

func NewIn(constraint []string) In {
	return In{
		Constraint: constraint,
	}
}
