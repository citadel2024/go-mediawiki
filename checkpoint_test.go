package mediawiki

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestNewCheckpointManager(t *testing.T) {
	tmpFile := "TestNewCheckpointManager.json"
	cm := NewCheckpointManagerWithConfig(
		&CheckpointConfig{
			SaveInterval:   time.Second,
			ItemsThreshold: 1,
			CheckpointFile: tmpFile,
		})
	defer os.Remove(tmpFile)

	assert.NotNil(t, cm)
	assert.NotNil(t, cm.currentCheckpoint)
	assert.Equal(t, int64(0), cm.currentCheckpoint.TotalItems)
}

func TestCheckpointManager_AutoSaveTicker(t *testing.T) {
	tmpFile := "TestCheckpointManager_AutoSaveTicker.json"
	cm := NewCheckpointManagerWithConfig(
		&CheckpointConfig{
			SaveInterval:   time.Second,
			ItemsThreshold: 1,
			CheckpointFile: tmpFile,
		})
	err := cm.UpdateProgressAndMaybeSave(1, "item1")
	assert.NoError(t, err)
	time.Sleep(time.Second * 2)

	data, err := os.ReadFile(tmpFile)
	assert.NoError(t, err)

	var checkpoint Checkpoint
	err = json.Unmarshal(data, &checkpoint)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), checkpoint.TotalItems)
	assert.Equal(t, int64(1), checkpoint.Position)
	assert.Equal(t, "item1", checkpoint.LastItemID)
}

func TestCheckpointManager_UpdateProgressAndMaybeSave(t *testing.T) {
	tmpFile := "TestCheckpointManager_UpdateProgressAndMaybeSave.json"
	cm := NewCheckpointManagerWithConfig(
		&CheckpointConfig{
			SaveInterval:   time.Second,
			ItemsThreshold: 2,
			CheckpointFile: tmpFile,
		})
	defer os.Remove(tmpFile)

	err := cm.UpdateProgressAndMaybeSave(1, "item1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), cm.currentCheckpoint.TotalItems)
	assert.Equal(t, int64(1), cm.itemsSinceLastCheckpoint)

	_, err = os.Stat(tmpFile)
	assert.True(t, os.IsNotExist(err))

	err = cm.UpdateProgressAndMaybeSave(2, "item2")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), cm.currentCheckpoint.TotalItems)

	data, err := os.ReadFile(tmpFile)
	assert.NoError(t, err)

	var checkpoint Checkpoint
	err = json.Unmarshal(data, &checkpoint)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), checkpoint.TotalItems)
	assert.Equal(t, "item2", checkpoint.LastItemID)
}

func TestCheckpointManager_Save(t *testing.T) {
	tmpFile := "TestCheckpointManager_Save.json"
	cm := &CheckpointManager{
		config: &CheckpointConfig{
			CheckpointFile: tmpFile,
		},
		currentCheckpoint: &Checkpoint{
			TotalItems: 100,
			LastItemID: "test_item",
			Position:   50,
		},
		dirty: true,
	}
	defer os.Remove(tmpFile)

	err := cm.Save()
	assert.NoError(t, err)
	assert.False(t, cm.dirty)
	assert.Equal(t, int64(0), cm.itemsSinceLastCheckpoint)

	data, err := os.ReadFile(tmpFile)
	assert.NoError(t, err)

	var checkpoint Checkpoint
	err = json.Unmarshal(data, &checkpoint)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), checkpoint.TotalItems)
	assert.Equal(t, "test_item", checkpoint.LastItemID)
	assert.Equal(t, int64(50), checkpoint.Position)
}

func TestCheckpointManager_LoadCheckpoint(t *testing.T) {
	tmpFile := "TestCheckpointManager_LoadCheckpoint.json"
	testCheckpoint := &Checkpoint{
		TotalItems:    200,
		LastItemID:    "last_item",
		Position:      150,
		SaveTimestamp: time.Now(),
	}

	data, err := json.MarshalIndent(testCheckpoint, "", "  ")
	assert.NoError(t, err)

	err = os.WriteFile(tmpFile, data, 0644)
	assert.NoError(t, err)
	defer os.Remove(tmpFile)

	cm := &CheckpointManager{
		config: &CheckpointConfig{
			CheckpointFile: tmpFile,
		},
	}

	err = cm.loadCheckpoint()
	assert.NoError(t, err)
	assert.Equal(t, testCheckpoint.TotalItems, cm.currentCheckpoint.TotalItems)
	assert.Equal(t, testCheckpoint.LastItemID, cm.currentCheckpoint.LastItemID)
	assert.Equal(t, testCheckpoint.Position, cm.currentCheckpoint.Position)
}

func TestCheckpointManager_CloseAndSave(t *testing.T) {
	tmpFile := "TestCheckpointManager_Close.json"
	cm := &CheckpointManager{
		config: &CheckpointConfig{
			CheckpointFile: tmpFile,
		},
		currentCheckpoint: &Checkpoint{
			TotalItems: 300,
			LastItemID: "final_item",
			Position:   250,
		},
		dirty: true,
	}
	defer os.Remove(tmpFile)

	err := cm.Close()
	assert.NoError(t, err)

	data, err := os.ReadFile(tmpFile)
	assert.NoError(t, err)

	var checkpoint Checkpoint
	err = json.Unmarshal(data, &checkpoint)
	assert.NoError(t, err)
	assert.Equal(t, int64(300), checkpoint.TotalItems)
	assert.Equal(t, "final_item", checkpoint.LastItemID)
	assert.Equal(t, int64(250), checkpoint.Position)
}
