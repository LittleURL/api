package engines

type CloudflareWorkers struct {}

func (engine *CloudflareWorkers) SupportsTracking() bool {
	return true
}

func (engine *CloudflareWorkers) SupportsOpaque() bool {
	return true
}