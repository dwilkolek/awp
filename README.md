# AWS Web proxy aka AWP

This application can be used to setup your local environment and expose you services that are accessible via bastion ssh tunnel.

# Requirements

- 1password cli https://1password.com/downloads/command-line/
- aws cli https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html

# How to setup

KILL ZSCALER!

1. Downlaod file from `awp` from [releases](https://github.com/tfmcdigital/aws-web-proxy/releases/latest)
2. Create `awp` directory in your home directory: `mkdir -p ~/awp && cp ~/Downloads/awp ~/awp/awp`
3. Setup cli 1password: `op account add --address team-dsi.1password.com --email your.email@technipfmc.com`
4. If you already had it setup just execute `eval $(op signin)`
5. Go to AWP directory: `cd ~/awp`
6. Execute: `chmod +x awp`
7. Open containing directory in finder `open ~/awp`, ctrl+click on file and select Open with terminal
8. Execute: `awp setup`
9. Execute: `sudo awp hosts`

Optionally add alias and exeucte it from anywhere by `awp [command]`:

- `echo "alias awp=\"~/awp/awp\"" >> ~/.zshrc && source ~/.zshrc`

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
