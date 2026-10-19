package io

import (
	"encoding/xml"
	"fmt"
	"io"
	"iter"
)

func XmlAllTokens(rc io.ReadCloser) iter.Seq2[xml.Token, error] {
	return func(yield func(xml.Token, error) bool) {
		if rc == nil {
			yield(nil, fmt.Errorf("Invalid input"))
			return
		}
		defer rc.Close()

		decoder := xml.NewDecoder(rc)
		for {
			token, err := decoder.Token()
			if err == io.EOF {
				return
			} else if err != nil {
				yield(nil, err)
				return
			} else if !yield(token, err) {
				return
			}
		}
	}
}
