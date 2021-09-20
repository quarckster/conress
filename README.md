# Cypress container image

This repository contains Dockerfiles for building a Fedora based container image with installed
Cypress binary, Google Chrome and Mozilla Firefox browsers as welll as minimia; amount of required
dependencies.

## Building arguemnts

`FIREFOX_VERSION` - desired version of Mozilla Firefox.

`CHROME_VERSION` - desired version of Google Chrome.

`CYPRESS_VERSION` - desired version of Cypress binary.

## Pull the image

You can pull the image from here:

`quay.io/redhatqe/cypress:latest`

All tags are available on this page:

<https://quay.io/repository/redhatqe/cypress?tab=tags>

## Usage

Mount a directory with the application source code and run Cypress using `startcypress`
command:

```sh
podman run -it --rm -p 5999:5999 --shm-size=2g -w /mnt -v .:/mnt quay.io/redhatqe/cypress:latest bash
startcypress open
```

`startcypress` starts and stop `Xvnc`, `fluxbox` and `Cypress` in the right order. It passes all
arguments to `Cypress`.
