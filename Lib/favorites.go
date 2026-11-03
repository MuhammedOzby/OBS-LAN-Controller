package lib

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type favoriteConfig struct {
	Collection string    `json:"collection"`
	Slots      [3]string `json:"slots"`
}

type favoriteStore struct {
	mu    sync.RWMutex
	path  string
	value *favoriteConfig
}

func loadFavorites(path string) (*favoriteStore, error) {
	s := &favoriteStore{path: path}
	if path == "" {
		return s, nil
	}
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	var v favoriteConfig
	if json.Unmarshal(b, &v) != nil || v.Collection == "" {
		return nil, errors.New("favori ayarları geçersiz")
	}
	s.value = &v
	return s, nil
}

func (s *favoriteStore) get() *favoriteConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.value == nil {
		return nil
	}
	v := *s.value
	return &v
}

func (s *favoriteStore) save(v favoriteConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.path != "" {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		f, err := os.CreateTemp(filepath.Dir(s.path), ".favorites-*")
		if err != nil {
			return err
		}
		defer os.Remove(f.Name())
		if err = f.Chmod(0600); err == nil {
			_, err = f.Write(b)
		}
		if err == nil {
			err = f.Sync()
		}
		if e := f.Close(); err == nil {
			err = e
		}
		if err == nil {
			err = os.Rename(f.Name(), s.path)
		}
		if err != nil {
			return err
		}
	}
	s.value = &v
	return nil
}

func (o *OBSClient) collectionLocked(ctx context.Context) (string, error) {
	var result struct {
		Name string `json:"currentSceneCollectionName"`
	}
	if err := o.request(ctx, "GetSceneCollectionList", map[string]any{}, &result); err != nil {
		return "", err
	}
	if result.Name == "" {
		return "", errors.New("sahne koleksiyonu bilinmiyor")
	}
	return result.Name, nil
}

func favoriteSlots(config *favoriteConfig, collection string, scenes []string) [3]string {
	var slots [3]string
	if config == nil {
		copy(slots[:], scenes)
		return slots
	}
	if config.Collection == collection {
		return config.Slots
	}
	return slots
}

func containsScene(scenes []string, name string) bool {
	if name == "" {
		return false
	}
	for _, s := range scenes {
		if s == name {
			return true
		}
	}
	return false
}

func favoriteRoutes(router *gin.Engine, obs *OBSClient, store *favoriteStore) {
	router.GET("/api/favorites", func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		if obs == nil {
			c.Status(503)
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		obs.mu.Lock()
		defer obs.mu.Unlock()
		status := obs.snapshotLocked(ctx, false, true)
		if !status.Connected || obs.conn == nil {
			c.Status(503)
			return
		}
		collection, err := obs.collectionLocked(ctx)
		if err != nil {
			c.Status(503)
			return
		}
		c.JSON(200, gin.H{"collection": collection, "scenes": status.Scenes, "slots": favoriteSlots(store.get(), collection, status.Scenes)})
	})
	router.POST("/api/favorites", func(c *gin.Context) {
		if !pairingSameOrigin(c.Request) {
			c.Status(403)
			return
		}
		media, _, err := mime.ParseMediaType(c.GetHeader("Content-Type"))
		if err != nil || media != "application/json" {
			c.Status(415)
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 4096)
		var request struct {
			Collection string   `json:"collection"`
			Slots      []string `json:"slots"`
		}
		dec := json.NewDecoder(c.Request.Body)
		dec.DisallowUnknownFields()
		if dec.Decode(&request) != nil || dec.Decode(&struct{}{}) != io.EOF || len(request.Slots) != 3 || request.Collection == "" {
			c.Status(400)
			return
		}
		if obs == nil {
			c.Status(503)
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		obs.mu.Lock()
		defer obs.mu.Unlock()
		status := obs.snapshotLocked(ctx, false, true)
		if !status.Connected || obs.conn == nil {
			c.Status(503)
			return
		}
		collection, err := obs.collectionLocked(ctx)
		if err != nil {
			c.Status(503)
			return
		}
		if collection != request.Collection {
			c.Status(409)
			return
		}
		seen := map[string]bool{}
		old := store.get()
		for i, name := range request.Slots {
			retained := old != nil && old.Collection == collection && old.Slots[i] == name
			if name != "" && ((!containsScene(status.Scenes, name) && !retained) || seen[name]) {
				c.Status(400)
				return
			}
			seen[name] = true
		}
		v := favoriteConfig{Collection: collection}
		copy(v.Slots[:], request.Slots)
		if err := store.save(v); err != nil {
			c.Status(500)
			return
		}
		c.Status(204)
	})
}
