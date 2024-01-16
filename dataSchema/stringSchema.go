package dataSchema

import (
	"errors"
	"regexp"

	zlog "github.com/rs/zerolog/log"
)

type String struct {
	MinLength *uint16 `json:"minLength,omitempty"` // (optional) Specifies the minimum length of a string.
	MaxLength *uint16 `json:"maxLength,omitempty"` // (optional) Specifies the maximum length of a string.package dataSchema
	Pattern   string  `json:"pattern,omitempty"`   // (optional) Provides a regular expressions to express constraints of the string value. The regular expression must follow the [ECMA-262] dialect.
	// TODO ContentEncoding  string `json:"contentEncoding,omitempty"`  // (optional) Specifies the encoding used to store the contents, as specified in RFC 2054. The values that are accepted: "7bit", "8bit", "binary", "quoted-printable" and "base64".
	// TODO ContentMediaType string `json:"contentMediaType,omitempty"` // (optional) Specifies the MIME type (e.g., image/png, audio/mpeg) of the contents of a string value, as described in RFC 2046.
	regexpPattern *regexp.Regexp `json:"-"`
}

func NewString(options ...StringOption) (Data, error) {
	opts := &StringOptions{
		Default: nil,
	}
	for _, option := range options {
		option(opts)
	}
	var regexpPattern *regexp.Regexp
	var err error

	if opts.Pattern != "" {
		// Add "^" and "$" to the pattern, for global matching
		if opts.Pattern[0] != '^' {
			opts.Pattern = "^" + opts.Pattern
		}
		if opts.Pattern[len(opts.Pattern)-1] != '$' {
			opts.Pattern = opts.Pattern + "$"
		}
		if regexpPattern, err = regexp.Compile(opts.Pattern); err != nil {
			zlog.Error().Err(err).Msg("invalid pattern")
			return Data{}, errors.New("invalid pattern")
		}
	}

	d := Data{
		Default: opts.Default,
		Unit:    opts.Unit,
		Type:    "string",
		DataSchema: String{
			MinLength:     opts.MinLength,
			MaxLength:     opts.MaxLength,
			Pattern:       opts.Pattern,
			regexpPattern: regexpPattern,
		},
	}
	if d.Default != nil {
		if err := d.Validate(d.Default); err != nil {
			return Data{}, errors.New("invalid default value: " + err.Error())
		}
	}
	return d, nil
}

type StringOption func(*StringOptions)

type StringOptions struct {
	Default   interface{}
	Unit      string
	MinLength *uint16
	MaxLength *uint16
	Pattern   string
}

func StringDefault(value string) StringOption {
	return func(opts *StringOptions) {
		opts.Default = value
	}
}

func StringUnit(unit string) StringOption {
	return func(opts *StringOptions) {
		opts.Unit = unit
	}
}

func StringMinLength(minLen uint16) StringOption {
	return func(opts *StringOptions) {
		opts.MinLength = &minLen
	}
}

func StringMaxLength(maxLen uint16) StringOption {
	return func(opts *StringOptions) {
		opts.MaxLength = &maxLen
	}
}

func StringPattern(pattern string) StringOption {
	return func(opts *StringOptions) {
		opts.Pattern = pattern
	}
}

func (s String) Validate(value interface{}) error {
	if _, ok := value.(string); !ok {
		return errors.New("incorrect string value type")
	}
	if s.MinLength != nil && uint16(len(value.(string))) < *s.MinLength {
		return errors.New("string too short")
	}
	if s.MaxLength != nil && uint16(len(value.(string))) > *s.MaxLength {
		return errors.New("string too long")
	}
	if s.regexpPattern != nil && !s.regexpPattern.MatchString(value.(string)) {
		return errors.New("string does not match pattern")
	}
	return nil
}

func (b String) FromString(value string) (interface{}, error) {
	return value, nil
}
