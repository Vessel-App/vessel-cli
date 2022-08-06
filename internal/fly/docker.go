package fly

// In previous iteration, we manually exec'ed
// `fly auth docker`, then `docker tag some-local-image registry.fly.io/$app:latest`,
// and finally `docker push registry.fly.io/$app:latest`

type PushToRegistryRequest struct {
}
