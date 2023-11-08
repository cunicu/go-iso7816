<!--
SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
SPDX-License-Identifier: Apache-2.0
-->

# go-iso7816: Go implementation of the ISO 7816 standard for smart card communication

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/cunicu/go-iso7816/build.yaml?style=flat-square)](https://github.com/cunicu/go-iso7816/actions)
[![goreportcard](https://goreportcard.com/badge/github.com/cunicu/go-iso7816?style=flat-square)](https://goreportcard.com/report/github.com/cunicu/go-iso7816)
[![Codecov branch](https://img.shields.io/codecov/c/github/cunicu/go-iso7816/main?style=flat-square&token=sLflBUZ2ky)](https://app.codecov.io/gh/cunicu/go-iso7816/tree/main)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square)](https://github.com/cunicu/go-iso7816/blob/main/LICENSES/Apache-2.0.txt)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/cunicu/go-iso7816?style=flat-square)
[![Go Reference](https://pkg.go.dev/badge/github.com/cunicu/go-iso7816.svg)](https://pkg.go.dev/github.com/cunicu/go-iso7816)

`go-iso7816` implements several helpers and utilities to communicate with [ISO7816](https://en.wikipedia.org/wiki/ISO/IEC_7816) compliant smart cards.

This includes:

- Abstract interface for smart card communication
- APDU parsing and serialization
  - Extended-length support
  - TLV en- & decoding variants
    - ASN.1 BER-TLV
    - Simple TLVs
    - Compact TLVs

- Constants of
  - Inter-industry instructions and status codes
  - Application Identifiers (AIDs)

- Basic card management and query support for:
  - [YubiKeys](https://www.yubico.com/)
  - [Nitrokeys](https://nitrokey.com/)

- Card enumeration and filters

- Testing utilities
  - Smartcard Mock Object
  - Tracing Wrapper

In the future we might want to add support for:

- More device support for existing cards
- Global Platform's Secure Channel Protocol (SCP03)
- Cross-platform transport implementations
  - Direct CCID
  - Apples CryptoTokenKit
  - CGo-less pcscd / libpcsc-lite

## Authors

- Steffen Vogel ([@stv0g](https://github.com/stv0g))

## License

go-iso7816 is licensed under the [Apache 2.0](./LICENSE) license.
