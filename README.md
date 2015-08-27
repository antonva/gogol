# Denton
A modular irc bot written in Go.

The idea is for a stable set of core functions while the majority of functionality
would be offloaded to drop-in golang plugins. 

### milestones:
- [ ] Core functionality definition.
- Core functionality:
  - [ ] A channel defined dictionary (by request) 
    - Requires additional message parser functionality, the bot should recognize it's own nickname.
    - Consider dropping trigger functionality for bot name utterance pattern instead.
  - [ ] Help functionality. 
  - [ ] Admin interface.
- Documentation:
  - [ ] Codebase.
  - [ ] Plugin interface.
  - [ ] Admin interface.
  - [ ] Plugin example.
- Tests:
  - [x] Here be dragons.
