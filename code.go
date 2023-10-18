// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package iso7816

import (
	"fmt"
)

// Code encapsulates (some) response codes from the spec
//
//nolint:errname
type Code [2]byte

var (
	ErrSuccess = Code{0x90, 0x00}

	ErrUnspecifiedWarning          = Code{0x62, 0x00} // No information given (warning)
	Err                            = Code{0x62, 0x02} // to '80' Triggering by the card (see 8.6.1)
	ErrResponseMayBeCorrupted      = Code{0x62, 0x81} // Part of returned data may be corrupted
	ErrEOF                         = Code{0x62, 0x82} // End of file or record reached before reading Ne bytes
	ErrSelectedFileDeactivated     = Code{0x62, 0x83} // Selected file deactivated
	ErrFileControlInfoNotFormatted = Code{0x62, 0x84} // File control information not formatted according to 5.3.3
	ErrSelectedFileInTermination   = Code{0x62, 0x85} // Selected file in termination state
	ErrNoSensorData                = Code{0x62, 0x86} // No input data available from a sensor on the card

	ErrUnspecifiedWarningModiefied = Code{0x63, 0x00} // No information given (warning)
	ErrFileFilledUp                = Code{0x63, 0x81} // File filled up by the last write
	// ErrCounter                     = Code{0x63, 0xC0} // Counter from 0 to 15 encoded by 'X' (exact meaning depending on the command)

	ErrExecution                 = Code{0x64, 0x00} // Execution error (error)
	ErrImmediateResponseRequired = Code{0x64, 0x01} // Immediate response required by the card
	ErrCard                      = Code{0x64, 0x02} // to '80' Triggering by the card (see 8.6.1)

	ErrUnspecified = Code{0x64, 0x00} // No information given (error)
	ErrMemory      = Code{0x64, 0x81} // Memory failure

	ErrWrongLength = Code{0x67, 0x00} // Wrong length; no further indication

	ErrFunctionInCLANotSupported   = Code{0x68, 0x00} // Function in CLA not supported
	ErrLogicalChannelNotSupported  = Code{0x68, 0x81} // Logical channel not supported
	ErrSecureMessagingNotSupported = Code{0x68, 0x82} // Secure messaging not supported
	ErrExpectedLastCommand         = Code{0x68, 0x83} // Last command of the chain expected
	ErrCommandChainingNotSupported = Code{0x68, 0x84} // Command chaining not supported

	ErrCommandNotAllowed                  = Code{0x69, 0x00} // Command not allowed
	ErrCommandIncompatibleWithFile        = Code{0x68, 0x81} // Command incompatible with file structure
	ErrSecurityStatusNotSatisfied         = Code{0x68, 0x82} // Security status not satisfied
	ErrAuthenticationMethodBlocked        = Code{0x68, 0x83} // Authentication method blocked
	ErrReferenceDataNotUsable             = Code{0x68, 0x84} // Reference data not usable
	ErrConditionsOfUseNotSatisfied        = Code{0x68, 0x85} // Conditions of use not satisfied
	ErrCommandNotAllowedNoCurrentEF       = Code{0x68, 0x86} // Command not allowed (no current EF)
	ErrExpectedSecureMessaging            = Code{0x68, 0x87} // Expected secure messaging data objects missing
	ErrIncorredSecureMessagingDataObjects = Code{0x68, 0x88} // Incorrect secure messaging data objects

	ErrWrongParamsUnspecified       = Code{0x6A, 0x00} // No information given (error)
	ErrIncorrectParamsInCommandData = Code{0x6A, 0x80} // Incorrect parameters in the command data field
	ErrFunctionNotSupported         = Code{0x6A, 0x81} // Function not supported
	ErrFileNotFound                 = Code{0x6A, 0x82} // File or application not found
	ErrRecordNotFound               = Code{0x6A, 0x83} // Record not found
	ErrNoSpace                      = Code{0x6A, 0x84} // Not enough memory space in the file
	ErrNcInconsistentWithTLV        = Code{0x6A, 0x85} // Nc inconsistent with TLV structure
	ErrIncorrectParams              = Code{0x6A, 0x86} // Incorrect parameters P1-P2
	ErrNcInconsistentWithParams     = Code{0x6A, 0x87} // Nc inconsistent with parameters P1-P2
	ErrNotFound                     = Code{0x6A, 0x88} // Referenced data or reference data not found (exact meaning depending on the command)
	ErrFileAlreadyExists            = Code{0x6A, 0x89} // File already exists
	ErrNameAlreadyExists            = Code{0x6A, 0x8A} // DF name already exists

	ErrWrongParams = Code{0x6B, 0x00} // Wrong parameters P1-P2

	ErrUnsupportedInstruction = Code{0x6D, 0x00} // Instruction code not supported or invalid
	ErrUnsupportedClass       = Code{0x6E, 0x00} // Class not supported
	ErrNoDiag                 = Code{0x6F, 0x00} // No precise diagnosis
)

// Error return the encapsulated error string
func (c Code) Error() string {
	return fmt.Sprintf("unknown (%x)", c[:])
}

// HasMore indicates more data that needs to be fetched
func (c Code) HasMore() (bool, int) {
	return c[0] == 0x61, int(c[1])
}

// IsSuccess indicates that all data has been successfully fetched
func (c Code) IsSuccess() bool {
	return c == ErrSuccess
}

func (c Code) IsCompleted() bool {
	return c.IsSuccess() || c[0] == 0x61 || c[0] == 0x62 || c[0] == 0x63
}

func (c Code) IsAborted() bool {
	return c[0] == 0x64 || c[0] == 0x66 || c[0] == 0x67 || c[0] == 0x6f
}
