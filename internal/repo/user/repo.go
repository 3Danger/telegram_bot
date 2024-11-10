package user

//go:generate gowrap gen -p ./postgres -i Querier -o ./wrappers/timeout.go -t timeout
//go:generate gowrap gen -p ./postgres -i Querier -o ./wrappers/skip_no_rows.go -t ../../../gowrap/skip_no_rows.tmpl
