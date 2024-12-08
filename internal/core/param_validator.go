package gocli

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	gt "github.com/vcharco/gocli/internal/types"
	gu "github.com/vcharco/gocli/internal/utils"
)

func GetClosestCommand(candidates []gt.Command, command string) (gt.Command, error) {
	words := strings.Fields(command)

	if len(words) == 0 {
		return gt.Command{}, fmt.Errorf("empty command")
	}

	candidate := gt.Command{Name: "", Params: []gt.Param{}}
	for _, cmd := range candidates {
		bestMatch, _ := gu.BestMatch(words[0], candidates)
		if cmd.Name == bestMatch {
			candidate = cmd
			break
		}
	}

	if len(candidate.Name) == 0 {
		return gt.Command{}, fmt.Errorf("invalid command")
	}

	return candidate, nil
}

func ValidateCommand(candidates []gt.Command, command string) (gt.Command, map[string]interface{}, error) {

	words := strings.Fields(command)

	if len(words) == 0 {
		return gt.Command{}, nil, errors.New("empty command")
	}

	candidate := gt.Command{Name: "", Params: []gt.Param{}}
	for _, cmd := range candidates {
		bestMatch, _ := gu.BestMatch(words[0], candidates)
		if cmd.Name == bestMatch {
			candidate = cmd
			break
		}
	}

	if len(candidate.Name) == 0 {
		return gt.Command{}, nil, errors.New("invalid command")
	}

	if len(words) == 1 {
		err := checkRequiredParams(nil, candidate.Params)
		if err != nil {
			return gt.Command{}, nil, err
		}
		return candidate, nil, nil
	}

	if len(words) > 1 && len(candidate.Params) == 0 {
		return candidate, nil, fmt.Errorf("parameters not supported for this command")
	}

	params, err := ValidateParams(candidate, words[1:])

	return candidate, params, err
}

func ValidateParams(candidate gt.Command, inputParams []string) (map[string]interface{}, error) {
	var parsedParams map[string]interface{}

	if len(inputParams) == 0 {
		parsedParams[candidate.Name] = nil
		return parsedParams, nil
	}

	defaultParam, err := getDefaultParam(candidate.Params)

	if err != nil {
		return nil, err
	}

	checkedDefaultParam := len(defaultParam.Name) == 0

	for i := 0; i < len(inputParams); i++ {
		param, err := getParamOrError(inputParams[i], candidate.Params)
		if err != nil {
			if !checkedDefaultParam {
				err := ValidateType(*defaultParam, inputParams[i])
				if err != nil {
					return nil, err
				}
				parsedParams[defaultParam.Name] = inputParams[i]
				checkedDefaultParam = true
				continue
			} else {
				return nil, fmt.Errorf("invalid parameter %v", inputParams[i])
			}
		}
		if param.Type == gt.None {
			parsedParams[param.Name] = nil
		} else {
			if i+1 < len(inputParams) {
				err := ValidateType(param, inputParams[i+1])
				if err != nil {
					return nil, err
				}
				if param.Type == gt.Number {
					toInt, err := strconv.Atoi(inputParams[i+1])
					if err != nil {
						return nil, errors.New("error when casting the Numeric param to an Integer")
					}
					parsedParams[param.Name] = toInt
				} else if param.Type == gt.FloatNumber {
					toFloat, err := strconv.ParseFloat(inputParams[i+1], 64)
					if err != nil {
						return nil, errors.New("error when casting the Numeric param to a Float")
					}
					parsedParams[param.Name] = toFloat
				} else {
					parsedParams[param.Name] = inputParams[i+1]
				}
				i++
			} else {
				return nil, errors.New("missing value")
			}
		}
	}

	err = checkRequiredParams(parsedParams, candidate.Params)

	if err != nil {
		return nil, err
	}

	return parsedParams, nil
}

func checkRequiredParams(parsedParams map[string]interface{}, params []gt.Param) error {
	for _, param := range params {
		_, exists := parsedParams[param.Name]
		if !exists && param.Modifier&gt.REQUIRED != 0 {
			if param.Modifier&gt.DEFAULT != 0 {
				return fmt.Errorf("default parameter <%v> is required", GetValidationTypeName(param.Type))
			} else {
				return fmt.Errorf("parameter %v is required", param.Name)
			}
		}
	}
	return nil
}

func getDefaultParam(params []gt.Param) (*gt.Param, error) {
	candidate := gt.Param{Name: ""}
	for _, param := range params {
		if param.Modifier&gt.DEFAULT != 0 {
			if len(candidate.Name) > 0 {
				return nil, fmt.Errorf("cannot exist more than one default param")
			}
			candidate = param
		}
	}
	return &candidate, nil
}

func getParamOrError(param string, params []gt.Param) (gt.Param, error) {
	for _, p := range params {
		if p.Name == param {
			return p, nil
		}
	}
	return gt.Param{}, fmt.Errorf("candidate not found")
}

func ValidateType(param gt.Param, inputParam string) error {
	switch param.Type {
	case gt.None:
		return nil
	case gt.Text:
		if inputParam == "" {
			return fmt.Errorf("text cannot be empty")
		}
		return nil

	case gt.Number:
		if _, err := strconv.Atoi(inputParam); err != nil {
			return fmt.Errorf("parameter %v must be a number", param.Name)
		}
		return nil

	case gt.FloatNumber:
		if _, err := strconv.ParseFloat(inputParam, 64); err != nil {
			return fmt.Errorf("parameter %v must be a float number", param.Name)
		}
		return nil

	case gt.Ipv4:
		re := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
		if !re.MatchString(inputParam) {
			return fmt.Errorf("parameter %v must be an IPv4", param.Name)
		}
		return nil

	case gt.Ipv6:
		re := regexp.MustCompile(`([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}`)
		if !re.MatchString(inputParam) {
			return fmt.Errorf("parameter %v must be an IPv6", param.Name)
		}
		return nil

	case gt.Email:
		re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !re.MatchString(inputParam) {
			return fmt.Errorf("parameter %v must be an email address", param.Name)
		}
		return nil

	case gt.Domain:
		re := regexp.MustCompile(`^(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
		if !re.MatchString(inputParam) {
			return fmt.Errorf("parameter %v must be a domain name", param.Name)
		}
		return nil

	case gt.Phone:
		re := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
		if !re.MatchString(inputParam) {
			return fmt.Errorf("parameter %v must be a phone number", param.Name)
		}
		return nil

	case gt.Date:
		re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
		if !re.MatchString(inputParam) {
			return fmt.Errorf("parameter %v must be a date (YYYY-MM-DD)", param.Name)
		}
		return nil

	case gt.Time:
		re := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)
		if !re.MatchString(inputParam) {
			return fmt.Errorf("parameter %v must be a time (HH:mm)", param.Name)
		}
		return nil

	case gt.Url:
		re := regexp.MustCompile(`^https?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\\(\\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+$`)
		if !re.MatchString(inputParam) {
			return fmt.Errorf("parameter %v must be a URL", param.Name)
		}
		return nil

	case gt.UUID:
		re := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89ab][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
		if !re.MatchString(inputParam) {
			return fmt.Errorf("parameter %v must be a UUID (v4)", param.Name)
		}
		return nil

	default:
		return fmt.Errorf("parameter %v has a unrecognized type", param.Name)
	}
}

func GetValidationTypeName(val gt.ParamType) string {
	switch val {
	case gt.None:
		return "None"
	case gt.Date:
		return "Date"
	case gt.Domain:
		return "Domain"
	case gt.Email:
		return "Email"
	case gt.Ipv4:
		return "Ipv4"
	case gt.Ipv6:
		return "Ipv6"
	case gt.Number:
		return "Number"
	case gt.FloatNumber:
		return "FloatNumber"
	case gt.Phone:
		return "Phone"
	case gt.Text:
		return "Text"
	case gt.Time:
		return "Time"
	case gt.Url:
		return "Url"
	case gt.UUID:
		return "UUID"
	}

	return ""
}
