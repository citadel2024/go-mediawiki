package mediawiki

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"sync"
	"time"
)

const (
	checkpointFile = "checkpoint.json"
	itemsThreshold = 1000
	saveInterval   = 5 * time.Minute
)

// CheckpointConfig represents the configuration for the CheckpointManager
// SaveInterval is the interval to save the checkpoint to the checkpoint file.
// Due to saving interval, the program may lose some items if it crashes before saving the checkpoint.
type CheckpointConfig struct {
	SaveInterval   time.Duration
	ItemsThreshold int
	CheckpointFile string
}

// Checkpoint represents the checkpoint data
// TotalItems is the total number of items we send into channel, we will increase this number based on the previous checkpoint.
// ProcessedPosition is the current position of the item we are processing, however, for multiple-consumer program, ProcessedPosition does not
// mean the items before ProcessedPosition were processed, because the items in the channel are not processed in order.
// Some items may be processed, the others may not.
// Hence, next time we start the program, to make sure the items in the channel were processed correctly,
// we need to retry these items in the channel, so we need to start from the beginning of the channel (ProcessedPosition - len(goroutines)),
// whatever the position is the first or the last item in the channel, we can make sure all items in the channel were processed.
// TotalItems >= ProcessedPosition - len(goroutines)
// The producer is single-threaded, so the rows channel is FIFO.
// Bitset is also a good choice to track the processed items, but it is not necessary. We only need a checkpoint.
type Checkpoint struct {
	TotalItems        int       `json:"total_items"`
	SaveTimestamp     time.Time `json:"timestamp"`
	LastItemID        string    `json:"last_item_id"` // This field is unused, user can just use TotalItems to skip items already processed
	ProcessedPosition int       `json:"position"`
	//LastProcessedThreadCount int `json:"last_processed_thread_count"` we can store the number of goroutines in the last run
}

type CheckpointManager struct {
	config                   *CheckpointConfig
	currentCheckpoint        *Checkpoint
	itemsSinceLastCheckpoint int
	mu                       sync.Mutex
	dirty                    bool
}

// NewCheckpointManager creates a new CheckpointManager
// For better performance, we save the checkpoint to the checkpoint file every saveInterval or when the number of
// items processed since the last checkpoint exceeds the threshold.
// It means that we don't save the checkpoint every time we process an item, so the program which uses this package need ability to
// 1. Recover from the last checkpoint, and skip items already processed.
// 2. The program should be able to handle the case when the checkpoint file is missing.
// 3. You may need to handle duplicate items if the program crashes after processing an item but before saving the checkpoint.
func NewCheckpointManager() *CheckpointManager {
	cm := &CheckpointManager{
		config: &CheckpointConfig{
			SaveInterval:   saveInterval,
			ItemsThreshold: itemsThreshold,
			CheckpointFile: checkpointFile,
		},
	}
	if err := cm.loadCheckpoint(); err != nil {
		if os.IsNotExist(err) {
			cm.currentCheckpoint = &Checkpoint{}
		} else {
			panic(fmt.Sprintf("Failed to load checkpoint: %v", err))
		}
	}
	go cm.autoSave()
	return cm
}

func NewCheckpointManagerWithConfig(config *CheckpointConfig) *CheckpointManager {
	cm := &CheckpointManager{
		config: config,
	}
	if err := cm.loadCheckpoint(); err != nil {
		if os.IsNotExist(err) {
			cm.currentCheckpoint = &Checkpoint{}
		} else {
			panic(fmt.Sprintf("Failed to load checkpoint: %v", err))
		}
	}
	go cm.autoSave()
	return cm
}

// autoSave saves the checkpoint to the checkpoint file every saveInterval
func (cm *CheckpointManager) autoSave() {
	ticker := time.NewTicker(cm.config.SaveInterval)
	defer ticker.Stop()

	for range ticker.C {
		if cm.dirty {
			if err := cm.Save(); err != nil {
				fmt.Println("Failed to auto save checkpoint:", err)
			}
		}
	}
}

// UpdateProgressAndMaybeSave updates the checkpoint with the current position and itemID
// We need to invoke this method every time we process an item
func (cm *CheckpointManager) UpdateProgressAndMaybeSave(position int, itemID string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.currentCheckpoint.ProcessedPosition = position
	cm.currentCheckpoint.LastItemID = itemID
	cm.currentCheckpoint.TotalItems++
	cm.itemsSinceLastCheckpoint++
	cm.dirty = true
	if cm.itemsSinceLastCheckpoint >= cm.config.ItemsThreshold {
		return cm.save()
	}
	return nil
}

// Save saves the current checkpoint to the checkpoint file
func (cm *CheckpointManager) Save() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.save()
}

// save saves the current checkpoint to the checkpoint file
// There are two conditions to save the checkpoint:
// 1. The number of items processed since the last checkpoint exceeds the threshold
// 2. The time since the last save exceeds the save interval
// So we use mutex to protect this method.
func (cm *CheckpointManager) save() error {
	if !cm.dirty {
		return nil
	}
	cm.currentCheckpoint.SaveTimestamp = time.Now()
	data, err := json.MarshalIndent(cm.currentCheckpoint, "", "  ")
	if err != nil {
		return errors.WithMessage(err, "marshal checkpoint")
	}
	tempFile := cm.config.CheckpointFile + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return errors.WithMessage(err, "write checkpoint")
	}
	if err := os.Rename(tempFile, cm.config.CheckpointFile); err != nil {
		return errors.WithMessage(err, "rename checkpoint file")
	}
	cm.itemsSinceLastCheckpoint = 0
	cm.dirty = false
	return nil
}

// loadCheckpoint loads the checkpoint from the checkpoint file
func (cm *CheckpointManager) loadCheckpoint() error {
	data, err := os.ReadFile(cm.config.CheckpointFile)
	if err != nil {
		return err
	}
	checkpoint := &Checkpoint{}
	if err := json.Unmarshal(data, checkpoint); err != nil {
		return errors.WithMessage(err, "unmarshal checkpoint")
	}
	fmt.Println("Loaded checkpoint", checkpoint)
	cm.currentCheckpoint = checkpoint
	return nil
}

func (cm *CheckpointManager) GetCheckpoint() *Checkpoint {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.currentCheckpoint
}

func (cm *CheckpointManager) Close() error {
	return cm.Save()
}
