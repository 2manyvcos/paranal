package dynamicconfig

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/2manyvcos/paranal/utils"
	"github.com/bep/debounce"
)

type configRevKey struct {
	Type string
	ID   string
}

type configRev struct {
	Config
	Rev string
}

type Loader struct {
	path             string
	debounceDuration time.Duration
	watcher          *utils.FileWatcher
	watchOnce        sync.Once
	lock             sync.Mutex
	revs             map[configRevKey]configRev
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
	revs := make(map[configRevKey]configRev, len(configs))

	for _, config := range configs {
		key := configRevKey{Type: config.Type(), ID: config.ID()}
		if _, ok := revs[key]; ok {
			return nil, fmt.Errorf("duplicate %q config %q", config.Type(), config.ID())
		}
		rev, err := utils.HashJSONValue(config)
		if err != nil {
			return nil, fmt.Errorf("deriving config revision: %w", err)
		}
		revs[key] = configRev{config, rev}
	}

	changes := make([]ConfigChange, 0, len(revs)+len(l.revs))
	for key, rev := range revs {
		if prev, ok := l.revs[key]; !ok || prev.Rev != rev.Rev {
			changes = append(changes, ConfigChange{Config: rev.Config})
		}
	}
	for key, rev := range l.revs {
		if _, ok := revs[key]; !ok {
			changes = append(changes, ConfigChange{Config: rev.Config, Removed: true})
		}
	}

	l.revs = revs

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
					log.Println(err)
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
