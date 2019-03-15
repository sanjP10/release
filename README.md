# Release

Release is a binary that validates and creates tags against bitbucket

The two subcommands for release are `validate` and `create`

The flags for these commands are 

```
-username <username>
-password <password/authroization token> 
-repo <owner/org>/<repo name>
-tag <tag name>
-hash <commit sha>
```

This is an example `validate` command

```
release validate -username $USER -password $ACCESS_TOKEN -repo cloudreach/release -tag initial -hash e1db5e6db25ec6a8592c879d3ff3435c5503d03d
```

This is an example `create` command

```
release validate -username $USER -password $ACCESS_TOKEN -repo cloudreach/release -tag initial -hash e1db5e6db25ec6a8592c879d3ff3435c5503d03d
```

To integrate the `validate` use this in bitbucket pipelines you can use the following as steps

```
- step:
    name: validate version
    image: golang
    script:
      - VERSION_FILE=$(pwd)/version.json
      - PACKAGE_PATH="${GOPATH}/src/bitbucket.org/cloudreach"
      - mkdir -pv "${PACKAGE_PATH}"
      - cd "${PACKAGE_PATH}"
      - git clone https://$USER:$ACCESS_TOKEN@bitbucket.org/cloudreach/release
      - cd release
      - go get -u github.com/golang/dep/cmd/dep
      - dep ensure
      - go install
      - apt-get update && apt-get install -y && apt-get install jq -y
      - VERSION=$(cat $VERSION_FILE | jq '.version'| tr -d \")
      # Test version does not exist
      - release validate -username $USER -password $ACCESS_TOKEN -repo $BITBUCKET_REPO_OWNER/$BITBUCKET_REPO_SLUG -tag $VERSION -hash $BITBUCKET_COMMIT

```

To integrte the `create` use this in the bitbucket pipeline after you merge to master

To integrate this into bitbucket pipelines you can use the following as steps

```
- step:
    name: create version
    image: golang
    script:
      - VERSION_FILE=$(pwd)/version.json
      - PACKAGE_PATH="${GOPATH}/src/bitbucket.org/cloudreach"
      - mkdir -pv "${PACKAGE_PATH}"
      - cd "${PACKAGE_PATH}"
      - git clone https://$USER:$ACCESS_TOKEN@bitbucket.org/cloudreach/release
      - cd release
      - go get -u github.com/golang/dep/cmd/dep
      - dep ensure
      - go install
      - apt-get update && apt-get install -y && apt-get install jq -y
      - VERSION=$(cat $VERSION_FILE | jq '.version'| tr -d \")
      - release create -username $USER -password $ACCESS_TOKEN -repo $BITBUCKET_REPO_OWNER/$BITBUCKET_REPO_SLUG -tag $VERSION -hash $BITBUCKET_COMMIT
```
