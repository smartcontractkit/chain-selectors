{
  description = "Chain selectors metadata";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    manifest.url = "git+ssh://git@github.com/smartcontractkit/manifest?ref=main&rev=13ceea7a9921669b975986207ccb322ea389a5da";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    manifest,
    ...
  } @ inputs:
  # it enables us to generate all outputs for each default system
  (
    flake-utils.lib.eachDefaultSystem
    (system: let
      # choose nixpkgs system?
      pkgs = import nixpkgs {
        inherit system;
      };

      # Reference to self (used to resolve model dependencies)
      # NOTE: this should only contain references to other models, anything else can affect build times
      mf = models;

      lib = manifest.lib.${system}.override {
        manifest = manifestArgs;
        pkgs = pkgs;
      };

      binArgs = {
        pkgs = pkgs;
        lib = lib;
      };

      mfcli = manifest.packages.${system}.mf;

      # To define your own bins and libs, rather than importing them from the manifest, you can use this
      # https://github.com/smartcontractkit/manifest/blob/5511f3663a6286074e67dd85c138756beddf9810/flake.nix#L43
      bins = {
      };
      schemas = {
        
      };

      # Manifest packages that host Manifest Nix templates which are used by the Manifest Nix lib to resolve models
      templates = lib.templates;

      # Metadata about the version of manifest. You can add more metadata here as an override
      metadata = {
        "manifest-metadata" = manifest.packages.${system}."manifest-metadata".override {
          pkgs = pkgs;
          metadataOverride = {
            description = "Chain selectors metadata";
          };
        };
      };

      manifestArgs = {
        inherit mf;
        inherit lib;
        inherit bins;
        inherit schemas;
        inherit templates;
      };

      # Resolve and flatten all .manifest models
      # You don't need the '// manifest.models.${system}' if you don't need access to the imported models
      models =
        lib.resolve ./. ./.manifest;

      # Group and flatten all pkgs for this Manifest
      manifestPackages = flake-utils.lib.flattenTree (lib.packagesTree models);
    in {
      # it outputs packages all packages defined in plugins
      packages = manifestPackages // metadata;

      # A mechanism to build all model pkgs via `nix flake check`
      # naming of check packages should match ones exported in `packages`
      # uses a subset of manifestPackages
      checks = flake-utils.lib.flattenTree (lib.packagesModels models);

      # access shell with mf CLI
      devShells.mf = pkgs.mkShell {
        buildInputs = [mfcli];
      };

      # shell with various linting tools
      devShells.default = pkgs.mkShell {
        buildInputs = [
          # nix tooling
          pkgs.alejandra
        ];
      };
    })
  );
}
