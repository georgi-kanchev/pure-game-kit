package tiled

import (
	"pure-game-kit/internal"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/flag"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/text"
)

type Object struct {
	Project    *Project
	Properties map[string]any
	Points     [][2]float32
}

func newObject(data *internal.LayerObject, project *Project) *Object {
	var result = Object{Project: project}
	result.initProperties(data)
	result.initPoints(data)
	return &result
}

//=================================================================

func (t *Object) initProperties(data *internal.LayerObject) {
	t.Properties = make(map[string]any)
	t.Properties[property.ObjectId] = data.Id
	t.Properties[property.ObjectClass] = data.Class
	t.Properties[property.ObjectTemplate] = data.Template
	t.Properties[property.ObjectName] = data.Name
	t.Properties[property.ObjectVisible] = data.Visible != "false"
	t.Properties[property.ObjectLocked] = data.Locked
	t.Properties[property.ObjectX] = data.X
	t.Properties[property.ObjectY] = data.Y
	t.Properties[property.ObjectWidth] = data.Width
	t.Properties[property.ObjectHeight] = data.Height
	t.Properties[property.ObjectRotation] = data.Rotation
	t.Properties[property.ObjectTileId] = data.Gid
	t.Properties[property.ObjectFlipX] = flag.IsOn(data.Gid, internal.FlipX)
	t.Properties[property.ObjectFlipY] = flag.IsOn(data.Gid, internal.FlipY)

	for _, prop := range data.Properties {
		t.Properties[prop.Name] = parseProperty(&prop, t.Project)
	}
}
func (t *Object) initPoints(data *internal.LayerObject) {
	var ptsData = ""
	if data.Polyline.Points != "" {
		ptsData = data.Polyline.Points
	}
	if data.Polygon.Points != "" {
		ptsData = data.Polygon.Points
	}
	if ptsData == "" {
		var w, h = data.Width, data.Height
		if data.Ellipse != nil {
			const segments = 16
			var rx, ry = w / 2, h / 2
			var step = 360.0 / float32(segments)
			var firstValue = ""

			for i := range segments {
				var cx, cy = point.MoveAtAngle(0, 0, float32(i)*step, 1)
				var value = text.New(cx*rx, ",", cy*ry, " ")
				ptsData += value

				if i == 0 {
					firstValue = value
				}
			}
			ptsData += firstValue
		} else { // rectangle
			ptsData = text.New(0, ",", 0, " ", w, ",", 0, " ", w, ",", h, " ", 0, ",", h)
		}
	}

	var points = [][2]float32{}
	var pts = text.Split(text.Trim(ptsData), " ")
	for _, pt := range pts {
		var xy = text.Split(pt, ",")
		if len(xy) == 2 {
			var x, y = text.ToNumber[float32](xy[0]), text.ToNumber[float32](xy[1])
			x, y = point.RotateAroundPoint(x, y, 0, 0, data.Rotation)
			points = append(points, [2]float32{x, y})
		}
	}

	t.Points = points
}
