package gocli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	gt "github.com/vcharco/gocli/internal/types"
)

func ValidateCommand(candidates []gt.Candidate, command string) (map[string]string, error) {

	words := strings.Fields(command)
	candidate := gt.Candidate{Name: "", DefaultOptionType: gt.None, Options: []gt.CandidateOption{}}
	for _, cmd := range candidates {
		if cmd.Name == words[0] {
			candidate = cmd
			break
		}
	}

	if len(candidate.Name) == 0 {
		return map[string]string{}, fmt.Errorf("invalid command")
	}

	return ValidateParams(candidate, words[1:])
}

func ValidateParams(candidate gt.Candidate, params []string) (map[string]string, error) {
	parsedParams := map[string]string{}

	if candidate.DefaultOptionType == gt.None && len(candidate.Options) == 0 {
		parsedParams[candidate.Name] = ""
		return parsedParams, nil
	}

	defaultParamAlreadyValidated := candidate.DefaultOptionType != gt.None
	for i := 0; i < len(params); i++ {
		candidateOption, err := getCandidateOrError(params[i], candidate.Options)
		if err != nil {
			if !defaultParamAlreadyValidated {
				candidateOption = gt.CandidateOption{Name: "default", Type: candidate.DefaultOptionType}
				err := ValidateType(candidateOption, params[i])
				if err != nil {
					return map[string]string{}, fmt.Errorf("type of default param value does not match")
				}
				parsedParams[candidateOption.Name] = params[i]
				defaultParamAlreadyValidated = true
			} else {
				return map[string]string{}, fmt.Errorf("too much default parameters or invalid parameters")
			}
		}
		if candidateOption.Type != gt.None {
			if i+1 < len(params)-1 {
				err := ValidateType(candidateOption, params[i+1])
				if err != nil {
					return map[string]string{}, fmt.Errorf("type of %v param value does not match", candidateOption.Name)
				}
				parsedParams[candidateOption.Name] = params[i+1]
				i++
			} else {
				return map[string]string{}, fmt.Errorf("missing parameter value")
			}
		} else {
			parsedParams[candidateOption.Name] = ""
		}
	}

	if !defaultParamAlreadyValidated {
		return map[string]string{}, fmt.Errorf("you must specify a default parameter")
	}

	return parsedParams, nil
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
	case gt.Text:
		if param == "" {
			return fmt.Errorf("text cannot be empty")
		}
		return nil

	case gt.Number:
		if _, err := strconv.Atoi(param); err != nil {
			return fmt.Errorf("invalid number")
		}
		return nil

	case gt.Ipv4:
		re := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
		if !re.MatchString(param) {
			return fmt.Errorf("invalid IPv4 address")
		}
		return nil

	case gt.Ipv6:
		re := regexp.MustCompile(`([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}`)
		if !re.MatchString(param) {
			return fmt.Errorf("invalid IPv6 address")
		}
		return nil

	case gt.Email:
		re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !re.MatchString(param) {
			return fmt.Errorf("invalid email address")
		}
		return nil

	case gt.Domain:
		re := regexp.MustCompile(`^(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
		if !re.MatchString(param) {
			return fmt.Errorf("invalid domain")
		}
		return nil

	case gt.Phone:
		re := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
		if !re.MatchString(param) {
			return fmt.Errorf("invalid phone number")
		}
		return nil

	case gt.Date:
		re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
		if !re.MatchString(param) {
			return fmt.Errorf("invalid date format (YYYY-MM-DD)")
		}
		return nil

	case gt.Time:
		re := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)
		if !re.MatchString(param) {
			return fmt.Errorf("invalid time format (HH:mm)")
		}
		return nil

	case gt.Url:
		re := regexp.MustCompile(`^https?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\\(\\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+$`)
		if !re.MatchString(param) {
			return fmt.Errorf("invalid URL")
		}
		return nil

	case gt.UUID:
		re := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89ab][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
		if !re.MatchString(param) {
			return fmt.Errorf("invalid UUID")
		}
		return nil

	default:
		return fmt.Errorf("unknown type")
	}
}
