/*
Contains constants, representing dynamic variables of:
  - Camera passed to the GUI when drawing - useful for anchoring elements or resizing them in percents.
  - Container owner of a specific widget - useful for using it as a pivot or reference.
  - Optional container target of a specific container - useful for using it as a pivot or reference.

Can be included in some XML properties alongside math expressions like so:

	// in XML:
	containerX="TargetX + TargetWidth + 50" containerWidth="CameraWidth / 4"
	// or in code:
	gui.Container("name", dynamic.CameraLeftX+"+10", dynamic.CameraTopY+"+10", dynamic.CameraWidth+"-20", "1100")
*/
package dynamic

const (
	CameraCenterX = "CameraCenterX" // [widget] [container]
	CameraCenterY = "CameraCenterY" // [widget] [container]
	CameraLeftX   = "CameraLeftX"   // [widget] [container]
	CameraRightX  = "CameraRightX"  // [widget] [container]
	CameraTopY    = "CameraTopY"    // [widget] [container]
	CameraBottomY = "CameraBottomY" // [widget] [container]
	CameraWidth   = "CameraWidth"   // [widget] [container]
	CameraHeight  = "CameraHeight"  // [widget] [container]

	OwnerCenterX = "OwnerCenterX" // [widget]
	OwnerCenterY = "OwnerCenterY" // [widget]
	OwnerLeftX   = "OwnerLeftX"   // [widget]
	OwnerRightX  = "OwnerRightX"  // [widget]
	OwnerTopY    = "OwnerTopY"    // [widget]
	OwnerBottomY = "OwnerBottomY" // [widget]
	OwnerWidth   = "OwnerWidth"   // [widget]
	OwnerHeight  = "OwnerHeight"  // [widget]

	TargetCenterX  = "TargetCenterX"  // [container]
	TargetCenterY  = "TargetCenterY"  // [container]
	TargetLeftX    = "TargetLeftX"    // [container]
	TargetRightX   = "TargetRightX"   // [container]
	TargetTopY     = "TargetTopY"     // [container]
	TargetBottomY  = "TargetBottomY"  // [container]
	TargetWidth    = "TargetWidth"    // [container]
	TargetHeight   = "TargetHeight"   // [container]
	TargetHidden   = "TargetHidden"   // [container]
	TargetDisabled = "TargetDisabled" // [container]
)
