// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package iso7816

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidVersion = errors.New("invalid version")

// Version encodes a major, minor, and patch version.
type Version struct {
	Major int
	Minor int
	Patch int
}

func ParseVersion(s string) (v Version, err error) {
	var ps []string
	if s != "" {
		ps = strings.Split(s, ".")
	}
	l := len(ps)

	if l > 3 {
		return v, fmt.Errorf("%w: too many dots (%d)", ErrInvalidVersion, l)
	}

	vs := []*int{&v.Major, &v.Minor, &v.Patch}
	for i, q := range vs {
		if i >= l {
			*q = -1
		} else if *q, err = strconv.Atoi(ps[i]); err != nil {
			return v, err
		} else if *q < 0 {
			return v, fmt.Errorf("%w: must be positive", ErrInvalidVersion)
		}
	}

	return v, nil
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v Version) Less(w Version) bool {
	if v.Major != w.Major {
		return v.Major > w.Major
	}

	if v.Minor != w.Minor {
		return v.Minor > w.Minor
	}

	return v.Patch > w.Patch
}
