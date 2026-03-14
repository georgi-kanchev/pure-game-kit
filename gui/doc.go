// The most complex package of all - handling graphical user interfaces by depending heavily on multiple
// other packages used for file loading, drawing graphics, accepting input, playing sounds, etc.
//
// The GUI topic is long & thorough and there are many designs, but the few main ones for
// games seem to be:
//   - Object-oriented (OOP) - Offers the most freedom but way too verbose.
//   - Immediate mode (Im) - The simplest one but lacks customization & serialization for re-usability.
//   - Data-oriented (css) - Reusable and customizable but hard to parse, create with code, and handle custom logic.
//
// This package takes benefits from each one of them and tries to solve their problems.
// It relies on a simple design idea but its usage remains a bit complex due to its sheer depth.
// It is constructed of 3 types of elements: Containers, Widgets & Themes.
//
// The GUI creation relies on the Data-oriented approach by parsing XML. Its problems and how they are solved:
//   - Hard to parse - by not allowing nesting, having a max depth of 2.
//   - Hard to create with code + no autocomplete - by optionally chaining function calls to construct the XML.
//   - Hard to handle custom logic - by mixing in the Immediate mode approach.
//
// The Immediate mode approach brings its problems as well, here is how they are solved:
//   - Lacking customization - by separating creation details and functionality
//   - Lacking serialization for re-usability - by the Data-oriented XML approach
//   - Relying on code structures as existing elements - by a single GUI structure & accessing everything through ids.
//
// Another huge problem is GUI elements respecting any window aspect ratios. This is solved by replacing
// certain dynamic variable keywords during the XML parsing while handling math expresions.
// Due to the nature of those dynamic values, scaling the GUI comes for free by zooming its provided camera.
//
// While loading an XML is handled, saving is not and this is a deliberate choice. Saving a GUI state
// has a risk of doing damage to the initial state and has to deal with versioning or multiple GUI states.
// Another reason not to do it is that fundamentally it does not make sense to save a GUI.
// It's rather better to save its data instead, then load the GUI in its initial state each time and
// have it react to the separately loaded data.
//
// Alongside solving all of those problems, here are some of the very useful features in no particular order:
//
//   - Widgets inheriting/reusing properties from their Themes or Container owners and optionally overwriting them.
//   - Elements supporting custom properties that only custom logic may rely on.
//   - Dividing long & complex GUI systems into multiple XMLs by optionally merging them during parsing.
//   - Containers handling scrolling, masking and ordering widgets out-of-the-box.
//   - Easy to reference themes, widgets, containers and assets due to the nature of ids.
//   - Rendering fallbacks to basic colored shapes in case no assets are provided.
//   - Out-of-the-box Z ordering for input & drawing.
//   - Having tooltips for all widgets, including text labels & images.
package gui
