package dockertest

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"strings"

	"github.com/cenk/backoff"
	dc "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
)

// Pool represents a connection to the docker API and is used to create and remove docker images.
type Pool struct {
	Client  *dc.Client
	MaxWait time.Duration
}

// Resource represents a docker container.
type Resource struct {
	Container *dc.Container
}

// GetPort returns a resource's published port. You can use it to connect to the service via localhost, e.g. tcp://localhost:1231/
func (r *Resource) GetPort(id string) string {
	if r.Container == nil {
		return ""
	} else if r.Container.NetworkSettings == nil {
		return ""
	}

	m, ok := r.Container.NetworkSettings.Ports[dc.Port(id)]
	if !ok {
		return ""
	} else if len(m) == 0 {
		return ""
	}

	return m[0].HostPort
}

func (r *Resource) GetBoundIP(id string) string {
	if r.Container == nil {
		return ""
	} else if r.Container.NetworkSettings == nil {
		return ""
	}

	m, ok := r.Container.NetworkSettings.Ports[dc.Port(id)]
	if !ok {
		return ""
	} else if len(m) == 0 {
		return ""
	}

	return m[0].HostIP
}

// NewTLSPool creates a new pool given an endpoint and the certificate path. This is required for endpoints that
// require TLS communication.
func NewTLSPool(endpoint, certpath string) (*Pool, error) {
	ca := fmt.Sprintf("%s/ca.pem", certpath)
	cert := fmt.Sprintf("%s/cert.pem", certpath)
	key := fmt.Sprintf("%s/key.pem", certpath)

	client, err := dc.NewTLSClient(endpoint, cert, key, ca)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &Pool{
		Client: client,
	}, nil
}

// NewPool creates a new pool. You can pass an empty string to use the default, which is taken from the environment
// variable DOCKER_URL, or from docker-machine if the environment variable DOCKER_MACHINE_NAME is set,
// or if neither is defined a sensible default for the operating system you are on.
func NewPool(endpoint string) (*Pool, error) {
	if endpoint == "" {
		if os.Getenv("DOCKER_URL") != "" {
			endpoint = os.Getenv("DOCKER_URL")
		} else if os.Getenv("DOCKER_MACHINE_NAME") != "" {
			client, err := dc.NewClientFromEnv()
			if err != nil {
				return nil, errors.Wrap(err, "")
			}

			return &Pool{Client: client}, nil
		} else if runtime.GOOS == "windows" {
			endpoint = "http://localhost:2375"
		} else {
			endpoint = "unix:///var/run/docker.sock"
		}
	}

	client, err := dc.NewClient(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &Pool{
		Client: client,
	}, nil
}

// RunOptions is used to pass in optional parameters when running a container.
type RunOptions struct {
	Repository   string
	Tag          string
	Env          []string
	Entrypoint   []string
	Cmd          []string
	Mounts       []string
	Links        []string
	ExposedPorts []string
	Auth         dc.AuthConfiguration
}

// RunWithOptions starts a docker container.
//
// pool.Run(&RunOptions{Repository: "mongo", Cmd: []string{"mongod", "--smallfiles"}})
func (d *Pool) RunWithOptions(opts *RunOptions) (*Resource, error) {
	repository := opts.Repository
	tag := opts.Tag
	env := opts.Env
	cmd := opts.Cmd
	ep := opts.Entrypoint
	var exp map[dc.Port]struct{}

	if len(opts.ExposedPorts) > 0 {
		exp = map[dc.Port]struct{}{}
		for _, p := range opts.ExposedPorts {
			exp[dc.Port(p)] = struct{}{}
		}
	}

	mounts := []dc.Mount{}

	for _, m := range opts.Mounts {
		sd := strings.Split(m, ":")
		if len(sd) == 2 {
			mounts = append(mounts, dc.Mount{
				Source:      sd[0],
				Destination: sd[1],
				RW:          true,
			})
		} else {
			return nil, errors.Wrap(fmt.Errorf("invalid mount format: got %s, expected <src>:<dst>", m), "")
		}
	}

	if tag == "" {
		tag = "latest"
	}

	_, err := d.Client.InspectImage(fmt.Sprintf("%s:%s", repository, tag))
	if err != nil {
		if err := d.Client.PullImage(dc.PullImageOptions{
			Repository: repository,
			Tag:        tag,
		}, opts.Auth); err != nil {
			return nil, errors.Wrap(err, "")
		}
	}

	c, err := d.Client.CreateContainer(dc.CreateContainerOptions{
		Config: &dc.Config{
			Image:        fmt.Sprintf("%s:%s", repository, tag),
			Env:          env,
			Entrypoint:   ep,
			Cmd:          cmd,
			Mounts:       mounts,
			ExposedPorts: exp,
		},
		HostConfig: &dc.HostConfig{
			PublishAllPorts: true,
			Binds:           opts.Mounts,
			Links:           opts.Links,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	if err := d.Client.StartContainer(c.ID, nil); err != nil {
		return nil, errors.Wrap(err, "")
	}

	c, err = d.Client.InspectContainer(c.ID)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &Resource{
		Container: c,
	}, nil
}

// Run starts a docker container.
//
// pool.Run("mysql", "5.3", []string{"FOO=BAR", "BAR=BAZ"})
func (d *Pool) Run(repository, tag string, env []string) (*Resource, error) {
	return d.RunWithOptions(&RunOptions{Repository: repository, Tag: tag, Env: env})
}

// Purge removes a container and linked volumes from docker.
func (d *Pool) Purge(r *Resource) error {
	if err := d.Client.KillContainer(dc.KillContainerOptions{ID: r.Container.ID}); err != nil {
		return errors.Wrap(err, "")
	}

	if err := d.Client.RemoveContainer(dc.RemoveContainerOptions{ID: r.Container.ID, Force: true, RemoveVolumes: true}); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

// Retry is an exponential backoff retry helper. You can use it to wait for e.g. mysql to boot up.
func (d *Pool) Retry(op func() error) error {
	if d.MaxWait == 0 {
		d.MaxWait = time.Minute
	}
	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = time.Second * 5
	bo.MaxElapsedTime = d.MaxWait
	return backoff.Retry(op, bo)
}
