module github.com/xoctopus/agents

go 1.26.5

tool github.com/xoctopus/agents/internal/cmd/skill-install

require (
	// +skill:concx
	github.com/xoctopus/concx v0.2.0
	// +skill:appx
	// +skill:kg
	github.com/xoctopus/confx v0.5.8
	// +skill:genx
	github.com/xoctopus/genx v0.3.7
	// +skill:logx
	github.com/xoctopus/logx v0.3.7
	// +skill:sqlx
	github.com/xoctopus/sqlx v0.4.2
	// +skill:testx
	github.com/xoctopus/x v0.5.7
)

require (
	github.com/go-think/openssl v1.23.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/xoctopus/typx v0.4.6 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	golang.org/x/mod v0.40.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/tools v0.49.0 // indirect
)
