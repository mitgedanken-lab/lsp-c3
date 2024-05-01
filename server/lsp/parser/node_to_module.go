package parser

import (
	"github.com/pherrymason/c3-lsp/lsp/document"
	sitter "github.com/smacker/go-tree-sitter"
)

/*
	module: $ => seq(
	'module',
	field('path', $.path_ident),
	optional(alias($.generic_module_parameters, $.generic_parameters)),
	optional($.attributes),
	';'

	attributes:
		@private

),
*/
func (p *Parser) nodeToModule(doc *document.Document, node *sitter.Node, sourceCode []byte) string {

	name := node.ChildByFieldName("path").Content(sourceCode)
	/*
		for i := 0; i < int(node.ChildCount()); i++ {
			n := node.Child(i)
			switch n.Type() {
			case "alias":

			}
		}*/

	return name
}

/*
		import_declaration: $ => seq(
	      'import',
	      field('path', commaSep1($.path_ident)),
	      optional($.attributes),
	      ';'
	    ),
*/
func (p *Parser) nodeToImport(doc *document.Document, node *sitter.Node, sourceCode []byte) []string {
	imports := []string{}

	for i := 0; i < int(node.ChildCount()); i++ {
		n := node.Child(i)

		switch n.Type() {
		case "path_ident":
			temp_mod := ""
			for m := 0; m < int(n.ChildCount()); m++ {
				sn := n.Child(m)
				if sn.Type() == "ident" || sn.Type() == "module_resolution" {
					temp_mod += sn.Content(sourceCode)
				}
			}
			imports = append(imports, temp_mod)
		}
	}

	return imports
}