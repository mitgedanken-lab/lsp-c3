package language

import (
	"fmt"
	"testing"

	"github.com/pherrymason/c3-lsp/lsp/document"
	"github.com/stretchr/testify/assert"
	"github.com/tliron/commonlog"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestLanguage_BuildCompletionList(t *testing.T) {
	parser := createParser()
	language := NewLanguage(commonlog.MockLogger{})

	//documents := installDocuments(&language, &parser)

	t.Run("Should suggest variable names defined in module", func(t *testing.T) {
		source := `
		int variable = 3;
		int xanadu = 10;`
		expectedKind := protocol.CompletionItemKindVariable
		cases := []struct {
			input    string
			expected protocol.CompletionItem
		}{
			{"v", protocol.CompletionItem{Label: "variable", Kind: &expectedKind}},
			{"va", protocol.CompletionItem{Label: "variable", Kind: &expectedKind}},
			{"x", protocol.CompletionItem{Label: "xanadu", Kind: &expectedKind}},
		}

		for n, tt := range cases {
			t.Run(fmt.Sprintf("Case #%d", n), func(t *testing.T) {

				doc := document.NewDocument("test.c3", source+"\n"+tt.input)
				language.RefreshDocumentIdentifiers(&doc, &parser)
				position := buildPosition(4, 1) // Cursor after `v|`

				completionList := language.BuildCompletionList(&doc, position)

				assert.Equal(t, 1, len(completionList))
				assert.Equal(t, tt.expected, completionList[0])
			})
		}
	})

	t.Run("Should suggest variable names defined in module and inside current function", func(t *testing.T) {
		source := `
		int variable = 3;
		fn void main() {
			int value = 4;`
		expectedKind := protocol.CompletionItemKindVariable
		cases := []struct {
			input    string
			expected []protocol.CompletionItem
		}{
			{"v", []protocol.CompletionItem{
				{Label: "variable", Kind: &expectedKind},
				{Label: "value", Kind: &expectedKind},
			}},
			{"val", []protocol.CompletionItem{
				{Label: "value", Kind: &expectedKind},
			}},
		}

		for n, tt := range cases {
			t.Run(fmt.Sprintf("Case #%d", n), func(t *testing.T) {

				doc := document.NewDocument("test.c3", source+`
`+tt.input+`
				}`)
				language.RefreshDocumentIdentifiers(&doc, &parser)
				position := buildPosition(5, uint(len(tt.input))) // Cursor after `<input>|`

				completionList := language.BuildCompletionList(&doc, position)

				assert.Equal(t, len(tt.expected), len(completionList))
				assert.Equal(t, tt.expected, completionList)
			})
		}
	})

	t.Run("Should suggest struct members", func(t *testing.T) {
		source := `
		struct Square { int width; int height; }
		fn void main() {
			Square inst;
			inst`
		expectedKind := protocol.CompletionItemKindField
		cases := []struct {
			input    string
			expected []protocol.CompletionItem
		}{
			{".", []protocol.CompletionItem{
				{Label: "width", Kind: &expectedKind},
				{Label: "height", Kind: &expectedKind},
			}},
			{".w", []protocol.CompletionItem{
				{Label: "width", Kind: &expectedKind},
			}},
		}
		t.Skip()
		for n, tt := range cases {
			t.Run(fmt.Sprintf("Case #%d", n), func(t *testing.T) {

				doc := document.NewDocument("test.c3", source+tt.input+`}`)
				language.RefreshDocumentIdentifiers(&doc, &parser)
				position := buildPosition(5, 7+uint(len(tt.input))) // Cursor after `<input>|`

				completionList := language.BuildCompletionList(&doc, position)

				assert.Equal(t, len(tt.expected), len(completionList))
				assert.Equal(t, tt.expected, completionList)
			})
		}
	})
}