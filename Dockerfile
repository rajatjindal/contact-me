#FROM ghcr.io/fermyon/spin:v3.0.0-rc.1-distroless
FROM rajatjindal/spin:2906@sha256:91bcca98b9bdbb42757bb6266b143bf2639d9af42fed3e1e7d63566c06dfcbbd

WORKDIR /app

COPY spin.toml spin.toml
COPY main.wasm main.wasm
COPY runtimeconfig.toml runtimeconfig.toml

ENTRYPOINT ["spin", "up"]
