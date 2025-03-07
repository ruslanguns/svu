package main

import (
	"testing"

	"github.com/Masterminds/semver"
	"github.com/matryer/is"
)

func TestBuildVersion(t *testing.T) {
	t.Run("dev", func(t *testing.T) {
		is.New(t).Equal("svu version dev", buildVersion("dev", "", "", ""))
	})

	t.Run("goreleaser built", func(t *testing.T) {
		is.New(t).Equal(`svu version v1.2.3
commit: a123cd
built at: 2021-01-02
built by: goreleaser`, buildVersion("v1.2.3", "a123cd", "2021-01-02", "goreleaser"))
	})
}

func TestCmd(t *testing.T) {
	ver := func() *semver.Version { return semver.MustParse("1.2.3-pre+123") }
	t.Run(currentCmd.FullCommand(), func(t *testing.T) {
		cmd := currentCmd.FullCommand()
		t.Run("version has meta", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, ver(), "v1.2.3", "", "", false)
			is.NoErr(err)
			is.Equal("1.2.3-pre+123", v.String())
		})
		t.Run("version is clean", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, semver.MustParse("v1.2.3"), "v1.2.3", "doesnt matter", "nope", true)
			is.NoErr(err)
			is.Equal("1.2.3", v.String())
		})
	})

	t.Run(minorCmd.FullCommand(), func(t *testing.T) {
		cmd := minorCmd.FullCommand()
		t.Run("no meta", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, ver(), "v1.2.3", "", "", false)
			is.NoErr(err)
			is.Equal("1.3.0", v.String())
		})
		t.Run("build", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, ver(), "v1.2.3", "", "124", false)
			is.NoErr(err)
			is.Equal("1.3.0+124", v.String())
		})
		t.Run("prerel", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, ver(), "v1.2.3", "alpha.1", "", false)
			is.NoErr(err)
			is.Equal("1.3.0-alpha.1", v.String())
		})
		t.Run("all meta", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, ver(), "v1.2.3", "alpha.2", "125", false)
			is.NoErr(err)
			is.Equal("1.3.0-alpha.2+125", v.String())
		})
	})

	t.Run(patchCmd.FullCommand(), func(t *testing.T) {
		cmd := patchCmd.FullCommand()
		t.Run("no meta", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, semver.MustParse("1.2.3"), "v1.2.3", "", "", false)
			is.NoErr(err)
			is.Equal("1.2.4", v.String())
		})
		t.Run("previous had meta", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, semver.MustParse("1.2.3-alpha.1+1"), "v1.2.3", "", "", false)
			is.NoErr(err)
			is.Equal("1.2.3", v.String())
		})
		t.Run("previous had meta, force", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, semver.MustParse("1.2.3-alpha.1+1"), "v1.2.3", "", "", true)
			is.NoErr(err)
			is.Equal("1.2.4", v.String())
		})
		t.Run("previous had meta, force, add meta", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, semver.MustParse("1.2.3-alpha.1+1"), "v1.2.3-alpha.1+1", "alpha.2", "10", true)
			is.NoErr(err)
			is.Equal("1.2.4-alpha.2+10", v.String())
		})
		t.Run("previous had meta, change it", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, semver.MustParse("1.2.3-alpha.1+1"), "v1.2.3-alpha.1+1", "alpha.2", "10", false)
			is.NoErr(err)
			is.Equal("1.2.3-alpha.2+10", v.String())
		})
		t.Run("build", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, semver.MustParse("1.2.3"), "v1.2.3", "", "124", false)
			is.NoErr(err)
			is.Equal("1.2.4+124", v.String())
		})
		t.Run("prerel", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, semver.MustParse("1.2.3"), "v1.2.3", "alpha.1", "", false)
			is.NoErr(err)
			is.Equal("1.2.4-alpha.1", v.String())
		})
		t.Run("all meta", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, semver.MustParse("1.2.3"), "v1.2.3", "alpha.2", "125", false)
			is.NoErr(err)
			is.Equal("1.2.4-alpha.2+125", v.String())
		})
	})

	t.Run(majorCmd.FullCommand(), func(t *testing.T) {
		cmd := majorCmd.FullCommand()
		t.Run("no meta", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, ver(), "v1.2.3", "", "", false)
			is.NoErr(err)
			is.Equal("2.0.0", v.String())
		})
		t.Run("build", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, ver(), "v1.2.3", "", "124", false)
			is.NoErr(err)
			is.Equal("2.0.0+124", v.String())
		})
		t.Run("prerel", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, ver(), "v1.2.3", "alpha.1", "", false)
			is.NoErr(err)
			is.Equal("2.0.0-alpha.1", v.String())
		})
		t.Run("all meta", func(t *testing.T) {
			is := is.New(t)
			v, err := nextVersion(cmd, ver(), "v1.2.3", "alpha.2", "125", false)
			is.NoErr(err)
			is.Equal("2.0.0-alpha.2+125", v.String())
		})
	})

	t.Run("errors", func(t *testing.T) {
		t.Run("invalid build", func(t *testing.T) {
			is := is.New(t)
			_, err := nextVersion(minorCmd.FullCommand(), semver.MustParse("1.2.3"), "v1.2.3", "", "+125", false)
			is.True(err != nil)
		})
		t.Run("invalid prerelease", func(t *testing.T) {
			is := is.New(t)
			_, err := nextVersion(minorCmd.FullCommand(), semver.MustParse("1.2.3"), "v1.2.3", "+aaa", "", false)
			is.True(err != nil)
		})
	})
}
