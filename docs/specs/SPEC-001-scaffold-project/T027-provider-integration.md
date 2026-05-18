# T027: Provider Integration Tests and Multi-Source Precedence

**Phase**: 3 - Remote Providers  
**Category**: Testing  
**Complexity**: Medium  
**Estimated Duration**: 3-4 hours

## Objective

Establish integration tests for all remote providers and multi-source conflict resolution.

## Current State

- Phase 2 has local provider tests

## Target State

- All providers pass contract tests
- Multi-source precedence is validated
- Caching and retry behavior is tested
- Error scenarios are covered

## Acceptance Criteria

- [ ] Contract tests pass for all providers (git-tags, github-releases, ghcr)
- [ ] Multi-source chain respects precedence order
- [ ] Same version from different sources is merged correctly
- [ ] Caching reduces API calls
- [ ] Retry behavior works correctly
- [ ] Error handling across provider chain
- [ ] Mock implementations for CI testing

## Context

### Files to Create

- `tests/integration/providers_multi_test.go` — multi-provider tests
- `tests/integration/caching_test.go` — caching behavior tests

## Testing

- [ ] Integration test: All providers work together
- [ ] Integration test: Precedence is respected
- [ ] Integration test: Caching works correctly
- [ ] Integration test: Retries work correctly
- [ ] Error test: Provider failures are handled gracefully

## Related Tickets

- T021: Integration test patterns
- T023-T024: Remote providers

## Notes

- Use mock HTTP servers for testing remote providers
- Test rate limit handling
