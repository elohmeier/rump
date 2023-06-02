{
  description = "Hot sync two Redis servers using dumps";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "nixpkgs/nixos-22.11";

  outputs = { self, nixpkgs }:
    let

      # to work with older version of flakes
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";

      # Generate a user-friendly version number.
      version = builtins.substring 0 8 lastModifiedDate;

      # System types to support.
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

    in
    {

      # Provide some binary packages for selected system types.
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          rump = pkgs.buildGoModule {
            pname = "rump";
            inherit version;
            src = ./.;
            vendorSha256 = "sha256-92CuSrD3lIVbfSfUu1zAYSbEEPaR3U12/a2XbhL1iC8=";
            checkFlags = [ "-skip" "pkg/file" ]; # end to end test requires a redis server
          };
        });

      # The default package for 'nix build'. This makes sense if the
      # flake provides only one package or there is a clear "main"
      # package.
      defaultPackage = forAllSystems (system: self.packages.${system}.rump);

      devShell = forAllSystems
        (system:
          let
            pkgs =
              nixpkgsFor.${system};
          in
          pkgs.mkShell {
            buildInputs = with pkgs; [
              go
            ];
          });
    };
}
