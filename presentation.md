% Buildpacks: Towards 1.0, AI and Other Things
% Cloud Native Buildpacks
% 2026

# Introduction

## About This Talk

This session will explore two main threads:

- **The Road to 1.0** — achieving stability and feature completeness for widespread adoption
- **AI & Machine Learning** — simplifying how we build and deploy AI-driven applications

. . .

Along the way, we'll do a live demo building a Java application with a custom builder.

## Agenda

| Section | Time |
|---------|------|
| Buildpacks Demo — Java and Custom Builder | ~10 min |
| AI & ML — CUDA Base Images | ~10 min |
| Towards 1.0 — Roadmap and Key Milestones | ~5 min |
| Q & A | ~5 min |

## What Are Cloud Native Buildpacks?

- buildpacks.io maintains a _specification_
- Transform application source code into OCI container images
- No Dockerfile required
- Provide a **structured**, **repeatable** build process
- CNCF Incubating project

## CLoud Native Buildpacks Implementations

### Vendors:

- Heroku: [https://elements.heroku.com/buildpacks](https://elements.heroku.com/buildpacks)
- Paketo: [https://paketo.io/](https://paketo.io/)
- Google: [https://docs.cloud.google.com/docs/buildpacks/overview](https://docs.cloud.google.com/docs/buildpacks/overview)

# Demo: Building a Java Application

## The Goal

Build a production-ready Java application image using:

* As an application developer
  - `pack build` and `docker run`
* As a platform operator
  - a custom builder

## Application Developer

* LIVE(ish) Demo
  - (yes, it's pre-scripted)
* Animated `gif` available at [demo/java/demo.gif](demo/java//demo.gif)

```
$ pack build example --builder cuda-java-builder
```

. . . 

```
$ docker run --rm -it example
Backend: CpuBackend
Learned: y = 2.03x + 0.93  (loss: 0.000911)
```

## Platform Operator

* We have control over the **`--builder`**

### Questions
* What language stacks do we support in cusom images?
* What are the JDKs/JREs that we support in production?
* How much flexibility do we provide to application developers?

## Step 0 - Custom build/run images

:::: {.columns}
::: {.column width="33%"}
* Define build and run images with a `cnb` user
* Subsequent builds do not require `root`!
:::
::: {.column width="66%"}
```Dockerfile
FROM nvidia/cuda:13.1.1-cudnn-devel-ubuntu24.04

ARG cnb_uid=1001
ARG cnb_gid=1001
ENV CNB_USER_ID=${cnb_uid} \
    CNB_GROUP_ID=${cnb_gid}

RUN groupadd --gid ${cnb_gid} cnb \
 && useradd --uid ${cnb_uid} \ 
    --gid ${cnb_gid} -m cnb

USER ${cnb_uid}:${cnb_gid}
```
:::
::::

## Step 1 — Create a Custom Builder

:::: {.columns}
::: {.column width="33%"}
Define a `builder.toml` that references the Paketo Java buildpack:
:::
::: {.column width="66%"}
```toml
[build]
  image = "cuda-build:latest"

[[run.images]]
  image = "cuda-run:latest"

[[buildpacks]]
  uri = "paketobuildpacks/java"

[[targets]]
  os = "linux"
  arch = "amd64"

[[order]]
  [[order.group]]
    id = "paketo-buildpacks/java"
```
:::
::::
## Step 2 — Create the Builder Image

```bash
pack builder create cuda-java-builder \
  --config builder.toml
```

. . .

Verify it was created:

```bash
pack builder inspect cuda-java-builder
```

## Step 3 — Build the Java Application

```bash
pack build my-java-app \
  --builder cuda-java-builder \
  --path ./demo/java
```

. . .

```
===> DETECTING
paketo-buildpacks/ca-certificates 3.8.3
paketo-buildpacks/bellsoft-liberica 10.8.1
paketo-buildpacks/maven          6.15.14
paketo-buildpacks/executable-jar 6.10.3
paketo-buildpacks/spring-boot    5.29.2
===> BUILDING
...
===> EXPORTING
Successfully built image my-java-app
```

## Step 4 — Run It

```bash
docker run --rm -p 8080:8080 my-java-app
```

. . .
DEMO OUTPUT HERE

Started Application in 2.3 seconds
```

## Platform Owners: You Own the Builder

As a platform operator, the custom builder gives you **complete control**:

- **Pin specific buildpack versions** — reproducible builds across your org
- **Curate the buildpack ecosystem** — only approved buildpacks in your builder
- **Control the base images** — your security team manages the build & run images
- **Enforce compliance** — embed SBOMs, signing policies, and vulnerability scanning

. . .

> Developers get a simple `pack build` command.
> Operators get full governance.

# AI & Machine Learning

## The AI/ML Challenge

Building AI/ML applications involves unique infrastructure requirements:

- GPU drivers (NVIDIA CUDA)
- Large frameworks (PyTorch, TensorFlow)
- Complex native dependencies
- Reproducible environments across dev and prod

Buildpacks can help tame this complexity.

## Approach: CUDA Base Image + Buildpacks

```
┌────────────────────────────────────────┐
│            App Image (OCI)             │
├────────────────────────────────────────┤
│   PyTorch Hello World Application     │
├────────────────────────────────────────┤
│   Python Buildpack Layers              │
├────────────────────────────────────────┤
│   CUDA-enabled Base Image             │
│   (nvidia/cuda:12.x-runtime-ubuntu22) │
└────────────────────────────────────────┘
```

## Step 1 — Build a CUDA Base Image

Use an **image extension** to generate a CUDA-capable run image:

```dockerfile
# run.Dockerfile — generated by the extension
ARG base_image
FROM ${base_image}

# Install CUDA runtime libraries
RUN apt-get update && apt-get install -y --no-install-recommends \
    cuda-cudart-12-4 \
    libcublas-12-4 \
    libcufft-12-4 \
    libcurand-12-4 \
    libcusparse-12-4 \
    libcudnn9-cuda-12 \
    && rm -rf /var/lib/apt/lists/*

ENV NVIDIA_VISIBLE_DEVICES=all
ENV NVIDIA_DRIVER_CAPABILITIES=compute,utility
ENV LD_LIBRARY_PATH=/usr/local/cuda/lib64:${LD_LIBRARY_PATH}
```

## Step 2 — The PyTorch Hello World App

```python
# app.py
import torch
import platform

def main():
    print(f"Python: {platform.python_version()}")
    print(f"PyTorch: {torch.__version__}")
    print(f"CUDA available: {torch.cuda.is_available()}")

    if torch.cuda.is_available():
        device = torch.device("cuda")
        print(f"GPU: {torch.cuda.get_device_name(0)}")
    else:
        device = torch.device("cpu")

    # Create a simple tensor and perform a computation
    x = torch.randn(3, 3, device=device)
    y = torch.randn(3, 3, device=device)
    z = torch.matmul(x, y)
    print(f"Result tensor (on {device}):\n{z}")

if __name__ == "__main__":
    main()
```

## Step 2 — The Requirements File

```
# requirements.txt
torch>=2.2.0
torchvision>=0.17.0
numpy>=1.26.0
```

## Step 3 — Build with Buildpacks

```bash
pack build my-pytorch-app \
  --builder polyglot-cuda-builder \
  --path ./pytorch-hello-world
```

. . .

```
===> DETECTING
paketo-buildpacks/cpython     2.x.x
paketo-buildpacks/pip          1.x.x
paketo-buildpacks/pip-install  1.x.x
cuda-extension                 0.1.0
===> BUILDING
  Installing CPython 3.14.x
  Installing pip dependencies via requirements.txt
  Downloading torch-2.2.x ...
===> EXPORTING
Successfully built image my-pytorch-app
```

## Step 4 — Run with GPU Access

```bash
docker run --rm --gpus all my-pytorch-app
```

. . .

```
Python: 3.14.3
PyTorch: 2.2.1+cu124
CUDA available: True
GPU: NVIDIA A100-SXM4-40GB
Result tensor (on cuda):
tensor([[ 0.4521, -1.2345,  0.8901],
        [ 1.1234, -0.5678,  0.2345],
        [-0.3456,  0.7890, -1.0123]], device='cuda:0')
```

## Why Buildpacks for AI/ML?

- **Reproducibility** — consistent CUDA + Python environments
- **Security** — automatic base image updates without rebuilding
- **Separation of concerns** — data scientists write Python, platform teams manage CUDA
- **Rebase** — update GPU drivers in the base image without rebuilding application layers
    * CUDA libraries are ABI compatible across minor versions 

. . .

> The CUDA layer and application layer are managed independently — just like any other buildpacks workflow.

# Towards 1.0

## What Does 1.0 Mean?

- **Stable APIs** — Buildpack API and Platform API reach 1.0
- **No more breaking changes** without a major version bump
- **A promise** to the ecosystem: safe to build on, safe to depend on

. . .

The current versions:

| Component | Current | Target |
|-----------|---------|--------|
| Buildpack API | 0.13 | 1.0 |
| Platform API | 0.16 | 1.0 |
| Lifecycle | 0.20.x | 1.0 |
| Pack CLI | 0.40.x | 1.0 |

## Key Components Driving 1.0

```
  ┌──────────────────────────────────────────┐
  │              Spec (spec repo)            │
  │  Buildpack API · Platform API · Distro   │
  └──────────┬───────────────┬───────────────┘
             │               │
    ┌────────▼────┐   ┌──────▼──────┐
    │  Lifecycle  │   │   Pack CLI  │
    │  (0.20.x)  │   │  (0.40.x)   │
    └─────────────┘   └─────────────┘
```

* The Spec defines contracts. The Lifecycle implements them.
* Pack CLI is the primary user-facing tool.

## Breaking Changes That Must Land Before 1.0

These approved RFCs represent **intentional breaking changes** — they must ship before the API is frozen:

- **RFC #0096 — Remove Stacks and Mixins**
  - Stacks are replaced by build/run image targets
  - Simplifies the mental model significantly

- **RFC #0093 — Remove Shell Processes**
  - All processes become direct processes
  - Improved security and signal handling

- **RFC #0105 — Dockerfiles (Image Extensions)**
  - Allows customizing build/run images via Dockerfiles
  - Critical for use cases like CUDA, as we just saw

## Active RFCs Shaping the Future

| RFC | Title | Impact |
|-----|-------|--------|
| #0134 | Execution Environments | Define runtime contract for apps |
| #0131 | Build Observability (OTEL) | Tracing and metrics for builds |
| #0130 | OCI Image Annotations | Richer metadata on output images |
| #0128 | Multi-arch Support | Build once, run on amd64 & arm64 |
| #0125 | Parallel Cache/Image Export | Faster builds |
| #0113 | Additional OCI Artifacts | SBOMs, signatures as OCI artifacts |

## Lifecycle 0.21.0 — The Next Milestone

8 open issues targeting the next lifecycle release:

- **Export run image metadata** — richer output image labels
- **Report.toml enhancements** — more build data for CI/CD
- **Corrupt cache recovery** — security hardening
- **containerd socket export** — direct export to containerd
- **Self-signed certificate support** — enterprise-friendly Kubernetes
- **FreeBSD support** — expanding platform reach
- **OCI layout + extensions** — run image extensions with OCI layout export

## Pack CLI 0.41.0 — What's Coming

24 open issues in the 0.41.0 milestone, including:

- **Platform API 0.14 support** — `-run` flag in restorer
- **Podman compatibility** — fixing docker host issues
- **`pack extension new`** — scaffolding for image extensions
- **Cosign image signing** — sign `buildpacksio/pack` images + SBOM
- **containerd daemon support** — publish-then-pull workaround
- **`project.toml` directory exclusion** — fix glob behavior
- **`try-always` pull policy** — new pull policy option

## Spec Stabilization Roadmap

```
 Now          ──────────▶          1.0
  │                                 │
  ├─ Buildpack API 0.13             │
  │   └── PATH delimiter fix       │
  │                                 │
  ├─ Platform API 0.16              │
  │   └── Cosign/SBOM spec         │
  │                                 │
  ├─ Distribution 0.3 / 0.4        │
  │   └── Builder order changes    │
  │   └── Multi-platform builders  │
  │                                 │
  └── Remove deprecated features ──┘
      (Stacks, Mixins, Shell Procs)
```

## What 1.0 Means for You

**For Application Developers:**

- Stable CLI and build experience — no surprises
- Multi-arch images out of the box
- Better error messages and observability

. . .

**For Buildpack Authors:**

- A frozen Buildpack API to target
- Image extensions for advanced customization
- First-class SBOM and signing support

. . .

**For Platform Operators:**

- Stable Platform API to integrate against
- Builder spec for governance
- OCI-native — containerd, Kubernetes-ready

## How to Get Involved

- **GitHub:** [github.com/buildpacks](https://github.com/buildpacks)
- **RFCs:** [github.com/buildpacks/rfcs](https://github.com/buildpacks/rfcs) — propose and discuss changes
- **Slack:** [slack.buildpacks.io](https://slack.buildpacks.io) — join the community
- **Working Group:** Weekly calls — all are welcome

. . .

> We're a CNCF project — contributions from every perspective make us stronger.

# Thank You

## Questions?

**Buildpacks: Towards 1.0, AI and Other Things**

- Docs: [buildpacks.io/docs](https://buildpacks.io/docs)
- GitHub: [github.com/buildpacks](https://github.com/buildpacks)
- Slack: [slack.buildpacks.io](https://slack.buildpacks.io)

# Appendix

## Control JDK/JRE Versions in Builder

* Paketo specific
* Repackage image with updated metadata
  - pull down packeto image
  - remove unwanted JDK/JREs from `buildpack.toml`
  - update buildpack id
  - push the image internally

## Buildpacks Rebase

* 