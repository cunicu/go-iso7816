# SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
# SPDX-License-Identifier: Apache-2.0
{
  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };
  outputs =
    {
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        frameworks = pkgs.darwin.apple_sdk.frameworks;
      in
      {
        devShell = pkgs.mkShell {
          buildInputs =
            with pkgs;
            [
              pkg-config
              clang
              go
              golangci-lint
              reuse
            ]
            ++ lib.optionals pkgs.stdenv.isLinux [
              pcsclite
              pcsctools
            ]
            ++ lib.optionals pkgs.stdenv.isDarwin [ frameworks.PCSC ];

          shellHook =
            if pkgs.stdenv.isDarwin then
              ''
                export CGO_LDFLAGS="-F${frameworks.PCSC}/Library/Frameworks";
              ''
            else
              "";

          hardeningDisable = [ "fortify" ];
        };

        formatter = pkgs.nixfmt-rfc-style;
      }
    );
}
