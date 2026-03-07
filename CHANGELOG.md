# Changelog

## [3.0.0](https://github.com/manaelproxy/manael/compare/v2.1.0...v3.0.0) (2026-03-07)


### ⚠ BREAKING CHANGES

* replace go-libwebp and go-avif with bimg (libvips) ([#1579](https://github.com/manaelproxy/manael/issues/1579))

### Features

* Add dedicated health check endpoint (/_health) ([#1632](https://github.com/manaelproxy/manael/issues/1632)) ([75d2091](https://github.com/manaelproxy/manael/commit/75d2091d57576cd5186b31645ec2010e6fb00829))
* do not forward manael transform query params to upstream ([#1651](https://github.com/manaelproxy/manael/issues/1651)) ([e1f71f2](https://github.com/manaelproxy/manael/commit/e1f71f272b2f410f1a8e12d0f6574e60a0139392))
* Expose Prometheus metrics on a dedicated internal port ([#1631](https://github.com/manaelproxy/manael/issues/1631)) ([b7d4a25](https://github.com/manaelproxy/manael/commit/b7d4a2500875fd1628325ff78f7a6cd615f10662))
* Format-specific encode quality control via `q` query parameter ([#1642](https://github.com/manaelproxy/manael/issues/1642)) ([a70f0ad](https://github.com/manaelproxy/manael/commit/a70f0ad305588c3bfe515fde59bca49a4936d4c5))
* Image resize and crop via query parameters with security limits ([#1636](https://github.com/manaelproxy/manael/issues/1636)) ([aefe363](https://github.com/manaelproxy/manael/commit/aefe363dc53f757f7ea1cfbfa3eb969ec4128256))
* implement graceful shutdown for the HTTP server ([#1608](https://github.com/manaelproxy/manael/issues/1608)) ([6c4bbc3](https://github.com/manaelproxy/manael/commit/6c4bbc35a3451c68e14f87eecdec1913ee213de4))
* Instrument HTTP proxy with OpenTelemetry ([#1610](https://github.com/manaelproxy/manael/issues/1610)) ([deecbc5](https://github.com/manaelproxy/manael/commit/deecbc54ac5b85a7e70d3b60aeacba901f00af50))
* make image resizing opt-in via MANAEL_ENABLE_RESIZE env var ([#1647](https://github.com/manaelproxy/manael/issues/1647)) ([3393ffc](https://github.com/manaelproxy/manael/commit/3393ffc02b79309d8442de9f24188640ac70e2c3))
* migrate Go vanity URL resolution to Cloudflare Pages Functions ([#1640](https://github.com/manaelproxy/manael/issues/1640)) ([15f1c24](https://github.com/manaelproxy/manael/commit/15f1c247140747bad3538a7032d9b0cb4714f780))
* migrate to structured logging using log/slog ([#1620](https://github.com/manaelproxy/manael/issues/1620)) ([5292627](https://github.com/manaelproxy/manael/commit/529262708a4efdf9ae9c646c6ce7272527fbf601))
* replace go-libwebp and go-avif with bimg (libvips) ([#1579](https://github.com/manaelproxy/manael/issues/1579)) ([9c3fe13](https://github.com/manaelproxy/manael/commit/9c3fe130570b91ea98261790a4b23e9368b8f41c))


### Bug Fixes

* pass ProxyOptions to scanAcceptHeader and check ([#1634](https://github.com/manaelproxy/manael/issues/1634)) ([03c1734](https://github.com/manaelproxy/manael/commit/03c1734413adf4783fc9a39bb62c20ccf493eb91))

## [2.1.0](https://github.com/manaelproxy/manael/compare/v2.0.6...v2.1.0) (2026-03-04)


### Features

* add APNG passthrough support ([#1570](https://github.com/manaelproxy/manael/issues/1570)) ([9e2de45](https://github.com/manaelproxy/manael/commit/9e2de45b75a44240c9f71fc957e574a6018a3f2b))
* convert static GIFs to WebP, pass through animated GIFs ([#1574](https://github.com/manaelproxy/manael/issues/1574)) ([438b633](https://github.com/manaelproxy/manael/commit/438b6330335c591fbd6d10fbd0064bbd398bdf3f))
* update Content-Disposition filename extension when converting images ([#1571](https://github.com/manaelproxy/manael/issues/1571)) ([f865760](https://github.com/manaelproxy/manael/commit/f865760467fdddea9390654c04c52ae94ea37c21))
* **website:** replace Docusaurus with Hugo and Docsy ([#1564](https://github.com/manaelproxy/manael/issues/1564)) ([87b40d2](https://github.com/manaelproxy/manael/commit/87b40d2a43136dbf9cf549310b3ea6006bce8318))

## [2.0.6](https://github.com/manaelproxy/manael/compare/v2.0.5...v2.0.6) (2024-12-08)


### Bug Fixes

* **docker:** fix syntax for dockerfile ([#1354](https://github.com/manaelproxy/manael/issues/1354)) ([5638f9d](https://github.com/manaelproxy/manael/commit/5638f9d0f7f1a4714e00ba8466bbcd193fb4578a))

## [2.0.5](https://github.com/manaelproxy/manael/compare/v2.0.4...v2.0.5) (2023-12-22)


### Bug Fixes

* **docker:** add missing `main.go` ([#1080](https://github.com/manaelproxy/manael/issues/1080)) ([45b6dd6](https://github.com/manaelproxy/manael/commit/45b6dd6617663846b3de3005ebe8be122bb6a9df))


### Reverts

* **deps:** downgrade module github.com/harukasan/go-libwebp ([#1078](https://github.com/manaelproxy/manael/issues/1078)) ([b586e5f](https://github.com/manaelproxy/manael/commit/b586e5fa38a332100cd2e40e21aad801342c0cd2))

## [2.0.4](https://github.com/manaelproxy/manael/compare/v2.0.3...v2.0.4) (2023-12-20)


### Bug Fixes

* always set vary header ([#1075](https://github.com/manaelproxy/manael/issues/1075)) ([ef3d47f](https://github.com/manaelproxy/manael/commit/ef3d47f4cfc7e06143c073b6200a7017df067e52))

## [2.0.3](https://github.com/manaelproxy/manael/compare/v2.0.2...v2.0.3) (2023-12-20)


### Bug Fixes

* add build workflow ([#1072](https://github.com/manaelproxy/manael/issues/1072)) ([112232c](https://github.com/manaelproxy/manael/commit/112232c72a8c5826607c61104ecb16e406ca4255))

## [2.0.2](https://github.com/manaelproxy/manael/compare/v2.0.1...v2.0.2) (2023-12-20)


### Bug Fixes

* **docker:** add `cmake` ([#1070](https://github.com/manaelproxy/manael/issues/1070)) ([bc6dd34](https://github.com/manaelproxy/manael/commit/bc6dd344ce8b0de8e02a5e3e7802f8cc788bc382))
* **docker:** add `cmake` ([#1071](https://github.com/manaelproxy/manael/issues/1071)) ([aa12a63](https://github.com/manaelproxy/manael/commit/aa12a63bc20f023400e342d19a8a1474f757a1d7))
* **docker:** add missing `&` ([#1068](https://github.com/manaelproxy/manael/issues/1068)) ([9523d2d](https://github.com/manaelproxy/manael/commit/9523d2d2746d8a9d0f3a737326b1f651205395b6))
* **docker:** add missing `&` (2) ([#1069](https://github.com/manaelproxy/manael/issues/1069)) ([fc5025c](https://github.com/manaelproxy/manael/commit/fc5025c5258188c1eeae2d0239fb874eb58e876a))
* **docker:** build on docker ([#1065](https://github.com/manaelproxy/manael/issues/1065)) ([57fb700](https://github.com/manaelproxy/manael/commit/57fb700e8350ed7cbb96a15f5e116ec774841a09))
* **docker:** remove `sudo` ([#1067](https://github.com/manaelproxy/manael/issues/1067)) ([631c24e](https://github.com/manaelproxy/manael/commit/631c24e5a3ecd59daaebd6f3e3bcecdfac9dd152))
* **github-actions:** fix release flow ([#1064](https://github.com/manaelproxy/manael/issues/1064)) ([52c6010](https://github.com/manaelproxy/manael/commit/52c6010df710b18b5d796326439d536ed1de4039))
* **release-please:** fallback release_created ([#1066](https://github.com/manaelproxy/manael/issues/1066)) ([b316e59](https://github.com/manaelproxy/manael/commit/b316e595e88555e33531210d738ae874a198c884))
* remove goreleaser ([#1062](https://github.com/manaelproxy/manael/issues/1062)) ([33bf0b2](https://github.com/manaelproxy/manael/commit/33bf0b201b8dbf2a842bf730995113b327ea47fb))

## [2.0.1](https://github.com/manaelproxy/manael/compare/v2.0.0...v2.0.1) (2023-12-20)


### Bug Fixes

* **goreleaser:** fix invalid syntax ([#1060](https://github.com/manaelproxy/manael/issues/1060)) ([f78bff3](https://github.com/manaelproxy/manael/commit/f78bff38975965cc223f4cad1a5b6975e3c311b1))

## [2.0.0](https://github.com/manaelproxy/manael/compare/v1.9.1...v2.0.0) (2023-12-20)


### ⚠ BREAKING CHANGES

* replace to `httputil.ReverseProxy` ([#1059](https://github.com/manaelproxy/manael/issues/1059))
* **docker:** remove docker hub ([#1058](https://github.com/manaelproxy/manael/issues/1058))
* replace to pnpm ([#1047](https://github.com/manaelproxy/manael/issues/1047))

### Features

* **docker:** remove docker hub ([#1058](https://github.com/manaelproxy/manael/issues/1058)) ([50d85c8](https://github.com/manaelproxy/manael/commit/50d85c8ec507b16dec88cd0c2c38068122aacd0e))
* replace to `httputil.ReverseProxy` ([#1059](https://github.com/manaelproxy/manael/issues/1059)) ([62a86b6](https://github.com/manaelproxy/manael/commit/62a86b6cf44d1c5e34f613cc3c73be80c516d9bf)), closes [#1054](https://github.com/manaelproxy/manael/issues/1054)


### Bug Fixes

* **release-please:** remove legacy property ([#1048](https://github.com/manaelproxy/manael/issues/1048)) ([515ca51](https://github.com/manaelproxy/manael/commit/515ca516b5e447126634bece4a34188fce71d53b))


### Code Refactoring

* replace to pnpm ([#1047](https://github.com/manaelproxy/manael/issues/1047)) ([0226430](https://github.com/manaelproxy/manael/commit/0226430a061f54e66db1b5e91d75ee4013d5a7fb))

### [1.9.1](https://github.com/manaelproxy/manael/compare/v1.9.0...v1.9.1) (2022-04-17)


### Bug Fixes

* **release:** fix path ([#695](https://github.com/manaelproxy/manael/issues/695)) ([1f3f36a](https://github.com/manaelproxy/manael/commit/1f3f36a8c962eb59f8fb891c17235e19a2c3e1aa))

## [1.9.0](https://github.com/manaelproxy/manael/compare/v1.8.5...v1.9.0) (2022-04-17)


### Features

* **deps:** update libwebp and libaom ([#693](https://github.com/manaelproxy/manael/issues/693)) ([cfbc541](https://github.com/manaelproxy/manael/commit/cfbc541604e3997eb6322d7e035c07cdeeff4aec))


### Bug Fixes

* **website:** disable trailing slash ([#681](https://github.com/manaelproxy/manael/issues/681)) ([5882d8a](https://github.com/manaelproxy/manael/commit/5882d8a5c7e6b2a086eddce2c684db8054501f1f))
* **website:** rename pkg url ([#684](https://github.com/manaelproxy/manael/issues/684)) ([24274a2](https://github.com/manaelproxy/manael/commit/24274a20bac64ecfa557f447fda5446abf0f563c))


### [1.8.5](https://www.github.com/manaelproxy/manael/compare/v1.8.4...v1.8.5) (2021-05-20)


### Bug Fixes

* **transport:** fix duplicate variables ([#464](https://www.github.com/manaelproxy/manael/issues/464)) ([dd1f3d5](https://www.github.com/manaelproxy/manael/commit/dd1f3d573e41d94653c1d1e9fbebdd177ce6c6ee))

### [1.8.4](https://www.github.com/manaelproxy/manael/compare/v1.8.3...v1.8.4) (2021-05-20)


### Bug Fixes

* **transport:** disable avif when png ([#462](https://www.github.com/manaelproxy/manael/issues/462)) ([c293444](https://www.github.com/manaelproxy/manael/commit/c293444dc83670a61d53f5c1f035ec9d649abaa2))

### [1.8.3](https://www.github.com/manaelproxy/manael/compare/v1.8.2...v1.8.3) (2021-05-15)


### Bug Fixes

* **release:** add -lm to ldflags ([#453](https://www.github.com/manaelproxy/manael/issues/453)) ([ae591af](https://www.github.com/manaelproxy/manael/commit/ae591afe12f97257dc18bd31030535451e8af760))

### [1.8.2](https://www.github.com/manaelproxy/manael/compare/v1.8.1...v1.8.2) (2021-05-14)


### Bug Fixes

* **transport:** change variable name ([#447](https://www.github.com/manaelproxy/manael/issues/447)) ([7b14d20](https://www.github.com/manaelproxy/manael/commit/7b14d203c38b3d9e1da98614efadadb2bed0c26e))

### [1.8.1](https://www.github.com/manaelproxy/manael/compare/v1.8.0...v1.8.1) (2021-05-14)


### Bug Fixes

* **release:** make directory for libaom ([#442](https://www.github.com/manaelproxy/manael/issues/442)) ([2790b2f](https://www.github.com/manaelproxy/manael/commit/2790b2f233d496eb21466329f3906e7b917add67))

## [1.8.0](https://www.github.com/manaelproxy/manael/compare/v1.7.1...v1.8.0) (2021-05-14)


### Features

* **avif:** add support avif ([#372](https://www.github.com/manaelproxy/manael/issues/372)) ([f2721d9](https://www.github.com/manaelproxy/manael/commit/f2721d99bb5f831237e49f8daa7994874e9efee6))
* **i18n:** add japanese translations ([#408](https://www.github.com/manaelproxy/manael/issues/408)) ([d4034b4](https://www.github.com/manaelproxy/manael/commit/d4034b4a4812d4fde952f4ffcef8900a28544e3b))
* **transprot:** add flag for avif ([#441](https://www.github.com/manaelproxy/manael/issues/441)) ([37cea0f](https://www.github.com/manaelproxy/manael/commit/37cea0fab3f45fb58fe90dbab103bc24e09aa3d8))
* **website:** enable docsearch ([#409](https://www.github.com/manaelproxy/manael/issues/409)) ([959d83a](https://www.github.com/manaelproxy/manael/commit/959d83a000458e0854c25666600bc23d823487b0))


### Bug Fixes

* **deps:** change release tag to latest ([#356](https://www.github.com/manaelproxy/manael/issues/356)) ([95b59cb](https://www.github.com/manaelproxy/manael/commit/95b59cb5426f7b0daee491ead9ad5a2eeb9e3c24))
* **website:** add missing scripts ([#407](https://www.github.com/manaelproxy/manael/issues/407)) ([d3c4bb1](https://www.github.com/manaelproxy/manael/commit/d3c4bb1f274ce5fd047106027bdc0ef354822bee))
