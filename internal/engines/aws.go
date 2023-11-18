package engines

type AmazonS3 struct {}

func (engine *AmazonS3) SupportsTracking() bool {
	return false
}

func (engine *AmazonS3) SupportsOpaque() bool {
	return true
}