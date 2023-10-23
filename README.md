# Society Voting

*Online voting designed for student groups*

---

## Deployment instructions

The only production deployment method that's supported is in Docker.

Pull the image from `ghcr.io/cssuob/society-voting:latest` and then run it with:
 - a valid configuration file (see below) mounted as `/run/config.yml` inside the container
 - the port 8080 inside the container mapped to wherever you need it outside of the container (unless otherwise changed in the config file)

## Sample configuration

Ensure that the following keys are changed:
 - `guild.sessionToken` should be set to the value of the `.ASPXAUTH` cookie from the Guild of Students website
 - `platform.discordWebhook.url` should be set to a Discord webhook URL. This is optional.
 - `platform.discordWebhook.threadID` should be set to a Discord thread ID that webhook messages should be sent in. This is optional. If specified, the thread should be in the same channel that the webhook is associated with.

```yaml
http:
  host: "0.0.0.0"
guild:
  sessionToken: "Your Guild of Students session token from .ASPXAUTH"
  societyID: "6531"
platform:
  societyName: "CSS"
  adminToken: "plzletmeadmin"
  discordWebhook:
    url: "Discord webhook URL (optional)"
    threadID: "Discord thread ID (optional)"
```
## Building the Docker image

### Remotely

Anything that's pushed to `build` will be automatically built - ie:

```
git switch build
git merge main
git push -u origin build
```

### Locally

```yaml
docker build -t socvot .
```