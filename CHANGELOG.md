# CHANGELOG

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## Unreleased

### Added

- Background colors codes and template helpers
- `AllowEdit` for prompt to prevent deletion of the default value by any key

### Fixed

- `<Enter>` key press on Windows
- `juju/ansiterm` dependency
- `chzyer/readline#136` new api with ReadCloser

## [0.2.1] - 2017-11-30

### Fixed

- `SelectWithAdd` panicking on `.Run` due to lack of keys setup
- Backspace key on Windows

## [0.2.0] - 2017-11-16

### Added

- `Select` items can now be searched

## [0.1.0] - 2017-11-02

### Added

- extract `promptui` from [torus](https://github.com/manifoldco/torus-cli) as a
  standalone lib.
- `promptui.Prompt` provides a single input line to capture user information.
- `promptui.Select` provides a list of options to choose from. Users can
  navigate through the list either one item at time or by pagination
