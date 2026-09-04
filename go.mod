module github.com/xoctopus/agents

go 1.27.0

tool github.com/xoctopus/agents/internal/cmd/gen

require (
	github.com/spf13/cobra v1.10.2
	// +skill:concx
	github.com/xoctopus/concx v0.2.3
	// +skill:appx
	// +skill:kg
	github.com/xoctopus/confx v0.6.0
	// +skill:genx
	github.com/xoctopus/genx v0.3.9
	// +skill:logx
	github.com/xoctopus/logx v0.3.9
	// +skill:sqlx
	github.com/xoctopus/sqlx v0.4.4
	// +skill:testx
	github.com/xoctopus/x v0.5.9
	golang.org/x/mod v0.40.0
)

require (
	github.com/go-think/openssl v1.25.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/xoctopus/pkgx v0.4.4 // indirect
	github.com/xoctopus/typx v0.4.7 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/tools v0.49.0 // indirect
)
