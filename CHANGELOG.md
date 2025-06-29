# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
- Resolved idx_location index conflicts in reviews table migrations
- Enhanced BDD test security with database URL consistency and debug logging
- Improved test environment setup with better error handling and security measures

### Changed
- All migration files synchronized between development and production
- CI stability improvements for GitHub Actions workflows

## [2025-06-28] - Security Update

### Fixed
- **Critical**: Fixed Huma v2 authentication middleware silent context propagation failure
- Protected API endpoints now properly authenticate users
- Implemented proper `huma.WithValue()` context handling pattern

### Security
- All authentication systems now fully functional and production-ready
- Enhanced security features for API endpoints