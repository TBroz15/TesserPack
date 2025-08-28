package instancechecker

import (
	"os"
	"path"
	"tesserpack/internal/helpers"

	"github.com/charmbracelet/log"

	"github.com/gofrs/flock"
)

// TesserPack can be resource-intensive, so warn the user if there is any instances

type InstanceChecker struct {
	fileLock *flock.Flock
}

func New() (InstanceChecker) {
	fileLockPath := path.Join(helpers.TempDir, "tsp.lock")
	fileLock := flock.New(fileLockPath)

	return InstanceChecker{fileLock}
}

func (checker *InstanceChecker) CheckLock()  {
	isLocked, err := checker.fileLock.TryLock()
	if err != nil {
		log.Error("Failed to check lock file", "err", err)
		os.Exit(1)
	}

	if !isLocked {
		log.Warn("There are multiple instances of TesserPack running right now. Running many instances of it can slow down your computer.")
	}
}

func (checker *InstanceChecker) Unlock() {
	err := checker.fileLock.Unlock()
	
	if err != nil {
		log.Error("Failed to unlock lock file", "err", err)
		os.Exit(1)
	}

	os.Remove(checker.fileLock.Path())
}