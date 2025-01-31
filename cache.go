package ebui

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// ImageCache provides a global cache for commonly used images
type ImageCache struct {
	mu sync.RWMutex
	// Cache for solid color images, keyed by "WxH-RGBA"
	colorCache map[string]*ebiten.Image
	// Cache for specific sizes, keyed by "WxH"
	sizeCache map[string]*ebiten.Image
}

var (
	globalCache = &ImageCache{
		colorCache: make(map[string]*ebiten.Image),
		sizeCache:  make(map[string]*ebiten.Image),
	}
)

// ImageWithColor returns a cached image of the specified size and color
func (c *ImageCache) ImageWithColor(width, height int, col color.Color) *ebiten.Image {
	r, g, b, a := col.RGBA()
	key := fmt.Sprintf("%dx%d-%d%d%d%d", width, height, r>>8, g>>8, b>>8, a>>8)

	c.mu.RLock()
	if img, ok := c.colorCache[key]; ok {
		c.mu.RUnlock()
		return img
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if img, ok := c.colorCache[key]; ok {
		return img
	}

	img := ebiten.NewImage(width, height)
	img.Fill(col)
	c.colorCache[key] = img
	return img
}

// Clear empties the cache
func (c *ImageCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Clear all references to allow garbage collection
	c.colorCache = make(map[string]*ebiten.Image)
	c.sizeCache = make(map[string]*ebiten.Image)
}

// GetCache returns the global image cache instance
func GetCache() *ImageCache {
	return globalCache
}
