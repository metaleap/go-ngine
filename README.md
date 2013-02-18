go-ngine
========

An in-development OpenGL-based "3D rendering engine" written in Go, focusing on interactive real-time apps such as games or demos, implementing a **multi-threaded render pipeline**.

Early work-in-progress, "progressing" at a rather leisurely pace too. Performance is a high priority, slowing development considerably. **Not really ready for developer use just yet...**

Current state:

- renders a scene of textured primitive meshes such as cubes, pyramids, planes, quads, with multiple cameras and render-to-texture pipeline
- a single mesh can be reused by multiple "models" (differently textured), which can be reused by multiple "nodes" (essentially a potentially recursive hierarchy of transformations applied to a model in world space).
- cubemapped skybox
- Permutated fx-shaders / uber-shader system (fairly simple-minded for now)
- post-processing effects -- but only a few very simple ones implemented yet (gamma, grayscale...)

Coming up soon:

- per face-rather than per-model materials/effects
- batched rendering
- geometry culling
- lighting models
- HDR (16-bit) pipeline, tone-mapping

... you name it.
