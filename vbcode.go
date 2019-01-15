package blackfriday_vbcode

import (
	"bytes"
	"io"
	"log"

	"gopkg.in/russross/blackfriday.v2"
)

var headingSize = map[int]string {
	1: "+4",
	2: "+3",
	3: "+2",
	4: "+1",
	5: "+0",
	6: "-1",
	7: "-2",
}

type VBulletinRenderer struct {
	w bytes.Buffer
}

func NewVBulletinRenderer() *VBulletinRenderer {
	return &VBulletinRenderer{
		w: *bytes.NewBuffer([]byte{}),
	}
}

func (r *VBulletinRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {

}

func (r *VBulletinRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.Text:
		r.literal(node.Literal)

	case blackfriday.Softbreak:
		r.lineBreak()

	case blackfriday.Hardbreak:
		r.lineBreak()

	case blackfriday.Emph:
		r.tag("highlight", entering)

	case blackfriday.Strong:
		r.tag("b", entering)

	case blackfriday.Del:
		r.tag("strike", entering)

	case blackfriday.HTMLSpan:
		break

	case blackfriday.Link:
		r.tagWithParameter("url", string(node.LinkData.Destination), entering)

	case blackfriday.Code:
		r.startTag("code")
		r.literal(node.Literal)
		r.endTag("code")

	case blackfriday.Document:
		break

	case blackfriday.Paragraph:
		if !entering {
			r.lineBreak()
		}

	case blackfriday.BlockQuote:
		r.tag("quote", entering)

	case blackfriday.HTMLBlock:
		break

	case blackfriday.Heading:
		if entering && r.w.Len() > 0 {
			r.lineBreak()
		}
		r.tagWithParameter("size", headingSize[node.Level], entering)
		if !entering {
			r.lineBreak()
			r.lineBreak()
		}

	case blackfriday.HorizontalRule:
		r.lineBreak()
		r.lineBreak()

	case blackfriday.List:
		if entering && r.w.Len() > 0 {
			r.lineBreak()
		}
		if node.ListFlags&blackfriday.ListTypeOrdered != 0 {
			r.tagWithParameter("list", "1", entering)
		} else if node.ListFlags&blackfriday.ListTypeDefinition != 0 {
			r.tag("list", entering)
		} else {
			r.tag("list", entering)
		}
		r.lineBreak()


	case blackfriday.Item:
		if entering {
			r.startTag("*")
		}

	case blackfriday.CodeBlock:
		r.lineBreak()
		r.startTag("code")
		r.literal(node.Literal)
		r.endTag("code")
		r.lineBreak()

	case blackfriday.Image:
		r.tagWithParameter("img", string(node.LinkData.Destination), entering)

	case blackfriday.Table:
		log.Printf("Warning: Tables are not supported by vBulletin code. Skipping table.")

	case blackfriday.TableCell:
		log.Printf("Warning: Tables are not supported by vBulletin code. Skipping table-cell.")

	case blackfriday.TableHead:
		log.Printf("Warning: Tables are not supported by vBulletin code. Skipping table-head.")

	case blackfriday.TableBody:
		log.Printf("Warning: Tables are not supported by vBulletin code. Skipping table-body.")

	case blackfriday.TableRow:
		log.Printf("Warning: Tables are not supported by vBulletin code. Skipping table-row.")

	default:
		log.Printf("Warning: Unkown node type %s ", node.Type.String())
	}

	return blackfriday.GoToNext
}

func (r *VBulletinRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {
	if _, err := w.Write(r.w.Bytes()); err != nil {
		panic(err)
	}
}

func (r *VBulletinRenderer) lineBreak() {
	r.w.Write([]byte{'\n'})
}

func (r *VBulletinRenderer) literal(text []byte) {
	r.w.WriteString(string(text))
}

func (r *VBulletinRenderer) tag(tag string, entering bool) {
	if entering {
		r.startTag(tag)
	} else {
		r.endTag(tag)
	}
}

func (r *VBulletinRenderer) tagWithParameter(tag string, param string, entering bool) {
	if entering {
		r.startTagWithParameter(tag, param)
	} else {
		r.endTag(tag)
	}
}

func (r *VBulletinRenderer) startTag(tag string) {
	r.w.WriteString("[" + tag + "]")
}

func (r *VBulletinRenderer) startTagWithParameter(tag string, param string) {
	r.w.WriteString("[" + tag + "=" + param + "]")
}

func (r *VBulletinRenderer) endTag(tag string) {
	r.w.WriteString("[/" + tag + "]")
}