package inn

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/weekly87/internal/characters/doodads"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/characters/players"
	"github.com/oakmound/weekly87/internal/joys"
)

// NewInnWalker creates a special character for the inn
func NewInnWalker(innSpace floatgeom.Rect2, scale float64, anim *render.Switch) *entities.Interactive {

	s := entities.NewInteractive(
		float64(oak.ScreenWidth/2),
		float64(oak.ScreenHeight/2)+40,
		16*scale,
		32*scale,
		anim.Copy().Modify(mod.Scale(scale, scale)).(*render.Switch),
		nil,
		0,
		0,
	)

	lowestID := joys.LowestID()

	s.Bind(func(id int, _ interface{}) int {
		p, ok := event.GetEntity(id).(*entities.Interactive)
		if !ok {
			dlog.Error("Non-player sent to player binding")
		}

		p.Delta.Zero()

		js := joys.StickState(lowestID)
		// Todo: support full analog control

		if oak.IsDown(key.UpArrow) || js.StickLY > 8000 {
			p.Delta.Add(physics.NewVector(0, -p.Speed.Y()))
		}
		if oak.IsDown(key.DownArrow) || js.StickLY < -8000 {
			p.Delta.Add(physics.NewVector(0, p.Speed.Y()))
		}
		if oak.IsDown(key.LeftArrow) || js.StickLX < -8000 {
			p.Delta.Add(physics.NewVector(-p.Speed.X(), 0))
		}
		if oak.IsDown(key.RightArrow) || js.StickLX > 8000 {
			p.Delta.Add(physics.NewVector(p.Speed.X(), 0))
		}

		p.Vector.Add(p.Delta)
		_, h := p.R.GetDims()
		hf := float64(h)
		// if p.Vector.Y() < float64(oak.ScreenHeight)*1/3 {
		// 	p.Vector.SetY(float64(oak.ScreenHeight) * 1 / 3)
		// }
		if p.Vector.Y() < 32 {
			p.Vector.SetY(32)
		} else if p.Vector.Y() > (float64(oak.ScreenHeight) - hf) {
			p.Vector.SetY((float64(oak.ScreenHeight) - hf))
		}
		if p.Vector.X() < 220 {
			p.Vector.SetX(220)
		} else if p.Vector.X()+p.W > float64(oak.ScreenWidth) {
			p.Vector.SetX(float64(oak.ScreenWidth) - p.W)
		}
		p.R.SetPos(p.Vector.X(), p.Vector.Y())
		p.RSpace.Update(p.Vector.X(), p.Vector.Y(), p.RSpace.GetW(), p.RSpace.GetH())
		<-s.RSpace.CallOnHits()
		if collision.HitLabel(s.RSpace.Space, labels.Blocking, labels.NPC) != nil {
			p.Vector.Sub(p.Delta)
			p.R.SetPos(p.Vector.X(), p.Vector.Y())
			p.RSpace.Update(p.Vector.X(), p.Vector.Y(), p.RSpace.GetW(), p.RSpace.GetH())
		}

		swch := p.R.(*render.Switch)
		if p.Delta.X() != 0 || p.Delta.Y() != 0 {
			if p.Delta.X() > 0 {
				swch.Set("walkRT")
			} else {
				swch.Set("walkLT")
			}
		} else {
			cur := swch.Get()
			err := swch.Set("stand" + string(cur[len(cur)-2:]))
			dlog.ErrorCheck(err)
		}
		return 0
	}, "EnterFrame")
	s.Speed = physics.NewVector(5, 5) // We actually allow players to move around in the inn!

	s.RSpace.Add(collision.Label(labels.Door), (func(s1, s2 *collision.Space) {
		d, ok := s2.CID.E().(*doodads.InnDoor)
		if !ok {
			dlog.Error("Non-door sent to inndoor binding")
			return
		}
		nextscene = d.NextScene
		stayInMenu = false
	}))

	render.Draw(s.R, 2, 1)
	return s
}

// NPC is a inn only construct
type NPC struct {
	*entities.Interactive
	Swtch *render.Switch
	Class int
}

// Init the npc so it has a CID!
func (n *NPC) Init() event.CID {
	return event.NextID(n)
}

// FaceLeft take wether the npc should face left
// Set the facing renderable appropriately
func (n NPC) FaceLeft(shouldFaceLeft bool) NPC {
	if shouldFaceLeft {
		n.R.(*render.Switch).Set("standLT")
	} else {
		n.R.(*render.Switch).Set("standRT")
	}
	return n
}

// NewInnNPC creates a npc to interact with for setting up party
func NewInnNPC(class int, scale, x, y float64) NPC {
	pcon := players.ClassConstructor([]int{class})[0]
	n := NPC{}
	n.Class = class
	n.Swtch = render.NewSwitch("standRT", pcon.AnimationMap).Copy().(*render.Switch)
	n.Swtch.Modify(mod.Scale(scale, scale))
	n.Interactive = entities.NewInteractive(
		x,
		y,
		pcon.Dimensions.X()*scale,
		pcon.Dimensions.Y()*scale,
		n.Swtch,
		nil,
		n.Init(),
		0,
	)

	return n
}

// Activate draws the npc and makes it collidable
func (n NPC) Activate() {
	n.RSpace.UpdateLabel(labels.NPC)
	render.Draw(n.R, 2, 1)
}

func (n NPC) Destroy() {
	n.Tree.Remove(n.RSpace.Space)

}
