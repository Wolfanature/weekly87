package doodads

import (
	"path/filepath"

	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"

	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/restrictor"
)

type Chest struct {
	*entities.Reactive
	Unmoving
	Value  int64
	Active bool
}

func (c *Chest) Init() event.CID {
	return event.NextID(c)
}

func (c *Chest) Destroy() {
	c.Active = false
	c.Reactive.Destroy()
}

func (c *Chest) Activate() {
	restrictor.Add(c)
	c.Active = true
}

func (c *Chest) GetDims() (int, int) {
	return c.Reactive.R.GetDims()
}

func NewChest(value int64) *Chest {
	ch := &Chest{}
	// r := render.NewColorBox(16, 16, color.RGBA{0, 255, 255, 255})
	r, _ := render.LoadSprite("", filepath.Join("", "/16x16/chest.png"))
	// Todo: calculate image based on value
	ch.Reactive = entities.NewReactive(0, 0, 16, 16,
		r, nil, ch.Init())

	ch.RSpace.UpdateLabel(labels.Chest)
	ch.Value = value
	return ch
}
