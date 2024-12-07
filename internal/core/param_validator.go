package gocli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	gt "github.com/vcharco/gocli/internal/types"
	gu "github.com/vcharco/gocli/internal/utils"
)

func GetClosestCommand(candidates []gt.Candidate, command string) (gt.Candidate, error) {
	words := strings.Fields(command)

	if len(words) == 0 {
		return gt.Candidate{}, fmt.Errorf("empty command")
	}

	candidate := gt.Candidate{Name: "", Options: []gt.CandidateOption{}}
	for _, cmd := range candidates {
		if cmd.Name == gu.BestMatch(words[0], candidates) {
			candidate = cmd
			break
		}
	}

	if len(candidate.Name) == 0 {
		return gt.Candidate{}, fmt.Errorf("invalid command")
	}

	return candidate, nil
}

func ValidateCommand(candidates []gt.Candidate, command string) (gt.Candidate, map[string]string, error) {

	words := strings.Fields(command)

	if len(words) == 0 {
		return gt.Candidate{}, nil, fmt.Errorf("empty command")
	}

	candidate := gt.Candidate{Name: "", Options: []gt.CandidateOption{}}
	for _, cmd := range candidates {
		if cmd.Name == gu.BestMatch(words[0], candidates) {
			candidate = cmd
			break
		}
	}

	if len(candidate.Name) == 0 {
		return gt.Candidate{}, nil, fmt.Errorf("invalid command")
	}

	if len(words) == 1 {
		err := checkRequiredParams(map[string]string{}, candidate.Options)
		if err != nil {
			return gt.Candidate{}, nil, err
		}
		return candidate, map[string]string{}, nil
	}

	if len(words) > 1 && len(candidate.Options) == 0 {
		return candidate, nil, fmt.Errorf("parameters not supported for this command")
	}

	params, err := ValidateParams(candidate, words[1:])

	return candidate, params, err
}

func ValidateParams(candidate gt.Candidate, params []string) (map[string]string, error) {
	parsedParams := map[string]string{}

	if len(params) == 0 {
		parsedParams[candidate.Name] = ""
		return parsedParams, nil
	}

	defaultParam, err := getDefaultParam(candidate.Options)

	if err != nil {
		return nil, err
	}

	checkedDefaultParam := len(defaultParam.Name) == 0

	for i := 0; i < len(params); i++ {
		candidateOption, err := getCandidateOrError(params[i], candidate.Options)
		if err != nil {
			if !checkedDefaultParam {
				err := ValidateType(*defaultParam, params[i])
				if err != nil {
					return nil, err
				}
				parsedParams[defaultParam.Name] = params[i]
				checkedDefaultParam = true
				continue
			} else {
				return nil, fmt.Errorf("invalid parameter %v", params[i])
			}
		}
		if candidateOption.Type == gt.None {
			parsedParams[candidateOption.Name] = ""
		} else {
			if i+1 < len(params) {
				err := ValidateType(candidateOption, params[i+1])
				if err != nil {
					return nil, err
				}
				parsedParams[candidateOption.Name] = params[i+1]
				i++
			} else {
				return nil, fmt.Errorf("missing value")
			}
		}
	}

	err = checkRequiredParams(parsedParams, candidate.Options)

	if err != nil {
		return nil, err
	}

	return parsedParams, nil
}

func checkRequiredParams(params map[string]string, options []gt.CandidateOption) error {
	for _, option := range options {
		_, exists := params[option.Name]
		if !exists && option.Modifier&gt.REQUIRED != 0 {
			if option.Modifier&gt.DEFAULT != 0 {
				return fmt.Errorf("default parameter is required")
			} else {
				return fmt.Errorf("parameter %v is required", option.Name)
			}
		}
	}
	return nil
}

func getDefaultParam(params []gt.CandidateOption) (*gt.CandidateOption, error) {
	candidate := gt.CandidateOption{Name: ""}
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

func getCandidateOrError(param string, candidateList []gt.CandidateOption) (gt.CandidateOption, error) {
	for _, candidateOption := range candidateList {
		if candidateOption.Name == param {
			return candidateOption, nil
		}
	}
	return gt.CandidateOption{}, fmt.Errorf("candidate not found")
}

func ValidateType(candidateOption gt.CandidateOption, param string) error {
	switch candidateOption.Type {
	case gt.None:
		return nil
	case gt.Text:
		if param == "" {
			return fmt.Errorf("text cannot be empty")
		}
		return nil

	case gt.Number:
		if _, err := strconv.Atoi(param); err != nil {
			return fmt.Errorf("parameter %v must be a number", candidateOption.Name)
		}
		return nil

	case gt.Ipv4:
		re := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
		if !re.MatchString(param) {
			return fmt.Errorf("parameter %v must be an IPv4", candidateOption.Name)
		}
		return nil

	case gt.Ipv6:
		re := regexp.MustCompile(`([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}`)
		if !re.MatchString(param) {
			return fmt.Errorf("parameter %v must be an IPv6", candidateOption.Name)
		}
		return nil

	case gt.Email:
		re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !re.MatchString(param) {
			return fmt.Errorf("parameter %v must be an email address", candidateOption.Name)
		}
		return nil

	case gt.Domain:
		re := regexp.MustCompile(`^(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
		if !re.MatchString(param) {
			return fmt.Errorf("parameter %v must be a domain name", candidateOption.Name)
		}
		return nil

	case gt.Phone:
		re := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
		if !re.MatchString(param) {
			return fmt.Errorf("parameter %v must be a phone number", candidateOption.Name)
		}
		return nil

	case gt.Date:
		re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
		if !re.MatchString(param) {
			return fmt.Errorf("parameter %v must be a date (YYYY-MM-DD)", candidateOption.Name)
		}
		return nil

	case gt.Time:
		re := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)
		if !re.MatchString(param) {
			return fmt.Errorf("parameter %v must be a time (HH:mm)", candidateOption.Name)
		}
		return nil

	case gt.Url:
		re := regexp.MustCompile(`^https?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\\(\\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+$`)
		if !re.MatchString(param) {
			return fmt.Errorf("parameter %v must be a URL", candidateOption.Name)
		}
		return nil

	case gt.UUID:
		re := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89ab][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
		if !re.MatchString(param) {
			return fmt.Errorf("parameter %v must be a UUID (v4)", candidateOption.Name)
		}
		return nil

	default:
		return fmt.Errorf("parameter %v has a unrecognized type", candidateOption.Name)
	}
}

func GetValidationTypeName(val gt.CandidateType) string {
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
