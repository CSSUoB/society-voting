# Society Voting

*Online voting designed for student groups*

---

## TODO

### System

- [x] Database setup
- [x] Account provisioning and login
- [x] Guild data scraping
- [x] Events
  - [ ] Discord webhook event notifier
- [ ] **Use database transactions!**
- [ ] Change package namespace
- [ ] Save election results to dedicated table

### API

- [ ] Allow admin access to non-user-specific sections of the normal API
- [x] Add `isRON` flag to `BallotEntry`
- [ ] Vote validation code

#### User

- [x] Change display name
- [x] Stand/withdraw from election
- [x] List all elections
- [x] Display currently running election in /api/elections
- [x] Vote endpoint
- [ ] Make only the main election list endpoint return the candidate list

#### Admin

- [x] Create election
- [x] Delete election
- [x] Run election
  - [x] Add `Ballot` table 
  - [x] Accept extra ballot options
  - [x] Create ballot in setup endpoint
  - [x] Store active election
- [x] Stop and finalise election
- [ ] Election status SSE

### Frontend
