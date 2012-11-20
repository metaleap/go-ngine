//	The *assets* package provides the (de)serializable, logic-less base data structures for all resource types
//	used in a go:ngine app.
//	
//	________
//	
//	
//	First off, the go:ngine 3 RULES of bugless(er) development:
//	
//	1. you do NOT set exported fields directly, those are provided for simplified read-access. If a field is
//	meant to be modifiable, there will be a *SetFoo()* method or it will be documented explicitly as directly
//	modifiable, but such instances will be rare.
//	
//	2. you do NOT instantiate exported struct types directly, as in "new(ImageDef)" or "&ImageDef {}". Many of
//	those are exported only for documentation, but are to be instantiated only inside the go:ngine package.
//	For types to be instantiated by package-external code (ie. your code), go:ngine packages provide constructor
//	functions -- often in collections, ie. img := ImageDefs.New("id") instead of, say, img := NewImageDef("id").
//	
//	3. Those 2 are *default* assumptions and modes of operation -- there are a few "if you know what you're
//	doing" exemptions and those are (or will be) explicitly documented as such.
//	
//	________
//	
//	
//	The *assets* package provides the (de)serializable, logic-less base data structures for all resource types
//	used in a go:ngine app.
//	These data structures only describe resource definitions and instances, but do not provide
//	any specific logic or algorithms, such as physics or visual rendering.
//	
//	Context: any go:ngine app uses a variety of both simple and complex resource types, including for
//	example geometry definitions, image textures, material and light descriptions and many others.
//	
//	While for all these resource types, the go:ngine *core* package manages their respective run-time
//	representation in its entirety,	some basic part of that representation is "merely descriptive" and
//	should be readable from or writable to a binary stream ("design-time"). That "merely descriptive"
//	sub-set of resource definitions is fully contained in this go:ngine *assets* package and thus decoupled
//	from the go:ngine *core* package.
//	
//	This *assets* package thus allows for things such as server-side procedural asset generators, networked
//	resource streaming or simple custom asset-import/export/conversion tools, all of which shouldn't have to
//	needlessly depend on the graphics, windowing etc. stacks.
//	
//	NOTE: there are essentially TWO distinct "modes" or use-cases in which the *assets* package is active:
//	
//	1. in graphics-independent, server-side or non-interactive "toolage", dealing only with the raw resource
//	data where it can be generated, loaded, stored, manipulated at will with no particular repercussions.
//	
//	2. in an interactive graphical go:ngine app that also has the *core* package loaded:
//	
//	All "Sync"-related functions pertain to use-case #2, where the *assets* package essentially becomes the
//	live repository for all resource definitions loaded, used, or manipulated by the *core* package at runtime.
//	So now every image definition in *assets* may have a corresponding GPU-bound texture object in *core*, every
//	*assets* mesh definition may be bound to a *core* *MeshBuffer*, etc.
//	
//	Structure: generally speaking, all resource types are organized in a consistent fashion as follows ---
//	users familiar with the COLLADA format will notice a stark resemblance in terminology and resource organization:
//	
//	1. First, there is a *FooDef* struct for the one-time definition of a unique resource:
//	GeometryMeshDef, ImageDef, LightDef, MaterialDef etc.
//	
//	2. Next, there is a smaller *FooInst* struct for handling many individual (sometimes parameterized) instances of
//	an *FooDef*: GeometryMeshInst, ImageInst, LightInst, MaterialInst etc.
//	
//	3. Finally, there is a light-weight *LibFooDefs* struct type (based on a *map* collection) containing Defs
//	associated with their ID: LibGeometryMeshDefs, LibImageDefs, LibLightDefs, LibMaterialDefs etc.
//	
//	4. The package also provides a pre-initialized global *FooDefs* variable for each such Lib type, in simple apps
//	considered the "default / main / only Lib you'll need": GeometryMeshDefs, ImageDefs, LightDefs, MaterialDefs etc.
//	
//	5. For more complex use-cases, you can also organize multiple libs for any given resource type in the global
//	*AllFooDefLibs* variable, essentially a *map* collection of Libs: AllGeometryMeshDefLibs (of type
//	LibsGeometryMeshDef), AllImageDefLibs (of type LibsImageDef), AllLightDefLibs (of type LibsLightDef),
//	AllMaterialDefLibs (of type LibsMaterialDef) etc.
//	
//	Any exported types in this package not following the above pattern (such as MeshData etc.) should be
//	considered "auxiliary helpers" rather than primary / "first-class" resource types.
package assets
