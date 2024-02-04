package nodes

import (
	"regexp"

	"github.com/hackborn/onefunc/pipeline"
)

var (
	regexpOperations = map[string]regexpOperationFn{
		"": func(re *regexp.Regexp, s string, n *RegexpNode) (string, error) {
			s = re.ReplaceAllString(s, n.Replace)
			return s, nil
		},
		"replace": func(re *regexp.Regexp, s string, n *RegexpNode) (string, error) {
			s = re.ReplaceAllString(s, n.Replace)
			return s, nil
		},
	}

	regexpTargets = map[string]regexpTargetFn{
		"content.name": func(pin any, fn regexpOperationFn, re *regexp.Regexp, n *RegexpNode) error {
			if cd, ok := pin.(*pipeline.ContentData); ok {
				s, err := fn(re, cd.Name, n)
				if err != nil {
					return err
				}
				cd.Name = s
			}
			return nil
		},
	}
)
