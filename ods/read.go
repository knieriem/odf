// This package implements rudimentary support
// for reading Open Document Spreadsheet files. At current
// stage table data can be accessed.
package ods

import (
	"encoding/xml"
	"errors"
	"github.com/knieriem/odf"
)

type Doc struct {
	XMLName xml.Name `xml:"document-content"`
	Table   []Table  `xml:"body>spreadsheet>table"`
}

type Table struct {
	Name   string   `xml:"name,attr"`
	Column []string `xml:"table-column"`
	Row    []Row    `xml:"table-row"`
}

type Row struct {
	RepeatedRows int `xml:"number-rows-repeated,attr"`

	Cell []Cell `xml:",any"` // use ",any" to match table-cell and covered-table-cell
}

// Return the contents of a row as a slice of strings. Cells that are
// covered by other cells will appear as empty strings.
func (r *Row) Strings() (row []string) {
	if len(r.Cell) == 0 {
		return
	}

	n := 0
	// calculate the real number of cells (including repeated)
	for _, c := range r.Cell {
		switch {
		case c.RepeatedCols != 0:
			n += c.RepeatedCols
		default:
			n++
		}
	}

	row = make([]string, n)
	w := 0
	for _, c := range r.Cell {
		cs := ""
		if c.XMLName.Local != "covered-table-cell" {
			cs = c.String()
		}
		row[w] = cs
		w++
		switch {
		case c.RepeatedCols != 0:
			for j := 1; j < c.RepeatedCols; j++ {
				row[w] = cs
				w++
			}
		}
	}
	return
}

type Cell struct {
	XMLName xml.Name

	// attributes
	ValueType    string `xml:"value-type,attr"`
	Value        string `xml:"value,attr"`
	Formula      string `xml:"formula,attr"`
	RepeatedCols int    `xml:"number-columns-repeated,attr"`
	ColSpan      int    `xml:"number-columns-spanned,attr"`

	P []Par `xml:"p"`
}

func (c *Cell) String() string {
	n := len(c.P)
	if n == 1 {
		return c.P[0].XML
	}
	s := ""
	for i := range c.P {
		if i != n-1 {
			s += c.P[i].XML + "\n"
		} else {
			s += c.P[i].XML
		}
	}
	return s
}

type Par struct {
	XML string `xml:",chardata"`
}

func (t *Table) Width() int {
	return len(t.Column)
}
func (t *Table) Height() int {
	return len(t.Row)
}
func (t *Table) Strings() (s [][]string) {
	if len(t.Row) == 0 {
		return
	}

	n := 0
	// calculate the real number of rows (including repeated rows)
	for _, r := range t.Row {
		switch {
		case r.RepeatedRows != 0:
			n += r.RepeatedRows
		default:
			n++
		}
	}

	s = make([][]string, n)
	w := 0
	for _, r := range t.Row {
		row := r.Strings()
		s[w] = row
		w++
		for j := 1; j < r.RepeatedRows; j++ {
			s[w] = row
			w++
		}
	}
	return
}

type File struct {
	*odf.File
}

// Open an ODS file. If the file doesn't exist or doesn't look
// like a spreadsheet file, an error is returned.
func Open(fileName string) (f *File, err error) {
	of, err := odf.Open(fileName)
	if err != nil {
		return
	}

	if of.MimeType != odf.MimeTypePfx+"spreadsheet" {
		of.Close()
		err = errors.New("not a spreadsheet")
	} else {
		f = &File{of}
	}
	return
}

// Parse the content.xml part of an ODS file. On Success
// the returned Doc will contain the data of the rows and cells
// of the table(s) contained in the ODS file.
func (f *File) ParseContent(doc *Doc) (err error) {
	content, err := f.Open("content.xml")
	if err != nil {
		return
	}
	defer content.Close()

	d := xml.NewDecoder(content)
	err = d.Decode(doc)
	return
}
