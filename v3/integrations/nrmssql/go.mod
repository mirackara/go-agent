module github.com/newrelic/go-agent/v3/integrations/nrmssql

go 1.20

require (
	github.com/microsoft/go-mssqldb v0.19.0
	github.com/newrelic/go-agent/v3 v3.33.1
)


replace github.com/newrelic/go-agent/v3 => ../..
