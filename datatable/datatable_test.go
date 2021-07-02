package datatable

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	var (
		err   error
		table = New(nil)
	)

	err = table.Set(1, 1, Plain("asd"))
	assert.Error(t, err, "cannot set out of range")
	err = table.Set(0, 1, Plain("asd"))
	assert.Error(t, err, "cannot set out of range")
	err = table.Set(1, 0, Plain("asd"))
	assert.Error(t, err, "cannot set out of range")
}

func TestSetRowReplacesExistingCellsWithoutAppending(t *testing.T) {
	titleLine := TitleLine{NewTitle("A"), NewTitle("B")}
	table := New(titleLine)

	if err := table.AppendRow(Line{Plain("a1"), Plain("b1")}); err != nil {
		t.Fatalf("AppendRow() error = %v", err)
	}
	if err := table.SetRow(0, Line{Plain("a2")}); err != nil {
		t.Fatalf("SetRow() error = %v", err)
	}

	row := table.GetLine(0)
	if len(row) != 2 {
		t.Fatalf("SetRow() row length = %d, want 2", len(row))
	}
	if row[0].String() != "a2" {
		t.Fatalf("SetRow() first cell = %q, want %q", row[0].String(), "a2")
	}
	if row[1].String() != "b1" {
		t.Fatalf("SetRow() second cell = %q, want %q", row[1].String(), "b1")
	}
}

func TestLineGetReturnsEmptyAtLenBoundary(t *testing.T) {
	line := Line{Plain("a")}

	got := line.Get(line.Len())

	if got != Empty {
		t.Fatalf("Get(Len()) = %q, want Empty", got.String())
	}
}

func TestLineGetsReturnsOnlyRequestedColumns(t *testing.T) {
	line := Line{Plain("a"), Plain("b")}

	got := line.Gets(1, 2, -1)

	if len(got) != 3 {
		t.Fatalf("Gets() length = %d, want 3", len(got))
	}
	if got[0].String() != "b" {
		t.Fatalf("Gets()[0] = %q, want %q", got[0].String(), "b")
	}
	if got[1] != Empty {
		t.Fatalf("Gets()[1] = %q, want Empty", got[1].String())
	}
	if got[2] != Empty {
		t.Fatalf("Gets()[2] = %q, want Empty", got[2].String())
	}
}

func Example() {
	const (
		TagID          = "ID"
		TagDescription = "Description"
		TagSetting     = "Setting"
		TagAnnotation  = "Annotation"
	)

	var (
		TitleID   = NewTitle("The ID", TagID)
		TitleDesc = NewTitle("Desc", TagDescription)
		TitleA    = NewTitle("A", TagSetting)
		TitleB    = NewTitle("B", TagSetting)
		TitleC    = NewTitle("C", TagSetting)
		TitleD    = NewTitle("D", TagSetting, TagAnnotation)
		TitleE    = NewTitle("E", TagSetting, TagAnnotation)
	)

	titleLst := TitleLine{
		TitleID,
		TitleDesc.DerivativeAddOn("_XX"), TitleDesc.DerivativeAddOn("_YY"),
		TitleA, TitleB, TitleC, TitleD, TitleE,
	}

	repo := New(titleLst)
	_ = repo.AppendRow(Line{
		Plain("id_1"), Plain("xx1"), Plain("yy1"), Plain("a1"), Plain("b1"), Plain("c1"), Plain("d1"), Plain("e1"),
	})
	_ = repo.AppendRow(Line{
		Plain("id_2"), Plain("xx2"), Plain("yy2"), Plain("a2"), Plain("b2"), Plain("c2"), Plain("d2"), Plain("e2"),
	})

	exportTitles := titleLst.SubTitleLineByTags(TagID, TagDescription, TagAnnotation)
	data := repo.Render(exportTitles, true)

	fmt.Println(data)
	// Output: [[The ID Desc_XX Desc_YY D E] [id_1 xx1 yy1 d1 e1] [id_2 xx2 yy2 d2 e2]]
}
