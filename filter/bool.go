// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package filter

import iso "cunicu.li/go-iso7816"

// And matches if all of the provided filters are matching.
func And(fs ...Filter) Filter {
	return func(reader string, card *iso.Card) (bool, error) {
		for _, f := range fs {
			if r, err := f(reader, card); err != nil {
				return false, err
			} else if !r {
				return false, nil
			}
		}

		return true, nil
	}
}

// Or matches if any of the provided filters are matching.
func Or(fs ...Filter) Filter {
	return func(reader string, card *iso.Card) (bool, error) {
		for _, f := range fs {
			if r, err := f(reader, card); err != nil {
				return false, err
			} else if r {
				return true, nil
			}
		}

		return false, nil
	}
}

// Not matches if the provided filters does not match.
func Not(f Filter) Filter {
	return func(reader string, card *iso.Card) (bool, error) {
		r, err := f(reader, card)
		if err != nil {
			return false, err
		}

		return !r, nil
	}
}
