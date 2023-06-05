// This package is inspired by this repo https://github.com/go-validator/validator.
// I implemented it myself mainly because I wanted to learn more
// about the reflect package.

// Copyright 2014 Roberto Teixeira <robteix@robteix.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ErrZeroValue     = errors.New("zero value")
	ErrInvalidRegexp = errors.New("invalid regexp")
	ErrRegexp        = errors.New("regexp mismatch")
	ErrUnknownTag    = errors.New("unknown tag")
	ErrNotSupported  = errors.New("not supported")
	ErrInvalidArg    = errors.New("invalid argument")
	ErrMaxLen        = errors.New("max length exceeded")
)

type ValidationFunc func(str, arg string) error

type Validator struct {
	tag   string
	funcs map[string]ValidationFunc
}

func New() Validator {
	return Validator{tag: "validate", funcs: map[string]ValidationFunc{
		"regexp":  ValidateRegexp,
		"nonzero": ValidateNonzero,
		"maxlen":  ValidateMaxLen,
	}}
}

type Tag struct {
	Name string
	Fn   ValidationFunc
	Arg  string
}

func NewTag(name string, fn ValidationFunc, arg string) Tag {
	return Tag{Name: name, Fn: fn, Arg: arg}
}

func (v *Validator) Validate(s interface{}) error {
	st := reflect.TypeOf(s)
	sv := reflect.ValueOf(&s).Elem().Elem() // wtf?

	for i := 0; i < st.NumField(); i++ {
		err := v.ValidateField(st.Field(i), sv.Field(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Validator) ValidateField(sf reflect.StructField, sv reflect.Value) error {
	tag := sf.Tag.Get(v.tag)
	if tag == "" {
		return nil
	} else if tag != "" && sf.Type.Kind() != reflect.String {
		return ErrNotSupported
	} else {
		tags, err := v.parseTags(tag)
		if err != nil {
			return err
		}
		for _, t := range tags {
			err := t.Fn(sv.String(), t.Arg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *Validator) parseTags(t string) ([]Tag, error) {
	var tags []Tag
	splitted := strings.Split(t, ",")
	for _, tagsRaw := range splitted {
		tag := strings.SplitN(tagsRaw, "=", 2)
		if fn, ok := v.funcs[tag[0]]; !ok {
			return []Tag{}, ErrUnknownTag
		} else {
			tags = append(tags, NewTag(tag[0], fn, tag[1]))
		}
	}
	return tags, nil
}

func ValidateRegexp(str, arg string) error {
	re, err := regexp.Compile(arg)
	if err != nil {
		return ErrInvalidRegexp
	}

	if !re.MatchString(str) {
		return ErrRegexp
	}
	return nil
}

func ValidateNonzero(str, arg string) error {
	valid := utf8.RuneCountInString(str) != 0
	if !valid {
		return ErrZeroValue
	}
	return nil
}

func ValidateMaxLen(str string, arg string) error {
	maxLen, err := strconv.Atoi(arg)
	if err != nil {
		return ErrInvalidArg
	}
	valid := utf8.RuneCountInString(str) <= maxLen
	if !valid {
		return ErrMaxLen
	}
	return nil
}
