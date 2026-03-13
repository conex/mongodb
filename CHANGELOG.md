# Changelog

## [Unreleased]

### Changed
- Updated default image from `mongo:3` to `mongo:5` (mongo:6+ not compatible with mgo.v2 driver)
- Added connection retry logic for MongoDB startup
- Added nil config check
- Added `Wait()` call to ensure MongoDB is ready before connecting
- Updated to use Go modules
