Getting going:
==============

First, `go get github.com/go-gl/glfw` --- this may not compile at first, until you have the http://glfw.org dev libs installed properly on your operating system.

This `go-gl/glfw` package needs to be fully installed in your `$GOPATH/pkg` before you can proceed!

Once it is, `go get github.com/go3d/go-ngine` should download and build everything including dependencies. This may well take a couple of minutes!

There shouldn't be any compilation errors if the above-mentioned `go-gl/glfw` package is fully installed.

Now you should be able to `go run` pretty much any `main.go` packages underneath `go3d/go-ngine/_examples` -- note the programs need write access to the `go3d/go-ngine/_examples/-app-data/_tmp` directory!



go-ngine
========

An in-development OpenGL-based "3D rendering engine" written in Go, focusing on interactive real-time apps such as games or demos.

Early work-in-progress, "progressing" at a rather leisurely pace too. Performance is a high priority, slowing development considerably. **Not really ready for any "production" use-cases just yet...**

Current state:

- parallel: while rendering the current frame N+0, *concurrently* (on multi-core) prepares rendering for the next frame N+1 *and* executes app logic for next-next frame N+2.
- renders a scene of textured primitive meshes such as cubes, pyramids, planes, quads, with multiple cameras and render-to-texture pipeline
- a single mesh can be reused by multiple "models" (differently textured), which can be reused by multiple "nodes" (essentially a potentially recursive hierarchy of transformations applied to a model in world space).
- cubemapped sky-mesh. Any mesh works in theory (try a sphere if you want a dome), tested so far cube and pyramid (which surprisingly looks just-as-good with a LOT less geometry...)
- Permutated fx-shaders / uber-shader system (fairly simple-minded for now)
- post-processing effects -- but only a few very simple ones implemented yet (gamma, grayscale...)
- per-face materials/effects
- multi-texturing (specifically, multiple fx of the same type in an effect)
- texture disk cache. Feed the go:ngine normal un-processed texture images -- it fetches instead from a temp dir (or creates in there only if missing) a file containing its equivalent, but readily pre-processed & re-aligned pixel data.
- new: batched rendering
- new: frustum culling

TODO:

- import/load actual 3D model assets (instead of the current pre-fab cubes, pyramids etc.)
- more cullings
- basic lighting models (per-vertex, per-pixel, deferred...)
- HDR (16-bit) pipeline, tone-mapping... but gonna need decent lighting (and shadows) first, obviously!

... you name it.
