go-ngine
========

An in-development OpenGL-based "3D rendering engine" written in Go, focusing on interactive real-time apps such as games or demos, implementing a **multi-threaded render pipeline**.

Early work-in-progress, "progressing" at a rather leisurely pace too. Performance is a high priority, slowing development considerably. **Not really ready for developer use just yet...**

Current state:

- renders a scene of textured primitive meshes such as cubes, pyramids, planes, quads, with multiple cameras and render-to-texture pipeline
- a single mesh can be reused by multiple "models" (differently textured), which can be reused by multiple "nodes" (essentially a potentially recursive hierarchy of transformations applied to a model in world space).
- cubemapped sky-mesh. Any mesh works in theory (try a sphere if you want a dome), tested so far cube and pyramid (which surprisingly looks just-as-good with a LOT less geometry...)
- Permutated fx-shaders / uber-shader system (fairly simple-minded for now)
- post-processing effects -- but only a few very simple ones implemented yet (gamma, grayscale...)
- new: per-face (rather than just per-model) materials/effects
- new: multi-texturing (specifically, multiple fx of the same type in an effect)
- new: texture disk cache. Feed the go:ngine normal un-processed texture images -- it fetches instead from a temp dir (or creates in there only if missing) a file containing its equivalent, but readily pre-processed & re-aligned pixel data.

Coming up soon:

- batched rendering
- geometry culling
- lighting models
- HDR (16-bit) pipeline, tone-mapping... but gonna need decent lighting (and shadows) first, obviously!

... you name it.
