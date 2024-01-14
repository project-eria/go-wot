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

func NewString(defaultValue string, minLength *uint16, maxLength *uint16, pattern string) (Data, error) {
	var regexpPattern *regexp.Regexp
	var err error
	if pattern != "" {
		// Add "^" and "$" to the pattern, for global matching
		if pattern[0] != '^' {
			pattern = "^" + pattern
		}
		if pattern[len(pattern)-1] != '$' {
			pattern = pattern + "$"
		}
		if regexpPattern, err = regexp.Compile(pattern); err != nil {
			zlog.Error().Err(err).Msg("invalid pattern")
			return Data{}, errors.New("invalid pattern")
		}
	}
	d := Data{
		Default: defaultValue,
		Type:    "string",
		DataSchema: String{
			MinLength:     minLength,
			MaxLength:     maxLength,
			Pattern:       pattern,
			regexpPattern: regexpPattern,
		},
	}
	if err := d.Validate(d.Default); err != nil {
		return Data{}, errors.New("invalid default value: " + err.Error())
	}
	return d, nil
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
