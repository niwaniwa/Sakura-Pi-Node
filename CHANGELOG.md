# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Added reed switch api [`#38`](https://github.com/niwaniwa/Sakura-Pi-Node/issues/38)
- Added door_switch_state_request setting(env) [`#40`](https://github.com/niwaniwa/Sakura-Pi-Node/issues/40)
- Added dev container setting. [`#43`](https://github.com/niwaniwa/Sakura-Pi-Node/issues/45)
- Added network setting. [`#46`](https://github.com/niwaniwa/Sakura-Pi-Node/issues/47)
- Added network setting. [`#54`](https://github.com/niwaniwa/Sakura-Pi-Node/issues/54)

### Changed
- Separate KeyState from DoorState[`#36`](https://github.com/niwaniwa/Sakura-Pi-Node/issues/36)
- MQTT logging mode[`#42`](https://github.com/niwaniwa/Sakura-Pi-Node/issues/44))
- Change hard-coded values to config file.[`#33`](https://github.com/niwaniwa/Sakura-Pi-Node/pull/51)
- Organize debug logs and optimize servo control.[`#52`](https://github.com/niwaniwa/Sakura-Pi-Node/pull/52)

### Fixed
- Fixed usb inaccessibility issue when using docker.[`#48`](https://github.com/niwaniwa/Sakura-Pi-Node/issues/48)

### Removed

## 0.2.0 - 2024-04-06
### Changed
- Changed and support new hardware. [`#35`](https://github.com/niwaniwa/Sakura-Pi-Node/issues/35)

## 0.1.0 - 2023-11-24
### Changed
- Changed to MQTT protocol. [`#15`](https://github.com/niwaniwa/Sakura-Pi-Node/issues/15)
- Changed to Clean architecture. [`#17`](https://github.com/niwaniwa/Sakura-Pi-Node/issues/17)

### Fixed
- Fixed access permission. [`#14`](https://github.com/niwaniwa/Sakura-Pi-Node/issues/14)

## 0.0.1 - 2023-11-14
### Added
- Migration from [OpenKeyByFelica](https://github.com/niwaniwa?tab=repositories)

### Changed
- Docker Containerization [**#9**](https://github.com/niwaniwa/Sakura-Pi-Node/issues/9)
