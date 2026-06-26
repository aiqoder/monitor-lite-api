package version

// Version 由 CI 通过 -ldflags 注入；本地开发默认为 dev。
var Version = "dev"
