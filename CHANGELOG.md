# Changelog

## 1.7

### Fixed
* Removed commented out code 

## 1.6

### Added
* Add `-host` flag to allow for use with self hosted bitbucket

## 1.5

# Updated
* For `create` subcommand validate tag can be created
* Updated unit tests to use gock, which records http calls

## 1.4

### Updated
* Updated documentation

## 1.3

### Updated
* Updated documentation


## 1.2

### Fixed
* If a tag is redeployed on creation ensure that the error message is correct


## 1.1

### Updated
* Updated readme documentation

## 1.0

### Added

* Created subcommands `validate` and `create`
* Retrieve and create tags with bitbucket
* Added pipeline to execute linting, unit tests and integration tests