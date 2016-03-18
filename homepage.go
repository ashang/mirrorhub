package main

import (
	"bytes"
	"html/template"
)

const homepageTemplate = `
<!doctype html>
<html>
<head>
	<title>mirrorhub</title>
	<style>
	body {
		margin: 16px 6.25%;
	}
	hr {
		border: none;
		height: 1px;
		color: #888;
		background-color: #888;
		margin: 1em 0;
	}
	@media screen and (min-width: 1024px) {
		body {
			width: 896px;
			margin: 16px auto;
		}
	}
	</style>
</head>
<body>
	<h1>mirrorhub</h1>
	<h2>Sites</h2>
	{{ range $name, $site := .Sites }}
		<a href="{{ $site.URL }}">{{ $name }}</a>
	{{ end}}
	<h2>Software</h2>
	{{ range $distro, $_ := .Distros }}
		<a href="/{{ $distro }}">{{ $distro }}</a>
	{{ end }}
	<hr/>
	<footer>
	A <a href="https://tuna.moe">TUNA</a> project. Source code on <a href="https://github.com/tuna/mirrorhub">Github</a>.
	</footer>
</body>
</html>
`

func (conf *Config) makeHomepage() []byte {
	t := template.Must(template.New("homepage").Parse(homepageTemplate))
	var buf bytes.Buffer
	t.Execute(&buf, conf)
	return buf.Bytes()
}
