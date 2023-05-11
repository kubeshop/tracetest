postgres:
  host: {{ .pHost }}
  user: {{ .pUser }}
  password: {{ .pPasswd }}
  port: 5432
  dbname: postgres
  params: sslmode=disable
