Update May 2013:

- The project is "on hold for the time being" and has been for the last 2 months
- All go+gfx enthusiasts are more than welcome to pick up where I left off

Builds in Go 1.1 just as well as it did in Go 1.0.3, but performance has roughly halved and GC durations have roughly doubled. Sad story. Not sure why that is, since for most other projects out there, Go 1.1 had the opposite effect.

I suspect the current code-base overdid it badly on the prematurely-optimizing, over-clever "custom allocators, as few pointers as possible" approach. Lesson learned?

That being said---I still think it's a pretty great library to base your GL gfx/engine efforts on, rather than starting from scratch.

(Of course, the real buffs won't be considering exactly that piece of advice ;)

Getting going:
==============

First, **go get github.com/go-gl/glfw** --- this may not compile at first, until you have the glfw.org libs installed properly on your system.

This go-gl/glfw package needs to be working on your system before you can proceed.

Once it does, **go get github.com/go3d/go-ngine** should in theory download and build everything. This may well take a couple of minutes!

There shouldn't be any compilation errors if the above glfw package is fully installed.

Now you should be able to **go run** pretty much any main.go packages underneath *go3d/go-ngine/_examples* -- note the programs need write access to *go3d/go-ngine/_examples/-app-data/_tmp*



go-ngine
========

An in-development OpenGL-based "3D rendering engine" written in Go, focusing on interactive real-time apps such as games or demos.

Early work-in-progress, "progressing" at a rather leisurely pace too. Performance is a high priority, slowing development considerably. **Not really ready for developer use just yet...**

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

Coming up soon:

- more cullings
- lighting models
- HDR (16-bit) pipeline, tone-mapping... but gonna need decent lighting (and shadows) first, obviously!

... you name it.
