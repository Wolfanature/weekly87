package players

import (
	"image"
	"image/color"
	"path/filepath"
	"strings"
	"io/ioutil"

	"github.com/solovev/gopsd"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/fileutil"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/weekly87/internal/abilities"
	"github.com/oakmound/weekly87/internal/recolor"
)

var mageConstructors = map[string]*Constructor{}

func mageInit() {
	var mageDefinitions = []ClassDefinition{
		{
			Name: "Blue",
			LayerColors: map[string]color.RGBA{
				"clothes": color.RGBA{47, 47, 200, 200},
			},
			Special1: abilities.FrostBolt,
			Special2: abilities.Blizzard,
		},
		{
			Name: "White",
			LayerColors: map[string]color.RGBA{
				"clothes": color.RGBA{190, 190, 190, 190},
			},
			Special1: abilities.Invulnerability,
			Special2: abilities.Rez,
		},
		{
			Name: "Red",
			LayerColors: map[string]color.RGBA{
				"clothes": color.RGBA{180, 70, 70, 180},
			},
			Special1: abilities.Fireball,
			Special2: abilities.FireStorm,
		},
		{
			Name: "Time",
			// LayerColors: map[string]color.RGBA{
			// 	"clothes": color.RGBA{100, 240, 100, 150},
			// },
			Special1: abilities.Slow,
			Special2: abilities.CooldownRework,
		},
	}

	for _, def := range mageDefinitions {

		psdFilePath := filepath.Join("assets", "images", "16x32", "mage.psd")
		rd, err := fileutil.Open(psdFilePath)
		dlog.ErrorCheck(err)
		data, err := ioutil.ReadAll(rd)
		dlog.ErrorCheck(err)
		psd, err := gopsd.ParseFromBuffer(data)

		combined := render.NewCompositeM()

		for _, layer := range psd.Layers {
			img, err := layer.GetImage()
			dlog.ErrorCheck(err)
			rgba, ok := img.(*image.RGBA)
			if !ok {
				dlog.Error("Image was not RGBA in underlying type")
			}
			sp := render.NewSprite(float64(layer.Rectangle.X), float64(layer.Rectangle.Y), rgba)
			if c, ok := def.LayerColors[strings.ToLower(layer.Name)]; ok {
				//fmt.Println("We found the right layer", layer.Name)
				// Recolor this layer
				sp.Filter(recolor.WithStrategy(recolor.ColorMix(c)))
			}
			// Add this layer to the combined image
			// Todo: bug with shoulder having some pixel flashing
			combined.Append(sp)
		}

		// flatten composite
		combinedSp := combined.ToSprite()

		sh, err := render.MakeSheet(combinedSp.GetRGBA(), 16, 32, 0)
		dlog.ErrorCheck(err)
		sheet := sh.ToSprites()

		standRT := sheet[0][0].Copy()
		standLT := sheet[0][0].Copy().Modify(mod.FlipX)
		standHold := sheet[0][1].Copy().Modify(mod.FlipX)

		walkRT, err := render.NewSheetSequence(sh, 8, []int{1, 0, 2, 0, 0, 0}...)
		dlog.ErrorCheck(err)
		walkLT := walkRT.Copy().Modify(mod.FlipX)

		walkHold, err := render.NewSheetSequence(sh, 8, []int{1, 1, 2, 1, 0, 1}...)
		dlog.ErrorCheck(err)
		walkHold = walkHold.Copy().Modify(mod.FlipX).(*render.Sequence)

		ghostFilePath := filepath.Join("16x32", "mageghost.png")
		dlog.ErrorCheck(err)

		deadRT, err := render.LoadSheetSequence(ghostFilePath, 16, 32, 0, 8, []int{0, 0, 1, 0}...)
		dlog.ErrorCheck(err)
		deadLT := deadRT.Copy().Modify(mod.FlipX)

		consume, err := render.NewSheetSequence(sh, 8, []int{0, 0, 0, 1, 0, 1}...)
		dlog.ErrorCheck(err)
		consume = consume.Copy().Modify(mod.FlipX).(*render.Sequence)

		mageCharMap := map[string]render.Modifiable{
			"walkRT":    walkRT,
			"walkLT":    walkLT,
			"standRT":   standRT,
			"standLT":   standLT,
			"deadRT":    deadRT,
			"deadLT":    deadLT,
			"walkHold":  walkHold,
			"standHold": standHold,
			"consume":   consume,
		}

		mageConstructors[def.Name] = &Constructor{
			AnimationMap: mageCharMap,
			Dimensions:   floatgeom.Point2{16, 32},
			Speed:        floatgeom.Point2{0, 5},
			RunSpeed:     4.0,
			Special1:     def.Special1,
			Special2:     def.Special2,
		}
	}
}
