# mr-pump
A modular irc bot written in Go.

The idea is for a stable set of core functions while the majority of functionality
would be offloaded to drop-in Lua modules/plugins. The decision to incorporate the
sed function as core stemmed from it being treated differently in the message parser.
This is not final and may very well be substituted for a Lua module instead.

### milestones:
- [ ] Core functionality definition.
- Core functionality:
  - [x] sed replace for channels.
  - [ ] A channel defined dictionary (by request) 
    - Requires additional message parser functionality, the bot should recognize it's own nickname.
    - Consider dropping trigger functionality for bot name utterance pattern instead.
  - [ ] help functionality. 
  - [ ] Admin interface.
  - [ ] Lua plugin interface.
    - Work has begun on implementation.
- Documentation:
  - [ ] Codebase.
  - [ ] Plugin interface.
  - [ ] Admin interface.
  - [ ] Plugin example.
- Tests:
  - [x] Here be dragons.
