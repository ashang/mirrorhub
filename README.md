# mirrorhub

Redirects the user to a nearby mirror site.

An instance is [mirror.moe](https://mirror.moe).

## Basics

* 302 based
* redirect user to the closest mirror site
* powered by golang

## Architecture

All components communicate with each other via database.

* user-redirector
* info-fetcher, which fetch ip and gpg data required from the
  open Internet.
* etc (see features below)

## Features (to be implemented)

* measure mirror state
    + passive monitor via 3-times redirect
    + active monitor via check_mk (maybe)

* user interface
    + restful API
    + human-friendly web interface based on restful API

## Roadmap (if not dream)

* decentralized
* multi-CDN
* (That's enough
