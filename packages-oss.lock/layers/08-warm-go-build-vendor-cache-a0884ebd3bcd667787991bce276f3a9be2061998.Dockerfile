# WARNING: Do not EDIT or MERGE this file, it is generated by 'packagespec lock'.
ARG BASE_IMAGE
FROM $BASE_IMAGE
COPY . ./
ENV GOOS=solaris
ENV GOARCH=amd64
ENV CGO_ENABLED=0
# Try to build vendored packages. We first filter out packages which report
# errors in 'go list', because trying to run go build ./vendor/... fails early
# if we include them. We also don't care about the exit code here, because
# some of the vendored packages may fail to build, but this won't necessarily
# mean that the final package will fail to build, and we will still get a
# usefully warmed cache.
RUN go list -f '{{.ImportPath}}{{if or .Error .DepsErrors}} ERROR{{end}}' ./vendor/... | grep -v ERROR | xargs go build -v || true