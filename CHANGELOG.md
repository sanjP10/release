# Changelog

## 3.1.0
### Added
* Support for non-api git providers. Allowing for tagging via git.

## 3.0.1
### Updated
* Separated each provider to a package
* Moved all private packages to the internal directory
* Moved from `dep` to `go mod` for dependency management

## 3.0.0

### Added
* Github and Gitlab support

## 2.0.4

### Changed

* Print version number

## 2.0.3
### Fixed
* Trim space on desired tag

## 2.0.2

### Updated
* Added hashicorp go-version to handle semantic version checks

## 2.0.1

### Fixed
* Host override was not being set as a flag

### Updated
* Host override allows for schema override as well

## 2.0.0
* Removed `-tag` flag
* Added `-changelog` flag
* Retrieve tag name from changelog
* Added semantic validations
* Added additional stderr messages

## 1.7.0

### Fixed
* Removed commented out code 

## 1.6.0

### Added
* Add `-host` flag to allow for use with self hosted bitbucket

## 1.5.0

### Updated
* For `create` subcommand validate tag can be created
* Updated unit tests to use gock, which records http calls

## 1.4.0

### Updated
* Updated documentation

## 1.3.0

### Updated
* Updated documentation


## 1.2.0

### Fixed
* If a tag is redeployed on creation ensure that the error message is correct


## 1.1.0

### Updated
* Updated readme documentation

## 1.0.0

### Added

* Created subcommands `validate` and `create`
* Retrieve and create tags with bitbucket
* Added pipeline to execute linting, unit tests and integration tests
