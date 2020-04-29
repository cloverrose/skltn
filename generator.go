package skltn

import (
	"go/format"
)

// Generate generates source code from template and method signature
func Generate(src []byte, template string, formatFlag bool) ([]byte, error) {
	target, err := newParser(src).parse()
	if err != nil {
		return nil, err
	}
	mft, err := renderTemplate(target, template)
	if err != nil {
		return nil, err
	}
	if !formatFlag {
		return mft, nil
	}
	ret, err := format.Source(mft)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
