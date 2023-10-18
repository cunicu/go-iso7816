// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package tlv

import "errors"

var (
	ErrInvalidTag   = errors.New("invalid tag")
	ErrTagToBig     = errors.New("tag is too big for this encoding")
	ErrValueToLarge = errors.New("valye to large for this encoding")
)

type Tag uint // In theory ISO 7816-4 also supports 3-byte tags
