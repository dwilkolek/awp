# AWS Service proxy

This application can be used to setup your local environment and expose you services that are accessible via bastion ssh tunnel.

# Requirements

- 1password cli https://1password.com/downloads/command-line/
- aws cli https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html

# How to use it

0. Setup cli 1password: `op account add --address team-dsi.1password.com --email your.email@technipfmc.com`
1. If you already had it setup just execute `eval $(op signin)`
2. Downlaod file from release
3. Execute: `./aws-service-proxy setup`
4. Execute: `sudo ./aws-service-proxy hosts`
5. Run application by `./aws-service-proxy start` or `./aws-service-proxy start dev|demo|prod`
6. All services are available locally on port 80 at `<servicename>.service`
