package templates

const DefaultTemplate = `import http from 'k6/http';
import { check, group } from 'k6';
import { uuidv4, randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { URL } from 'https://jslib.k6.io/url/1.0.0/index.js';

export const options = {
  thresholds: {
    {{- range .Paths }}
      {{- $pathname := .Pathname }}
    'http_req_duration{group:::{{ $pathname }}}': [
      {{- range .Detail.Latency -}}
        {{- if .IsK6Supported -}}
          '{{- . -}}', 
        {{- end -}}
      {{- end -}}
    ],
        
      {{- range .Detail.ErrorRate }}
        {{- if .IsK6Supported }}
    'http_req_failed{group:::{{ $pathname }}}': ['{{ . }}'],
        {{- end }}
      {{- end }}
    {{- end }}
  },
  {{- if .Config.Users }}
  vus: {{ .Config.Users }},
  {{- end }}
  {{- if .Config.Duration }}
  duration: "{{ .Config.Duration }}s"
  {{- end }}
  {{- if .Config.Stages }}
  stages: [
    {{- range .Config.Stages }}
    { target: {{ .Target }}, duration: "{{ .Duration }}" }
    {{- end }}
  ]
  {{- end }}
}

const pickRandom = (array) => {
  return array[Math.floor(Math.random() * array.length)]
}

const randomBool = () => {
  return Math.random() < 0.5;
}

const applyPathParams = (path, params) => {
  return path
    .split('/')
    .map(s => s.startsWith('{') && s.endsWith('}')
      ? params[s.substr(1, s.length - 2)]
      : s)
    .join('/');
}

const generateFromRange = (min, max, step = 1) => {
  return Math.round(Math.random() * (max - min) + min);
}

const generateFromPattern = (pattern) => {
  const stringPattern = /string\((\d+)\)/

  if (pattern === 'uuid') {
    return uuidv4();
  } else if (stringPattern.test(pattern)) {
    const length = parseInt(pattern.match(stringPattern)[1]);

    return randomString(length);
  } else if (pattern === 'bool') {
    return randomBool();
  }
}

{{ define "params" }}
  {{- range $name, $value := . }}
    {{- if $value.Examples }}
      {{ $name }}: pickRandom([
        {{- range $value.Examples -}}
          '{{- . -}}',
        {{- end -}}
      ])
    {{- else if $value.Range }}
      {{ $name }}: generateFromRange({{ $value.Range.Min }}, {{ $value.Range.Max }}),
    {{- else if $value.Pattern }}
      {{ $name }}: generateFromPattern('{{ $value.Pattern }}'),
    {{- end }}
  {{- end }}
{{- end }}

export default function () {
  {{- range .Paths }}
  group("{{ .Pathname }}", function () {
    const pathParams = {
      {{- template "params" .Detail.Params.Path }}
    }

    const queryParams = {
      {{- template "params" .Detail.Params.Query }}
    }

    const url = new URL('{{ $.BaseUrl }}' + applyPathParams('{{ .Pathname }}', pathParams));

    for (const param of Object.keys(queryParams)) {
        url.searchParams.append(param, queryParams[param]);
    }

    http.
    {{- if eq .Method "post" -}}
    post
    {{- else if eq .Method "put" -}}
    put
    {{- else if eq .Method "patch" -}}
    patch
    {{- else if eq .Method "delete" -}}
    delete
    {{- else -}}
    get
    {{- end -}}
    (url.href);
  });
  {{ end }}
}`
