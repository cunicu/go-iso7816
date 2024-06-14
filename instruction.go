// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package iso7816

type Instruction byte

const (
	// MaxLenCmdDataStandard defines the maximum command data length of a standard length command APDU.
	MaxLenCmdDataStandard = 255

	// MaxLenRespDataStandard defines the maximum response data length of a standard length response APDU.
	MaxLenRespDataStandard = 256

	// MaxLenCmdDataExtended defines the maximum command data length of an extended length command APDU.
	MaxLenCmdDataExtended = 65535

	// MaxLenRespDataExtended defines the maximum response data length of an extended length response APDU.
	MaxLenRespDataExtended = 65536
)

const (
	InsDeactivateFile                 Instruction = 0x04 // Part 9
	InsEraseRecord                    Instruction = 0x0C // Section 7.3.8
	InsEraseBinary                    Instruction = 0x0E // Section 7.2.7
	InsEraseBinaryEven                Instruction = 0x0F // Section 7.2.7
	InsPerformSCQLOperation           Instruction = 0x10 // Part 7
	InsPerformTransactionOperation    Instruction = 0x12 // Part 7
	InsPerformUserOperation           Instruction = 0x14 // Part 7
	InsVerify                         Instruction = 0x20 // Section 7.5.6
	InsVerifyOdd                      Instruction = 0x21 // Section 7.5.6
	InsManageSecurityEnvironment      Instruction = 0x22 // Section 7.5.11
	InsChangeReferenceData            Instruction = 0x24 // Section 7.5.7
	InsDisableVerificationRequirement Instruction = 0x26 // Section 7.5.9
	InsEnableVerificationRequirement  Instruction = 0x28 // Section 7.5.8
	InsPerformSecurityOperation       Instruction = 0x2A // Part 8
	InsResetRetryCounter              Instruction = 0x2C // Section 7.5.10
	InsActivateFile                   Instruction = 0x44 // Part 9
	InsGenerateAsymmetricKeyPair      Instruction = 0x46 // Part 8
	InsManageChannel                  Instruction = 0x70 // Section 7.1.2
	InsExternalOrMutualAuthenticate   Instruction = 0x82 // Section 7.5.4
	InsGetChallenge                   Instruction = 0x84 // Section 7.5.3
	InsGeneralAuthenticate            Instruction = 0x87 // Section 7.5.5
	InsInternalAuthenticate           Instruction = 0x88 // Section 7.5.2
	InsSearchBinary                   Instruction = 0xA0 // Section 7.2.6
	InsSearchBinaryOdd                Instruction = 0xA1 // Section 7.2.6
	InsSearchRecord                   Instruction = 0xA2 // Section 7.3.7
	InsSelect                         Instruction = 0xA4 // Section 7.1.1
	InsReadBinary                     Instruction = 0xB0 // Section 7.2.3
	InsReadBinaryOdd                  Instruction = 0xB1 // Section 7.2.3
	InsReadRecord                     Instruction = 0xB3 // Section 7.3.3
	InsGetResponse                    Instruction = 0xC0 // Section 7.6.1
	InsEnvelope                       Instruction = 0xC2 // Section 7.6.2
	InsEnvelopeOdd                    Instruction = 0xC3 // Section 7.6.2
	InsGetData                        Instruction = 0xCA // Section 7.4.2
	InsGetDataOdd                     Instruction = 0xCB // Section 7.4.2
	InsWriteBinary                    Instruction = 0xD0 // Section 7.2.6
	InsWriteBinaryOdd                 Instruction = 0xD1 // Section 7.2.6
	InsWriteRecord                    Instruction = 0xD2 // Section 7.3.4
	InsUpdateBinary                   Instruction = 0xD7 // Section 7.2.5
	InsPutData                        Instruction = 0xDA // Section 7.4.3
	InsPutDataOdd                     Instruction = 0xDB // Section 7.4.3
	InsUpdateRecord                   Instruction = 0xDC // Section 7.3.5
	InsUpdateRecordOdd                Instruction = 0xDD // Section 7.3.5
	InsCreateFile                     Instruction = 0xE0 // Part 9
	InsAppendRecord                   Instruction = 0xE2 // Section 7.3.6
	InsDeleteFile                     Instruction = 0xE4 // Part 9
	InsTerminateDF                    Instruction = 0xE6 // Part 9
	InsTerminateEF                    Instruction = 0xE8 // Part 9
	InsTerminateCardUsage             Instruction = 0xFE // Part 9
)
