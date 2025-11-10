# Changelog

## [1.1.0](https://github.com/Firecrown-Media/stax/compare/v1.0.0...v1.1.0) (2025-11-10)


### Features

* Enhanced error messaging and UX improvements ([#34](https://github.com/Firecrown-Media/stax/issues/34)) ([84ef266](https://github.com/Firecrown-Media/stax/commit/84ef2662d90c9dc62ce67408634dacdd60300c68))

## [1.0.0](https://github.com/Firecrown-Media/stax/compare/v0.5.0...v1.0.0) (2025-11-10)


### ⚠ BREAKING CHANGES

* Default project type changed from wordpress-multisite to wordpress

### Documentation

* clarify single site and multisite support ([3a2414b](https://github.com/Firecrown-Media/stax/commit/3a2414b294503ebb852d76d7edc380ddf6a78db3))

## [0.5.0](https://github.com/Firecrown-Media/stax/compare/v0.4.2...v0.5.0) (2025-11-10)


### Features

* **list:** add global list command for WPEngine installs ([#30](https://github.com/Firecrown-Media/stax/issues/30)) ([45f73c2](https://github.com/Firecrown-Media/stax/commit/45f73c224feb70037f91fc74a5f34a9099c14bfd))

## [0.4.2](https://github.com/Firecrown-Media/stax/compare/v0.4.1...v0.4.2) (2025-11-10)


### Bug Fixes

* credential storage for CGO-disabled builds (Homebrew) ([#28](https://github.com/Firecrown-Media/stax/issues/28)) ([3379b3a](https://github.com/Firecrown-Media/stax/commit/3379b3ae3fafedeb17457dd7bfe3706446de63e7)), closes [#27](https://github.com/Firecrown-Media/stax/issues/27)

## [0.4.1](https://github.com/Firecrown-Media/stax/compare/v0.4.0...v0.4.1) (2025-11-09)


### Bug Fixes

* resolve platform-specific keychain build issues for releases ([#25](https://github.com/Firecrown-Media/stax/issues/25)) ([e39f9ab](https://github.com/Firecrown-Media/stax/commit/e39f9ab1e638ff8fc7d585b5395767104634cad2))

## [0.4.0](https://github.com/Firecrown-Media/stax/compare/v0.3.0...v0.4.0) (2025-11-09)


### Features

* complete codebase refactor with build system, tests, and release automation ([e09284f](https://github.com/Firecrown-Media/stax/commit/e09284f0f73dc00ce77eb92202b3764bd663f34c))


### Bug Fixes

* disable CGO for Darwin builds to enable cross-compilation ([d65db36](https://github.com/Firecrown-Media/stax/commit/d65db36556b401216cb7f9f1b13510daf2e6a245))

## [0.3.0](https://github.com/Firecrown-Media/stax/compare/v0.2.0...v0.3.0) (2025-11-09)


### Features

* complete codebase refactor with build system, tests, and CI fixes ([9de3d0d](https://github.com/Firecrown-Media/stax/commit/9de3d0dd4bfb4b7862055d2a7c72d02e6004d4f2))


### Bug Fixes

* update GoReleaser action to v6 for version 2 config support ([6865c4f](https://github.com/Firecrown-Media/stax/commit/6865c4fe7d901f090daa29e7cc4cbbbe3f1a73ef))

## [0.2.0](https://github.com/Firecrown-Media/stax/compare/v0.1.1...v0.2.0) (2025-11-09)


### Features

* add automated release process and reorganize documentation ([6d605c5](https://github.com/Firecrown-Media/stax/commit/6d605c51038b0c765000d52f66b57d9e9267f97e))


### Bug Fixes

* add wp_blogs table to mock database dump for multisite tests ([07faa2d](https://github.com/Firecrown-Media/stax/commit/07faa2d2a20533ade31b6a632c3c8d612e370b1d))
* correct mock types and remove omitempty from Build config booleans ([e39a618](https://github.com/Firecrown-Media/stax/commit/e39a6185ec5c4e7046e2b38bb781b87f315bdf1a)), closes [#9](https://github.com/Firecrown-Media/stax/issues/9)
* format code with gofmt to pass CI checks ([c2654f7](https://github.com/Firecrown-Media/stax/commit/c2654f73db9a20c8966797a955611766722d5616))
* format test files with gofmt ([559c111](https://github.com/Firecrown-Media/stax/commit/559c1111fc37c059df89ef8fa24778f938e88359))
* update .gitignore to include pkg/build directory ([bac9821](https://github.com/Firecrown-Media/stax/commit/bac982110fa36d62974cb47d5d95fdd680dad3b4))
