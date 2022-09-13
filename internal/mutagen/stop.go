package mutagen

import "fmt"

// StopSession will attempt to stop a syncing and forwarding session.
// Even if one fails, it still attempts to do the other.
func StopSession(name string) error {
	errSync := StopSync(name)
	errForward := StopForward(name)

	if errSync != nil || errForward != nil {
		return fmt.Errorf("stop syncing error: %w, stop forwarding error: %w", errSync, errForward)
	}

	return nil
}
