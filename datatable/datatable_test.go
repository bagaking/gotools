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

func ExampleCreateAndRender() {
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
