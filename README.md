# AWS Service proxy

This application can be used to setup your local environment and expose you services that are accessible via bastion ssh tunnel.

# Requirements

- 1password cli https://1password.com/downloads/command-line/
- aws cli https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html

# How to use it

KILL ZSCALER!

0. Setup cli 1password: `op account add --address team-dsi.1password.com --email your.email@technipfmc.com`
1. If you already had it setup just execute `eval $(op signin)`
2. Downlaod file from release
3. Execute: `chmod +x ./aws-service-proxy`
4. Open containing directory in finder and Open with terminal
5. Execute: `./aws-service-proxy setup`
6. Execute: `sudo ./aws-service-proxy hosts`
7. Run application by `./aws-service-proxy start` or `./aws-service-proxy start dev|demo|prod`
8. All services are available locally on port 80 at `<servicename>.service` eg. http://material-match.service
