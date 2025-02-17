package config

import (
	"testing"
)

func TestNoRedis(t *testing.T) {
	_, err := validate("/s.rump", "/t.rump", false, false, 60*1024)
	if err == nil {
		t.Error("file-only operations should not be supported")
	}
}

func TestNoFrom(t *testing.T) {
	_, err := validate("", "redis://t", false, false, 60*1024)
	if err == nil {
		t.Error("from should be required")
	}
}

func TestNoTo(t *testing.T) {
	_, err := validate("redis://s", "", false, false, 60*1024)
	if err == nil {
		t.Error("to should be required")
	}
}

func TestFromRedisToRedis(t *testing.T) {
	cfg, err := validate("redis://s", "redis://t", false, false, 60*1024)
	if err != nil {
		t.Error("from redis to redis should work")
	}

	if !cfg.Source.IsRedis {
		t.Error("wrong from")
	}

	if !cfg.Target.IsRedis {
		t.Error("wrong to")
	}

	if cfg.Source.URI != "redis://s" {
		t.Error("wrong source")
	}

	if cfg.Target.URI != "redis://t" {
		t.Error("wrong target")
	}
}

func TestFromRedisToFile(t *testing.T) {
	cfg, err := validate("redis://s", "/t.rump", false, false, 60*1024)
	if err != nil {
		t.Error("from redis to file should work")
	}

	if !cfg.Source.IsRedis {
		t.Error("wrong from")
	}

	if cfg.Target.IsRedis {
		t.Error("wrong to")
	}

	if cfg.Source.URI != "redis://s" {
		t.Error("wrong source")
	}

	if cfg.Target.URI != "/t.rump" {
		t.Error("wrong target")
	}
}

func TestFromFileToRedis(t *testing.T) {
	cfg, err := validate("/s.rump", "redis://t", false, false, 60*1024)
	if err != nil {
		t.Error("from file to redis should work")
	}

	if cfg.Source.IsRedis {
		t.Error("wrong from")
	}

	if !cfg.Target.IsRedis {
		t.Error("wrong to")
	}

	if cfg.Source.URI != "/s.rump" {
		t.Error("wrong source")
	}

	if cfg.Target.URI != "redis://t" {
		t.Error("wrong target")
	}
}
