package fpth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func wrapAdapt(t *testing.T, expected string, pth string, opts ...Option) {
	result, err := Adapt(pth, opts...)
	assert.Nil(t, err, "parse `%s` should not failed", pth)
	assert.Equal(t, expected, result)
}

func testError(t *testing.T, pth string) {
	_, err := Adapt(pth)
	assert.NotNil(t, err, "parse `%s` should failed", pth)
}

func TestAdapt(t *testing.T) {
	testError(t, "")

	wrapAdapt(t, "/", "/")
	// ..
	wrapAdapt(t, "..", "..")
	wrapAdapt(t, "/a", "/a/b/..")
	wrapAdapt(t, "/", "/a/..")
	wrapAdapt(t, "/", "/a/b/../..")
	wrapAdapt(t, "/", "/a/b/../../../..")
	// .
	wrapAdapt(t, ".", ".")
	wrapAdapt(t, "a", "./a")
	wrapAdapt(t, "a", "a/.")
	wrapAdapt(t, "/a/b/c", "/a/b/./c")
	// sym
	wrapAdapt(t, "~/a", "~/a")
	wrapAdapt(t, ":/a/b/c", ":/a/b/./c")
	wrapAdapt(t, "++@#$%^", "++@#$%^")

	// home dir
	wrapAdapt(t, cachedHomePath, "~/", OEnableHomeDir())
	wrapAdapt(t, Join(cachedHomePath, "a", "b"), "~/a/b", OEnableHomeDir())
	wrapAdapt(t, Clean(Join(cachedHomePath, "..")), "~/a/../..", OEnableHomeDir())

	// relative dir
	wrapAdapt(t, cachedPWDPath, "./", ORelativePWDPath())
	wrapAdapt(t, Join(cachedPWDPath, "a", "b"), "./a/b", ORelativePWDPath())
	wrapAdapt(t, Clean(Join(cachedPWDPath, "..")), "./..", ORelativePWDPath())

	// relative dir
	wrapAdapt(t, "/a/b", "./", ORelativeGivenPath("/a/b"))
	wrapAdapt(t, "/a/b/a/b", "./a/b", ORelativeGivenPath("/a/b"))
	wrapAdapt(t, "/a", "./..", ORelativeGivenPath("/a/b"))
	wrapAdapt(t, "/a/c/d", "../c/d", ORelativeGivenPath("/a/b"))

	// relative header
	wrapAdapt(t, "/a/b/c/d", "@/c/d", ORelativeHeader("@", "/a/b"))
	wrapAdapt(t, "/b/c/d", "#/c/d", ORelativeHeader("@", "/a"), ORelativeHeader("#", "/b"))
}
