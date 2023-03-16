# Superbchat

- Self-hosted, noncustodial and minimalist Bitcoin Cash (BCH) superchat system written in Go.
- Provides an admin view page to see donations with corresponding comments.
- Provides notification methods usable in OBS with an HTML page.

To see a working instance of superbchat, see [superbchat-vulkan0n.fly.dev](https://superbchat-vulkan0n.fly.dev/).

# Installation

1. ```apt install golang```
2. ```git clone https://github.com/vulkan0n/superbchat.git```
3. ```cd superbchat```
4. ```go install github.com/skip2/go-qrcode@latest```
5. edit ```config.json```
6. ```go run main.go```

A webserver at 127.0.0.1:8900 is running.

# Usage
- The BCH address should have the format `"bitcoincash:..."`
- Visit 127.0.0.1:8900/view to view your superchat history
- Visit 127.0.0.1:8900/alert?auth=adminadmin to see notifications
- The default username is `admin` and password `adminadmin`. Change these in `config.json`
- Edit web/index.html and style/style.css to customize your front page!

# OBS

- Add a Browser source in obs and point it to `https://example.com/alert?auth=adminadmin`

# Fly.io
You can try to deploy it fly.io for free. You will need to download their [CLI app](https://fly.io/docs/hands-on/install-flyctl/) `flyctl`.
But you will need to create a volume for the app with :
- ```fly volumes create example_data --size 1```
- Just 1GB should be more than enough.
- Then add the corresponding [mounts section](https://fly.io/docs/reference/volumes/#using-volumes) to the ```fly.toml``` file

# License

GPLv3

### Origin

This comes from [https://git.sr.ht/~anon_/shadowchat](https://git.sr.ht/~anon_/shadowchat) and is not Vulkan0n's original
work.

### Donate

sir,thank you
`bitcoincash:qpwcmmp5636akgpz8m07zhefgykczegekv5hvwsd24`
