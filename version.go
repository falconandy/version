package version

import (
	"regexp"
)

const (
	zero = byte('0')
)

var (
	partRe = regexp.MustCompile(`^(\d+)(([A-Za-z]+)(\d?))?(\.|-|\+|$)`)
)

type Version struct {
	rawVersion string

	Parts   []Part // "1" and "8rc2" in "1.8rc2-wheezy"
	Postfix string // "-wheezy" in "1.8rc2-wheezy"
}

type Part struct {
	Number        string // "1" and "8" in 2-part "1.8rc2-wheezy"
	Postfix       string // "" and "rc"  in 2-part "1.8rc2-wheezy"
	NumberPostfix string // "" and "2" in 2-part "1.8rc2-wheezy"
}

func NewVersion(rawVersion string) Version {
	parts, postfix := parseRawVersion(rawVersion)
	return Version{rawVersion: rawVersion, Parts: parts, Postfix: postfix}
}

func (v Version) String() string {
	return v.rawVersion
}

func (v Version) Compare(other Version) int {
	result := v.CompareParts(other)
	if result != 0 {
		return result
	}

	return v.ComparePostfix(other)
}

func (v Version) CompareParts(other Version) int {
	if len(v.Parts) == 0 && len(other.Parts) > 0 {
		return -1
	}

	if len(v.Parts) > 0 && len(other.Parts) == 0 {
		return 1
	}

	for i := 0; i < len(v.Parts) && i < len(other.Parts); i++ {
		part1, part2 := v.Parts[i], other.Parts[i]

		compare := compareNumbers(part1.Number, part2.Number)
		if compare != 0 {
			return compare
		}

		if part1.Postfix != "" && part2.Postfix == "" {
			return -1
		}

		if part1.Postfix == "" && part2.Postfix != "" {
			return 1
		}

		if part1.Postfix < part2.Postfix {
			return -1
		}

		if part1.Postfix > part2.Postfix {
			return 1
		}

		compare = compareNumbers(part1.NumberPostfix, part2.NumberPostfix)
		if compare != 0 {
			return compare
		}
	}

	if len(v.Parts) < len(other.Parts) {
		return -1
	}

	if len(v.Parts) > len(other.Parts) {
		return 1
	}

	return 0
}

func (v Version) ComparePostfix(other Version) int {
	if v.Postfix < other.Postfix {
		return -1
	}

	if v.Postfix > other.Postfix {
		return 1
	}

	return 0
}

func parseRawVersion(rawVersion string) (parts []Part, postfix string) {
	postfix = rawVersion

	for {
		match := partRe.FindStringSubmatch(postfix)
		if match == nil {
			break
		}

		part := Part{Number: match[1], Postfix: match[3], NumberPostfix: match[4]}
		parts = append(parts, part)

		if match[5] != "." {
			postfix = postfix[len(match[0])-len(match[5]):]
			break
		}

		postfix = postfix[len(match[0]):]
	}

	return parts, postfix
}

func compareNumbers(n1, n2 string) int {
	var maxLen int
	if len(n2) > len(n1) {
		maxLen = len(n2)
	} else {
		maxLen = len(n1)
	}

	for i := maxLen; i > 0; i-- {
		i1, i2 := len(n1)-i, len(n2)-i
		d1, d2 := zero, zero
		if i1 >= 0 {
			d1 = n1[i1]
		}
		if i2 >= 0 {
			d2 = n2[i2]
		}

		if d1 < d2 {
			return -1
		}

		if d1 > d2 {
			return 1
		}
	}

	return 0
}
