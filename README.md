# Release

Release is a binary that validates and creates tags against bitbucket by reading your changelog file

Requires a markdown formatted changelog, with the most recent changes at the top.

The tht consists of a version must start with a `h2` markup and have a number afterwards.

An example changelog would be 

```
# Changelog

[//]: <> (Spaces an no spaces on version number lines are for checking regex in unit tests)
## 1.1.0

### Updated
* An update happened

## 1.0.0

### Added

* Initial release

```

the version numbers can be of a format with decimals separating them.

Example formats tha can be used are

```
major
major.minor
major.minor.patch
major.minor.patch.micro
```

***Note the format must be consistent within the changelog***

The two subcommands for release are `validate` and `create`

The flags for these commands are 

```
-username <username>
-password <password/authroization token> 
-repo <owner/org>/<repo name>
-changelog <changelog md file>
-hash <commit sha>
-host <bitbucket dns> (optional) (default is bitbucket.org)
```

This is an example `validate` command

```
release validate -username $USER -password $ACCESS_TOKEN -repo cloudreach/release -changelog changelog.md -hash e1db5e6db25ec6a8592c879d3ff3435c5503d03d
```

This is an example `create` command

```
release validate -username $USER -password $ACCESS_TOKEN -repo cloudreach/release -changelog changelog.md -hash e1db5e6db25ec6a8592c879d3ff3435c5503d03d
```

This is an example of `validate` command against a self-hosted bitbucket
```
release validate -username $USER -password $ACCESS_TOKEN -repo cloudreach/release -changelog changelog.md -hash e1db5e6db25ec6a8592c879d3ff3435c5503d03d -host api.mybitbucket.com
```

To integrate the `validate` use this in bitbucket pipelines you can use the following as steps

```
- step:
    name: validate version
    image: golang
    script:
      - CHANGELOG_FILE=$(pwd)/Changelog.md
      - PACKAGE_PATH="${GOPATH}/src/bitbucket.org/cloudreach"
      - mkdir -pv "${PACKAGE_PATH}"
      - cd "${PACKAGE_PATH}"
      - git clone https://$USER:$ACCESS_TOKEN@bitbucket.org/cloudreach/release
      - cd release
      - go get -u github.com/golang/dep/cmd/dep
      - dep ensure
      - go install
      # Test version does not exist
      - release validate -username $USER -password $ACCESS_TOKEN -repo $BITBUCKET_REPO_OWNER/$BITBUCKET_REPO_SLUG -changelog $CHANGELOG_FILE -hash $BITBUCKET_COMMIT

```

To integrte the `create` use this in the bitbucket pipeline after you merge to master

To integrate this into bitbucket pipelines you can use the following as steps

```
- step:
    name: create version
    image: golang
    script:
      - CHANGELOG_FILE=$(pwd)/Changelog.md
      - PACKAGE_PATH="${GOPATH}/src/bitbucket.org/cloudreach"
      - mkdir -pv "${PACKAGE_PATH}"
      - cd "${PACKAGE_PATH}"
      - git clone https://$USER:$ACCESS_TOKEN@bitbucket.org/cloudreach/release
      - cd release
      - go get -u github.com/golang/dep/cmd/dep
      - dep ensure
      - go install
      - release create -username $USER -password $ACCESS_TOKEN -repo $BITBUCKET_REPO_OWNER/$BITBUCKET_REPO_SLUG -changelog $CHANGELOG_FILE -hash $BITBUCKET_COMMIT
```
