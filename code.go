// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package iso7816

import (
	"fmt"
)

// Code encapsulates (some) response codes from the spec
// See: ISO 7816-4 Section 5.1.3 Status Bytes
//
//nolint:errname
type Code [2]byte

var (
	// General or unspecified warnings & errors
	ErrSuccess                    = Code{0x90, 0x00}
	ErrUnspecifiedWarning         = Code{0x62, 0x00} // No information given (warning)
	ErrUnspecifiedWarningModified = Code{0x63, 0x00} // No information given (warning), on-volatile memory has changed
	ErrUnspecifiedError           = Code{0x64, 0x00} // No information given (error)
	ErrUnspecifiedErrorModified   = Code{0x65, 0x00} // No information given (error), on-volatile memory has changed
	ErrWrongLength                = Code{0x67, 0x00} // Wrong length; no further indication
	ErrUnsupportedFunction        = Code{0x68, 0x00} // Function in CLA not supported
	ErrCommandNotAllowed          = Code{0x69, 0x00} // Command not allowed
	ErrWrongParamsNoInfo          = Code{0x6A, 0x00} // No information given (error)
	ErrWrongParams                = Code{0x6B, 0x00} // Wrong parameters P1-P2
	ErrUnsupportedInstruction     = Code{0x6D, 0x00} // Instruction code not supported or invalid
	ErrUnsupportedClass           = Code{0x6E, 0x00} // Class not supported
	ErrNoDiag                     = Code{0x6F, 0x00} // No precise diagnosis

	// Specific warnings & errors
	ErrResponseMayBeCorrupted              = Code{0x62, 0x81} // Part of returned data may be corrupted
	ErrEOF                                 = Code{0x62, 0x82} // End of file or record reached before reading Ne bytes
	ErrSelectedFileDeactivated             = Code{0x62, 0x83} // Selected file deactivated
	ErrInvalidFileControlInfo              = Code{0x62, 0x84} // File control information not formatted according to 5.3.3
	ErrSelectedFileInTermination           = Code{0x62, 0x85} // Selected file in termination state
	ErrNoSensorData                        = Code{0x62, 0x86} // No input data available from a sensor on the card
	ErrFileFilledUp                        = Code{0x63, 0x81} // File filled up by the last write
	ErrImmediateResponseRequired           = Code{0x64, 0x01} // Immediate response required by the card
	ErrMemory                              = Code{0x65, 0x81} // Memory failure
	ErrLogicalChannelNotSupported          = Code{0x68, 0x81} // Logical channel not supported
	ErrSecureMessagingNotSupported         = Code{0x68, 0x82} // Secure messaging not supported
	ErrExpectedLastCommand                 = Code{0x68, 0x83} // Last command of the chain expected
	ErrCommandChainingNotSupported         = Code{0x68, 0x84} // Command chaining not supported
	ErrCommandIncompatibleWithFile         = Code{0x69, 0x81} // Command incompatible with file structure
	ErrSecurityStatusNotSatisfied          = Code{0x69, 0x82} // Security status not satisfied
	ErrAuthenticationMethodBlocked         = Code{0x69, 0x83} // Authentication method blocked
	ErrReferenceDataNotUsable              = Code{0x69, 0x84} // Reference data not usable
	ErrConditionsOfUseNotSatisfied         = Code{0x69, 0x85} // Conditions of use not satisfied
	ErrCommandNotAllowedNoCurrentEF        = Code{0x69, 0x86} // Command not allowed (no current EF)
	ErrExpectedSecureMessaging             = Code{0x69, 0x87} // Expected secure messaging data objects missing
	ErrIncorrectSecureMessagingDataObjects = Code{0x69, 0x88} // Incorrect secure messaging data objects
	ErrIncorrectData                       = Code{0x6A, 0x80} // Incorrect parameters in the command data field
	ErrFunctionNotSupported                = Code{0x6A, 0x81} // Function not supported
	ErrFileOrAppNotFound                   = Code{0x6A, 0x82} // File or application not found
	ErrRecordNotFound                      = Code{0x6A, 0x83} // Record not found
	ErrNoSpace                             = Code{0x6A, 0x84} // Not enough memory space in the file
	ErrInvalidNcWithTLV                    = Code{0x6A, 0x85} // Nc inconsistent with TLV structure
	ErrIncorrectParams                     = Code{0x6A, 0x86} // Incorrect parameters P1-P2
	ErrInvalidNcWithParams                 = Code{0x6A, 0x87} // Nc inconsistent with parameters P1-P2
	ErrReferenceNotFound                   = Code{0x6A, 0x88} // Referenced data or reference data not found (exact meaning depending on the command)
	ErrFileAlreadyExists                   = Code{0x6A, 0x89} // File already exists
	ErrNameAlreadyExists                   = Code{0x6A, 0x8A} // DF name already exists
)

// Error return the encapsulated error string
func (c Code) Error() string {
	switch c {
	case ErrSuccess:
		return "success"
	case ErrUnspecifiedWarning:
		return "unspecified warning"
	case ErrUnspecifiedWarningModified:
		return "unspecified warning; on-volatile memory has changed"
	case ErrUnspecifiedError:
		return "unspecified error"
	case ErrUnspecifiedErrorModified:
		return "unspecified error; on-volatile memory has changed"
	case ErrWrongLength:
		return "wrong length; no further indication"
	case ErrUnsupportedFunction:
		return "function in CLI not supported"
	case ErrCommandNotAllowed:
		return "command not allowed"
	case ErrWrongParamsNoInfo:
		return "wrong parameters"
	case ErrWrongParams:
		return "wrong parameters p1-p2"
	case ErrUnsupportedInstruction:
		return "instruction code not supported or invalid"
	case ErrUnsupportedClass:
		return "class not supported"
	case ErrNoDiag:
		return "no precise diagnosis"

	case ErrResponseMayBeCorrupted:
		return "part of returned data may be corrupted"
	case ErrEOF:
		return "end of file or record reached before reading Ne bytes"
	case ErrSelectedFileDeactivated:
		return "selected file deactivated"
	case ErrInvalidFileControlInfo:
		return "file control information not formatted correctly"
	case ErrSelectedFileInTermination:
		return "selected file in termination state"
	case ErrNoSensorData:
		return "no input data available from a sensor on the card"
	case ErrFileFilledUp:
		return "file filled up by the last write"
	case ErrImmediateResponseRequired:
		return "immediate response required by the card"
	case ErrMemory:
		return "memory failure"
	case ErrLogicalChannelNotSupported:
		return "logical channel not supported"
	case ErrSecureMessagingNotSupported:
		return "secure messaging not supported"
	case ErrExpectedLastCommand:
		return "last command of the chain expected"
	case ErrCommandChainingNotSupported:
		return "command chaining not supported"
	case ErrCommandIncompatibleWithFile:
		return "command incompatible with file structure"
	case ErrSecurityStatusNotSatisfied:
		return "security status not satisfied"
	case ErrAuthenticationMethodBlocked:
		return "authentication method blocked"
	case ErrReferenceDataNotUsable:
		return "reference data not usable"
	case ErrConditionsOfUseNotSatisfied:
		return "conditions of use not satisfied"
	case ErrCommandNotAllowedNoCurrentEF:
		return "command not allowed (no current EF)"
	case ErrExpectedSecureMessaging:
		return "expected secure messaging data objects missing"
	case ErrIncorrectSecureMessagingDataObjects:
		return "incorrect secure messaging data objects"
	case ErrIncorrectData:
		return "incorrect parameters in the command data field"
	case ErrFunctionNotSupported:
		return "function not supported"
	case ErrFileOrAppNotFound:
		return "file or application not found"
	case ErrRecordNotFound:
		return "record not found"
	case ErrNoSpace:
		return "not enough memory space in the file"
	case ErrInvalidNcWithTLV:
		return "Nc inconsistent with tlv structure"
	case ErrIncorrectParams:
		return "incorrect parameters P1-P2"
	case ErrInvalidNcWithParams:
		return "Nc inconsistent with parameters P1-P2"
	case ErrReferenceNotFound:
		return "referenced data or reference data not found (exact meaning depending on the command)"
	case ErrFileAlreadyExists:
		return "file already exists"
	case ErrNameAlreadyExists:
		return "DF name already exists"
	}
	return fmt.Sprintf("unknown (%x)", c[:])
}

// HasMore indicates more data that needs to be fetched
func (c Code) HasMore() bool {
	return c[0] == 0x61
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
