package gfxw

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"../path"
)

// var serververzeichnis string = os.Getenv("GOPATH")+"\\src\\gfxw\\gfxwserver\\"
var path_to_server string = path.Give_Path()
var serververzeichnis string = (path_to_server + "\\gfxw\\gfxwserver\\")

// const serververzeichnis string = "/home/lewein/go/bin/"

var zielIP string = "127.0.0.1"
var portnummer uint16 = 55555
var anfragekanal net.Conn  // für nichtblockierende Anfragen an den Server !!!
var mauskanal net.Conn     // für ggf. blockierende Mausanfragen an den Server
var tastaturkanal net.Conn // für ggf. blockierende Tastaturanfragen an den Server
var fensterschloss = make(chan int, 1)
var kommunikationsschloss = make(chan int, 1)
var mausschloss = make(chan int, 1)
var tastaturschloss = make(chan int, 1)
var serverprozess *os.Process

// intern
func start(args ...string) (p *os.Process, err error) {
	if args[0], err = exec.LookPath(args[0]); err == nil {
		var procAttr os.ProcAttr
		procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
		p, err := os.StartProcess(args[0], args, &procAttr)
		if err == nil {
			return p, nil
		}
	}
	return nil, err
}

func gfxServerAnfrage(anfrage string) string {
	var laenge int32
	var b = make([]byte, 1024)     //Lese-Puffer
	var l []byte = make([]byte, 4) //kodierte Länge der Nachricht
	var nachricht []byte
	var kanal net.Conn

	laenge = int32(len(anfrage))
	nachricht = make([]byte, laenge)
	for i := 0; i < 4; i++ {
		l[i] = byte(laenge % 256)
		laenge = laenge / 256
	}
	copy(nachricht, anfrage)
	nachricht = append(l, nachricht...) //4-Byte-Präfix gibt die Länge der Nachricht an

	switch anfrage {
	case "MAL1", "MPL1":
		mausschloss <- 1 // WICHTIG: Mauskanal sperren, nur ein Prozess darf die Maus zur gleichen Zeit lesen!
		kanal = mauskanal
	case "TAL1", "TPL1":
		tastaturschloss <- 1 // WICHTIG: Tastaturkanal sperren, nur ein Prozess darf die Tastatur zur gleichen Zeit lesen!
		kanal = tastaturkanal
	default:
		kommunikationsschloss <- 1 // WICHTIG: Kanal sperren
		kanal = anfragekanal       //         Nur ein Prozess darf mit dem Server kommunizieren!!
		//         Sonst gibt es ein durcheinander auf dem Kanal!!
	}
	n, err := kanal.Write(nachricht)
	if err != nil {
		panic("Übertragungsfehler bzgl. der Verbindung beim Senden!")
	}
	if n != len(nachricht) {
		panic("Es sind Bytes beim Senden verloren gegangen!")
	}
	// NUN WIRD DIE ANTWORT VOM SERVER EMPFANGEN
	var erwartet, angekommen int32
	n, err = kanal.Read(l)
	if n < 4 || err != nil {
		if err.Error() == "EOF" || err.Error()[0:4] == "read" { // Die Gegenseite existiert nicht mehr!
			panic("Das Grafikfenster wurde geschlossen! Programmabbruch!!")
		} else {
			panic("Fehler beim Empfangen des Nachrichtenbeginns")
		}
	}
	for i := 0; i < 4; i++ {
		erwartet = erwartet*256 + int32(l[3-i])
	}
	nachricht = make([]byte, erwartet)
	for angekommen < erwartet {
		n, err = kanal.Read(b)
		if err != nil {
			panic("Fehler beim Empfangen der Nachricht!")
		}
		copy(nachricht[angekommen:], b[0:n])
		angekommen = angekommen + int32(n)
	}
	switch anfrage {
	case "MAL1", "MPL1":
		<-mausschloss // Kanal wieder entsperren
	case "TAL1", "TPL1":
		<-tastaturschloss
	default:
		<-kommunikationsschloss // Kanal wieder entsperren
	}
	return string(nachricht)
}

func split(text string) []string {
	var erg []string = make([]string, 0)
	var teil string = ""
	for _, z := range text {
		if z == ':' {
			erg = append(erg, teil)
			teil = ""
		} else {
			teil = teil + string(z)
		}
	}
	erg = append(erg, teil)
	return erg
}

// NEUE Funktionen unter Windows
func GfxPortnummer() uint16 {
	return portnummer
}

func SetzeGfxPortnummer(p uint16) {
	portnummer = p
}

func SetzeServerprotokoll(w bool) {
	if w {
		if gfxServerAnfrage("SPAN") != "OK" {
			panic("Fehler!!")
		}
	} else {
		if gfxServerAnfrage("SPAU") != "OK" {
			panic("Fehler!!")
		}
	}
}

//Ab hier kommt alt bekanntes

func Fenster(breite, hoehe uint16) {
	b, h, p := fmt.Sprint(breite), fmt.Sprint(hoehe), fmt.Sprint(portnummer)
	prozess, err := start(serververzeichnis+"gfxwserver", b, h, p, zielIP)
	if err != nil {
		fmt.Println(err)
		panic("Der Gfx-Server konnte so nicht gestartet werden!")
	}
	time.Sleep(5e8) // Zeit zum Starten lassen!
	// Jetzt wird der Kanal zum Server aufgebaut
	conn, err := net.Dial("tcp", zielIP+":"+fmt.Sprint(portnummer))
	if err != nil {
		panic("TCP-Verbindung zum Gfx-Server konnte nicht aufgebaut werden!")
	}
	anfragekanal = conn
	conn2, err2 := net.Dial("tcp", zielIP+":"+fmt.Sprint(portnummer))
	if err2 != nil {
		panic("TCP-Verbindung zum Gfx-Server konnte nicht aufgebaut werden!")
	}
	mauskanal = conn2
	conn3, err3 := net.Dial("tcp", zielIP+":"+fmt.Sprint(portnummer))
	if err3 != nil {
		panic("TCP-Verbindung zum Gfx-Server konnte nicht aufgebaut werden!")
	}
	tastaturkanal = conn3
	serverprozess = prozess
}

func FensterOffen() bool {
	if anfragekanal == nil {
		return false
	}
	return gfxServerAnfrage("FEOF") == "true"
}

func FensterAus() {
	gfxServerAnfrage("FEAU")
	anfragekanal = nil
	mauskanal = nil
	tastaturkanal = nil
	// ACHTUNG: "gfxwserver" läuft noch und muss abgewürgt werden!!
	err := serverprozess.Kill()
	if err != nil {
		fmt.Println("Fehler beim Serverkill:", err)
	}
	time.Sleep(2e8) // Kurz warten, um sicherzustellen, dass das Fenster aus
	// und der Server down sind
}

func Grafikzeilen() uint16 {
	a := split(gfxServerAnfrage("GRZE"))
	if len(a) != 1 {
		panic("Fehler!!")
	}
	ze, err := strconv.Atoi(a[0])
	if err != nil {
		panic("Fehler beim 'Grafikzeilen'-Aufruf!")
	}
	return uint16(ze)
}

func Grafikspalten() uint16 {
	a := split(gfxServerAnfrage("GRSP"))
	if len(a) != 1 {
		panic("Fehler!!")
	}
	sp, err := strconv.Atoi(a[0])
	if err != nil {
		panic("Fehler beim 'Grafikspalten'-Aufruf!")
	}
	return uint16(sp)
}

func Fenstertitel(s string) {
	gfxServerAnfrage("FETI:" + s)
}

func Cls() {
	gfxServerAnfrage("CLSC")
}

func Stiftfarbe(r, g, b uint8) {
	rot, green, blue := fmt.Sprint(r), fmt.Sprint(g), fmt.Sprint(b)
	gfxServerAnfrage("STFA:" + rot + ":" + green + ":" + blue)
}

func Transparenz(t uint8) {
	gfxServerAnfrage("TRAN:" + fmt.Sprint(t))
}

func Punkt(x, y uint16) {
	xk, yk := fmt.Sprint(x), fmt.Sprint(y)
	gfxServerAnfrage("PNKT:" + xk + ":" + yk)
}

func GibPunktfarbe(x, y uint16) (r, g, b uint8) {
	xk, yk := fmt.Sprint(x), fmt.Sprint(y)
	antw := gfxServerAnfrage("GPTF:" + xk + ":" + yk)
	a := split(antw)
	if len(a) != 3 {
		panic("Fehler!!")
	}
	red, err := strconv.Atoi(a[0])
	green, err2 := strconv.Atoi(a[1])
	blue, err3 := strconv.Atoi(a[2])
	if err != nil || err2 != nil || err3 != nil {
		panic("Fehler beim 'GibPunktFarbe'-Aufruf!")
	}
	return uint8(red), uint8(green), uint8(blue)
}

func Linie(x1, y1, x2, y2 uint16) {
	x1k, y1k, x2k, y2k := fmt.Sprint(x1), fmt.Sprint(y1), fmt.Sprint(x2), fmt.Sprint(y2)
	gfxServerAnfrage("LINE:" + x1k + ":" + y1k + ":" + x2k + ":" + y2k)
}

func Kreis(x, y, r uint16) {
	xk, yk, ra := fmt.Sprint(x), fmt.Sprint(y), fmt.Sprint(r)
	gfxServerAnfrage("KREI:" + xk + ":" + yk + ":" + ra)
}

func Vollkreis(x, y, r uint16) {
	xk, yk, ra := fmt.Sprint(x), fmt.Sprint(y), fmt.Sprint(r)
	gfxServerAnfrage("VOKR:" + xk + ":" + yk + ":" + ra)
}

func Ellipse(x, y, rx, ry uint16) {
	xk, yk, rxk, ryk := fmt.Sprint(x), fmt.Sprint(y), fmt.Sprint(rx), fmt.Sprint(ry)
	gfxServerAnfrage("ELLI:" + xk + ":" + yk + ":" + rxk + ":" + ryk)
}

func Vollellipse(x, y, rx, ry uint16) {
	xk, yk, rxk, ryk := fmt.Sprint(x), fmt.Sprint(y), fmt.Sprint(rx), fmt.Sprint(ry)
	gfxServerAnfrage("VOEL:" + xk + ":" + yk + ":" + rxk + ":" + ryk)
}

func Kreissektor(x, y, r, w1, w2 uint16) {
	xk, yk, ra, wi1, wi2 := fmt.Sprint(x), fmt.Sprint(y), fmt.Sprint(r), fmt.Sprint(w1), fmt.Sprint(w2)
	gfxServerAnfrage("KRSE:" + xk + ":" + yk + ":" + ra + ":" + wi1 + ":" + wi2)
}

func Vollkreissektor(x, y, r, w1, w2 uint16) {
	xk, yk, ra, wi1, wi2 := fmt.Sprint(x), fmt.Sprint(y), fmt.Sprint(r), fmt.Sprint(w1), fmt.Sprint(w2)
	gfxServerAnfrage("VKSE:" + xk + ":" + yk + ":" + ra + ":" + wi1 + ":" + wi2)
}

func Rechteck(x1, y1, b, h uint16) {
	x1k, y1k, b1, h1 := fmt.Sprint(x1), fmt.Sprint(y1), fmt.Sprint(b), fmt.Sprint(h)
	gfxServerAnfrage("RECH:" + x1k + ":" + y1k + ":" + b1 + ":" + h1)
}

func Vollrechteck(x1, y1, b, h uint16) {
	x1k, y1k, b1, h1 := fmt.Sprint(x1), fmt.Sprint(y1), fmt.Sprint(b), fmt.Sprint(h)
	gfxServerAnfrage("VORE:" + x1k + ":" + y1k + ":" + b1 + ":" + h1)
}

func Dreieck(x1, y1, x2, y2, x3, y3 uint16) {
	x1k, y1k, x2k, y2k, x3k, y3k := fmt.Sprint(x1), fmt.Sprint(y1), fmt.Sprint(x2), fmt.Sprint(y2), fmt.Sprint(x3), fmt.Sprint(y3)
	gfxServerAnfrage("DREI:" + x1k + ":" + y1k + ":" + x2k + ":" + y2k + ":" + x3k + ":" + y3k)
}

func Volldreieck(x1, y1, x2, y2, x3, y3 uint16) {
	x1k, y1k, x2k, y2k, x3k, y3k := fmt.Sprint(x1), fmt.Sprint(y1), fmt.Sprint(x2), fmt.Sprint(y2), fmt.Sprint(x3), fmt.Sprint(y3)
	gfxServerAnfrage("VODR:" + x1k + ":" + y1k + ":" + x2k + ":" + y2k + ":" + x3k + ":" + y3k)
}

func Schreibe(x, y uint16, s string) {
	x1, y1 := fmt.Sprint(x), fmt.Sprint(y)
	gfxServerAnfrage("SCHR:" + x1 + ":" + y1 + ":" + s)
}

func SetzeFont(s string, groesse int) bool {
	g := fmt.Sprint(groesse)
	antw := gfxServerAnfrage("SEFO:" + g + ":" + s)
	if antw[0:1] == "E" {
		panic("Fehler!!")
	}
	return antw == "true"
}

func GibFont() string {
	return gfxServerAnfrage("GIFO:")
}

func SchreibeFont(x, y uint16, s string) {
	x1, y1 := fmt.Sprint(x), fmt.Sprint(y)
	gfxServerAnfrage("SCFO:" + x1 + ":" + y1 + ":" + s)
}

func LadeBild(x, y uint16, s string) {
	x1, y1 := fmt.Sprint(x), fmt.Sprint(y)
	gfxServerAnfrage("LABI:" + x1 + ":" + y1 + ":" + s)
}

func LadeBildMitColorKey(x, y uint16, s string, r, g, b uint8) {
	x1, y1 := fmt.Sprint(x), fmt.Sprint(y)
	gfxServerAnfrage("LBMC:" + x1 + ":" + y1 + ":" + s + ":" + fmt.Sprint(r) + ":" + fmt.Sprint(g) + ":" + fmt.Sprint(b))
}

func LadeBildInsClipboard(s string) {
	gfxServerAnfrage("LBIC:" + s)
}

func Archivieren() {
	gfxServerAnfrage("ARCH")
}

func Restaurieren(x, y, b, h uint16) {
	xk, yk, b1, h1 := fmt.Sprint(x), fmt.Sprint(y), fmt.Sprint(b), fmt.Sprint(h)
	gfxServerAnfrage("REST:" + xk + ":" + yk + ":" + b1 + ":" + h1)
}

func Clipboard_kopieren(x, y, b, h uint16) {
	xk, yk, b1, h1 := fmt.Sprint(x), fmt.Sprint(y), fmt.Sprint(b), fmt.Sprint(h)
	gfxServerAnfrage("CLKO:" + xk + ":" + yk + ":" + b1 + ":" + h1)
}

func Clipboard_einfuegen(x, y uint16) {
	xk, yk := fmt.Sprint(x), fmt.Sprint(y)
	gfxServerAnfrage("CLEI:" + xk + ":" + yk)
}

func Clipboard_einfuegenMitColorKey(x, y uint16, r, g, b uint8) {
	xk, yk := fmt.Sprint(x), fmt.Sprint(y)
	gfxServerAnfrage("CEMC:" + xk + ":" + yk + ":" + fmt.Sprint(r) + ":" + fmt.Sprint(g) + ":" + fmt.Sprint(b))
}

func Sperren() {
	fensterschloss <- 1
}

func Entsperren() {
	<-fensterschloss
}

func UpdateAus() {
	gfxServerAnfrage("UPAU")
}

func UpdateAn() {
	gfxServerAnfrage("UPAN")
}

func TastaturLesen1() (taste uint16, gedrueckt uint8, tiefe uint16) {
	a := split(gfxServerAnfrage("TAL1"))
	if len(a) != 3 {
		panic("Fehler!!")
	}
	ta, err := strconv.Atoi(a[0])
	ge, err2 := strconv.Atoi(a[1])
	ti, err3 := strconv.Atoi(a[2])
	if err != nil || err2 != nil || err3 != nil {
		panic("Fehler beim 'TastaturLesen1'-Aufruf!")
	}
	taste, gedrueckt, tiefe = uint16(ta), uint8(ge), uint16(ti)
	return
}

func Tastaturzeichen(taste, tiefe uint16) rune {
	ta, ti := fmt.Sprint(taste), fmt.Sprint(tiefe)
	a := split(gfxServerAnfrage("TAZE:" + ta + ":" + ti))
	if len(a) != 1 {
		panic("Fehler!!")
	}
	ru, err := strconv.Atoi(a[0])
	if err != nil {
		panic("Fehler beim 'Tastaturzeichen'-Aufruf!")
	}
	return rune(ru)
}

func TastaturpufferAn() {
	gfxServerAnfrage("TPAN")
}

func TastaturpufferAus() {
	gfxServerAnfrage("TPAU")
}

func TastaturpufferLesen1() (taste uint16, gedrueckt uint8, tiefe uint16) {
	a := split(gfxServerAnfrage("TPL1"))
	if len(a) != 3 {
		panic("Fehler!!")
	}
	ta, err := strconv.Atoi(a[0])
	ge, err2 := strconv.Atoi(a[1])
	ti, err3 := strconv.Atoi(a[2])
	if err != nil || err2 != nil || err3 != nil {
		panic("Fehler beim 'TastaturpufferLesen1'-Aufruf!")
	}
	taste, gedrueckt, tiefe = uint16(ta), uint8(ge), uint16(ti)
	return
}

func MausLesen1() (taste uint8, status int8, mausX, mausY uint16) {
	a := split(gfxServerAnfrage("MAL1"))
	if len(a) != 4 {
		panic("Fehler!!")
	}
	ta, err := strconv.Atoi(a[0])
	st, err2 := strconv.Atoi(a[1])
	mx, err3 := strconv.Atoi(a[2])
	my, err4 := strconv.Atoi(a[3])
	if err != nil || err2 != nil || err3 != nil || err4 != nil {
		panic("Fehler beim 'MausLesen1'-Aufruf!")
	}
	taste, status, mausX, mausY = uint8(ta), int8(st), uint16(mx), uint16(my)
	return
}

func MauspufferAn() {
	gfxServerAnfrage("MPAN")
}

func MauspufferAus() {
	gfxServerAnfrage("MPAU")
}

func MauspufferLesen1() (taste uint8, status int8, mausX, mausY uint16) {
	a := split(gfxServerAnfrage("MPL1"))
	if len(a) != 4 {
		panic("Fehler!!")
	}
	ta, err := strconv.Atoi(a[0])
	st, err2 := strconv.Atoi(a[1])
	mx, err3 := strconv.Atoi(a[2])
	my, err4 := strconv.Atoi(a[3])
	if err != nil || err2 != nil || err3 != nil || err4 != nil {
		panic("Fehler beim 'MauspufferLesen1'-Aufruf!")
	}
	taste, status, mausX, mausY = uint8(ta), int8(st), uint16(mx), uint16(my)
	return
}

func SpieleSound(s string) {
	gfxServerAnfrage("SPSO:" + s)
}

func GibNotenTempo() uint8 {
	a := split(gfxServerAnfrage("GNTE"))
	if len(a) != 1 {
		panic("Fehler!!")
	}
	tempo, err := strconv.Atoi(a[0])
	if err != nil {
		panic("Fehler beim 'GibNotenTempo'-Aufruf!")
	}
	return uint8(tempo)
}

func SetzeNotenTempo(t uint8) {
	a := gfxServerAnfrage("SNTE:" + fmt.Sprint(t))
	if a[0:1] == "E" {
		panic("Fehler bei 'SetzeNotenTempo' !")
	}
}

func GibHuellkurve() (float64, float64, float64, float64) {
	a := split(gfxServerAnfrage("GHUE"))
	if len(a) != 4 {
		panic("Fehler!!")
	}
	a1, err := strconv.ParseFloat(a[0], 64)
	a2, err2 := strconv.ParseFloat(a[1], 64)
	a3, err3 := strconv.ParseFloat(a[2], 64)
	a4, err4 := strconv.ParseFloat(a[3], 64)
	if err != nil || err2 != nil || err3 != nil || err4 != nil {
		panic("Fehler beim 'GibHuellkurve'-Aufruf!")
	}
	return a1, a2, a3, a4
}

func SetzeHuellkurve(a, d, s, r float64) {
	antw := gfxServerAnfrage("SHUE:" + fmt.Sprint(a) + ":" + fmt.Sprint(d) + ":" + fmt.Sprint(s) + ":" + fmt.Sprint(r))
	if antw[0:1] == "E" {
		panic("Fehler bei 'SetzeHuellkurve' !")
	}
}

func GibKlangparameter() (r uint32, b uint8, k uint8, s uint8, p float64) {
	a := split(gfxServerAnfrage("GKPA"))
	if len(a) != 5 {
		panic("Fehler!!")
	}
	rate, err := strconv.Atoi(a[0])
	bits, err2 := strconv.Atoi(a[1])
	kanaele, err3 := strconv.Atoi(a[2])
	signalform, err4 := strconv.Atoi(a[3])
	pulsweite, err5 := strconv.ParseFloat(a[4], 64)
	if err != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		panic("Fehler beim 'GibKlangparameter'-Aufruf!")
	}
	r, b, k, s, p = uint32(rate), uint8(bits), uint8(kanaele), uint8(signalform), pulsweite
	return
}

func SetzeKlangparameter(rate uint32, aufloesung, kanaele, signal uint8, p float64) {
	a := gfxServerAnfrage("SKPA:" + fmt.Sprint(rate) + ":" + fmt.Sprint(aufloesung) + ":" + fmt.Sprint(kanaele) + ":" + fmt.Sprint(signal) + ":" + fmt.Sprint(p))
	if a[0:1] == "E" {
		panic("Fehler bei 'SetzeKlangparameter' !")
	}
}

func SpieleNote(tonname string, laenge float64, wartedauer float64) {
	gfxServerAnfrage("SPNO:" + tonname + ":" + fmt.Sprint(laenge) + ":" + fmt.Sprint(wartedauer))
}
