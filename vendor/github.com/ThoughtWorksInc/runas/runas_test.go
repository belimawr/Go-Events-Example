package runas

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

const linhaLetraA = "0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;"

const linhas3Da43 = `
003D;EQUALS SIGN;Sm;0;ON;;;;;N;;;;;
003E;GREATER-THAN SIGN;Sm;0;ON;;;;;Y;;;;;
003F;QUESTION MARK;Po;0;ON;;;;;N;;;;;
0040;COMMERCIAL AT;Po;0;ON;;;;;N;;;;;
0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;
0042;LATIN CAPITAL LETTER B;Lu;0;L;;;;;N;;;;0062;
0043;LATIN CAPITAL LETTER C;Lu;0;L;;;;;N;;;;0063;
`

func TestAnalisarLinha(t *testing.T) {
	runa, nome, palavras := AnalisarLinha(linhaLetraA) // ➊
	if runa != 'A' {
		t.Errorf("Esperado: 'A'; recebido: %q", runa)
	}
	const nomeA = "LATIN CAPITAL LETTER A"
	if nome != nomeA {
		t.Errorf("Esperado: %q; recebido: %q", nomeA, nome)
	}
	palavrasA := []string{"LATIN", "CAPITAL", "LETTER", "A"} // ➋
	if !reflect.DeepEqual(palavras, palavrasA) {             // ➌
		t.Errorf("\n\tEsperado: %q\n\trecebido: %q", palavrasA, palavras) // ➍
	}
}

func TestAnalisarLinhaComHífenECampo10(t *testing.T) {
	var casos = []struct { // ➊
		linha    string
		runa     rune
		nome     string
		palavras []string
	}{ // ➋
		{"0021;EXCLAMATION MARK;Po;0;ON;;;;;N;;;;;",
			'!', "EXCLAMATION MARK", []string{"EXCLAMATION", "MARK"}},
		{"002D;HYPHEN-MINUS;Pd;0;ES;;;;;N;;;;;",
			'-', "HYPHEN-MINUS", []string{"HYPHEN", "MINUS"}},
		{"0027;APOSTROPHE;Po;0;ON;;;;;N;APOSTROPHE-QUOTE;;;",
			'\'', "APOSTROPHE (APOSTROPHE-QUOTE)", []string{"APOSTROPHE", "QUOTE"}},
	}
	for _, caso := range casos { // ➌
		runa, nome, palavras := AnalisarLinha(caso.linha) // ➍
		if runa != caso.runa || nome != caso.nome ||
			!reflect.DeepEqual(palavras, caso.palavras) {
			t.Errorf("\nAnalisarLinha(%q)\n-> (%q, %q, %q)", // ➎
				caso.linha, runa, nome, palavras)
		}
	}
}

func TestContém(t *testing.T) {
	casos := []struct { // ➊
		fatia     []string
		procurado string
		esperado  bool
	}{ // ➋
		{[]string{"A", "B"}, "B", true},
		{[]string{}, "A", false},
		{[]string{"A", "B"}, "Z", false}, // ➌
	} // ➍
	for _, caso := range casos { // ➎
		recebido := contém(caso.fatia, caso.procurado) // ➏
		if recebido != caso.esperado {                 // ➐
			t.Errorf("contém(%#v, %#v) esperado: %v; recebido: %v",
				caso.fatia, caso.procurado, caso.esperado, recebido) // ➑
		}
	}
}

func TestContémTodos(t *testing.T) {
	casos := []struct { // ➊
		fatia      []string
		procurados []string
		esperado   bool
	}{ // ➋
		{[]string{"A", "B"}, []string{"B"}, true},
		{[]string{}, []string{"A"}, false},
		{[]string{"A"}, []string{}, true}, // ➌
		{[]string{"A", "B"}, []string{"Z"}, false},
		{[]string{"A", "B", "C"}, []string{"A", "C"}, true},
		{[]string{"A", "B", "C"}, []string{"A", "Z"}, false},
		{[]string{"A", "B"}, []string{"A", "B", "C"}, false},
	}
	for _, caso := range casos {
		obtido := contémTodos(caso.fatia, caso.procurados) // ➍
		if obtido != caso.esperado {
			t.Errorf("contémTodos(%#v, %#v)\nesperado: %v; recebido: %v",
				caso.fatia, caso.procurados, caso.esperado, obtido) // ➎
		}
	}
}

func TestSeparar(t *testing.T) {
	casos := []struct {
		texto    string
		esperado []string
	}{
		{"A", []string{"A"}},
		{"A B", []string{"A", "B"}},
		{"A B-C", []string{"A", "B", "C"}},
	}
	for _, caso := range casos {
		obtido := separar(caso.texto)
		if !reflect.DeepEqual(obtido, caso.esperado) {
			t.Errorf("separar(%q)\nesperado: %#v; recebido: %#v",
				caso.texto, caso.esperado, obtido)
		}
	}
}

func TestListar(t *testing.T) {
	texto := strings.NewReader(linhas3Da43)
	m := Listar(texto, "MARK")

	if len(m) != 1 {
		t.Error("Esperando 1 runa")
	}

	for k := range m {
		if k != "U+003F" {
			t.Errorf("Não esperava a runa %q", k)
		}
	}
	// Output: U+003F	?	QUESTION MARK
}

func TestListar_doisResultados(t *testing.T) {
	texto := strings.NewReader(linhas3Da43)
	m := Listar(texto, "SIGN")

	if len(m) != 2 {
		t.Error("Esperando 2 runas")
	}

	for k := range m {
		if !(k == "U+003D" || k == "U+003E") {
			t.Errorf("Não esperava a runa %q", k)
		}
	}
	// Output:
	// U+003D	=	EQUALS SIGN
	// U+003E	>	GREATER-THAN SIGN
}

func TestListar_duasPalavras(t *testing.T) {
	texto := strings.NewReader(linhas3Da43)
	m := Listar(texto, "CAPITAL LATIN")

	if len(m) != 3 {
		t.Error("Esperando 3 runas")
	}

	for k := range m {
		if !(k == "U+0041" || k == "U+0042" || k == "U+0043") {
			t.Errorf("Não esperava a runa %q", k)
		}
	}
	// Output:
	// U+0041	A	LATIN CAPITAL LETTER A
	// U+0042	B	LATIN CAPITAL LETTER B
	// U+0043	C	LATIN CAPITAL LETTER C
}

func restaurar(nomeVar, valor string, existia bool) {
	if existia {
		os.Setenv(nomeVar, valor)
	} else {
		os.Unsetenv(nomeVar)
	}
}

func TestObterCaminhoUCD_setado(t *testing.T) {
	caminhoAntes, existia := os.LookupEnv("UCD_PATH")
	defer restaurar("UCD_PATH", caminhoAntes, existia)
	caminhoUCD := fmt.Sprintf("./TEST%d-UnicodeData.txt", time.Now().UnixNano())
	os.Setenv("UCD_PATH", caminhoUCD)
	obtido := ObterCaminhoUCD()
	if obtido != caminhoUCD {
		t.Errorf("ObterCaminhoUCD() [setado]\nesperado: %q; recebido: %q", caminhoUCD, obtido)
	}
}

func TestObterCaminhoUCD_default(t *testing.T) {
	caminhoAntes, existia := os.LookupEnv("UCD_PATH")
	defer restaurar("UCD_PATH", caminhoAntes, existia)
	os.Unsetenv("UCD_PATH")
	sufixoCaminhoUCD := "/UnicodeData.txt"
	obtido := ObterCaminhoUCD()
	if !strings.HasSuffix(obtido, sufixoCaminhoUCD) {
		t.Errorf("ObterCaminhoUCD() [default]\nesperado (sufixo): %q; recebido: %q", sufixoCaminhoUCD, obtido)
	}
}

func TestBaixarUCD(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(linhas3Da43))
		}))
	defer srv.Close()

	caminhoUCD := fmt.Sprintf("./TEST%d-UnicodeData.txt", time.Now().UnixNano())
	feito := make(chan bool)
	go baixarUCD(srv.URL, caminhoUCD, feito)
	_ = <-feito
	ucd, err := os.Open(caminhoUCD)
	if os.IsNotExist(err) {
		t.Errorf("baixarUCD não gerou:%v\n%v", caminhoUCD, err)
	}
	ucd.Close()
	os.Remove(caminhoUCD)
}

func TestAbrirUCD_local(t *testing.T) {
	caminhoUCD := ObterCaminhoUCD()
	ucd, err := AbrirUCD(caminhoUCD)
	if err != nil {
		t.Errorf("AbrirUCD(%q):\n%v", caminhoUCD, err)
	}
	ucd.Close()
}

func TestAbrirUCD_remoto(t *testing.T) {
	if testing.Short() {
		t.Skip("teste ignorado [opção -test.short]")
	}
	caminhoUCD := fmt.Sprintf("./TEST%d-UnicodeData.txt", time.Now().UnixNano())
	ucd, err := AbrirUCD(caminhoUCD)
	if err != nil {
		t.Errorf("AbrirUCD(%q):\n%v", caminhoUCD, err)
	}
	ucd.Close()
	os.Remove(caminhoUCD)
}
