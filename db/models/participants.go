package models

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	attrs = []string{"pname", "pid", "pclub", "pweightCat", "pbelt", "pshirt", "extraReq"}
)

// Participant represents a participant in the checkout form.
type Participant struct {
	Suffix            string `bigquery:"-"`
	OrdinalPos        int    `bigquery:"-"`
	ItemID            int    `bigquery:"item_id"`
	Name              string `id:"pname" validate:"nonzero=true,maxlen=30" bigquery:"name"`
	ID                string `id:"pid" validate:"nonzero=true,maxlen=20" bigquery:"id"`
	Club              string `id:"pclub" validate:"maxlen=30" bigquery:"club"`
	WeightCategory    string `id:"pweightCat" validate:"maxlen=10" bigquery:"weight_category"`
	Belt              string `id:"pbelt" validate:"maxlen=20" bigquery:"belt"`
	ShirtSize         string `id:"pshirt" validate:"nonzero=true,maxlen=10" bigquery:"shirt_size"`
	ExtraRequirements string `id:"extraReq" validate:"maxlen=100" bigquery:"extra_requirements"`
}

// newParticipant returns a new Participant with the given suffix and form.
// The suffix is used to match the attributes in the form.
func newParticipant(suffix string, form url.Values) (Participant, error) {
	splittedSuffix := strings.Split(suffix, "-")
	if len(splittedSuffix) != 3 {
		return Participant{}, fmt.Errorf("models: invalid suffix %s", suffix)
	}

	ordinalPos, err := strconv.Atoi(splittedSuffix[1])
	if err != nil {
		return Participant{}, err
	}

	itemID, err := strconv.Atoi(splittedSuffix[2])
	if err != nil {
		return Participant{}, err
	}

	return Participant{
		Suffix:            suffix,
		OrdinalPos:        ordinalPos,
		ItemID:            itemID,
		Name:              form.Get("pname" + suffix),
		ID:                form.Get("pid" + suffix),
		Club:              form.Get("pclub" + suffix),
		WeightCategory:    form.Get("pweightCat" + suffix),
		Belt:              form.Get("pbelt" + suffix),
		ShirtSize:         form.Get("pshirt" + suffix),
		ExtraRequirements: form.Get("extraReq" + suffix),
	}, nil
}

// decodeParticipants returns a slice of Participants from the given form.
func decodeParticipants(form url.Values) ([]Participant, error) {
	var participants []Participant

	suffixes, err := attrSuffixes(form)
	if err != nil {
		return participants, err
	}

	for _, s := range suffixes {
		p, err := newParticipant(s, form)
		if err != nil {
			return participants, err
		}
		participants = append(participants, p)
	}
	return participants, nil
}

// attrSuffixes returns a slice of suffixes that match the prefix
// of the attributes in the form.
// Example:
//   - map[string][]string{"pname-1-1", "pweightCat-2-1"} -> []string{"1-1", "2-1"}
func attrSuffixes(form url.Values) ([]string, error) {
	var matchedAttrs, suffixes []string

	re := regexp.MustCompile("^(?P<Attr>[a-zA-Z]+)(?P<Suffix>-[0-9]+-[0-9]+)$")

	// get all attributes tha match the prefix
	for k := range form {
		if attrMatch(k, attrs) {
			matchedAttrs = append(matchedAttrs, k)
		}
	}

	// get all suffixes
	for _, t := range matchedAttrs {
		if re.MatchString(t) {
			matches := re.FindStringSubmatch(t)
			suffix := re.SubexpIndex("Suffix")
			suffixes = append(suffixes, matches[suffix])
		} else {
			return []string{}, fmt.Errorf("models: invalid attribute %s", t)
		}
	}
	// remove duplicates
	m := make(map[string]struct{})

	for _, s := range suffixes {
		m[s] = struct{}{}
	}

	// reset suffixes
	suffixes = []string{}

	for k := range m {
		suffixes = append(suffixes, k)
	}

	sort.Strings(suffixes)

	return suffixes, nil
}

// attrMatch returns true if the attribute matches any of the prefixes
// in the slice.
func attrMatch(attr string, slice []string) bool {
	for _, s := range slice {
		if strings.HasPrefix(attr, s) {
			return true
		}
	}
	return false
}
