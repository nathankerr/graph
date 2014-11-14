package graph_test

import (
	"github.com/nathankerr/graph"
	"math"
	"testing"
)

func TestTileGraph(t *testing.T) {
	tg := graph.NewTileGraph(4, 4, false)
	if tg == nil || tg.String() != "▀▀▀▀\n▀▀▀▀\n▀▀▀▀\n▀▀▀▀" {
		t.Fatal("Tile graph not generated correctly")
	}

	tg.SetPassability(0, 1, true)
	if tg == nil || tg.String() != "▀ ▀▀\n▀▀▀▀\n▀▀▀▀\n▀▀▀▀" {
		t.Fatal("Passability set incorrectly")
	}

	tg.SetPassability(0, 1, false)
	if tg == nil || tg.String() != "▀▀▀▀\n▀▀▀▀\n▀▀▀▀\n▀▀▀▀" {
		t.Fatal("Passability set incorrectly")
	}

	tg.SetPassability(0, 1, true)
	if tg == nil || tg.String() != "▀ ▀▀\n▀▀▀▀\n▀▀▀▀\n▀▀▀▀" {
		t.Fatal("Passability set incorrectly")
	}

	tg.SetPassability(0, 2, true)
	if tg == nil || tg.String() != "▀  ▀\n▀▀▀▀\n▀▀▀▀\n▀▀▀▀" {
		t.Fatal("Passability set incorrectly")
	}

	tg.SetPassability(1, 2, true)
	if tg == nil || tg.String() != "▀  ▀\n▀▀ ▀\n▀▀▀▀\n▀▀▀▀" {
		t.Fatal("Passability set incorrectly")
	}

	tg.SetPassability(2, 2, true)
	if tg == nil || tg.String() != "▀  ▀\n▀▀ ▀\n▀▀ ▀\n▀▀▀▀" {
		t.Fatal("Passability set incorrectly")
	}

	tg.SetPassability(3, 2, true)
	if tg == nil || tg.String() != "▀  ▀\n▀▀ ▀\n▀▀ ▀\n▀▀ ▀" {
		t.Fatal("Passability set incorrectly")
	}

	if tg2, err := graph.GenerateTileGraph("▀  ▀\n▀▀ ▀\n▀▀ ▀\n▀▀ ▀"); err != nil {
		t.Error("Tile graph errored on interpreting valid template string\n▀  ▀\n▀▀ ▀\n▀▀ ▀\n▀▀ ▀")
	} else if tg2.String() != "▀  ▀\n▀▀ ▀\n▀▀ ▀\n▀▀ ▀" {
		t.Error("Tile graph failed to generate properly with input string\n▀  ▀\n▀▀ ▀\n▀▀ ▀\n▀▀ ▀")
	}

	if tg.CoordsToID(0, 0) != 0 {
		t.Error("Coords to ID fails on 0,0")
	} else if tg.CoordsToID(3, 3) != 15 {
		t.Error("Coords to ID fails on 3,3")
	} else if tg.CoordsToID(0, 3) != 3 {
		t.Error("Coords to ID fails on 0,3")
	} else if tg.CoordsToID(3, 0) != 12 {
		t.Error("Coords to ID fails on 3,0")
	}

	if r, c := tg.IDToCoords(0); r != 0 || c != 0 {
		t.Error("ID to Coords fails on 0,0")
	} else if r, c := tg.IDToCoords(15); r != 3 || c != 3 {
		t.Error("ID to Coords fails on 3,3")
	} else if r, c := tg.IDToCoords(3); r != 0 || c != 3 {
		t.Error("ID to Coords fails on 0,3")
	} else if r, c := tg.IDToCoords(12); r != 3 || c != 0 {
		t.Error("ID to Coords fails on 3,0")
	}

	if succ := tg.Successors(graph.GonumNode(0)); succ != nil || len(succ) != 0 {
		t.Error("Successors for impassable tile not 0")
	}

	if succ := tg.Successors(graph.GonumNode(2)); succ == nil || len(succ) != 2 {
		t.Error("Incorrect number of successors for (0,2)")
	} else {
		for _, s := range succ {
			if s.ID() != 1 && s.ID() != 6 {
				t.Error("Successor for (0,2) neither (0,1) nor (1,2)")
			}
		}
	}

	if tg.Degree(graph.GonumNode(2)) != 4 {
		t.Error("Degree returns incorrect number for (0,2)")
	}
	if tg.Degree(graph.GonumNode(1)) != 2 {
		t.Error("Degree returns incorrect number for (0,2)")
	}
	if tg.Degree(graph.GonumNode(0)) != 0 {
		t.Error("Degree returns incorrect number for impassable tile (0,0)")
	}

}

func TestSimpleAStar(t *testing.T) {
	tg, err := graph.GenerateTileGraph("▀  ▀\n▀▀ ▀\n▀▀ ▀\n▀▀ ▀")
	if err != nil {
		t.Fatal("Couldn't generate tilegraph")
	}

	path, cost, _ := graph.AStar(graph.GonumNode(1), graph.GonumNode(14), tg, nil, nil)
	if math.Abs(cost-4.0) > .00001 {
		t.Errorf("A* reports incorrect cost for simple tilegraph search")
	}

	if path == nil {
		t.Fatalf("A* fails to find path for for simple tilegraph search")
	} else {
		correctPath := []int{1, 2, 6, 10, 14}
		if len(path) != len(correctPath) {
			t.Fatalf("Astar returns wrong length path for simple tilegraph search")
		}
		for i, node := range path {
			if node.ID() != correctPath[i] {
				t.Errorf("Astar returns wrong path at step", i, "got:", node, "actual:", correctPath[i])
			}
		}
	}
}

func TestHarderAStar(t *testing.T) {
	tg := graph.NewTileGraph(3, 3, true)

	path, cost, _ := graph.AStar(graph.GonumNode(0), graph.GonumNode(8), tg, nil, nil)

	if math.Abs(cost-4.0) > .00001 || !graph.IsPath(path, tg) {
		t.Error("Non-optimal or impossible path found for 3x3 grid")
	}

	tg = graph.NewTileGraph(1000, 1000, true)
	path, cost, _ = graph.AStar(graph.GonumNode(00), graph.GonumNode(999*1000+999), tg, nil, nil)
	if !graph.IsPath(path, tg) || cost != 1998.0 {
		t.Error("Non-optimal or impossible path found for 100x100 grid; cost:", cost, "path:\n"+tg.PathString(path))
	}
}
