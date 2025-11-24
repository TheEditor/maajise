# Maajise - High-Level Development Roadmap

**Version:** 2.0.0
**Date:** 2025-11-23
**Target:** Feature-Complete Multi-Language CLI Tool

---

## Vision

Transform Maajise from a monolithic single-purpose tool into a modular, extensible CLI that supports project initialization across multiple languages and frameworks with consistent tooling (Git, Beads, UBS).

---

## Current State Summary

- Architecture: 50% migrated (interfaces defined, not integrated)
- Features: Core init working, multi-language templates blocked
- Code Quality: Security fixed, tests in place
- Technical Debt: Main.go still monolithic (~700 lines)

---

## Development Phases

### Phase 1: Complete Architecture Migration
**Priority:** HIGH  
**Status:** In Progress (50%)  
**Estimated Effort:** 2-3 focused sessions

#### Goal
Finish the subcommand refactor started but not completed. Migrate main.go from monolithic structure to thin dispatcher using cmd/ packages.

#### Tasks

1. **Create cmd/init.go**
   - Move initialization logic from main.go
   - Implement Command interface
   - Register as "init" command
   - Handle all flags and options
   - Estimated: 1 session

2. **Create cmd/help.go**
   - Implement help command
   - Auto-discover registered commands
   - Display command list with descriptions
   - Show usage for specific commands
   - Estimated: 0.5 session

3. **Create cmd/version.go**
   - Simple version display command
   - Show build info if available
   - Estimated: 0.25 session

4. **Refactor main.go**
   - Reduce to thin dispatcher (~100 lines)
   - Route args to appropriate command
   - Implement convenience fallback (bare arg = init)
   - Handle errors gracefully
   - Estimated: 0.5 session

5. **Wire up internal packages**
   - Use internal/git instead of direct exec calls
   - Use internal/beads instead of direct exec calls
   - Move Config to internal/config
   - Use internal/ui for all output
   - Estimated: 0.5 session

#### Success Criteria
- [ ] main.go < 150 lines
- [ ] `maajise help` works
- [ ] `maajise version` works
- [ ] `maajise init my-project` works
- [ ] `maajise my-project` (convenience) works
- [ ] All tests pass
- [ ] Build succeeds

#### Why First
- Establishes clean foundation
- Removes technical debt
- Makes all future work easier
- Currently half-done creates confusion

---

### Phase 2: Template System Implementation
**Priority:** MEDIUM  
**Status:** Not Started (0%)  
**Estimated Effort:** 2-3 focused sessions

#### Goal
Enable multi-language project initialization with appropriate configurations for TypeScript, Python, Rust, PHP, and Go projects.

#### Tasks

1. **Implement Base Template**
   - Complete templates/base.go
   - Implement Template interface
   - Move current file generation logic
   - Create language-agnostic configs
   - Register in template registry
   - Estimated: 0.5 session

2. **Create TypeScript Template**
   - New file: templates/typescript.go
   - TypeScript-specific .gitignore
   - package.json structure
   - tsconfig.json
   - README with npm commands
   - Estimated: 0.5 session

3. **Create Python Template**
   - New file: templates/python.go
   - Python-specific .gitignore
   - requirements.txt / pyproject.toml
   - Virtual environment setup
   - README with pip commands
   - Estimated: 0.5 session

4. **Create Rust Template**
   - New file: templates/rust.go
   - Rust-specific .gitignore
   - Cargo.toml structure
   - src/main.rs skeleton
   - README with cargo commands
   - Estimated: 0.5 session

5. **Create PHP Template**
   - New file: templates/php.go
   - PHP-specific .gitignore
   - composer.json structure
   - README with composer commands
   - Estimated: 0.5 session

6. **Add Template Selection**
   - Add --template flag to init command
   - Default to "base" if not specified
   - Validate template exists
   - Pass template to file creation logic
   - Estimated: 0.25 session

7. **Create Template List Command**
   - New file: cmd/templates.go
   - List all registered templates
   - Show descriptions
   - Estimated: 0.25 session

#### Success Criteria
- [ ] `maajise init my-ts-app --template=typescript` works
- [ ] `maajise init my-py-app --template=python` works
- [ ] `maajise init my-rust-app --template=rust` works
- [ ] `maajise init my-php-app --template=php` works
- [ ] `maajise templates` lists all available templates
- [ ] Each template creates appropriate files
- [ ] All tests pass

#### Why Second
- You explicitly need multi-language support
- Clean architecture makes this straightforward
- High value-add for daily workflow
- Extends tool's utility significantly

---

### Phase 3: Additional Commands
**Priority:** MEDIUM  
**Status:** Not Started (0%)  
**Estimated Effort:** 1-2 focused sessions

#### Goal
Add utility commands for repo maintenance and validation.

#### Tasks

1. **Create Update Command**
   - New file: cmd/update.go
   - Refresh .gitignore, .ubsignore, README template
   - Respect --no-overwrite by default
   - Allow --force to overwrite
   - Estimated: 0.5 session

2. **Create Validate Command**
   - New file: cmd/validate.go
   - Check: Git initialized
   - Check: Beads initialized
   - Check: Required files present
   - Check: go.mod valid (if Go project)
   - Report findings with colored output
   - Estimated: 0.5 session

3. **Update Documentation**
   - Update README with all commands
   - Add usage examples for each command
   - Document all flags and options
   - Estimated: 0.25 session

#### Success Criteria
- [ ] `maajise update` refreshes config files
- [ ] `maajise validate` reports repo health
- [ ] README documents all commands
- [ ] Help text shows all commands

#### Why Third
- Adds polish and utility
- Not blocking core functionality
- Nice-to-have vs must-have
- Quick wins once architecture solid

---

### Phase 4: Advanced Features
**Priority:** LOW  
**Status:** Future Work  
**Estimated Effort:** 3-5 focused sessions

#### Goal
Add convenience features and extend customization.

#### Tasks

1. **Configuration File Support**
   - Support ~/.maajiserc
   - YAML/JSON format
   - Default template selection
   - Default flags
   - User preferences
   - Estimated: 1 session

2. **Custom Templates**
   - Support user-defined templates
   - Template directory: ~/.maajise/templates/
   - Template specification format
   - Template validation
   - Estimated: 1 session

3. **Pre-commit Hook Setup**
   - Optional --setup-hooks flag
   - Install pre-commit hooks
   - Configure for UBS scanning
   - Configure for tests
   - Estimated: 0.5 session

4. **CI/CD Configuration**
   - Generate .github/workflows/ configs
   - Support GitLab CI
   - Support Jenkins
   - Template-aware CI configs
   - Estimated: 1 session

5. **Docker Support**
   - Docker template
   - Generate Dockerfile
   - Generate docker-compose.yml
   - .dockerignore file
   - Estimated: 0.5 session

#### Success Criteria
- [ ] ~/.maajiserc configures defaults
- [ ] Custom templates loadable
- [ ] Pre-commit hooks installable
- [ ] CI configs generated
- [ ] Docker support available

#### Why Last
- Nice-to-haves, not essential
- Build on solid foundation
- Can be added incrementally
- User requests may change priorities

---

## Implementation Strategy

### Approach: Phased Delivery

**Phase 1 → Phase 2 → Phase 3 → Phase 4**

Each phase delivers complete, usable functionality before moving to next.

### Parallel Work Considerations

Some tasks can be done in parallel:
- Templates (Phase 2) are independent of each other
- Commands (Phase 3) are independent of each other
- Advanced features (Phase 4) are independent

But Phase 1 must complete first.

### Testing Strategy

Each phase includes:
- Unit tests for new functions
- Integration tests for new commands
- Manual testing checklist
- UBS scan after changes

### Documentation Strategy

Update docs after each phase:
- README.md for user-facing changes
- Code comments for complex logic
- Task specs archived in repo

---

## Milestones

### Milestone 1: Architecture Complete
**Target:** End of Phase 1  
**Deliverable:** Clean, modular subcommand architecture

**Definition of Done:**
- main.go < 150 lines
- All commands registered and working
- All internal packages in use
- Tests pass
- Build succeeds

### Milestone 2: Multi-Language Support
**Target:** End of Phase 2  
**Deliverable:** TypeScript, Python, Rust, PHP templates working

**Definition of Done:**
- All 5 templates (base + 4 languages) implemented
- Template selection via flag works
- Templates command lists all options
- Each template creates appropriate files
- Documentation updated

### Milestone 3: Feature Complete
**Target:** End of Phase 3  
**Deliverable:** Full CLI with update and validate commands

**Definition of Done:**
- Update command working
- Validate command working
- All commands documented
- README comprehensive
- No critical technical debt

### Milestone 4: Advanced Features
**Target:** End of Phase 4  
**Deliverable:** Configuration files, custom templates, hooks, CI

**Definition of Done:**
- Config file support working
- Custom templates loadable
- Hooks installable
- CI generation working
- Docker support available

---

## Resource Requirements

### Development Time
- Phase 1: 2-3 focused sessions
- Phase 2: 2-3 focused sessions
- Phase 3: 1-2 focused sessions
- Phase 4: 3-5 focused sessions

**Total:** 8-13 focused sessions

### Tools Required
- Go 1.23+
- Git
- Beads (bd)
- UBS (optional)
- Text editor / IDE

### No External Dependencies
- Pure Go standard library
- No npm packages
- No pip packages
- No cargo crates

---

## Risk Management

### Risk: Architecture Changes Break Features
**Mitigation:** Comprehensive test suite before refactoring

### Risk: Template Complexity
**Mitigation:** Start simple, iterate based on usage

### Risk: Scope Creep in Phase 4
**Mitigation:** Phase 4 is optional enhancements, can defer

### Risk: Breaking Changes for Users
**Mitigation:** Semantic versioning, maintain backwards compatibility

---

## Success Metrics

### Quantitative
- Lines of code in main.go: 700 → < 150
- Test coverage: 65% → 80%
- Number of templates: 0 → 5
- Number of commands: 1 → 7+

### Qualitative
- Code maintainability: Easier to add features
- User experience: Intuitive multi-language support
- Developer experience: Clean contribution path
- Documentation: Comprehensive and current

---

## Decision Points

### After Phase 1
**Decision:** Proceed to Phase 2 or address technical debt?

**Factors:**
- Are tests comprehensive?
- Is architecture clean?
- Any critical bugs?

### After Phase 2
**Decision:** Proceed to Phase 3 or add more templates?

**Factors:**
- Are 5 templates sufficient?
- User requests for other languages?
- Template quality acceptable?

### After Phase 3
**Decision:** Proceed to Phase 4 or consider project complete?

**Factors:**
- Is current feature set sufficient?
- User requests for advanced features?
- ROI on additional development?

---

## Maintenance Plan

### Post-Roadmap
- Monitor for bugs via GitHub issues
- Update templates as languages evolve
- Add templates for new languages as requested
- Keep dependencies updated (Go version)
- Maintain UBS scan compliance

### Community Involvement
- Accept template contributions
- Review feature requests
- Maintain issue tracker
- Keep documentation current

---

## Conclusion

This roadmap provides a structured path from current state (50% architecture migration) to feature-complete multi-language CLI tool. The phased approach ensures each stage delivers value while building toward the complete vision.

**Recommended Next Step:** Begin Phase 1 - Complete Architecture Migration

**Timeline:** With focused effort, Phases 1-3 achievable in 5-8 sessions, reaching feature-complete state. Phase 4 can be deferred or implemented based on user demand.
