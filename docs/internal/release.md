# Release Procedure

This document outlines the steps to be followed for releasing a new version of the Nuxeo Go Client.

## Pre-release: Every Contribution

1. **Update Changelog**: Ensure that all notable changes are documented in the `CHANGELOG.md` file following the [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format.
2. **Run Tests**: Execute all tests to ensure that the codebase is stable.
3. **Update Documentation**: Review and update any relevant documentation to reflect the changes in the new release.
4. **Commit Changes**: Commit all changes with a message that follows [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) guidelines.

## Release Manager: Finalizing the Release

1. **Determine Version**: Decide the new version number based on the changes made (major, minor, patch) following [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
2. **Update Changelog**: Move the changes from the "Unreleased" section to a new section with the determined version number and release date.
3. **Tag the Release**:
   ```bash
   git tag -a vX.Y.Z -m "Release version X.Y.Z"
   git push origin vX.Y.Z
   ```
4. The Github Actions workflow will automatically trigger an re-indexing of pkg.go.dev upon pushing the new tag. Check the [Nuxeo Go Client page on pkg.go.dev](https://pkg.go.dev/github.com/anselm94/nuxeo-go-client) to ensure the new version is listed correctly after about 20 mins.