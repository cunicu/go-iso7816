// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package iso7816

func concat(prefix []byte, rest ...byte) (r []byte) {
	r = append(r, prefix...)
	return append(r, rest...)
}

type RID [5]byte

func (r RID) String() string {
	switch r {
	case RidNIST:
		return "NIST"
	case RidFSFE:
		return "FSFE"
	case RidYubico:
		return "Yubico"
	case RidFIDO:
		return "FIDO"
	case RidSolokeys:
		return "Solokeys"
	case RidGlobalPlatform:
		return "GlobalPlatform"
	case RidNXPNFC:
		return "NXP NFC"
	default:
		return "<unknown>"
	}
}

//nolint:gochecknoglobals
var (
	// https://www.eftlab.com/knowledge-base/complete-list-of-registered-application-provider-identifiers-rid
	RidNIST           = RID{0xa0, 0x00, 0x00, 0x03, 0x08}
	RidFSFE           = RID{0xd2, 0x76, 0x00, 0x01, 0x24}
	RidYubico         = RID{0xa0, 0x00, 0x00, 0x05, 0x27}
	RidFIDO           = RID{0xa0, 0x00, 0x00, 0x06, 0x47}
	RidSolokeys       = RID{0xA0, 0x00, 0x00, 0x08, 0x47}
	RidGlobalPlatform = RID{0xa0, 0x00, 0x00, 0x01, 0x51}
	RidNXPNFC         = RID{0xD2, 0x76, 0x00, 0x00, 0x85}
)

//nolint:gochecknoglobals
var (
	// https://nvlpubs.nist.gov/nistpubs/specialpublications/nist.sp.800-73-4.pdf
	AidPIV = concat(RidNIST[:], 0x00, 0x00, 0x10, 0x00)

	// https://gnupg.org/ftp/specs/OpenPGP-smart-card-application-3.4.1.pdf
	AidOpenPGP = concat(RidFSFE[:], 0x01)

	// https://fidoalliance.org/specs/fido-v2.1-ps-20210615/fido-client-to-authenticator-protocol-v2.1-ps-20210615.html#nfc-applet-selection
	AidFIDO = concat(RidFIDO[:], 0x2f, 0x00, 0x01)

	// https://github.com/Yubico/yubikey-manager/blob/6496393f9269e86fb7b4b67907b397db33b50c2d/yubikit/core/smartcard.py#L66
	AidYubicoOTP           = concat(RidYubico[:], 0x20, 0x01)
	AidYubicoManagement    = concat(RidYubico[:], 0x47, 0x11, 0x17)
	AidYubicoOATH          = concat(RidYubico[:], 0x21, 0x01)
	AidYubicoHSMAuth       = concat(RidYubico[:], 0x21, 0x07, 0x01)
	AidSolokeysAdmin       = concat(RidSolokeys[:], 0x00, 0x00, 0x00, 0x01)
	AidSolokeysProvisioner = concat(RidSolokeys[:], 0x01, 0x00, 0x00, 0x01)
	AidCardManager         = concat(RidGlobalPlatform[:], 0x00, 0x00, 0x00)
	AidNDEF                = concat(RidNXPNFC[:], 0x01, 0x01)
)
