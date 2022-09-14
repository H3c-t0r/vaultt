package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/docker"
)

func PrepareTestContainer(t *testing.T, version string) (func(), string) {
	env := []string{
		"POSTGRES_PASSWORD=secret",
		"POSTGRES_DB=database",
	}

	_, cleanup, url, _ := prepareTestContainer(t, "postgres", "postgres", version, "secret", "database", true, false, false, env)

	return cleanup, url
}

func PrepareTestContainerWithPassword(t *testing.T, version, password string) (func(), string) {
	env := []string{
		"POSTGRES_PASSWORD=" + password,
		"POSTGRES_DB=database",
	}

	_, cleanup, url, _ := prepareTestContainer(t, "postgres", "postgres", version, password, "database", true, false, false, env)

	return cleanup, url
}

func PrepareTestContainerRepmgr(t *testing.T, name, version string, envVars []string) (*docker.Runner, func(), string, string) {
	return prepareTestContainer(t, name, "bitnami/postgresql-repmgr", version, "secret", "", false, true, true, envVars)
}

func StopContainer(t *testing.T, ctx context.Context, runner *docker.Runner, containerID string) {
	if err := runner.Stop(ctx, containerID); err != nil {
		t.Fatalf("Could not stop docker Postgres: %s", err)
	}
}

func RestartContainer(t *testing.T, ctx context.Context, runner *docker.Runner, containerID string) {
	if err := runner.Restart(ctx, containerID); err != nil {
		t.Fatalf("Could not restart docker Postgres: %s", err)
	}
}

func prepareTestContainer(t *testing.T, name, repo, version, password, db string, addSuffix, forceLocalAddr, doNotAutoRemove bool, envVars []string) (*docker.Runner, func(), string, string) {
	if os.Getenv("PG_URL") != "" {
		return nil, func() {}, "", os.Getenv("PG_URL")
	}

	if version == "" {
		version = "11"
	}

	runOpts := docker.RunOptions{
		ContainerName:   name,
		ImageRepo:       repo,
		ImageTag:        version,
		Env:             envVars,
		Ports:           []string{"5432/tcp"},
		DoNotAutoRemove: doNotAutoRemove,
	}
	if repo == "bitnami/postgresql-repmgr" {
		runOpts.NetworkID = os.Getenv("POSTGRES_MULTIHOST_NET")
	}

	runner, err := docker.NewServiceRunner(runOpts)
	if err != nil {
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	svc, containerID, err := runner.StartNewService(context.Background(), addSuffix, forceLocalAddr, connectPostgres(password, repo))
	if err != nil {
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	return runner, svc.Cleanup, svc.Config.URL().String(), containerID
}

func connectPostgres(password, repo string) docker.ServiceAdapter {
	return func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		u := url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword("postgres", password),
			Host:     fmt.Sprintf("%s:%d", host, port),
			Path:     "postgres",
			RawQuery: "sslmode=disable",
		}

		db, err := sql.Open("pgx", u.String())
		if err != nil {
			return nil, err
		}
		defer db.Close()

		if err = db.Ping(); err != nil {
			return nil, err
		}
		return docker.NewServiceURL(u), nil
	}
}
