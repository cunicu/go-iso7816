// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package iso7816

import (
	"encoding/binary"
	"log/slog"

	"cunicu.li/go-iso7816/encoding/tlv"
)

// Compact TLV tags used in historical bytes
const (
	ctagCountryCode   tlv.Tag = 0x1 // ISO 7816-4 Section 8.1.1.2.1 Country or issuer indicator
	ctagIssuerID      tlv.Tag = 0x2 // ISO 7816-4 Section 8.1.1.2.1 Country or issuer indicator
	ctagAID           tlv.Tag = 0xf // ISO 7816-4 Section 8.1.1.2.2 Application identifier
	ctagCardService   tlv.Tag = 0x3 // ISO 7816-4 Section 8.1.1.2.2 Application identifier
	ctagInitialAccess tlv.Tag = 0x4 // ISO 7816-4 Section 8.1.1.2.4 Initial access data
	ctagCardIssuer    tlv.Tag = 0x5 // ISO 7816-4 Section 8.1.1.2.5 Card issuer's data
	ctagPreIssuing    tlv.Tag = 0x6 // ISO 7816-4 Section 8.1.1.2.6 Pre-issuing data
	ctagCapabilities  tlv.Tag = 0x7 // ISO 7816-4 Section 8.1.1.2.7 Card capabilities
)

type CardService byte

const (
	CardServiceMF                    CardService = (1 << 0) // Card with MF
	CardServiceBERTLVInATR           CardService = (1 << 4) // BER-TLV data objects available in EF.ATR (see 8.2.1.1)
	CardServiceBERTLVInDIR           CardService = (1 << 5) // BER-TLV data objects available in EF.DIR (see 8.2.1.1)
	CardServiceAppSelectionPartialDF CardService = (1 << 6) // Application Selection by partial DF name
	CardServiceAppSelectionFullDF    CardService = (1 << 7) // Application Selection by full DF name (AID)

	// EF.DIR and EF.ATR access services
	// Access via Cardservice.AccessServices()
	CardServiceAccessByGetDataCmd    CardService = (1 << 2) // By the GET DATA command (TLV structure)
	CardServiceAccessByReadRecordCmd CardService = 0        // By the READ RECORD (S) command (record structure)
	CardServiceAccessByReadBinaryCmd CardService = (1 << 3) // By the READ BINARY command (transparent structure)
)

func (cs CardService) AccessServices() byte {
	return byte(cs) & 0xe
}

func (cs *CardService) Decode(b []byte) error {
	if len(b) != 1 {
		return errInvalidLength
	}

	*cs = CardService(b[0])

	return nil
}

// Card capabilities
// See 8.1.1.2.7 Card capabilities
type (
	CardCapabilities uint32
)

const (
	// Table 1
	CardCapRecordIdentifier         CardCapabilities = (1 << iota) // Record identifier supported
	CardCapRecordNumber                                            // Record number supported
	CardCapShortEFIdentifier                                       // Short EF identifier supported
	CardCapImplicitDFSelection                                     // Implicit DF selection
	CardCapDFSelectByFileIdentifier                                // DF selection by file identifier
	CardCapDFSelectByPath                                          // DF selection by path
	CardCapDFSelectByPartialDFName                                 // DF selection by partial DF name
	CardCapDFSelectByFullDFName                                    // DF selection by full DF name

	// Table 2
	_                          // Data unit size in quartets (Byte 1)
	_                          // Data unit size in quartets (Byte 2)
	_                          // Data unit size in quartets (Byte 3)
	_                          // Data unit size in quartets (Byte 4)
	CardCapFirstTagByteValidFF // Value 'FF' is valid for the first byte of BER-TLV tag fields (see 5.2.2.1)
	_                          // Behavior of write functions
	_                          // Behavior of write functions
	CardCapEFsTLVStructure     // EFs of TLV structure supported

	// Table 3
	_                                            // Maximum number of logical channels (see 5.1.1 and 5.1.1.2) (Byte 1)
	_                                            // Maximum number of logical channels (see 5.1.1 and 5.1.1.2) (Byte 2)
	_                                            // Maximum number of logical channels (see 5.1.1 and 5.1.1.2) (Byte 3)
	CardCapLogicalChannelAssignByInterfaceDevice // Logical channel number assignment by the interface device (see 7.1.2)
	CardCapLogicalChannelAssignByCard            // Logical channel number assignment by the card (see 7.1.2)
	CardCapExtendedLengthInfoInEFATR             // Extended Length Information in EF.ATR/ INFO (OpenPGP specific?)
	CardCapExtendedLength                        // Extended Lc and Le fields (see 5.1)
	CardCapCommandChaining                       // Command chaining (see 5.1.1.1)
)

// DataUnitSize returns the data unit size in quartets
// Note: 2 quartets are one byte!
func (cc CardCapabilities) DataUnitSize() int {
	shift := int((cc >> 8) & 0xf)
	return 1 << shift
}

func (cc CardCapabilities) WriteBehaviour() CardCapabilities {
	return (cc >> 8) & 0x60
}

func (cc CardCapabilities) LogicalChannelCount() int {
	return int((cc>>16)&0x7 + 1)
}

func (cc *CardCapabilities) Decode(b []byte) error {
	if len(b) < 1 || len(b) > 3 {
		return errInvalidLength
	}

	c := make([]byte, 4)
	copy(c, b)

	*cc = CardCapabilities(binary.LittleEndian.Uint32(c))

	return nil
}

// See: ISO-7816-4 - Section 8.1.1 Historical bytes
type HistoricalBytes struct {
	CategoryIndicator byte
	LifeCycleStatus   byte // 8.1.1.3 Status indicator
	ProcessingStatus  Code // 8.1.1.3 Status indicator

	CountryCode      []byte           // 8.1.1.2.1 Country or issuer indicator
	IssuerID         []byte           // 8.1.1.2.1 Country or issuer indicator
	AID              []byte           // 8.1.1.2.2 Application identifier
	CardService      CardService      // 8.1.1.2.3 Card service data
	InitialAccess    []byte           // 8.1.1.2.4 Initial access data
	CardIssuer       []byte           // 8.1.1.2.5 Card issuer's data
	PreIssuing       []byte           // 8.1.1.2.6 Pre-issuing data
	CardCapabilities CardCapabilities // 8.1.1.2.7 Card capabilities
}

func (h *HistoricalBytes) Decode(b []byte) (err error) {
	h.CategoryIndicator = b[0]

	switch h.CategoryIndicator {
	case 0x10:
		// Not supported

	case 0x00:
		lb := len(b)
		h.LifeCycleStatus = b[lb-3]
		h.ProcessingStatus = Code{b[lb-2], b[lb-1]}
		b = b[:lb-3]
		fallthrough

	case 0x80:
		tvs, err := tlv.DecodeCompact(b)
		if err != nil {
			return err
		}

		for _, tv := range tvs {
			switch tv.Tag {
			case ctagCountryCode:
				h.CountryCode = tv.Value
			case ctagIssuerID:
				h.IssuerID = tv.Value
			case ctagAID:
				h.AID = tv.Value
			case ctagInitialAccess:
				h.InitialAccess = tv.Value
			case ctagCardIssuer:
				h.CardIssuer = tv.Value
			case ctagPreIssuing:
				h.PreIssuing = tv.Value
			case ctagCapabilities:
				if err := h.CardCapabilities.Decode(tv.Value); err != nil {
					return err
				}
			case ctagCardService:
				if err := h.CardService.Decode(tv.Value); err != nil {
					return err
				}
			}
		}

	default:
		slog.Warn("Received unknown category indicator",
			slog.String("do", "historical bytes"),
			slog.Int("category_indicator", int(h.CategoryIndicator)))
	}

	return nil
}
