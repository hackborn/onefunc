package nodes

import (
	"regexp"

	"github.com/hackborn/onefunc/pipeline"
)

var (
	regexpOperations = map[string]regexpOperationFn{
		"": func(re *regexp.Regexp, s string, data *regexpData) (string, error) {
			s = re.ReplaceAllString(s, data.Replace)
			return s, nil
		},
		"replace": func(re *regexp.Regexp, s string, data *regexpData) (string, error) {
			s = re.ReplaceAllString(s, data.Replace)
			return s, nil
		},
	}

	regexpTargets = map[string]regexpTargetFn{
		"content.name": func(pin any, fn regexpOperationFn, re *regexp.Regexp, data *regexpData) error {
			if cd, ok := pin.(*pipeline.ContentData); ok {
				s, err := fn(re, cd.Name, data)
				if err != nil {
					return err
				}
				cd.Name = s
			}
			return nil
		},
	}
)
