package rpcregistry

type Config = map[int64]ChainConfig

type ChainConfig struct {
	HTTP string
}
