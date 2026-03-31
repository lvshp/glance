package reader

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEpubReaderLoadReadsSpineOrder(t *testing.T) {
	tempDir := t.TempDir()
	bookPath := filepath.Join(tempDir, "book.epub")

	files := map[string]string{
		"META-INF/container.xml": `<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>`,
		"OEBPS/content.opf": `<?xml version="1.0" encoding="UTF-8"?>
<package version="3.0" xmlns="http://www.idpf.org/2007/opf">
  <manifest>
    <item id="chapter-1" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
    <item id="chapter-2" href="chapter2.xhtml" media-type="application/xhtml+xml"/>
  </manifest>
  <spine>
    <itemref idref="chapter-1"/>
    <itemref idref="chapter-2"/>
  </spine>
</package>`,
		"OEBPS/chapter1.xhtml": `<?xml version="1.0" encoding="UTF-8"?>
<html xmlns="http://www.w3.org/1999/xhtml"><body><h1>Chapter 1</h1><p>Hello world.</p></body></html>`,
		"OEBPS/chapter2.xhtml": `<?xml version="1.0" encoding="UTF-8"?>
<html xmlns="http://www.w3.org/1999/xhtml"><body><h1>Chapter 2</h1><p>Goodbye world.</p></body></html>`,
	}

	if err := writeTestZip(bookPath, files); err != nil {
		t.Fatalf("write epub: %v", err)
	}

	r := NewEpubReader()
	if err := r.Load(bookPath); err != nil {
		t.Fatalf("load epub: %v", err)
	}

	if got := r.First(); got != "Chapter 1" {
		t.Fatalf("first line = %q, want %q", got, "Chapter 1")
	}

	if got := r.Next(); got != "Hello world." {
		t.Fatalf("second line = %q, want %q", got, "Hello world.")
	}

	r.Last()
	if got := r.Current(); got != "Goodbye world." {
		t.Fatalf("last line = %q, want %q", got, "Goodbye world.")
	}
}

func TestReadChapterTextStripsMarkup(t *testing.T) {
	tempDir := t.TempDir()
	bookPath := filepath.Join(tempDir, "chapter.epub")

	files := map[string]string{
		"chapter.xhtml": `<?xml version="1.0" encoding="UTF-8"?>
<html xmlns="http://www.w3.org/1999/xhtml"><body><p>Hello &amp; welcome</p><div>Keep <span>reading</span>.</div></body></html>`,
	}

	if err := writeTestZip(bookPath, files); err != nil {
		t.Fatalf("write zip: %v", err)
	}

	zr, err := zip.OpenReader(bookPath)
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	defer zr.Close()

	got, err := readChapterText(zr.File[0])
	if err != nil {
		t.Fatalf("read chapter text: %v", err)
	}

	want := strings.Join([]string{"Hello & welcome", "Keep reading."}, "\n")
	if got != want {
		t.Fatalf("text = %q, want %q", got, want)
	}
}

func writeTestZip(zipPath string, files map[string]string) error {
	f, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil {
			zw.Close()
			return err
		}

		if _, err := w.Write([]byte(content)); err != nil {
			zw.Close()
			return err
		}
	}

	return zw.Close()
}
