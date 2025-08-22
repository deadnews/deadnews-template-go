# Changelog

## [unreleased]

### Features

- migrate from `echo` to `chi` ([#127](https://github.com/deadnews/deadnews-template-go/issues/127)) - ([27207dc](https://github.com/deadnews/deadnews-template-go/commit/27207dc9b7496d9ba0b65d29deefc38630fab685))
- add built-in `healthcheck` functionality ([#123](https://github.com/deadnews/deadnews-template-go/issues/123)) - ([9a5f097](https://github.com/deadnews/deadnews-template-go/commit/9a5f0978bf5ba802f8c30412cbe724378c432322))

### Documentation

- _(changelog)_ update `git-cliff` config - ([f3fa2b9](https://github.com/deadnews/deadnews-template-go/commit/f3fa2b9974896905c571ba06baa97cf2410ef5cb))

### Build

- _(docker)_ add `healthcheck` ([#88](https://github.com/deadnews/deadnews-template-go/issues/88)) - ([2619c07](https://github.com/deadnews/deadnews-template-go/commit/2619c0772d289698f722ea307e8ebe73d7219c93))

### Chores

- _(config)_ migrate renovate config ([#100](https://github.com/deadnews/deadnews-template-go/issues/100)) - ([4cd7c4c](https://github.com/deadnews/deadnews-template-go/commit/4cd7c4c46339e9f048734862c5bf86982aeda3bd))
- _(github)_ replace manual `govulncheck` with action ([#126](https://github.com/deadnews/deadnews-template-go/issues/126)) - ([bf4d388](https://github.com/deadnews/deadnews-template-go/commit/bf4d388b9e242db94dd884bec8778e7911b04321))
- _(github)_ remove `macos` from tests matrix - ([561b4c0](https://github.com/deadnews/deadnews-template-go/commit/561b4c036a3c09b719841b371f9acd40839c56d6))
- _(golangci)_ remove dups - ([0bad907](https://github.com/deadnews/deadnews-template-go/commit/0bad907a8a413ebbcb6ae55ae0aa3c87f953082c))
- _(renovate)_ move regex to custom config - ([1ce89e0](https://github.com/deadnews/deadnews-template-go/commit/1ce89e0f8fee00a8050f3fea56ada6136a00207b))
- _(renovate)_ detect docker images in test files ([#129](https://github.com/deadnews/deadnews-template-go/issues/129)) - ([0f59215](https://github.com/deadnews/deadnews-template-go/commit/0f59215a1a2224e0c4f4cc61c163db02d40e3ff9))
- standardize repository references ([#132](https://github.com/deadnews/deadnews-template-go/issues/132)) - ([e97fec4](https://github.com/deadnews/deadnews-template-go/commit/e97fec4eeaf0ef90cedec12d1017359556dce4f0))
- update `golangci` configuration ([#99](https://github.com/deadnews/deadnews-template-go/issues/99)) - ([a82df2e](https://github.com/deadnews/deadnews-template-go/commit/a82df2e327aeea56ef19ab3b35f4639150bd085c))

### Dependencies

- prepare for `golangci:v2` ([#118](https://github.com/deadnews/deadnews-template-go/issues/118)) - ([0a12e36](https://github.com/deadnews/deadnews-template-go/commit/0a12e367ba08b43093f1159d3690e2586a1dd6a4))

## [2.0.6](https://github.com/deadnews/deadnews-template-go/compare/v2.0.5...v2.0.6) - 2024-05-18

### Documentation

- _(readme)_ update badges - ([1a64c93](https://github.com/deadnews/deadnews-template-go/commit/1a64c93b726de87e029206403ef3af30b8e65fe0))

### Chores

- _(github)_ update `goreleaser` job ([#75](https://github.com/deadnews/deadnews-template-go/issues/75)) - ([7e57d9c](https://github.com/deadnews/deadnews-template-go/commit/7e57d9cbe6d4156631b2cc7b30399f144b1ebe5a))
- _(github)_ update `goreleaser` job ([#70](https://github.com/deadnews/deadnews-template-go/issues/70)) - ([f703317](https://github.com/deadnews/deadnews-template-go/commit/f703317d71f90aeb45fa87210624440e9908244f))
- _(github)_ update `goreleaser` job ([#69](https://github.com/deadnews/deadnews-template-go/issues/69)) - ([961743a](https://github.com/deadnews/deadnews-template-go/commit/961743a14257db69f6d63ec4cb693bef42f980d3))
- _(typos)_ ignore short words - ([9bff8a9](https://github.com/deadnews/deadnews-template-go/commit/9bff8a9a769315168791358f84ae94e1937bc19f))

## [2.0.5](https://github.com/deadnews/deadnews-template-go/compare/v2.0.4...v2.0.5) - 2024-04-06

### Chores

- _(github)_ update `goreleaser` job ([#68](https://github.com/deadnews/deadnews-template-go/issues/68)) - ([0b6064b](https://github.com/deadnews/deadnews-template-go/commit/0b6064b460c765316104ef350c252708ebbef741))

## [2.0.4](https://github.com/deadnews/deadnews-template-go/compare/v2.0.3...v2.0.4) - 2024-04-06

### Documentation

- _(changelog)_ add `git-cliff` ([#67](https://github.com/deadnews/deadnews-template-go/issues/67)) - ([1f93052](https://github.com/deadnews/deadnews-template-go/commit/1f930527660452e569d10576c3b22a82ecc49089))
- _(readme)_ update badges - ([e171584](https://github.com/deadnews/deadnews-template-go/commit/e171584c9c59e14a67747491b0804189ec3b242a))
- _(readme)_ add badges - ([b0e89dc](https://github.com/deadnews/deadnews-template-go/commit/b0e89dc0b9c1191933101d2ba7f450a731988f1f))

### Build

- _(dockerfile)_ explicitly disable `cgo` - ([d863cfe](https://github.com/deadnews/deadnews-template-go/commit/d863cfe6533db65d19c5afc3759fe08b621616c6))
- _(dockerfile)_ explicitly disable `cgo` - ([44bc4e3](https://github.com/deadnews/deadnews-template-go/commit/44bc4e3ccc5c4363af529485ca709ab2e9435f16))
- _(dockerfile)_ explicitly disable `cgo` - ([9acff34](https://github.com/deadnews/deadnews-template-go/commit/9acff34cd3bd4f065d519ddbc1c95e3f41566aaa))

### Chores

- _(makefile)_ update - ([aa7e699](https://github.com/deadnews/deadnews-template-go/commit/aa7e69997481479e4f0d62a03798d739409e58a8))
- _(pre-commit)_ replace `hadolint-docker` with `hadolint-py` - ([944f617](https://github.com/deadnews/deadnews-template-go/commit/944f617e6f5ac6913739c760979ae31ce5473d82))

### Dependencies

- update module github.com/stretchr/testify to v1.9.0 ([#62](https://github.com/deadnews/deadnews-template-go/issues/62)) - ([3f2b941](https://github.com/deadnews/deadnews-template-go/commit/3f2b9418c6f51f0622184bb5fd31a9c1a1ca7da8))

## [2.0.3](https://github.com/deadnews/deadnews-template-go/compare/v2.0.2...v2.0.3) - 2024-01-24

### Build

- _(docker)_ update `Dockerfile` - ([83b98a4](https://github.com/deadnews/deadnews-template-go/commit/83b98a41fbaa9d12745ce357c469f807cb1ba121))
- _(docker)_ add `docker-compose` - ([fc59f04](https://github.com/deadnews/deadnews-template-go/commit/fc59f04a8c95e230dc2157f9285a0f295ffb0cac))
- _(goreleaser)_ update config ([#54](https://github.com/deadnews/deadnews-template-go/issues/54)) - ([b85c149](https://github.com/deadnews/deadnews-template-go/commit/b85c1492d90a1b2fef35f57711628c1088deba79))

### Chores

- _(pre-commit)_ add `checkmake` hook - ([16c1e3f](https://github.com/deadnews/deadnews-template-go/commit/16c1e3fd97f15ef4e4b4f1344649dba56cb29090))

## [2.0.2](https://github.com/deadnews/deadnews-template-go/compare/v2.0.1...v2.0.2) - 2024-01-07

### Documentation

- _(readme)_ update badges - ([6c6735a](https://github.com/deadnews/deadnews-template-go/commit/6c6735a1d53cc4a02ae9c59de4fe3f92eeea178a))

### Build

- _(docker)_ update `Dockerfile` - ([1e825c5](https://github.com/deadnews/deadnews-template-go/commit/1e825c5aa9ee4c5b8f27b5ace518c726dab848c4))
- _(dockerfile)_ update pathes - ([2ac9e66](https://github.com/deadnews/deadnews-template-go/commit/2ac9e66809d7cf82750c0df3324408deec7d7ec5))

### Dependencies

- update module github.com/labstack/echo/v4 to v4.11.4 ([#51](https://github.com/deadnews/deadnews-template-go/issues/51)) - ([cb037ae](https://github.com/deadnews/deadnews-template-go/commit/cb037ae81263f354fa2fe054255d6aaf7cfdef89))

## [2.0.1](https://github.com/deadnews/deadnews-template-go/compare/v2.0.0...v2.0.1) - 2023-09-21

### Documentation

- _(readme)_ update badges - ([f5d121a](https://github.com/deadnews/deadnews-template-go/commit/f5d121a5112a3c0308920d0bc038eb6739b2efc6))
- _(readme)_ update badge - ([d04944b](https://github.com/deadnews/deadnews-template-go/commit/d04944b7cd4301d7553baca3ebea3c83a49dd289))
- _(readme)_ update badge - ([079a2af](https://github.com/deadnews/deadnews-template-go/commit/079a2af6ab083bb992c0dad6f02fd709898a33ee))

### Build

- _(docker)_ use more explicit tags - ([882ca05](https://github.com/deadnews/deadnews-template-go/commit/882ca0589a7dab7aa2f01dd295ef44bef50f8d31))
- _(docker)_ update stage alias - ([bf55221](https://github.com/deadnews/deadnews-template-go/commit/bf55221ff3d76404549d321c3e0355bbec73ff31))

### Chores

- use `reusable workflow` ([#35](https://github.com/deadnews/deadnews-template-go/issues/35)) - ([b718e3f](https://github.com/deadnews/deadnews-template-go/commit/b718e3ff9362fedf7ad1f2a05490243ce0b29756))

## [2.0.0](https://github.com/deadnews/deadnews-template-go/compare/v1.0.0...v2.0.0) - 2023-09-19

### Features

- change the sample application ([#33](https://github.com/deadnews/deadnews-template-go/issues/33)) - ([c22ecdf](https://github.com/deadnews/deadnews-template-go/commit/c22ecdf0fca8be184ddc461528334cee0fd8d39f))

### Chores

- _(dockerfile)_ add `label` - ([48bc20d](https://github.com/deadnews/deadnews-template-go/commit/48bc20dfbb5301c8d02c6da7ee53b69fc43fa605))
- updete `makefile` ([#34](https://github.com/deadnews/deadnews-template-go/issues/34)) - ([0252d60](https://github.com/deadnews/deadnews-template-go/commit/0252d602ef9acd17f0e99f233f05e57611a65c8c))
- enable more linters ([#31](https://github.com/deadnews/deadnews-template-go/issues/31)) - ([cb5b2cf](https://github.com/deadnews/deadnews-template-go/commit/cb5b2cfae46db812d3582ddff6788e24e3e8d07b))
- use `stable/oldstable` aliases in tests matrix ([#28](https://github.com/deadnews/deadnews-template-go/issues/28)) - ([9938420](https://github.com/deadnews/deadnews-template-go/commit/99384209c0b8b2247c24725b9b2a412258a0d587))

## [1.0.0](https://github.com/deadnews/deadnews-template-go/compare/v0.0.7...v1.0.0) - 2023-07-24

### Documentation

- fix `workflow` name - ([19759e7](https://github.com/deadnews/deadnews-template-go/commit/19759e7b07743ee9873ad7427d8773594d032a76))
- fix `workflow` name - ([54ea0d4](https://github.com/deadnews/deadnews-template-go/commit/54ea0d4c8c4180241a4286fa67b0e267b2271878))

### Chores

- _(renovate)_ adjust schedule - ([4ddda1e](https://github.com/deadnews/deadnews-template-go/commit/4ddda1e5c0670a53844145806aca62567ff92279))
- _(renovate)_ use shared config - ([6b46c46](https://github.com/deadnews/deadnews-template-go/commit/6b46c46320e33ce46330ae880c779a49954e66e2))
- _(renovate)_ adjust schedule - ([41b7867](https://github.com/deadnews/deadnews-template-go/commit/41b7867002f6f56f8085dbba4a786ebb1f81f3d2))
- clean up - ([8693b3b](https://github.com/deadnews/deadnews-template-go/commit/8693b3b2ea89cef678ca66b222b46900dacacb5f))
- disable `codeql` schedule - ([7e56036](https://github.com/deadnews/deadnews-template-go/commit/7e56036c15fed9972cb15c95ef4c857b13a4c946))
- enable more linters - ([4e83528](https://github.com/deadnews/deadnews-template-go/commit/4e83528a128ee8800672e864a0619e5540c1df61))
- use `digest pinning` ([#23](https://github.com/deadnews/deadnews-template-go/issues/23)) - ([79e7d5a](https://github.com/deadnews/deadnews-template-go/commit/79e7d5a5e97d29c5f150da9e377d551fe71bb287))
- update `workflows` ([#19](https://github.com/deadnews/deadnews-template-go/issues/19)) - ([46d7970](https://github.com/deadnews/deadnews-template-go/commit/46d7970a069287a5b44c06214f3dfa4f8f437993))

### Dependencies

- update module github.com/stretchr/testify to v1.8.4 ([#21](https://github.com/deadnews/deadnews-template-go/issues/21)) - ([ebb7b19](https://github.com/deadnews/deadnews-template-go/commit/ebb7b190cd8431de541ed4df60d810db4b333729))

## [0.0.7](https://github.com/deadnews/deadnews-template-go/compare/v0.0.6...v0.0.7) - 2023-05-04

### Features

- rename project - ([e74bbef](https://github.com/deadnews/deadnews-template-go/commit/e74bbef1c7853e88ecf1c533cee2fdd7a1d2470c))

### Chores

- _(pre-commit)_ comment `golangci-lint` hook - ([1554393](https://github.com/deadnews/deadnews-template-go/commit/155439337ebe1b981e60909a528de4904bed8d96))
- _(pre-commit)_ update hooks - ([b0f895d](https://github.com/deadnews/deadnews-template-go/commit/b0f895dd242c68bc8453ddf96ff08c9d4252df51))
- _(renovate)_ replace `dependabot` with `renovate` - ([7205a41](https://github.com/deadnews/deadnews-template-go/commit/7205a41b02a937d32add946e4b5cb5a97a4261b7))
- test `govulncheck` - ([fe39d19](https://github.com/deadnews/deadnews-template-go/commit/fe39d1933ee84df7abf08d10f92526feac770a50))

## [0.0.1](https://github.com/deadnews/deadnews-template-go/commit/v0.0.1) - 2022-09-06

<!-- generated by git-cliff -->
