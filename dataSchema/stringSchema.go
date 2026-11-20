package dataSchema

import (
	"errors"
	"regexp"

	zlog "github.com/rs/zerolog/log"
)

type String struct {
	MinLength        *uint16        `json:"minLength,omitempty"`        // (optional) Specifies the minimum length of a string.
	MaxLength        *uint16        `json:"maxLength,omitempty"`        // (optional) Specifies the maximum length of a string.package dataSchema
	Pattern          string         `json:"pattern,omitempty"`          // (optional) Provides a regular expressions to express constraints of the string value. The regular expression must follow the [ECMA-262] dialect.
	ContentEncoding  string         `json:"contentEncoding,omitempty"`  // (optional) Specifies the encoding used to store the contents, as specified in RFC 2054. The values that are accepted: "7bit", "8bit", "binary", "quoted-printable" and "base64".
	ContentMediaType string         `json:"contentMediaType,omitempty"` // (optional) Specifies the MIME type (e.g., image/png, audio/mpeg) of the contents of a string value, as described in RFC 2046.
	regexpPattern    *regexp.Regexp `json:"-"`
}

func NewString(options ...StringOption) (Data, error) {
	opts := &StringOptions{
		Default: nil,
	}
	for _, option := range options {
		if option != nil {
			option(opts)
		}
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
		Enum:    opts.Enum,
		OneOf:   opts.OneOf,
		DataSchema: String{
			MinLength:        opts.MinLength,
			MaxLength:        opts.MaxLength,
			Pattern:          opts.Pattern,
			ContentEncoding:  opts.ContentEncoding,
			ContentMediaType: opts.ContentMediaType,
			regexpPattern:    regexpPattern,
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
	Enum      []interface{}
	OneOf     []ConstSchema
	MinLength *uint16
	MaxLength *uint16
	Pattern   string

	ContentEncoding  string
	ContentMediaType string
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

func StringEnum(enum []string) StringOption {
	return func(opts *StringOptions) {
		opts.Enum = []interface{}{}
		for _, e := range enum {
			opts.Enum = append(opts.Enum, e)
		}
	}
}

// StringOneOf restricts the values to the given ones, each with a title.
func StringOneOf(values ...ConstSchema) StringOption {
	return func(opts *StringOptions) {
		opts.OneOf = values
	}
}

// StringContentMediaType declares the media type of the content (image/jpeg,
// text/html…). A media type other than JSON makes the HTTP binding serve the
// value as is, with that Content-Type, instead of a JSON string.
func StringContentMediaType(mediaType string) StringOption {
	return func(opts *StringOptions) {
		opts.ContentMediaType = mediaType
	}
}

// StringContentEncoding declares how binary content is stored in the string.
// With "base64", handlers return a []byte: JSON encodes it as base64 and the
// HTTP binding serves the raw bytes.
func StringContentEncoding(encoding string) StringOption {
	return func(opts *StringOptions) {
		opts.ContentEncoding = encoding
	}
}

func (s String) Validate(value interface{}) error {
	if b, ok := value.([]byte); ok && s.ContentEncoding != "" {
		value = string(b) // binary content: the length constraints apply to the bytes
	}
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
