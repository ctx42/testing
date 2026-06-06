package randkit

import (
	"strings"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_WithChars(t *testing.T) {
	// --- Given ---
	def := &options{}

	// --- When ---
	WithChars("ab", "cd", "ef")(def)

	// --- Then ---
	want := &options{
		chars:  "abcdef",
		n:      0,
		prefix: "",
		suffix: "",
	}
	assert.Equal(t, want, def)
}

func Test_WithPrefix(t *testing.T) {
	// --- Given ---
	def := &options{}

	// --- When ---
	WithPrefix("abc")(def)

	// --- Then ---
	want := &options{
		chars:  "",
		n:      0,
		prefix: "abc",
		suffix: "",
	}
	assert.Equal(t, want, def)
}

func Test_WithSuffix(t *testing.T) {
	// --- Given ---
	def := &options{}

	// --- When ---
	WithSuffix("abc")(def)

	// --- Then ---
	want := &options{
		chars:  "",
		n:      0,
		prefix: "",
		suffix: "abc",
	}
	assert.Equal(t, want, def)
}

func Test_WithExt(t *testing.T) {
	// --- Given ---
	def := &options{}

	// --- When ---
	WithExt(".exe")(def)

	// --- Then ---
	want := &options{
		chars:  "",
		n:      0,
		prefix: "",
		suffix: ".exe",
	}
	assert.Equal(t, want, def)
}

func Test_WithLen(t *testing.T) {
	// --- Given ---
	def := &options{}

	// --- When ---
	WithLen(42)(def)

	// --- Then ---
	want := &options{
		chars:  "",
		n:      42,
		prefix: "",
		suffix: "",
	}
	assert.Equal(t, want, def)
}

func Test_Str(t *testing.T) {
	t.Run("does not repeat", func(t *testing.T) {
		// --- Given ---
		history := make(map[string]struct{})

		// --- Then ---
		for i := 0; i < 1000; i++ {
			str := Str()
			if _, ok := history[str]; ok {
				t.Error("did not expect string to repeat")
			}
			history[str] = struct{}{}
			want := 10
			have := len(str)
			if want != have {
				t.Errorf("expected len(s)=%d got %d", want, have)
			}
		}
	})

	t.Run("set prefix", func(t *testing.T) {
		// --- When ---
		str := Str(WithPrefix("pref-"))

		// --- Then ---
		assert.True(t, strings.HasPrefix(str, "pref-"))
	})

	t.Run("set suffix", func(t *testing.T) {
		// --- When ---
		str := Str(WithSuffix("-suf"))

		// --- Then ---
		assert.True(t, strings.HasSuffix(str, "-suf"))
	})

	t.Run("set prefix and suffix", func(t *testing.T) {
		// --- When ---
		str := Str(WithPrefix("pref-"), WithSuffix("-suf"))

		// --- Then ---
		assert.True(t, strings.HasPrefix(str, "pref-"))
		assert.True(t, strings.HasSuffix(str, "-suf"))
	})

	t.Run("has only lowercase letters", func(t *testing.T) {
		// --- Then ---
		for i := 0; i < 500; i++ {
			str := Str(WithChars(Lowercase))
			want := strings.ToLower(str)
			have := str
			if want != have {
				t.Errorf("expected only lowercase letters got %s", have)
			}
		}
	})

	t.Run("has only uppercase letters", func(t *testing.T) {
		// --- Then ---
		for i := 0; i < 500; i++ {
			str := Str(WithChars(Uppercase))
			want := strings.ToUpper(str)
			have := str
			if want != have {
				t.Errorf("expected only uppercase letters got %s", have)
			}
		}
	})
}

func Test_FileName(t *testing.T) {
	t.Run("no dir prefix ext", func(t *testing.T) {
		// --- When ---
		name := FileName("")

		// --- Then ---
		assert.True(t, strings.HasPrefix(name, "file-"))
		assert.True(t, strings.HasSuffix(name, ".txt"))
		assert.Len(t, len("file-")+7+len(".txt"), name)
	})

	t.Run("no dir prefix with ext", func(t *testing.T) {
		// --- When ---
		name := FileName("", WithExt(".my"))

		// --- Then ---
		assert.True(t, strings.HasPrefix(name, "file-"))
		assert.True(t, strings.HasSuffix(name, ".my"))
		assert.Len(t, len("file-")+7+len(".my"), name)
	})

	t.Run("no dir with prefix ext", func(t *testing.T) {
		// --- When ---
		name := FileName("", WithPrefix("prefix-"), WithExt(".my"))

		// --- Then ---
		assert.True(t, strings.HasPrefix(name, "prefix-"))
		assert.True(t, strings.HasSuffix(name, ".my"))
		assert.Len(t, len("prefix-")+7+len(".my"), name)
	})

	t.Run("with dir prefix ext", func(t *testing.T) {
		// --- When ---
		name := FileName("/dir", WithPrefix("prefix-"), WithExt(".my"))

		// --- Then ---
		assert.True(t, strings.HasPrefix(name, "/dir/prefix-"))
		assert.True(t, strings.HasSuffix(name, ".my"))
		assert.Len(t, len("/dir/prefix-")+7+len(".my"), name)
	})

	t.Run("has only lowercase letters", func(t *testing.T) {
		// --- Then ---
		for i := 0; i < 500; i++ {
			name := FileName(
				"/dir",
				WithPrefix(""),
				WithSuffix(""),
				WithChars(Lowercase),
			)
			want := strings.ToLower(name)
			have := name
			if want != have {
				t.Errorf("expected only lowercase letters got %s", have)
			}
		}
	})

	t.Run("has only uppercase letters", func(t *testing.T) {
		// --- Then ---
		for i := 0; i < 500; i++ {
			name := FileName(
				"/DIR",
				WithPrefix(""),
				WithSuffix(""),
				WithChars(Uppercase),
			)
			want := strings.ToUpper(name)
			have := name
			if want != have {
				t.Errorf("expected only uppercase letters got %s", have)
			}
		}
	})
}

func Test_Int(t *testing.T) {
	const count = 100_000
	const maxVal = 999
	for i := 0; i < count; i++ {
		num := Int(maxVal)
		if num < 1 || num > maxVal {
			t.Errorf("expected num [1, %d]: %d", maxVal, num)
		}
	}
}

func Test_Password(t *testing.T) {
	const count = 100_000

	for i := 0; i < count; i++ {
		pass := Password(16)
		if len(pass) != 16 {
			t.Errorf("expected password 16 characters long got: %s", pass)
		}
	}
}
