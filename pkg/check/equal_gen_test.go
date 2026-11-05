// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"math/rand/v2"
	"reflect"
	"slices"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/dump"
)

// genSeeds is the number of generated values each generated test checks.
const genSeeds = 100

// genDepth limits how many levels of new nodes the generator creates.
const genDepth = 2

// genDumpDepth keeps failure dumps of generated cyclic values small.
const genDumpDepth = 2

// genNode is the type of generated values. Its fields hold the shapes that
// hide false passes: slices sharing one array, zero-size elements, nil and
// empty values, unexported fields, and cycles.
type genNode struct {
	Ints   []int
	Alias  []int
	Zero   []struct{}
	Map    map[string]int
	Ptr    *genNode
	Kids   []*genNode
	Any    any
	hidden []int
	hidMap map[string]int
	hidAny any
}

// genSite is one way to change a generated value.
type genSite struct {
	name  string
	apply func()
}

// genBuilder builds one generated value. The same seed builds the same shape.
type genBuilder struct {
	rng   *rand.Rand
	root  *genNode
	nodes []*genNode
	path  []*genNode
	pool  [][]int
	sites []genSite
	seen  map[string]bool
}

// genBuild returns a builder holding the value generated for the seed.
func genBuild(seed uint64) *genBuilder {
	bld := &genBuilder{
		rng:  rand.New(rand.NewPCG(seed, 0)), //nolint:gosec
		seen: make(map[string]bool),
	}
	bld.root = bld.node(genDepth)
	return bld
}

// node builds a node, adding new child nodes while depth is above zero.
func (bld *genBuilder) node(depth int) *genNode {
	nod := new(genNode)
	bld.nodes = append(bld.nodes, nod)
	bld.path = append(bld.path, nod)
	defer func() { bld.path = bld.path[:len(bld.path)-1] }()

	nod.Ints = bld.ints()
	nod.Alias = bld.alias()
	nod.Zero = bld.zeros()
	nod.Map = bld.strMap()
	nod.hidden = bld.ints()
	nod.hidMap = bld.strMap()
	nod.Any = bld.anyVal()
	nod.hidAny = bld.anyVal()
	nod.Ptr = bld.ptr(depth)
	for range bld.rng.IntN(3) {
		nod.Kids = append(nod.Kids, bld.ptr(depth))
	}

	bld.intSites("Ints", &nod.Ints)
	bld.intSites("Alias", &nod.Alias)
	bld.intSites("hidden", &nod.hidden)
	genSliceSites(bld, "Zero", &nod.Zero, struct{}{})
	genSliceSites(bld, "Kids", &nod.Kids, nil)
	bld.mapSites("Map", &nod.Map)
	bld.mapSites("hidMap", &nod.hidMap)
	bld.ptrSites(&nod.Ptr)
	bld.anySites("Any", &nod.Any)
	bld.anySites("hidAny", &nod.hidAny)
	return nod
}

// ints returns a nil, empty, or short slice. A non-empty one joins the alias
// pool.
func (bld *genBuilder) ints() []int {
	var ints []int
	switch num := bld.rng.IntN(5); num {
	case 0:
		return nil

	case 1:
		return make([]int, 0)

	default:
		for range num - 1 {
			ints = append(ints, bld.rng.IntN(3))
		}
	}
	bld.pool = append(bld.pool, ints)
	return ints
}

// alias returns a part of a slice from the pool, or a new slice.
func (bld *genBuilder) alias() []int {
	if len(bld.pool) == 0 || bld.rng.IntN(4) == 0 {
		return bld.ints()
	}
	src := bld.pool[bld.rng.IntN(len(bld.pool))]
	beg := bld.rng.IntN(len(src))
	end := beg + bld.rng.IntN(len(src)-beg+1)
	if beg == 0 && end > 0 && end < len(src) {
		bld.seen["shorter alias"] = true
	}
	return src[beg:end]
}

// zeros returns a nil, empty, or short slice of zero-size elements.
func (bld *genBuilder) zeros() []struct{} {
	switch num := bld.rng.IntN(4); num {
	case 0:
		return nil

	default:
		return make([]struct{}, num-1)
	}
}

// strMap returns a nil, empty, or short map.
func (bld *genBuilder) strMap() map[string]int {
	switch bld.rng.IntN(4) {
	case 0:
		return nil

	case 1:
		return make(map[string]int)

	case 2:
		return map[string]int{"a": bld.rng.IntN(3)}

	default:
		return map[string]int{"a": bld.rng.IntN(3), "b": bld.rng.IntN(3)}
	}
}

// ptr returns nil, a new node, or a node built before, which can form a cycle.
func (bld *genBuilder) ptr(depth int) *genNode {
	switch bld.rng.IntN(3) {
	case 0:
		return nil

	case 1:
		if depth > 0 {
			return bld.node(depth - 1)
		}
		return nil

	default:
		return bld.known()
	}
}

// known returns a node built before and records what it forms: a cycle when
// the node is still being built, and sharing when it is finished.
func (bld *genBuilder) known() *genNode {
	nod := bld.nodes[bld.rng.IntN(len(bld.nodes))]
	if slices.Contains(bld.path, nod) {
		bld.seen["cycle"] = true
		return nod
	}
	bld.seen["shared node"] = true
	return nod
}

// anyVal returns a value for an interface field.
func (bld *genBuilder) anyVal() any {
	switch bld.rng.IntN(7) {
	case 0:
		return nil

	case 1:
		return bld.rng.IntN(3)

	case 2:
		return bld.alias()

	case 3:
		return bld.strMap()

	case 4:
		return bld.zeros()

	case 5:
		return bld.known()

	default:
		// A slice holding a shorter view of itself.
		cyc := make([]any, 2)
		cyc[0] = cyc[:1]
		cyc[1] = bld.rng.IntN(3)
		bld.seen["any-cycle"] = true
		bld.add("any-cycle element", func() { cyc[1] = "changed" })
		return cyc
	}
}

// add records one way to change the value.
func (bld *genBuilder) add(name string, apply func()) {
	bld.sites = append(bld.sites, genSite{name, apply})
}

// intSites records ways to change an int slice.
func (bld *genBuilder) intSites(name string, ptr *[]int) {
	genSliceSites(bld, name, ptr, 7)
	for i := range *ptr {
		bld.add(name+" element", func() { (*ptr)[i] += 10 })
	}
}

// genSliceSites records ways to change a slice.
func genSliceSites[T any](bld *genBuilder, name string, ptr *[]T, one T) {
	switch {
	case *ptr == nil:
		bld.add(name+" nil to empty", func() { *ptr = make([]T, 0) })
		bld.add(name+" nil to one", func() { *ptr = []T{one} })

	case len(*ptr) == 0:
		bld.add(name+" empty to nil", func() { *ptr = nil })
		bld.add(name+" empty to one", func() { *ptr = []T{one} })

	default:
		bld.add(name+" to nil", func() { *ptr = nil })
		bld.add(name+" shorter", func() { *ptr = (*ptr)[:len(*ptr)-1] })
		bld.add(name+" longer", func() {
			*ptr = append(slices.Clip(*ptr), one)
		})
	}
}

// mapSites records ways to change a map.
func (bld *genBuilder) mapSites(name string, ptr *map[string]int) {
	switch {
	case *ptr == nil:
		bld.add(name+" nil to empty", func() {
			*ptr = make(map[string]int)
		})

	case len(*ptr) == 0:
		bld.add(name+" empty to nil", func() { *ptr = nil })
		bld.add(name+" empty to one", func() { (*ptr)["a"] = 1 })

	default:
		bld.add(name+" to nil", func() { *ptr = nil })
		bld.add(name+" new key", func() { (*ptr)["z"] = 1 })
		bld.add(name+" other value", func() { (*ptr)["a"] += 10 })
		bld.add(name+" other key", func() {
			(*ptr)["z"] = (*ptr)["a"]
			delete(*ptr, "a")
		})
	}
}

// ptrSites records ways to change a pointer.
func (bld *genBuilder) ptrSites(ptr **genNode) {
	if *ptr == nil {
		bld.add("Ptr nil to node", func() { *ptr = new(genNode) })
		return
	}
	bld.add("Ptr to nil", func() { *ptr = nil })
}

// anySites records ways to change an interface field.
func (bld *genBuilder) anySites(name string, ptr *any) {
	bld.add(name+" other type", func() { *ptr = "changed" })
	if *ptr != nil {
		bld.add(name+" to nil", func() { *ptr = nil })
	}
}

// unshare replaces about half of the Ptr and Kids references with shallow
// copies. The value stays equal, but those nodes are no longer shared.
func (bld *genBuilder) unshare(rng *rand.Rand) {
	for _, nod := range bld.nodes {
		nod.Ptr = genCopy(rng, nod.Ptr)
		for i, kid := range nod.Kids {
			nod.Kids[i] = genCopy(rng, kid)
		}
	}
}

// genCopy returns a shallow copy of a non-nil node for about half of the calls
// and the node itself otherwise.
func genCopy(rng *rand.Rand, nod *genNode) *genNode {
	if nod == nil || rng.IntN(2) == 0 {
		return nod
	}
	dup := *nod
	return &dup
}

func Test_genBuild(t *testing.T) {
	t.Run("same seed builds equal values", func(t *testing.T) {
		// --- When ---
		have := genBuild(42)

		// --- Then ---
		want := genBuild(42)
		affirm.Equal(t, true, reflect.DeepEqual(want.root, have.root))
		affirm.Equal(t, len(want.sites), len(have.sites))
	})

	t.Run("seeds reach the listed shapes", func(t *testing.T) {
		// --- Given ---
		seen := make(map[string]bool)

		// --- When ---
		for seed := range uint64(genSeeds) {
			bld := genBuild(seed)
			for shape := range bld.seen {
				seen[shape] = true
			}
			for _, site := range bld.sites {
				seen[site.name] = true
			}
		}

		// --- Then ---
		wants := []string{
			"shorter alias",
			"cycle",
			"shared node",
			"any-cycle",
			"any-cycle element",
			"Ints nil to empty",
			"Ints empty to nil",
			"Ints element",
			"Alias shorter",
			"hidden nil to empty",
			"hidden element",
			"Zero nil to empty",
			"Zero shorter",
			"Zero longer",
			"Map nil to empty",
			"Map empty to nil",
			"Map other key",
			"hidMap nil to empty",
			"hidMap other key",
			"hidAny to nil",
			"Kids shorter",
			"Ptr to nil",
		}
		for _, want := range wants {
			if !seen[want] {
				t.Errorf("shape not reached: %s", want)
			}
		}
	})
}

func Test_genBuilder_unshare(t *testing.T) {
	// --- Given ---
	var ptrCopies, kidCopies int
	for seed := range uint64(genSeeds) {
		want := genBuild(seed).root
		have := genBuild(seed)
		ptrs := make([]*genNode, 0, len(have.nodes))
		kids := make([][]*genNode, 0, len(have.nodes))
		for _, nod := range have.nodes {
			ptrs = append(ptrs, nod.Ptr)
			kids = append(kids, slices.Clone(nod.Kids))
		}

		// --- When ---
		have.unshare(rand.New(rand.NewPCG(seed, 1))) //nolint:gosec

		// --- Then ---
		if !reflect.DeepEqual(want, have.root) {
			t.Fatalf("seed %d: unshared value is not equal", seed)
		}
		for i, nod := range have.nodes {
			if nod.Ptr != ptrs[i] {
				ptrCopies++
			}
			for j, kid := range nod.Kids {
				if kid != kids[i][j] {
					kidCopies++
				}
			}
		}
	}
	affirm.Equal(t, true, ptrCopies > 0)
	affirm.Equal(t, true, kidCopies > 0)
}

func Test_genBuilder_known(t *testing.T) {
	t.Run("node being built forms a cycle", func(t *testing.T) {
		// --- Given ---
		nod := new(genNode)
		bld := genBuild(1)
		bld.seen = make(map[string]bool)
		bld.nodes = []*genNode{nod}
		bld.path = []*genNode{nod}

		// --- When ---
		have := bld.known()

		// --- Then ---
		affirm.Equal(t, true, have == nod)
		affirm.DeepEqual(t, map[string]bool{"cycle": true}, bld.seen)
	})

	t.Run("finished node is shared", func(t *testing.T) {
		// --- Given ---
		nod := new(genNode)
		bld := genBuild(1)
		bld.seen = make(map[string]bool)
		bld.nodes = []*genNode{nod}
		bld.path = nil

		// --- When ---
		have := bld.known()

		// --- Then ---
		affirm.Equal(t, true, have == nod)
		affirm.DeepEqual(t, map[string]bool{"shared node": true}, bld.seen)
	})
}

func Test_genCopy(t *testing.T) {
	t.Run("nil stays nil", func(t *testing.T) {
		for seed := range uint64(genSeeds) {
			// --- Given ---
			rng := rand.New(rand.NewPCG(seed, 1)) //nolint:gosec

			// --- When ---
			have := genCopy(rng, nil)

			// --- Then ---
			affirm.Equal(t, true, have == nil)
		}
	})

	t.Run("copy has the same content", func(t *testing.T) {
		// --- Given ---
		want := &genNode{Ints: []int{1}}

		// --- When ---
		have := want
		for seed := range uint64(genSeeds) {
			if have == want {
				rng := rand.New(rand.NewPCG(seed, 1)) //nolint:gosec
				have = genCopy(rng, want)
			}
		}

		// --- Then ---
		affirm.Equal(t, true, have != want)
		affirm.Equal(t, true, reflect.DeepEqual(want, have))
	})
}

func Test_Equal_generated_copies(t *testing.T) {
	opt := WithDumper(dump.WithMaxDepth(genDumpDepth))
	for seed := range uint64(genSeeds) {
		// --- Given ---
		want := genBuild(seed).root
		have := genBuild(seed).root

		// --- When ---
		err := Equal(want, have, opt)

		// --- Then ---
		if err != nil {
			t.Fatalf("seed %d: copies are not equal: %v", seed, err)
		}
	}
}

func Test_Equal_generated_single_change(t *testing.T) {
	opt := WithDumper(dump.WithMaxDepth(genDumpDepth))
	for seed := range uint64(genSeeds) {
		for idx := range genBuild(seed).sites {
			// --- Given ---
			want := genBuild(seed).root
			bld := genBuild(seed)
			bld.sites[idx].apply()
			name := bld.sites[idx].name

			// --- When ---
			err := Equal(want, bld.root, opt)
			rev := Equal(bld.root, want, opt)

			// --- Then ---
			if reflect.DeepEqual(want, bld.root) {
				t.Fatalf("seed %d: %s: no real change", seed, name)
			}
			if err == nil {
				t.Fatalf("seed %d: %s: change not found", seed, name)
			}
			if rev == nil {
				t.Fatalf("seed %d: %s: reverse not found", seed, name)
			}
		}
	}
}

func Test_Equal_generated_differential(t *testing.T) {
	opt := WithDumper(dump.WithMaxDepth(genDumpDepth))
	verdicts := make(map[bool]int)
	for seed := range uint64(genSeeds) {
		// --- Given ---
		rng := rand.New(rand.NewPCG(seed, 1)) //nolint:gosec
		want := genBuild(seed).root
		bld := genBuild(seed)
		bld.unshare(rng)
		if rng.IntN(2) == 0 {
			bld.sites[rng.IntN(len(bld.sites))].apply()
		}

		// --- When ---
		err := Equal(want, bld.root, opt)
		rev := Equal(bld.root, want, opt)

		// --- Then ---
		same := reflect.DeepEqual(want, bld.root)
		verdicts[same]++
		if same != (err == nil) {
			t.Fatalf("seed %d: differs from DeepEqual: %v", seed, err)
		}
		if same != (rev == nil) {
			t.Fatalf("seed %d: reverse differs from DeepEqual: %v", seed, rev)
		}
	}
	affirm.Equal(t, true, verdicts[true] > 0)
	affirm.Equal(t, true, verdicts[false] > 0)
}
