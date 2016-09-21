# mirrorhub

Redirects the user to a nearby mirror site.

An instance is [mirrors.moe](https://mirrors.moe).

## Configuration

The configuration is a JSON file and should be specified with the `-conf` option when running mirrorhub.

It should be an object with the following members:

* `sites`: an object mapping site names to site configurations, where a site configuration is an object with the following members:

  * `url`: the full URL for the site. A trailing slash is not needed. Example: `https://mirrors.ustc.edu.cn`

  * `distros`: a list of distros available on the site.

* `routes`: a list of routes, where every route is an object with the following members:

  * `ipnet`: an IP subnet in CIDR notation.

  * `ordering`: an ordering of site names appropriate for this subnet.

* `default-ordering`: the default ordering of sites.
