load("@io_bazel_rules_docker//container:pull.bzl", "container_pull")

def containers():
    container_pull(
        name = "alpine_linux_amd64",
        registry = "index.docker.io",
        repository = "library/alpine",
        tag = "3.14.2",
    )
