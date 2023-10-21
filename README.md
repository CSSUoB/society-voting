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

### API

- [ ] Allow admin access to non-user-specific sections of the normal API
- [ ] Add `isRON` flag to `BallotEntry`

#### User

- [x] Change display name
- [x] Stand/withdraw from election
- [x] List all elections
- [ ] Display currently running election in /api/elections

#### Admin

- [x] Create election
- [x] Delete election

- [ ] Run election
  - [ ] Add `Ballot` table 
  - [ ] Accept extra ballot options
  - [ ] Store active election

### Frontend
