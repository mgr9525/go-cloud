select
{{ if .getCount}}
count(1)
{{ else }}
*
{{ end }}
from user where
1=1
{{ if .name}}
and name=?name
{{ end }}