package pdf

import (
	"bytes"
	"log"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func NewPdf() {
	// pdf := gofpdf.New("P", "mm", "A4", "")
	// pdf.AddPage()
	// pdf.SetFont("Arial", "B", 16)
	// pdf.Cell(40, 10, "Hello, world")

	// pdf.Write(80, "ggggggg")
	// err := pdf.OutputFileAndClose("hello.pdf")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}
	// pdfg.AddPage(wkhtmltopdf.NewPage("https://www.baidu.com"))

	reader := bytes.NewReader([]byte("aaaaaaa"))
	pdfg.AddPage(wkhtmltopdf.NewPageReader(reader))
	pdfg.WriteFile("a.pdf")
}
