package mgorm

// NewUpdater func
func NewUpdater(updater interface{}) interface{} {
	return map[string]interface{}{
		"$set": updater,
	}
}
