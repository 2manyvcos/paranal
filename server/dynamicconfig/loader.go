package dynamicconfig

import (
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"path/filepath"
	"sync"
	"time"

	"github.com/2manyvcos/paranal/utils"
	"github.com/bep/debounce"
)

type collectorWithRev struct {
	CollectorConfig
	Rev string
}

type Loader struct {
	path             string
	debounceDuration time.Duration
	watcher          *utils.FileWatcher
	watchOnce        sync.Once
	lock             sync.Mutex
	collectorRevs    map[string]collectorWithRev
	changes          chan<- []ConfigChange
	done             chan struct{}
	doneOnce         sync.Once
}

func (l *Loader) Watch() {
	l.watchOnce.Do(func() {
		go l.watch()
	})
}

func (l *Loader) Shutdown() {
	l.doneOnce.Do(func() {
		close(l.done)
	})

	l.watcher.Shutdown()
}

func (l *Loader) loadConfigs() ([]Config, error) {
	var configs []Config
	err := filepath.WalkDir(l.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		bundle, err := loadConfigFile(path)
		if err != nil {
			return fmt.Errorf("loading %q: %w", path, err)
		}

		configs = append(configs, bundle...)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func (l *Loader) applyConfigs(configs []Config) ([]ConfigChange, error) {
	collectorRevs := make(map[string]collectorWithRev, len(configs))

	for _, config := range configs {
		switch config.Kind() {
		case "collector":
			collectorConfig := config.(CollectorConfig)
			id := url.QueryEscape(collectorConfig.ID) + "@file"
			if _, ok := collectorRevs[id]; ok {
				return nil, fmt.Errorf("duplicate collector %q", id)
			}
			rev, err := utils.HashJSONValue(collectorConfig)
			if err != nil {
				return nil, fmt.Errorf("deriving collector revision: %w", err)
			}
			collectorRevs[id] = collectorWithRev{collectorConfig, rev}

		default:
			return nil, fmt.Errorf("unsupported kind %q", config.Kind())
		}
	}

	changes := make([]ConfigChange, 0, len(collectorRevs)+len(l.collectorRevs))
	for id, c := range collectorRevs {
		if prev, ok := l.collectorRevs[id]; !ok || prev.Rev != c.Rev {
			changes = append(changes, ConfigChange{Config: c.CollectorConfig})
		}
	}
	for id, c := range l.collectorRevs {
		if _, ok := collectorRevs[id]; !ok {
			changes = append(changes, ConfigChange{Config: c.CollectorConfig, Removed: true})
		}
	}

	l.collectorRevs = collectorRevs

	return changes, nil
}

func (l *Loader) loadChanges() error {
	configs, err := l.loadConfigs()
	if err != nil {
		return err
	}

	changes, err := l.applyConfigs(configs)
	if err != nil {
		return fmt.Errorf("applying config: %w", err)
	}

	if len(changes) > 0 {
		l.changes <- changes
	}

	return nil
}

func (l *Loader) watch() {
	debounced := debounce.New(l.debounceDuration)

	for {
		select {
		case <-l.done:
			return

		case _, ok := <-l.watcher.Events():
			if !ok {
				return
			}

			debounced(func() {
				select {
				case <-l.done:
					return
				default:
				}

				l.lock.Lock()
				defer l.lock.Unlock()

				err := l.loadChanges()
				if err != nil {
					log.Printf("config watcher: %v\n", err)
				}
			})

		case err, ok := <-l.watcher.Errors():
			if !ok {
				return
			}

			log.Printf("config watcher: %v\n", err)
		}
	}
}

func NewLoader(path string, debounceDuration time.Duration, changes chan<- []ConfigChange) (*Loader, []Config, error) {
	var l Loader
	l.path = path
	l.debounceDuration = debounceDuration
	l.collectorRevs = make(map[string]collectorWithRev)
	l.changes = changes
	l.done = make(chan struct{})

	configs, err := l.loadConfigs()
	if err != nil {
		return nil, nil, err
	}
	_, err = l.applyConfigs(configs)
	if err != nil {
		return nil, nil, fmt.Errorf("applying config: %w", err)
	}

	l.watcher, err = utils.NewFileWatcher(path)
	if err != nil {
		return nil, nil, err
	}

	return &l, configs, nil
}
