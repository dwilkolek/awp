# AWS Web proxy aka AWP

This application can be used to setup your local environment and expose you services that are accessible via bastion ssh tunnel.

# Requirements

- 1password cli https://1password.com/downloads/command-line/
- aws cli https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html

# How to setup

KILL ZSCALER!

0. Setup cli 1password: `op account add --address team-dsi.1password.com --email your.email@technipfmc.com`
1. If you already had it setup just execute `eval $(op signin)`
2. Downlaod file from release
3. Execute: `chmod +x ./awp`
4. Open containing directory in finder and Open with terminal
5. Execute: `./awp setup`
6. Execute: `sudo ./awp hosts`

# How to start

All services are available locally on port 80 at `<servicename>.service` eg. http://material-match.service

- `./awp start` - starts proxy to dev cluster
- `./awp start dev` - start proxy to dev
- `./awp start demo` - start proxy to demo
- `./awp start prod` - start proxy to prod

# Commands

- `./awp setup` - creates aws profile to use `hosts` command and updates bastion keys from 1password
- `./awp hosts` - requires sudo, updates `/etc/hosts` with service list from aws
- `./awp update-keys` - updates keys to bastion from 1password
