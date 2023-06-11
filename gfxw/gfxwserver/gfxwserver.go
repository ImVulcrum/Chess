package main
// Autor: St. Schmidt
// Datum: 04.02.2019-07.02.2019; letzte Änderung: 18.10.2020
// Zweck: TCP/IP-Server, der ein gfx-Grafikfenster verwaltet

// Autor: Stefan Schmidt (Kontakt: St.Schmidt@online.de)
// Datum: 07.03.2016 ; letzte Änderung: 18.10.2020
// Zweck: - Grafik- und Soundausgabe und Eingabe per Tastatur und Maus
//          mit Go unter Windows und unter Linux
//        - 06.11.2020  Bug bei 'Clipboard_einfuegen' bzgl. der Transparenz entfernt
//        - 18.10.2020: Bug in 'FensterAus' entfernt,
//                      neue Funktion 'Fenstertitel' zur Festlegung eines eigenen
//                      Fenstertitels für das Grafikfenster,
//                      neue Funktionen 'LadeBildMitColorKey' und 
//                      'Clipboard_einfuegenMitColorKey'bei der Pixel einer
//                      bestimmten Farbe trasparent dargestellt werden (gut für "Sprites"!)
//                      neue Funktion 'Trasparenz', damit sich überdeckende Grafik-
//                      objekte erkennen kann
//                      'Bug entfernt': Mit 'defer' angemeldete Funktionsaufrufe werden
//                      nun nach dem Schließen des Grafikfensters auch noch ausgeführt.
//        - 01.09.2019: Bugfix: Maus- und Tastaturabfragen funktionieren nun nebenläufig
//                      zu Änderungen im Grafikfenster ; Einbau der neuen Musikbefehle
//        - 07.02.2019: Umbau zu einem Server, der genau ein Grafikfenster verwaltet
//        - 03.03.2018: Die Funktion 'SetzeFont' liefert nun einen Rückgabewert,
//                      der den Erfolg/Misserfolg angibt.
//        - 07.10.2017: neue Funktion 'Tastaturzeichen'
//        - 07.10.2017: 'Bug' in Funktion 'Cls()' entfernt - KEIN FLACKERN MEHR 
//                       bei 'double-buffering' mit UpdateAus() und UpdateAn()

/*
#cgo LDFLAGS: -lSDL -lSDL_gfx -lSDL_ttf
#include <SDL/SDL.h>
#include <SDL/SDL_ttf.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <SDL/SDL_gfxPrimitives.h>

// Structure for loaded sounds. 
typedef struct sound_s {
    Uint8 *samples;		// raw PCM sample data 
    Uint32 length;		// size of sound data in bytes 
} sound_t, *sound_p;

// Structure for a currently playing sound. 
typedef struct playing_s {
    int active;                 // 1 if this sound should be played 
    sound_p sound;              // sound data to play 
    Uint32 position;            // current position in the sound buffer 
} playing_t, *playing_p;

// Array for all active sound effects. 
#define MAX_PLAYING_SOUNDS      10 
playing_t playing[MAX_PLAYING_SOUNDS];

// The higher this is, the louder each currently playing sound will be.
// However, high values may cause distortion if too many sounds are
// playing. Experiment with this. 
#define VOLUME_PER_SOUND        SDL_MIX_MAXVOLUME / 2

static SDL_Surface *screen;
static SDL_Surface *archiv;
static SDL_Surface *clipboard = NULL;
static Uint8 updateOn = 1;
static Uint8 red,green,blue, alpha;
static Uint8 ck_red, ck_green, ck_blue;
static Uint8 colorkey;
static SDL_Event event; 
static Uint8 gedrueckt;
static Uint16 taste,tiefe;
static Uint8 tasteLesen = 0;
static Uint8 tastaturpuffer = 0;
static Uint32 t_puffer[256];
static Uint8 t_pufferkopf;
static Uint8 t_pufferende;
static Uint16 mausX, mausY;
static Uint8 mausLesen = 0;
static Uint8 mausTaste;
static Uint8 mauspuffer = 0;
static Uint32 m_puffer[256];
static Uint8 m_pufferkopf;
static Uint8 m_pufferende;
static Uint8 fensteroffen = 0;
static Uint8 fensterzu = 1;
static char aktFont[256];
static int aktFontSize;
static SDL_AudioSpec desired, obtained; // Audio format specifications.
static sound_t s[10];                   // Our loaded sounds and their formats. 
   
//------------------------------------------------------------------
// This function is called by SDL whenever the sound card
// needs more samples to play. It might be called from a
// separate thread, so we should be careful what we touch. 
static void AudioCallback(void *user_data, Uint8 *audio, int length)
{
    int i;
    // Avoid compiler warning. 
    user_data += 0;
    // Clear the audio buffer so we can mix samples into it. 
    memset(audio, 0, length);
    // Mix in each sound. 
    for (i = 0; i < MAX_PLAYING_SOUNDS; i++) {
	  if (playing[i].active) {
	    Uint8 *sound_buf;
	    Uint32 sound_len;
	    // Locate this sound's current buffer position. 
	    sound_buf = playing[i].sound->samples;
	    sound_buf += playing[i].position;
	    // Determine the number of samples to mix. 
	    if ((playing[i].position + length) > playing[i].sound->length) {
		sound_len = playing[i].sound->length - playing[i].position;
	    } else {
		sound_len = length;
	    }
	    // Mix this sound into the stream. 
	    SDL_MixAudio(audio, sound_buf, sound_len, VOLUME_PER_SOUND);
	    // Update the sound buffer's position. 
	    playing[i].position += length;
	    // Have we reached the end of the sound? 
	    if (playing[i].position >= playing[i].sound->length) {
	    free(s[i].samples);      //zugehörigen Soundstruktur-Samplespeicher wieder freigeben 
		playing[i].active = 0;	 // und anschließend als inaktiv markieren
	    }
	  }
    }
}
//----------------------------------------------------------------
// This function loads a sound with SDL_LoadWAV and converts
// it to the specified sample format. Returns 0 on success
// and 1 on failure. 
static int LoadAndConvertSound(char *filename, SDL_AudioSpec *spec,
			sound_p sound)
{
    SDL_AudioCVT cvt;           // audio format conversion structure 
    SDL_AudioSpec loaded;       // format of the loaded data 
    Uint8 *new_buf;
    // Load the WAV file in its original sample format. 
    if (SDL_LoadWAV(filename,
		    &loaded, &sound->samples, &sound->length) == NULL) {
	//printf("Unable to load sound: %s\n", SDL_GetError());
	return 1;
    }
    // Build a conversion structure for converting the samples.
    // This structure contains the data SDL needs to quickly
    // convert between sample formats. 
    if (SDL_BuildAudioCVT(&cvt, loaded.format,
			  loaded.channels,
			  loaded.freq,
			  spec->format, spec->channels, spec->freq) < 0) {
	// printf("Unable to convert sound: %s\n", SDL_GetError());
	return 1;
    }
    // Since converting PCM samples can result in more data
    //   (for instance, converting 8-bit mono to 16-bit stereo),
    //   we need to allocate a new buffer for the converted data.
    //   Fortunately SDL_BuildAudioCVT supplied the necessary
    //   information. 
    cvt.len = sound->length;
    new_buf = (Uint8 *) malloc(cvt.len * cvt.len_mult);
    if (new_buf == NULL) {
	//printf("Memory allocation failed.\n");
	SDL_FreeWAV(sound->samples);
	return 1;
    }
    // Copy the sound samples into the new buffer.
    memcpy(new_buf, sound->samples, sound->length);
    // Perform the conversion on the new buffer. 
    cvt.buf = new_buf;
    if (SDL_ConvertAudio(&cvt) < 0) {
	//printf("Audio conversion error: %s\n", SDL_GetError());
	free(new_buf);
	SDL_FreeWAV(sound->samples);
	return 1;
    }
    // Swap the converted data for the original. 
    SDL_FreeWAV(sound->samples);
    sound->samples = new_buf;
    sound->length = sound->length * cvt.len_mult;
    // Success! 
    //printf("'%s' was loaded and converted successfully.\n", filename);
    return 0;
}
//----------------------------------------------------------------
// Diese Funktion übernimmt eine Bytefolge aus dem RAM ab der Adresse addr
// mit der Länge laenge, die dem Inhalt einer WAV-DAtei entspricht und konvertiert
// sie , damit es abgespielt werden kann. Die Funktion liefert 0 bei Erfolg
// 1 bei Misserfolg. 
static int LadeUndKonvertiereRAMWAV(const void* addr, int laenge, SDL_AudioSpec *spec,
			sound_p sound)
{
    SDL_AudioCVT cvt;           // audio format conversion structure 
    SDL_AudioSpec loaded;       // format of the loaded data 
    Uint8 *new_buf;
    // Lade 'RAMWAV' im Originalformat:  
    if (SDL_LoadWAV_RW(SDL_RWFromConstMem(addr,laenge),0,
		    &loaded, &sound->samples, &sound->length) == NULL) {
	printf("Unable to load sound: %s\n", SDL_GetError());
	return 1;
    }
    // Build a conversion structure for converting the samples.
    // This structure contains the data SDL needs to quickly
    // convert between sample formats. 
    if (SDL_BuildAudioCVT(&cvt, loaded.format,
			  loaded.channels,
			  loaded.freq,
			  spec->format, spec->channels, spec->freq) < 0) {
	// printf("Unable to convert sound: %s\n", SDL_GetError());
	return 1;
    }
    // Since converting PCM samples can result in more data
    //   (for instance, converting 8-bit mono to 16-bit stereo),
    //   we need to allocate a new buffer for the converted data.
    //   Fortunately SDL_BuildAudioCVT supplied the necessary
    //   information. 
    cvt.len = sound->length;
    new_buf = (Uint8 *) malloc(cvt.len * cvt.len_mult);
    if (new_buf == NULL) {
	//printf("Memory allocation failed.\n");
	SDL_FreeWAV(sound->samples);
	return 1;
    }
    // Copy the sound samples into the new buffer.
    memcpy(new_buf, sound->samples, sound->length);
    // Perform the conversion on the new buffer. 
    cvt.buf = new_buf;
    if (SDL_ConvertAudio(&cvt) < 0) {
	//printf("Audio conversion error: %s\n", SDL_GetError());
	free(new_buf);
	SDL_FreeWAV(sound->samples);
	return 1;
    }
    // Swap the converted data for the original. 
    SDL_FreeWAV(sound->samples);
    sound->samples = new_buf;
    sound->length = sound->length * cvt.len_mult;
    // Success! 
    //printf("'%s' was loaded and converted successfully.\n", filename);
    return 0;
}
//-----------------------------------------------------------------
static int LoadAndPlaySound (char *filename) 
{
	int i;
	//Finde einen freien Index (Bereich 0 <= index < MAX_PLAYING_SOUND
	for (i = 0; i < MAX_PLAYING_SOUNDS; i++) {
	if (playing[i].active == 0)
	    break;
    } 
    if (i == MAX_PLAYING_SOUNDS)
	return 1; //Fehler: Es werden schon die max. Anzahl an Dateien abgespielt.

	//Lade und konvertiere den Sound in die entsprechende Soundstruktur
	if (LoadAndConvertSound(filename, &obtained, &s[i]) != 0) {
	  return 2; //Laden fehlgeschlagen!
    }
    //Abspielen starten
    // The 'playing' structures are accessed by the audio callback,
    // so we should obtain a lock before we access them. 
    SDL_LockAudio();
    playing[i].active = 1;
    playing[i].sound = &s[i];
    playing[i].position = 0;
    SDL_UnlockAudio();
    return 0;
}    
//-----------------------------------------------------------------
static int LadeUndSpieleNote (const void* addr, int laenge) 
{
	int i;
	//Finde einen freien Index (Bereich 0 <= index < MAX_PLAYING_SOUND
	for (i = 0; i < MAX_PLAYING_SOUNDS; i++) {
	if (playing[i].active == 0)
	    break;
    } 
    if (i == MAX_PLAYING_SOUNDS)
	return 1; //Fehler: Es werden schon die max. Anzahl an Dateien abgespielt.

	//Lade und konvertiere den Sound in die entsprechende Soundstruktur
	if (LadeUndKonvertiereRAMWAV(addr, laenge, &obtained, &s[i]) != 0) {
	  return 2; //Laden fehlgeschlagen!
    }
    //Abspielen starten
    // The 'playing' structures are accessed by the audio callback,
    // so we should obtain a lock before we access them. 
    SDL_LockAudio();
    playing[i].active = 1;
    playing[i].sound = &s[i];
    playing[i].position = 0;
    SDL_UnlockAudio();
    return 0;
}   
//------------------------------------------------------------------
static int setFont (char *fontfile, int groesse) {
	strcpy (aktFont,fontfile);
	aktFontSize = groesse;
	TTF_Font *font = TTF_OpenFont(aktFont, aktFontSize);
	if (!font) {
	  //printf("TTF_OpenFont: %s\n", TTF_GetError());
	  return 1;
    }
    TTF_CloseFont(font);
	return 0;
}
//------------------------------------------------------------------
static char *getFont () {
	return aktFont;
}
//------------------------------------------------------------------
static int write (Sint16 x, Sint16 y, char *text) {
	TTF_Font *font = TTF_OpenFont(aktFont, aktFontSize);
	if (!font) {
	  //printf("TTF_OpenFont: %s\n", TTF_GetError());
	  return 1;
    }
	SDL_Color clrFg = {red,green,blue,alpha};  
	SDL_Surface *sText = TTF_RenderUTF8_Solid(font,text,clrFg);
	SDL_Rect rcDest = {x,y,0,0};
	SDL_BlitSurface(sText,NULL, screen,&rcDest);
	SDL_FreeSurface(sText);
	if (updateOn)
	  SDL_UpdateRect(screen,0,0,0,0);
	TTF_CloseFont(font);
	return 0;
}
//------------------------------------------------------------------
static void clearscreen () {
  SDL_FillRect(screen, NULL, SDL_MapRGB(screen->format, red, green, blue));
  if (updateOn)
    SDL_UpdateRect (screen,0,0,0,0);
}
//------------------------------------------------------------
static int GrafikfensterAn (Uint16 breite, Uint16 hoehe) 
{   
    if ( fensteroffen == 1) return 1;  //Es kann nur ein Grafikfenster geben!

	//1. SDL muss initialisiert werden.
	if (SDL_Init(SDL_INIT_VIDEO | SDL_INIT_AUDIO) != 0) {
		//printf ("Kann SDL nicht initialisieren: %s\n", SDL_GetError ());
		return 1;
	}		
	//2. Bekanntmachung: Diese Funktion soll mit dem Programmende aufgerufen werden.
	// atexit (SDL_Quit);
	//3. Bildschirm: Hier kann man auch SDL_DOUBLEBUF sagen!
	screen = SDL_SetVideoMode (breite,hoehe, 32, SDL_DOUBLEBUF); //SDL_FULLSCREEN);
	if (screen == NULL) {
		//printf ("Bildschirm-Modus nicht setzbar: %s\n",SDL_GetError ());
		return 1;
	}
	SDL_WM_SetCaption( "LWB FU-Berlin: GO-Grafikfenster", 0 );
	
	TTF_Init();
	
	red = 255;
	green = 255;
	blue = 255;
	alpha = 255; 
	clearscreen ();
	red   = 0;
	green = 0;
	blue  = 0;
	
	//Archiv-Surface erstellen
	archiv = SDL_ConvertSurface (screen, screen->format, SDL_HWSURFACE);
	if (archiv == NULL) {
		//printf ("Archiv-Surface konnte nicht erzeugt werden!\n");
		return 1;
	}
	
    // Open the audio device. The sound driver will try to give us
    // the requested format, but it might not succeed. The 'obtained'
    // structure will be filled in with the actual format data. 
    desired.freq = 44100;	// desired output sample rate 
    desired.format = AUDIO_S16;	// request signed 16-bit samples 
    desired.samples = 4096;	// this is more or less discretionary 
    desired.channels = 2;	// ask for stereo 
    desired.callback = AudioCallback;
    desired.userdata = NULL;	// we don't need this 
    if (SDL_OpenAudio(&desired, &obtained) < 0) {
    	//printf("Unable to open audio device: %s\n", SDL_GetError());
	    return 1;
    }
    // Initialisiere die Liste der möglichen Sounds (keiner aktiv zu Beginn) 
    int i;
    for (i = 0; i < MAX_PLAYING_SOUNDS; i++) {
	playing[i].active = 0;
    }

    // SDL's audio is initially paused. Start it. 
    SDL_PauseAudio(0);

	fensteroffen = 1; 
	fensterzu = 0;

	//Jetzt kommt die Event-Loop
	
	while (fensteroffen == 1 && SDL_WaitEvent(&event) != 0 ) {
		switch (event.type) {
			case SDL_KEYDOWN:
				if (tasteLesen)
				{
					gedrueckt = 1;                //Taste ist gerade heruntergedrückt.
					taste = event.key.keysym.sym;  //Das ist der Code der Taste auf der Tastatur.
					tiefe = event.key.keysym.mod;  //Gleichzeitig Steuerungstaste(n) gedrückt??
					//printf("%i,%i,%i\n",taste, gedrueckt, tiefe);
					tasteLesen = 0;
				}
				if (tastaturpuffer)
				{
					if (t_pufferende + 1 != t_pufferkopf)
					{
						t_puffer[t_pufferende] = ((Uint32) event.key.keysym.sym)*256*256 + (Uint32) 256*256*256*128 + ((Uint32) event.key.keysym.mod);
						t_pufferende++; //Umschlag auf 0 automatisch, da Uint8
					} 
				}
				break;
			case SDL_KEYUP:
				if (tasteLesen)
				{
					gedrueckt = 0; //Taste wurde gerade losgelassen.
					taste = event.key.keysym.sym;
					tiefe = event.key.keysym.mod;  //Gleichzeitig Steuerungstaste(n) gedrückt??
					//printf("%i,%i,%i\n",taste, gedrueckt, tiefe);
					tasteLesen = 0;
				}
				if (tastaturpuffer)
				{
					if (t_pufferende + 1 != t_pufferkopf)
					{
						t_puffer[t_pufferende] = ((Uint32) event.key.keysym.sym)*256*256 + ((Uint32) event.key.keysym.mod);
						t_pufferende++; //Umschlag auf 0 automatisch, da Uint8
					} 
				}
				break;
			case SDL_MOUSEMOTION:
				if (mausLesen)
				{   //BEi MOUSEMOTION GIBT ES NUR 3 MÖGLICHKEITEN FÜR EINE GEDRÜCKT-GEHALTENE TASTE: 1,2 oder 3
					// Dummerweise ist bei 3 der Tastenwert 4, daher Korrektur:
					mausTaste = (Uint8) event.button.button;
					if (mausTaste == 4)
						mausTaste--;
					mausX     = (Uint16) event.motion.x;
					mausY     = (Uint16) event.motion.y;
					mausLesen = 0;
				}
				if (mauspuffer)
				{   
					mausTaste = (Uint8) event.button.button;
					if (mausTaste == 4)
						mausTaste--;
					mausX     = (Uint16) event.motion.x;
					mausY     = (Uint16) event.motion.y;
					if (m_pufferende + 1 != m_pufferkopf)
					{
						m_puffer[m_pufferende] = ((Uint32) mausTaste)*256*256*256 + (((Uint32) mausX) <<12) + (Uint32) mausY;
						m_pufferende++;
					}
				}
				break;
			case SDL_MOUSEBUTTONDOWN:
				if (mausLesen)
				{
					mausTaste = (Uint8) event.button.button + 128; //+128: "pressed"
					mausX = (Uint16) event.motion.x;
					mausY = (Uint16) event.motion.y;
					mausLesen = 0;
				}
				if (mauspuffer)
				{
					mausTaste = (Uint8) event.button.button + 128; //+128: "pressed"
					mausX = (Uint16) event.motion.x;
					mausY = (Uint16) event.motion.y;
					if (m_pufferende + 1 != m_pufferkopf)
					{
						m_puffer[m_pufferende] = ((Uint32) mausTaste)*256*256*256 + (((Uint32) mausX) <<12) + (Uint32) mausY;
						m_pufferende++;
					}
				}
				break;
			case SDL_MOUSEBUTTONUP:
				if (mausLesen)
				{
					mausTaste = (Uint8) event.button.button + 64; //+64: "released"
					mausX = (Uint16) event.motion.x;
					mausY = (Uint16) event.motion.y;
					mausLesen = 0;
				}
				if (mauspuffer)
				{
					mausTaste = (Uint8) event.button.button + 64; //+64: "released"
					mausX = (Uint16) event.motion.x;
					mausY = (Uint16) event.motion.y;
					if (m_pufferende + 1 != m_pufferkopf)
					{
						m_puffer[m_pufferende] = ((Uint32) mausTaste)*256*256*256 + (((Uint32) mausX) <<12) + (Uint32) mausY;
						m_pufferende++;
					} 
				}
				break;
			case SDL_QUIT:
				//printf("Das Grafikfenster wurde geschlossen. Bye.\n");
				fensteroffen = 0;
				break;
		}
		
	}
	
	// Die event-Loop wurde beendet, also wird nun das Fenster geschlossen!
	TTF_Quit ();
    // Pause and lock the sound system so we can safely delete our sound data. 
    SDL_PauseAudio(1);
    SDL_LockAudio();
    // Free our sounds before we exit, just to be safe.
    for (i=0; i < MAX_PLAYING_SOUNDS;i++) {
		if (playing[i].active ==1) {
			free(s[i].samples);
		}
	}
    // At this point the output is paused and we know for certain that the
    // callback is not active, so we can safely unlock the audio system. 
    SDL_UnlockAudio();
	SDL_CloseAudio();
	SDL_Quit ();
	
	//screen wird automatisch wieder freigegeben
	SDL_FreeSurface(archiv); archiv = NULL;
	if (clipboard != NULL)
	{
		SDL_FreeSurface(clipboard); clipboard = NULL;
	}
	fensterzu = 1;
	
	return 0;
}
//------------------------------------------------------------------
static void Fenstertitel (const char *titel)
{
  SDL_WM_SetCaption(titel,0);
}
//------------------------------------------------------------------
static Uint8 FensterOffen ()
{
  return fensteroffen;
}
//-------------------------------------------------------------------
static Uint8 FensterZu ()
{
  return fensterzu;
}
//-------------------------------------------------------------------
static void GrafikfensterAus ()
{
  SDL_Event user_event;
	user_event.type=SDL_QUIT;
	SDL_PushEvent(&user_event);
}
//-------------------------------------------------------------------
static void updateAus ()
{
	updateOn = 0;
}
//-------------------------------------------------------------------
static void updateAn ()
{
	updateOn = 1;
	SDL_Flip (screen);
}
//-------------------------------------------------------------------
static void zeichnePunkt (Sint16 x, Sint16 y)
{
  pixelRGBA (screen, x, y ,red, green,blue,alpha);
  if (updateOn) 
	SDL_UpdateRect (screen, x, y, 1, 1);
}
//--------------------------------------------------------------
static Uint32 gibPixel(Sint16 x, Sint16 y)
{
    int bpp = screen->format->BytesPerPixel;
    // Here p is the address to the pixel we want to retrieve 
    Uint8 *p = (Uint8 *)screen->pixels + y * screen->pitch + x * bpp;

    switch(bpp) {
    case 1:
        return *p;
        break;
    case 2:
        return *(Uint16 *)p;
        break;
    case 3:
        if(SDL_BYTEORDER == SDL_BIG_ENDIAN)
            return p[0] << 16 | p[1] << 8 | p[2];
        else
            return p[0] | p[1] << 8 | p[2] << 16;
        break;
    case 4:
        return *(Uint32 *)p;
        break;
    default:
        return 0;       // shouldn't happen, but avoids warnings 
    }
}
//--------------------------------------------------------------
static void zeichneKreis (Sint16 x, Sint16 y, Sint16 r, Uint8 full)
{ 
	if (full)
		filledCircleRGBA(screen,x,y,r,red,green,blue,alpha);
	else
		circleRGBA (screen, x,y,r,red, green, blue,alpha);
	if (updateOn)
		SDL_UpdateRect (screen,x-r,y-r,2*r+1,2*r+1);
}
//---------------------------------------------------------------
static void zeichneEllipse (Sint16 x, Sint16 y, Sint16 rx, Sint16 ry, Uint8 filled)
{	
	if (filled)
		filledEllipseRGBA (screen, x, y, rx, ry, red, green, blue, alpha);
	else
		ellipseRGBA (screen, x, y, rx, ry, red, green,blue, alpha);
	if (updateOn)
		SDL_UpdateRect (screen, x-rx, y-ry,2*rx+1,2*ry+1);
}
//---------------------------------------------------------------
static void stiftfarbe (Uint8 r, Uint8 g, Uint8 b)
{
    red = r;
    green = g;
    blue = b;
}
//---------------------------------------------------------------
static void transparenz (Uint8 t)
{
	alpha = t;
}
//---------------------------------------------------------------
static void zeichneStrecke (Sint16 x1, Sint16 y1, Sint16 x2, Sint16 y2)
{ 
	int upx,upy,breite,hoehe;
	
	lineRGBA (screen, x1,y1,x2,y2, red, green, blue, alpha);
	if (x1 <= x2)
	{
		upx    = x1;
		breite = x2 - x1 + 1;
	} 
	else
	{
		upx    = x2;
		breite = x1 - x2 + 1;
	}
	if (y1 <= y2)
	{
		upy   = y1;
		hoehe = y2 - y1 + 1;
	}
	else
	{
		upy   = y2;
		hoehe = y1 - y2 + 1;
	}
	if (updateOn)
		SDL_UpdateRect (screen,upx,upy,breite,hoehe);
}
//--------------------------------------------------------------
static void rechteck (Sint16 x1, Sint16 y1, Sint16 b, Sint16 h, Uint8 filled)
{	
	if (filled)
		boxRGBA (screen, x1, y1 ,x1+b-1, y1+h-1, red, green, blue, alpha);
	else
		rectangleRGBA (screen, x1, y1 , x1+b-1, y1+h-1, red, green,blue,alpha);
	if (updateOn)
		SDL_UpdateRect (screen, x1, y1, b, h);
}	
//--------------------------------------------------------------
static void kreissektor (Sint16 x, Sint16 y, Sint16 r, Sint16 w1, Sint16 w2, Uint8 filled)
{
	if (filled)
		filledPieRGBA (screen, x, y , r, w1, w2, red, green, blue, alpha);
	else
		pieRGBA (screen, x, y , r, w1, w2, red, green, blue, alpha);
	if (updateOn)
		SDL_UpdateRect (screen, x-r, y-r, 2*r+1, 2*r+1);
}
//---------------------------------------------------------------
Sint16 minimum (Sint16 x, Sint16 y, Sint16 z)
{
  if ((x <= y) && (x <=z))
    return x;
  else if ((y<=x) && (y<=z))
    return y;
  else
    return z;
}
//-------------------------------------------------------------
Sint16 maximum (Sint16 x, Sint16 y, Sint16 z)
{
  if ((x >= y) && (x >=z))
    return x;
  else if ((y>=x) && (y>=z))
    return y;
  else
    return z;
}
//---------------------------------------------------------------
static void dreieck (Sint16 x1, Sint16 y1, Sint16 x2, Sint16 y2, Sint16 x3, Sint16 y3, Uint8 filled)
{
	int upx,upy,breite,hoehe;

	upx = minimum (x1, x2, x3);
	upy = minimum (y1, y2, y3);
	breite = maximum (x1, x2, x3) - upx + 1;
	hoehe  = maximum (y1, y2, y3) - upy + 1;
	if (filled)
		filledTrigonRGBA(screen, x1,y1,x2,y2,x3,y3,red,green,blue,alpha);
	else
		trigonRGBA (screen, x1,y1,x2,y2,x3,y3,red,green,blue,alpha);
	if (updateOn)
		SDL_UpdateRect (screen, upx,upy,breite,hoehe);
}
//----------------------------------------------------------------
static void setcolorkey (Uint8 r, Uint8 g, Uint8 b, Uint8 key)
{
  ck_red = r;
  ck_green = g;
  ck_blue = b;
  colorkey = key;
}
//----------------------------------------------------------------
static void ladeBild (Sint16 x, Sint16 y, char *cs)
{
	SDL_Surface *image;
	SDL_Rect src, dest;
	
	image = SDL_LoadBMP(cs);
	//printf ("Dateiname: %s\n",cs);
	if (image == NULL) {
		//printf("Bild konnte nicht geladen werden!\n");
		return;
	}
	src.x = 0;
	src.y = 0;
	src.w = image->w;
	src.h = image->h;
	
	dest.x = x;
	dest.y = y;
	dest.w = image->w;
	dest.h = image->h;

	if (colorkey)
		SDL_SetColorKey(image, SDL_SRCCOLORKEY, SDL_MapRGBA(screen->format, ck_red, ck_green, ck_blue, alpha));
	SDL_SetAlpha(image, SDL_SRCALPHA, alpha);
	SDL_BlitSurface(image, &src, screen, &dest);
	SDL_FreeSurface (image);  
	if (updateOn)
		SDL_UpdateRect(screen, x, y, dest.w, dest.h);
}
//---------------------------------------------------------
static void schreibe (Sint16 x, Sint16 y, char *cs)
{
	gfxPrimitivesSetFont(NULL, 0 ,0);
	stringRGBA (screen, x,y,cs,red, green, blue, alpha);
	if (updateOn)
		SDL_UpdateRect (screen,0,0,0,0);
}
//---------------------------------------------------------
static void ladeBildInsClipboard (char *cs)
{
	SDL_Surface *image;
	image = SDL_LoadBMP(cs);
	//printf ("Dateiname: %s\n",cs);
	if (image == NULL) {
		// printf("Bild konnte nicht geladen werden!\n");
		return;
	}
	SDL_FreeSurface (clipboard); //altes Clipboard freigeben
	clipboard = SDL_DisplayFormat (image);
	SDL_FreeSurface (image);
}
//---------------------------------------------------------
static void clipboardKopieren (Sint16 x, Sint16 y, Uint16 b, Uint16 h) 
{
	SDL_Surface *image;
	SDL_Rect src, dest;
	Uint32 rmask, gmask, bmask, amask;
	if (clipboard != NULL)
		SDL_FreeSurface (clipboard);
	#if SDL_BYTEORDER == SDL_BIG_ENDIAN
		rmask = 0xff000000;
		gmask = 0x00ff0000;
		bmask = 0x0000ff00;
		amask = 0x000000ff;
	#else
		rmask = 0x000000ff;
		gmask = 0x0000ff00;
		bmask = 0x00ff0000;
		amask = 0xff000000;
	#endif
	image = SDL_CreateRGBSurface(SDL_HWSURFACE | SDL_SRCALPHA, (int) b, (int) h, 32, rmask, gmask, bmask, amask);
	if (image == NULL) {
		// printf("Neues Clipboard konnte nicht erzeugt werden!\n");
		return;
	}
	src.x = x;
	src.y = y;
	src.w = b;
	src.h = h;
	dest.x = 0;
	dest.y = 0;
	dest.w = b;
	dest.h = h; 
	SDL_BlitSurface(screen, &src, image, &dest);
	SDL_UpdateRect (image, 0, 0, 0, 0);
	clipboard = SDL_DisplayFormat (image);
	SDL_FreeSurface (image);
}
//---------------------------------------------------------
static void clipboardEinfuegen (Sint16 x, Sint16 y)
{
	SDL_Rect src, dest;
	src.x = 0;
	src.y = 0;
	src.w = clipboard->w;
	src.h = clipboard->h;
	dest.x = x;
	dest.y = y;
	dest.w = clipboard->w;
	dest.h = clipboard->h; 
	if (colorkey)
	{
		SDL_SetColorKey(clipboard, SDL_SRCCOLORKEY, SDL_MapRGBA(screen->format, ck_red, ck_green, ck_blue, alpha));
	} else {
		SDL_SetColorKey(clipboard, 0, SDL_MapRGBA(screen->format, ck_red, ck_green, ck_blue, alpha));
	}
	SDL_SetAlpha(clipboard, SDL_SRCALPHA, alpha); 
	SDL_BlitSurface(clipboard, &src, screen, &dest);
	if (updateOn)
		SDL_UpdateRect (screen, x, y, dest.w, dest.h);
}
//---------------------------------------------------------
static void archivieren ()
{
	SDL_Rect src, dest;
	src.x = 0;
	src.y = 0;
	src.w = screen->w;
	src.h = screen->h;
	dest = src;
	SDL_BlitSurface(screen, &src, archiv, &dest);
	SDL_UpdateRect(archiv, 0,0,0,0);
}
//----------------------------------------------------------
static void restaurieren (Sint16 x, Sint16 y, Uint16 b, Uint16 h)
{
	SDL_Rect src, dest;
	src.x = x;
	src.y = y;
	src.w = b;
	src.h = h;
	dest = src;
	SDL_BlitSurface(archiv, &src, screen, &dest);
	if (updateOn)
		SDL_UpdateRect (screen, x, y, b, h);
}
//---------------------------------------------------------------
static Uint32 tastaturLesen1 () 
{
	tasteLesen = 1;
	while (tasteLesen && fensteroffen) 
	{
		SDL_Delay (5);
	}
	return ((Uint32) taste)*256*256 + ((Uint32) gedrueckt)*256*256*256*128+ ((Uint32) tiefe);
}
//-------------------------------------------------------------
static void tastaturpufferAn () {
	t_pufferkopf = 0;
	t_pufferende = 0;
	tastaturpuffer = 1;
}
//-------------------------------------------------------------
static void tastaturpufferAus () {
	tastaturpuffer = 0;
}
//-------------------------------------------------------------
static Uint32 tastaturpufferLesen1 ()
{
	Uint32 erg;
	while (t_pufferende == t_pufferkopf && fensteroffen)
	{
		SDL_Delay (5);
	}
	erg = t_puffer[t_pufferkopf];
	t_pufferkopf++; //Überlauf von 255 auf 0 automatisch, da Uint8
	return erg;
}
//-------------------------------------------------------------
static Uint32 mausLesen1 ()
{
	mausLesen = 1;
	while (mausLesen && fensteroffen) 
	{
		SDL_Delay (5);
	}
	return ((Uint32) mausTaste)*256*256*256 + (((Uint32) mausX) << 12) + ((Uint32) mausY);
}
//--------------------------------------------------------------
static void mauspufferAn () {
	m_pufferkopf = 0;
	m_pufferende = 0;
	mauspuffer = 1;
}
//-------------------------------------------------------------
static void mauspufferAus () {
	mauspuffer = 0;
}
//-------------------------------------------------------------
static Uint32 mauspufferLesen1 ()
{
	Uint32 erg;
	while (m_pufferende == m_pufferkopf && fensteroffen)
	{
		SDL_Delay (5);
	}
	erg = m_puffer[m_pufferkopf];
	m_pufferkopf++; //Überlauf von 255 auf 0 automatisch, da Uint8
	return erg;
}
//-------------------------------------------------------------
*/
import "C"

import ( "time" ; "unsafe" ; "net" ; "fmt" ; "os" ; "strconv" ; "math" )

const (             // Konstanten für die Signalform der gespielten Töne
Sinusform uint8 = iota
Rechteckform
Dreieckform
Sägezahnform
)

var r uint32= 44100 // Abtastrate: 11025 oder 22050 oder 44100 - Standard 44100 Hz
var b uint8 = 2     // Auflösung:    1: 8 Bit ; 2: 16 Bit      - Standard 16 Bit
var k uint8 = 2     // Kanalanzahl:  1: mono  ; 2: stereo      - Standard stereo

var s = Rechteckform // aktuelle Signalform:                 - Standard Rechteck

var anschlagzeit  float64 = 0.002 // Standard auf 2 ms gesetzt
var abschwellzeit float64 = 0.750 // Standard auf 750 ms gesetzt
var haltepegel    float64 = 0     // Standard auf 0 % gesetzt ; 1 = 100 %
var ausklingzeit  float64 = 0.006 // Standard auf 6 ms gesetzt
var pulsweite     float64 = 0.375 // nur wichtig für Rechteck-Signale, hier: Prozentsatz HIGH (Pulsweite)

var tempo uint8 = 120 // "Schläge" = Viertelnoten pro Minute"
var tVollnote uint16 = 4 * uint16(((60 * 1000 + uint32(tempo)/2 )/uint32(tempo))) // Zeit einer vollen Note in ms

var frequenzen  = map[string]float64{ // Frequenzen der Noten - 7. Oktave
	"C" :2093.00, "C#":2217.46,
	"D" :2349.32, "D#":2489.02,
	"E" :2637.02,
	"F" :2793.83, "F#":2959.96,
	"G" :3135.96, "G#":3322.44,
	"A" :3520.00, "A#":3729.31,
	"H" :3951.07}

var grafikschloss   = make (chan int,1)
var tastaturschloss = make (chan int,1)
var mausschloss     = make (chan int,1)
var fensterbreite,fensterhoehe uint16
var serverprotokoll bool // Standard: false 
var serverLäuft bool

// Es gibt 4 Tastenbelegungen: Standard, SHIFT, ALT GR, ALT GR mit SHIFT.

var z1 [4]string = [4]string{ ",-.", ";_:",  "·–…", "×—÷"}
var z2 [4]string = [4]string{"0123456789", "=!\"§$%&/()", "}¹²³¼½¬{[]", "°¡⅛£¤⅜⅝⅞™±"}
var z3 [4]string = [4]string{"abcdefghijklmnopqrstuvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "æ“¢ð€đŋħ→̣ĸłµ”øþ@¶ſŧ↓„ł«»←",  "Æ‘©Ð€ªŊĦı˙&Łº’ØÞΩ®ẞŦ↑‚Ł‹›¥"}
var z4 [4]string = [4]string{",/*-+",",/*-+",",/*-+",",/*-+"} //NUM-Block
var z5 [4]string = [4]string{"0123456789","0123456789","0123456789","0123456789"} //NUM-BLOCK
var taste_belegung [4][320]rune //vier Belegungen pro Taste
//-----------------------------------------------------------------------------

func lock () { grafikschloss <- 1 }
func unlock () { <- grafikschloss }
func t_lock () { tastaturschloss <- 1 }
func t_unlock () { <- tastaturschloss }
func m_lock () { mausschloss <- 1 }
func m_unlock () { <- mausschloss }

func Fenster (breite, hoehe uint16) { // terminiert nur, wenn das Grafikfenster geschlossen ist!
	if fensterZu () {
		if breite > 1920 {breite = 1920}
		if hoehe > 1200 {hoehe = 1200}
		fensterhoehe = hoehe
		fensterbreite = breite
		C.GrafikfensterAn (C.Uint16(breite), C.Uint16(hoehe))
		//for !FensterOffen () {
		//	time.Sleep (100 * 1000 * 1000) //Unter Windows notwendig!!
		//}
	}
}

func FensterOffen () bool {
	return uint8(C.FensterOffen ()) == 1
}

func fensterZu () bool { //interne Hilfsfunktion
	return uint8(C.FensterZu ()) == 1
}

func FensterAus () {
	lock ()
	if FensterOffen () {
		C.GrafikfensterAus ()
		for !fensterZu () {
			time.Sleep (100 * 1000 * 1000) 
		}
	}
	unlock ()
}

func Grafikzeilen () uint16 {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	return fensterhoehe
}

func Grafikspalten () uint16 {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	return fensterbreite
}

func Fenstertitel (s string) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:=C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.Fenstertitel (cs)
	unlock ()
}

func Cls () {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.clearscreen ()
	unlock ()
}

func Stiftfarbe (r,g,b uint8) {
	lock ()
	C.stiftfarbe (C.Uint8 (r), C.Uint8 (g), C.Uint8 (b))
	unlock ()
}

func Transparenz (t uint8) {
	lock ()
	C.transparenz (C.Uint8(255-t))
	unlock ()
}

func Punkt (x,y uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.zeichnePunkt (C.Sint16(x), C.Sint16(y))
	unlock ()
}

func GibPunktfarbe (x,y uint16) (r,g,b uint8) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	pixel:=uint32 (C.gibPixel(C.Sint16(x),C.Sint16(y)))
	r = uint8(pixel >> 16)
	g = uint8(pixel >> 8)
	b = uint8(pixel)
	unlock ()
	return
}

func Linie (x1,y1,x2,y2 uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.zeichneStrecke (C.Sint16(x1),C.Sint16(y1),C.Sint16(x2),C.Sint16(y2))
	unlock ()
}

func Kreis (x,y,r uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.zeichneKreis(C.Sint16(x),C.Sint16(y),C.Sint16(r),0)
	unlock ()
}

func Vollkreis (x,y,r uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.zeichneKreis(C.Sint16(x),C.Sint16(y),C.Sint16(r),1)
	unlock ()
}

func Ellipse (x,y,rx,ry uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.zeichneEllipse(C.Sint16(x),C.Sint16(y),C.Sint16(rx),C.Sint16(ry),0)
	unlock ()
}

func Vollellipse (x,y,rx,ry uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.zeichneEllipse(C.Sint16(x),C.Sint16(y),C.Sint16(rx),C.Sint16(ry),1)
	unlock ()
}

func Vollkreissektor (x,y,r,w1,w2 uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.kreissektor (C.Sint16(x),C.Sint16(y),C.Sint16(r),360-C.Sint16(w2),360-C.Sint16(w1),1)
	unlock ()
}

func Kreissektor (x,y,r,w1,w2 uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.kreissektor (C.Sint16(x),C.Sint16(y),C.Sint16(r),360-C.Sint16(w2),360-C.Sint16(w1),0)
	unlock ()
}

func Rechteck (x1,y1,b,h uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.rechteck (C.Sint16(x1),C.Sint16(y1),C.Sint16(b),C.Sint16(h),0)
	unlock ()
}

func Vollrechteck (x1,y1,b,h uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.rechteck (C.Sint16(x1),C.Sint16(y1),C.Sint16(b),C.Sint16(h),1)
	unlock ()
}

func Dreieck (x1,y1,x2,y2,x3,y3 uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.dreieck (C.Sint16(x1),C.Sint16(y1),C.Sint16(x2),C.Sint16(y2),C.Sint16(x3),C.Sint16(y3),0)
	unlock ()
}

func Volldreieck (x1,y1,x2,y2,x3,y3 uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.dreieck (C.Sint16(x1),C.Sint16(y1),C.Sint16(x2),C.Sint16(y2),C.Sint16(x3),C.Sint16(y3),1)
	unlock ()
}

func Schreibe (x,y uint16, s string) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:= C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.schreibe (C.Sint16(x), C.Sint16 (y), cs)
	unlock ()
}

func SetzeFont (s string, groesse int) (erg bool) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:=C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	if int(C.setFont(cs,C.int(groesse)))==0 {
		erg = true 
	} else {
		erg = false
	}
	unlock()
	return
}

func GibFont () (erg string) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:=C.getFont ()
	//defer C.free(unsafe.Pointer(cs))
	erg = C.GoString(cs)
	unlock()
	return
}

func SchreibeFont (x,y uint16, s string) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:=C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	if int(C.write (C.Sint16(x),C.Sint16(y),cs)) == 1 {
		println ("FEHLER: Kein aktueller Font: ", C.GoString(C.getFont()))
	}
	unlock()
}
	
	
func LadeBild (x,y uint16, s string) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:= C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.setcolorkey (C.Uint8(0),C.Uint8(0),C.Uint8(0),C.Uint8(0))
	C.ladeBild (C.Sint16(x),C.Sint16(y),cs)
	unlock ()
}

func LadeBildMitColorKey (x,y uint16, s string,r,g,b uint8) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:= C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.setcolorkey (C.Uint8(r),C.Uint8(g),C.Uint8(b),C.Uint8(1))
	C.ladeBild (C.Sint16(x),C.Sint16(y),cs)
	unlock ()
}


func LadeBildInsClipboard (s string) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:= C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.ladeBildInsClipboard (cs)
	unlock ()
}

func Archivieren () {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.archivieren ()
	unlock ()
}

func Restaurieren (x1,y1,b,h uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.restaurieren (C.Sint16(x1),C.Sint16(y1),C.Uint16(b),C.Uint16(h))
	unlock ()
}

func Clipboard_kopieren (x,y,b,h uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.clipboardKopieren (C.Sint16(x), C.Sint16(y), C.Uint16(b), C.Uint16 (h))
	unlock ()
}

func Clipboard_einfuegen(x,y uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.setcolorkey (C.Uint8(0),C.Uint8(0),C.Uint8(0),C.Uint8(0))
	C.clipboardEinfuegen(C.Sint16(x), C.Sint16(y))
	unlock ()
}

func Clipboard_einfuegenMitColorKey (x,y uint16, r,g,b uint8) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.setcolorkey (C.Uint8(r),C.Uint8(g),C.Uint8(b),C.Uint8(1))
	C.clipboardEinfuegen(C.Sint16(x), C.Sint16(y))
	unlock ()
}

// Sperren und Entsperren in gfx2impl.go

func UpdateAus () {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.updateAus ()
	unlock ()
}

func UpdateAn () {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.updateAn ()
	unlock ()
}


func TastaturLesen1 () (taste uint16, gedrueckt uint8, tiefe uint16) {
	var tastenwert uint32
	t_lock ()
	tastenwert = uint32(C.tastaturLesen1 ())
	t_unlock ()
	tiefe = uint16 (tastenwert % 65536)
	tastenwert = tastenwert >> 16
	gedrueckt = uint8(tastenwert >> 15)
	taste = uint16(tastenwert % 32768) //oberstes Bit rausschieben
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	return
}

func Tastaturzeichen (taste, tiefe uint16) rune {
	switch tiefe {
		case 0,4096, 8192+1, 8192+2, 8192+3,4096+8192+1,4096+8192+2,4096+8192+3:	// kein SHIFT, kein ALT GR, NUMLOCK an oder aus, CAPSLOCK an mit SHIFT
			return taste_belegung[0][taste]
		case 1,2,3,4096+1, 4096+2, 4096+3, 8192, 4096+8192:  // SHIFT, kein ALT GR, NUMLOCK an oder aus, CAPSLOCK an ohne SHIFT
			return taste_belegung[1][taste]
		case 16384, 16384 + 4096, 16384+8192+1,16384+8192+2,16384+8192+3,16384+8192+4096+1,16384+8192+4096+2,16384+8192+4096+3:  // kein SHIFT, ALT GR, NUMLOCK an oder aus, CAPSLOCK an mit SHIFT
			return taste_belegung[2][taste]
		case 16384+1, 16384+2, 16384+3, 16384+4096+1, 16384+4096+2, 16384+4096+3, 16384+8192, 16384+8192+4096: // ALT GR und SHIFT, NUMLOCK an oder aus, CAPSLOCK an ohne SHIFT
			return taste_belegung[3][taste]
		default:
		return 0
	}
}

func TastaturpufferAn () {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	t_lock ()
	C.tastaturpufferAn ()
	t_unlock ()
}

func TastaturpufferAus () {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	t_lock ()
	C.tastaturpufferAus ()
	t_unlock ()
}

func TastaturpufferLesen1 () (taste uint16, gedrueckt uint8, tiefe uint16) {
	var tastenwert uint32
	t_lock ()
	tastenwert = uint32(C.tastaturpufferLesen1 ())
	t_unlock ()
	tiefe = uint16 (tastenwert % 65536)
	tastenwert = tastenwert >> 16
	gedrueckt = uint8(tastenwert >> 15)
	taste = uint16(tastenwert % 32768)
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	return
}


func MausLesen1 () (taste uint8,status int8, mausX, mausY uint16) {
	var tastenwert uint32
	m_lock ()
	tastenwert = uint32(C.mausLesen1 ())
	m_unlock ()
	taste = uint8 (tastenwert >> 24)
	if taste < 64 {
		status=0   //Zustand wird gehalten
	} else if taste > 128 {
		status = 1 //gerade gedrückt
		taste = taste - 128
	} else  {//zwischen 64 und 128
		status = -1 //gerade losgelassen
		taste = taste-64
	} 
	mausY = uint16 (tastenwert % 4096)
	tastenwert = tastenwert >> 12
	mausX = uint16 (tastenwert % 4096)
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	return
} 

func MauspufferAn () {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	m_lock()
	C.mauspufferAn ()
	m_unlock ()
}

func MauspufferAus (){
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	m_lock ()
	C.mauspufferAus ()
	m_unlock ()
}

func MauspufferLesen1 () (taste uint8,status int8, mausX, mausY uint16) {
	var tastenwert uint32
	m_lock ()
	tastenwert = uint32(C.mauspufferLesen1 ())
	m_unlock ()
	taste = uint8 (tastenwert >> 24)
	if taste < 64 {
		status=0   //Zustand wird gehalten
	} else if taste > 128 {
		status = 1 //gerade gedrückt
		taste = taste - 128
	} else  {//zwischen 64 und 128
		status = -1 //gerade losgelassen
		taste = taste-64
	} 
	mausY = uint16 (tastenwert % 4096)
	tastenwert = tastenwert >> 12
	mausX = uint16 (tastenwert % 4096)
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	return
} 

func SpieleSound(s string) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:=C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	erg:= int(C.LoadAndPlaySound(cs))
	if erg == 1 {
		println("Es werden schon die max. Anzahl an Sounds abgespielt!")
	}
	if erg == 2 {
		println ("Konnte Sounddatei nicht laden! --> ", s)
	}
	unlock()
}

// Erg.: Das aktuelle Tempo ist geliefert, d.h. die Anzahl der Viertelnoten pro Minute.
func GibNotenTempo () uint8 { 
	return tempo 
}

// Vor.: 30 <= tempo <= 240 
// Eff.: Das Tempo ist auf den Wert t gesetzt.
func SetzeNotenTempo (t uint8) {
	lock ()
	if t >= 30  && t <= 240 {
		tempo = t
		tVollnote = 4 * uint16(((60 * 1000 + uint32(tempo)/2 )/uint32(tempo)))
	}
	unlock ()
}
		
// Vor.: -
// Erg.: Geliefert sind:
//       Abtastrate der WAV-Daten in Hz, z.B. 44100,
//       Auflösung der Klänge (1: 8 Bit; 2: 16 Bit),
//       die Anzahl der Kanäle (1: mono, 2:stereo),
//       die Signalform (0: Sinus, 1: Rechteck, 2:Dreieck, 3: Sägezahn) und
//       die Pulsweite HIGH bei Rechteckform als Prozentsatz zw. 0 und 1.
func GibKlangparameter () (uint32,uint8,uint8,uint8,float64) {
	return r,b,k,s,pulsweite
}
		
// Vor.: rate ist die Abtastrate, z.B. 11025, 22050 oder 44100.
//       auflösung ist 1 für 8 Bit oder 2 für 16 Bit.
//       kanaele ist 1 für mono oder 2 für stereo.
//       signal gibt die Signalform an: 0: Sinus, 1: Rechteck, 2:Dreieck, 3: Sägezahn
//       p ist die Pulsweite für Rechtecksignale und gibt den Prozentsatz (0<=p<=1) für den HIGH-Teil an.
// Eff.: Die klangparameter sind auf die angegebenen Werte gesetzt.
func SetzeKlangparameter(rate uint32, aufloesung,kanaele,signal uint8, p float64) {
	lock ()
	r = rate
	b = aufloesung
	k = kanaele
	s = signal
	pulsweite = p
	unlock ()
}
		
// Vor.: -
// Ergebnis: Anschlagzeit, Abschwellzeit, Haltepegel und Ausklingzeit sind geliefert.
func GibHuellkurve () (float64,float64,float64,float64) {
	return anschlagzeit,abschwellzeit,haltepegel,ausklingzeit
}

// Vor.: a ist die Anschlagzeit in s mit 0 <= a <= 1,
//       d ist die Abschwellzeit in s mit 0<= d <= 5,
//       s ist der Haltepegel in Prozent vom Maximum mit 0<= s <= 1.0,
//       r ist die Ausklingzeit in s mit 0< =r <= 5.
// Eff.: Für die Hüllkurve zukünftig zu spielender Töne bzw. Noten sind
//       die Parameter gesetzt.
func SetzeHuellkurve (a,d,s,r float64) {
	lock ()
	if 0<=a && a<= 1 && 0<=d && d<=5 && 0<=s && s<=1 && 0<=r && r<=5 {
		anschlagzeit  = a
		abschwellzeit = d
		haltepegel    = s
		ausklingzeit  = r
	}
	unlock ()
}


// INTERN
// Erg.: Ein Slice bestehend aus 4 Bytes ist geliefert, die den Wert x 
//       darstellen. Das erste Byte ist das LSB.
func uint32toSlice (x uint32) []byte {
	var erg []byte = make ([]byte,4)
	for i:=0; i<4;i++ {
		erg[i] = byte(x % 256)
		x = x / 256
	}
	return erg
}

// INTERN
// Vor.: t ist der echte Zeitpunkt innerhalb des Tons.
// Erg.: Der Maximalausschlag zwischen 0 und 1 ist gemäß
//       der aktuell festgelegten Hüllkurve geliefert. 
func amplitude (t,tges float64) float64 {
	switch {
		case t <= anschlagzeit:
		return t/anschlagzeit
		case t > tges-ausklingzeit:
		return haltepegel-(t-(tges-ausklingzeit))*haltepegel/ausklingzeit
		default: //Abschwell- und Haltepegelzeit
		return haltepegel+(1-haltepegel)*math.Pow(2,-(t-anschlagzeit)*6/abschwellzeit)
	}
}

// INTERN
// Vor.: tges gibt die Tondauer in Millisekunden an.
//       f ist die Tonfrequenz in Hertz.
// Erg.: Ein Byte-Slice ist geliefert, dass der entsprechenden WAV-Datei entspricht.
func ton (tges uint16, f float64) []byte {
	var laenge uint32           = r*uint32(tges)*uint32(b*k)/1000
	var dateigrößeMinus8 uint32 = laenge + 44 - 8 
	var bytes []byte            = make ([]byte,44 + laenge)
	var w float64
	
	// "Dateikopf gemäß RIFF-WAVE-Format
	copy (bytes,"RIFF")
	copy (bytes[4:], uint32toSlice(dateigrößeMinus8)) //DATEIGRÖSSE - 8 ----------------------------------------------------
	copy (bytes[8:],"WAVEfmt ")
	bytes[16] = 16 // Die Größe des fmt-Abschnitts ist 16 Bytes (uint32)
	bytes[17] = 0
	bytes[18] = 0
	bytes[19] = 0
	bytes[20] = 1 // Das verwendete Format: 01 = PCM (uint16)
	bytes[21] = 0
	bytes[22] = k // Wir verwenden k Kanal (1:mono; 2:stereo).
	bytes[23] = 0
	copy (bytes[24:],uint32toSlice(r)) // Eintrag der Abtastrate
	copy (bytes[28:],uint32toSlice(r * uint32(b*k))) // Übertragungsbandbreite (Bytes pro Sekunde): rate*b*k
	bytes[32] = k   // uint16 - 1: mono ; 2: stereo
	bytes[33] = 0
	bytes[34] = b*8 // uint16 - Auflösung: 8 oder 16 Bit
	bytes[35] = 0
	copy(bytes[36:],"data")
	copy(bytes[40:], uint32toSlice(laenge))//DATEIGRÖSSE - 44----------------------------------------------------
	// Es folgen die Daten - ein Frame = b*k Byte
	// Es sind rate Frames pro Sekunde
	for i:=uint32(0);i<laenge-uint32(b*k)+1;i=i+uint32(b*k) {
		t:= float64(i)/float64(r*uint32(b*k)) // echter Zeitpunkt 
		t2:= t-float64(uint64(t*f))/f         // Zeitpunkt innerhalb der aktuellen Schwingung 
		switch s { // nach Signalform
			case Sinusform:
			w= math.Sin(2*math.Pi*f*t2) // float64 aus [-1;1]
			case Rechteckform:
			if t2 <= pulsweite/f {
				w = 1
			} else {
				w = -1
			}
			case Dreieckform:
			if t2 <= 1/(2*f) {
				w = -1 + 4*f*t2
			} else {
				w = 1 - 4 * f * (t2-1/(2*f))
			}
			case Sägezahnform:
			w=-1+2*f*t2
			default:
			panic ("unbekannte Signalform!!")
		}
		w = amplitude(t,float64(tges)/1000) * w  // Einarbeiten der Hüllkurve
		switch b {
			case 1:
			bytes[44+i] = uint8 (128 + w * 127)
			if k == 2 { 
				bytes[45+i] = bytes[44+i]
			}
			case 2:
			bytes[44+i] = byte(uint16(w*32767) % 256)
			bytes[45+i] = byte(uint16(w*32767) / 256)
			if k == 2 {
				bytes[46+i] = bytes[44+i]
				bytes[47+i] = bytes[45+i]
			}
		}
	}	
	return bytes
}

// INTERN
// Vor.: data stellt die Bytefolge einer WAVE-Datei dar.
//       wartezeit ist die Abwartezeit nach dem Anspielen der WAV-Datei in ms.
// Eff.: Die 'WAV-Datei' wird bzw. ist gerade abgespielt. Der Programmablauf
//       ist dafür um wartezeit ms verzögert worden.
func spieleRAMWAV (data []byte,wartezeit uint32) {
	erg:= int (C.LadeUndSpieleNote(unsafe.Pointer(&data[0]),C.int(len(data))))
	if erg == 1 {
		println("Es werden schon die max. Anzahl an Sounds abgespielt!")
	}
	if erg == 2 {
		println ("Die Daten entsprechen keiner WAV-Datei! Daten nicht geladen!")
	}
	time.Sleep (time.Duration(int64(wartezeit)*1e6))
}

func SpieleNote (tonname string, laenge float64, wartedauer float64) {
	var o uint8 = byte(tonname[0])-48
	var freq float64 =frequenzen[tonname[1:]]
	lock ()
	for i:=uint8(7);i>o;i--{ freq = freq / 2 }
	for i:=uint8(7);i<o;i++{ freq = freq * 2 }  
	bytes:= ton(uint16(float64(tVollnote)*laenge),freq)
	spieleRAMWAV(bytes,uint32(wartedauer*float64(tVollnote)))
	unlock ()
}


func init_Tastatur_Deutsch ()  {
	// Es folgt die Initalisierung der Tastaturbelegung auf Deutsch.
	// Das wird für die Funktion 'Tastaturzeichen(taste, tiefe) rune' benötigt.
	for i:=0; i < 4; i++ {
		index:= 0; for _,e:= range z1[i] {taste_belegung[i][index+44] = e; index++}
		index = 0; for _,e:= range z2[i] {taste_belegung[i][index+48] = e; index++}
		index = 0; for _,e:= range z5[i] {taste_belegung[i][index+256] = e; index++} //Num-Block
		index = 0; for _,e:= range z4[i] {taste_belegung[i][index+266] = e; index++} //Num-Block
		index = 0; for _,e:= range z3[i] {taste_belegung[i][index+97] = e; index++}
	}
	// kein SHIFT, kein ALT GR
	taste_belegung[0][43]='+'  ; taste_belegung[0][35] ='#' ; taste_belegung[0][252]='ü'
	taste_belegung[0][246]='ö' ; taste_belegung[0][228]='ä' ; taste_belegung[0][223]='ß'
	taste_belegung[0][180]='´' ; taste_belegung[0][94] ='^' ; taste_belegung[0][60] ='<'
	taste_belegung[0][32]=' '
	// SHIFT, kein ALT GR
	taste_belegung[1][43] ='*' ; taste_belegung[1][35] ='\'' ; taste_belegung[1][252]='Ü'
	taste_belegung[1][246]='Ö' ; taste_belegung[1][228]='Ä'  ; taste_belegung[1][223]='?'
	taste_belegung[1][180]='`' ; taste_belegung[1][94] ='°'  ; taste_belegung[1][60]='>'
	taste_belegung[1][32]=' '
	// kein SHIFT, ALT GR
	taste_belegung[2][43] ='~' ; taste_belegung[2][35] ='`' ; taste_belegung[2][252]='¨'
	taste_belegung[2][246]='˝' ; taste_belegung[2][228]='^' ; taste_belegung[2][223]='\\'
	taste_belegung[2][180]='¸' ; taste_belegung[2][94] ='¬' ; taste_belegung[2][60] ='|'
	taste_belegung[2][32] =' '
	// SHIFT, ALT GR
	taste_belegung[3][43] ='¯' ; taste_belegung[3][35] ='`' ; taste_belegung[3][252]='¨'
	taste_belegung[3][246]='˝' ; taste_belegung[3][228]='^' ; taste_belegung[3][223]='¿'
	taste_belegung[3][180]='¸' ; taste_belegung[3][94] ='¬' ; taste_belegung[3][60] ='¦'
	taste_belegung[3][32] =' '
}

// Es folgen die Netzfunktionen.
func setzeServerprotokoll (wert bool) {
	serverprotokoll = wert
}

func starteGfxServer (quellIP string, portnummer uint16, f func (string) string) {
	server,err := net.Listen ("tcp",quellIP + ":" + fmt.Sprint(portnummer))
	if err != nil {
		fmt.Println (err)
		panic ("Gfx-Server konnte nicht gestartet werden!")
	}
	serverLäuft = true
	if serverprotokoll {
		fmt.Println ("gfx-Server wurde gestartet und wartet auf Kontaktaufnahme...")
		fmt.Println ("============================================================")
	}
	//ÖFFNE DIE DREI KANÄLE ZUM PROGRAMM!
	conn, err := server.Accept () // Standardanfragen
	if err != nil { panic ("Eine Verbindung ist fehlgeschlagen!") }
	if serverprotokoll { fmt.Println ("Verbunden mit:",conn.RemoteAddr ()) }
	connMaus, err2 := server.Accept () // Mausanfragen
	if err2 != nil { panic ("Verbindung für Mausanfragen ist fehlgeschlagen!") }
	if serverprotokoll { fmt.Println ("Maus-Verbindung über:",conn.RemoteAddr ()) }
	connTast, err3 := server.Accept () // Tastaturanfragen
	if err3 != nil { panic ("Verbindung für Tastaturanfragen ist fehlgeschlagen!") }
	if serverprotokoll { fmt.Println ("Tastatur-Verbunden über:",conn.RemoteAddr ()) }
	for !FensterOffen() { time.Sleep (1e8) } //Solange das Fenster noch nicht offen ist, warte ...
	go func () {
		for FensterOffen() {
			bearbeiteAnfrage (connMaus,f)
		}
	} ()
	go func () { 
		for FensterOffen() {
			bearbeiteAnfrage (connTast,f)
		}
	} ()
	for FensterOffen() { // Solange das Fenster offen ist, bearbeite Anfragen ...
		bearbeiteAnfrage (conn, f)
	} 
	// Wenn das Fenster nun zu ist ...
	//WIRD DER KANAL GESCHLOSSEN!
	conn.Close ()
	connMaus.Close()
	connTast.Close ()
	serverLäuft = false  // Hiermit wird das Hauptprogramm 'main' und damit auch der Server beendet!
}



// interne Funktion, die für jede Anfrage nebenläufig ausgeführt wird
func bearbeiteAnfrage (conn net.Conn, f func (string) string) {
	var l []byte = make ([]byte,4)     //Byte-Slice für die Länge der Nachricht
	var b []byte = make ([]byte,1024)  //Lese-Puffer
	var erwartet, angekommen, laenge int32
	n,err := conn.Read(l)
	if n < 4 || err != nil {
		if err.Error() == "EOF" || err.Error()[0:4]=="read" { // Die Gegenseite hat den Kanal geschlossen!
			FensterAus()
			return
		} else {
			fmt.Println("PANIK:",err)
			panic ("Fehler beim Empfangen einer Nachricht!")
		}
	}
	for i:=0; i < 4; i++ {
		erwartet = erwartet*256 + int32(l[3-i])
	}
	nachricht := make ([]byte, erwartet)
	for angekommen < erwartet {
		n,err = conn.Read(b)
		if err != nil {
			panic ("Fehler beim Empfangen einer Nachricht!")
		}
		copy(nachricht[angekommen:],b[0:n])
		angekommen = angekommen + int32(n)
	}
	if serverprotokoll {
		fmt.Print ("Empfangene Nachricht:")
		fmt.Println (string('"') + string(nachricht) + string ('"'))
	}
	//in -nachricht- ist die Anfrage
	//
	//NUN WIRD DIE ANTWORT GENERIERT ...
	antwort:= f(string(nachricht))
	//NUN WIRD ZURÜCK GESENDET ...
	laenge = int32(len(antwort))
	nachricht = make ([]byte,laenge)
	for i:=0; i < 4; i++ {
		l[i] = byte(laenge % 256)
		laenge = laenge / 256
	}
	copy(nachricht,antwort)
	nachricht = append (l,nachricht...) //4-Byte-Präfix gibt die Länge der Nachricht an
	n,err = conn.Write(nachricht)
	if err != nil {
		panic ("Übertragungsfehler bzgl. der Verbindung! beim Senden")
	}
	if n != len(nachricht) {
		panic ("Es sind Bytes beim Senden verloren gegangen!")
	}
	if serverprotokoll {
		fmt.Print ("Versende Antwort an ",conn.RemoteAddr(),":")
		fmt.Println (string ('"') + string(antwort) + string('"'))
		fmt.Println ("=======================================================")
	}
}

func split (text string) []string {
	var erg []string = make ([]string,0)
	var teil string = ""
	for _,z := range text {
		if z == ':' {
			erg = append(erg,teil)
			teil = ""
		} else {
			teil = teil + string(z)
		}
	}
	erg = append (erg, teil)
	return erg
}
		
// f ist die Funktion, die den Anfragestring der Netzwerkanfrage bekommt und eine
// Antwort generiert.
func f (anfrage string) string {
	a:= split(anfrage)
	switch a[0] {
		case "FEAU":
		if FensterOffen () {
			FensterAus()
		}
		return "OK"
		case "FEOF":
		if FensterOffen () {
			return "true"
		} else {
			return "false"
		}
		case "GRZE":
		if len (a) != 1 {return "ERROR1" }
		return fmt.Sprint(Grafikzeilen())
		case "GRSP":
		if len (a) != 1 {return "ERROR1" }
		return fmt.Sprint(Grafikspalten())
		case "FETI":
		if len(a) != 2 {return "ERROR1" }
		Fenstertitel(a[1])
		return "OK"
		case "CLSC":
		Cls ()
		return "OK"
		case "STFA":
		if len(a) != 4 { return "ERROR1" }
		rot,err:=strconv.Atoi(a[1])
		green,err2:= strconv.Atoi(a[2])
		blue,err3 := strconv.Atoi(a[3])
		if err != nil || err2 != nil || err3 != nil {return "ERROR2"}
		r,g,b:=uint8(rot),uint8(green),uint8(blue)
		Stiftfarbe(r,g,b)
		return "OK"
		case "TRAN":
		if len (a) != 2 { return "ERROR1" }
		t,err:=strconv.Atoi(a[1])
		if err != nil {return "ERROR2"}
		tr := uint8(t)
		Transparenz(tr)
		return "OK"
		case "PNKT":
		if len (a) != 3 {return "ERROR1" }
		xk,err:=strconv.Atoi(a[1])
		yk,err2:=strconv.Atoi(a[2])
		if err != nil || err2 != nil {return "ERROR2"}
		x,y:= uint16(xk),uint16(yk)
		Punkt(x,y)
		return "OK"
		case "GPTF":
		if len(a) != 3 {return "ERROR1" }
		xk,err:=strconv.Atoi(a[1])
		yk,err2:=strconv.Atoi(a[2])
		if err != nil || err2 != nil {return "ERROR2"}
		x,y:= uint16(xk),uint16(yk)
		r,g,b:= GibPunktfarbe (x,y)
		return fmt.Sprint(r)+":"+fmt.Sprint(g)+":"+fmt.Sprint(b)
		case "LINE":
		if len (a) != 5 {return "ERROR1" }
		x1k,err:=strconv.Atoi(a[1])
		y1k,err2:=strconv.Atoi(a[2])
		x2k,err3:=strconv.Atoi(a[3])
		y2k,err4:=strconv.Atoi(a[4])
		if err != nil || err2 != nil || err3 != nil || err4 != nil {return "ERROR2"}
		x1,y1,x2,y2:= uint16(x1k),uint16(y1k),uint16(x2k),uint16(y2k)
		Linie(x1,y1,x2,y2)
		return "OK"
		case "KREI":
		if len (a) != 4 {return "ERROR1" }
		xk,err:=strconv.Atoi(a[1])
		yk,err2:=strconv.Atoi(a[2])
		ra,err3:=strconv.Atoi(a[3])
		if err != nil || err2 != nil || err3 != nil {return "ERROR2"}
		x,y,r:= uint16(xk),uint16(yk),uint16(ra)
		Kreis(x,y,r)
		return "OK"
		case "VOKR":
		if len (a) != 4 {return "ERROR1" }
		xk,err:=strconv.Atoi(a[1])
		yk,err2:=strconv.Atoi(a[2])
		ra,err3:=strconv.Atoi(a[3])
		if err != nil || err2 != nil || err3 != nil {return "ERROR2"}
		x,y,r:= uint16(xk),uint16(yk),uint16(ra)
		Vollkreis(x,y,r)
		return "OK"
		case "ELLI":
		if len (a) != 5 {return "ERROR1" }
		xk,err:=strconv.Atoi(a[1])
		yk,err2:=strconv.Atoi(a[2])
		rxk,err3:=strconv.Atoi(a[3])
		ryk,err4:=strconv.Atoi(a[4])
		if err != nil || err2 != nil || err3 != nil || err4 != nil {return "ERROR2"}
		x,y,rx,ry:= uint16(xk),uint16(yk),uint16(rxk),uint16(ryk)
		Ellipse(x,y,rx,ry)
		return "OK"
		case "VOEL":
		if len (a) != 5 {return "ERROR1" }
		xk,err:=strconv.Atoi(a[1])
		yk,err2:=strconv.Atoi(a[2])
		rxk,err3:=strconv.Atoi(a[3])
		ryk,err4:=strconv.Atoi(a[4])
		if err != nil || err2 != nil || err3 != nil || err4 != nil {return "ERROR2"}
		x,y,rx,ry:= uint16(xk),uint16(yk),uint16(rxk),uint16(ryk)
		Vollellipse(x,y,rx,ry)
		return "OK"
		case "KRSE":
		if len (a) != 6 {return "ERROR1" }
		xk,err:=strconv.Atoi(a[1])
		yk,err2:=strconv.Atoi(a[2])
		ra,err3:=strconv.Atoi(a[3])
		wi1,err4:=strconv.Atoi(a[4])
		wi2,err5:=strconv.Atoi(a[5])
		if err != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {return "ERROR2"}
		x,y,r,w1,w2:= uint16(xk),uint16(yk),uint16(ra),uint16(wi1),uint16(wi2)
		Kreissektor(x,y,r,w1,w2)
		return "OK"
		case "VKSE":
		if len (a) != 6 {return "ERROR1" }
		xk,err:=strconv.Atoi(a[1])
		yk,err2:=strconv.Atoi(a[2])
		ra,err3:=strconv.Atoi(a[3])
		wi1,err4:=strconv.Atoi(a[4])
		wi2,err5:=strconv.Atoi(a[5])
		if err != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {return "ERROR2"}
		x,y,r,w1,w2:= uint16(xk),uint16(yk),uint16(ra),uint16(wi1),uint16(wi2)
		Vollkreissektor(x,y,r,w1,w2)
		return "OK"
		case "RECH":
		if len (a) != 5 {return "ERROR1" }
		x1k,err:=strconv.Atoi(a[1])
		y1k,err2:=strconv.Atoi(a[2])
		b1,err3:=strconv.Atoi(a[3])
		h1,err4:=strconv.Atoi(a[4])
		if err != nil || err2 != nil || err3 != nil || err4 != nil {return "ERROR2"}
		x1,y1,b,h:= uint16(x1k),uint16(y1k),uint16(b1),uint16(h1)
		Rechteck(x1,y1,b,h)
		return "OK"
		case "VORE":
		if len (a) != 5 {return "ERROR1" }
		x1k,err:=strconv.Atoi(a[1])
		y1k,err2:=strconv.Atoi(a[2])
		b1,err3:=strconv.Atoi(a[3])
		h1,err4:=strconv.Atoi(a[4])
		if err != nil || err2 != nil || err3 != nil || err4 != nil {return "ERROR2"}
		x1,y1,b,h:= uint16(x1k),uint16(y1k),uint16(b1),uint16(h1)
		Vollrechteck(x1,y1,b,h)
		return "OK"
		case "DREI":
		if len (a) != 7 {return "ERROR1" }
		x1k,err:=strconv.Atoi(a[1])
		y1k,err2:=strconv.Atoi(a[2])
		x2k,err3:=strconv.Atoi(a[3])
		y2k,err4:=strconv.Atoi(a[4])
		x3k,err5:=strconv.Atoi(a[5])
		y3k,err6:=strconv.Atoi(a[6])
		if err != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {return "ERROR2"}
		x1,y1,x2,y2,x3,y3:= uint16(x1k),uint16(y1k),uint16(x2k),uint16(y2k),uint16(x3k),uint16(y3k)
		Dreieck(x1,y1,x2,y2,x3,y3)
		return "OK"
		case "VODR":
		if len (a) != 7 {return "ERROR1" }
		x1k,err:=strconv.Atoi(a[1])
		y1k,err2:=strconv.Atoi(a[2])
		x2k,err3:=strconv.Atoi(a[3])
		y2k,err4:=strconv.Atoi(a[4])
		x3k,err5:=strconv.Atoi(a[5])
		y3k,err6:=strconv.Atoi(a[6])
		if err != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {return "ERROR2"}
		x1,y1,x2,y2,x3,y3:= uint16(x1k),uint16(y1k),uint16(x2k),uint16(y2k),uint16(x3k),uint16(y3k)
		Volldreieck(x1,y1,x2,y2,x3,y3)
		return "OK"
		case "SCHR":
		if len (a) < 4 {return "ERROR1" }
		xk,err:=strconv.Atoi(a[1])
		yk,err2:=strconv.Atoi(a[2])
		if err != nil || err2 != nil {return "ERROR2"}
		text:=a[3]
		for i:=4;i<len(a);i++ { text=text+":"+a[i]}
		x,y:= uint16(xk),uint16(yk)
		Schreibe(x,y,text)
		return "OK"
		case "SEFO":
		if len (a) < 3 {return "ERROR1" }
		gr,err:=strconv.Atoi(a[1])
		if err != nil {return "ERROR2"}
		s:=a[2]
		for i:=3;i<len(a);i++ { s=s+":"+a[i]}
		g:= int(gr)
		if SetzeFont(s,g) {
			return "true"
		} else {
			return "false"
		}
		case "GIFO":
		return GibFont ()
		case "SCFO":
		if len (a) < 4 {return "ERROR1" }
		xk,err:=strconv.Atoi(a[1])
		yk,err2:=strconv.Atoi(a[2])
		if err != nil || err2 != nil {return "ERROR2"}
		s:=a[3]
		for i:=4;i<len(a);i++ { s=s+":"+a[i]}
		x,y:= uint16(xk),uint16(yk)
		SchreibeFont(x,y,s)
		return "OK"
		case "LABI":
		if len (a) < 4 {return "ERROR1" }
		xk,err:=strconv.Atoi(a[1])
		yk,err2:=strconv.Atoi(a[2])
		if err != nil || err2 != nil {return "ERROR2"}
		s:=a[3]
		for i:=4;i<len(a);i++ { s=s+":"+a[i]}
		x,y:= uint16(xk),uint16(yk)
		LadeBild(x,y,s)
		return "OK"
		case "LBMC":
		if len (a) < 7 {return "ERROR1" }
		xk,err:=strconv.Atoi(a[1])
		yk,err2:=strconv.Atoi(a[2])
		if err != nil || err2 != nil {return "ERROR2"}
		s:=a[3]
		for i:=4;i<len(a)-3;i++ { s=s+":"+a[i]}
		red,err:= strconv.Atoi(a[len(a)-3])
		green,err2:= strconv.Atoi(a[len(a)-2])
		blue,err3:= strconv.Atoi(a[len(a)-1])
		if err != nil || err2 != nil || err3 != nil {return "ERROR2"}
		x,y := uint16(xk), uint16(yk)
		r,g,b:= uint8(red),uint8(green),uint8(blue)
		LadeBildMitColorKey(x,y,s,r,g,b)
		return "OK"
		case "LBIC":
		if len (a) < 2 {return "ERROR1" }
		name:=a[1]
		for i:=2;i<len(a);i++ { name=name+":"+a[i]}
		LadeBildInsClipboard(name)
		return "OK"
		case "ARCH":
		if len (a) != 1 {return "ERROR1"}
		Archivieren()
		return "OK"
		case "REST":
		if len (a) != 5 {return "ERROR1" }
		x1k,err:=strconv.Atoi(a[1])
		y1k,err2:=strconv.Atoi(a[2])
		b1,err3:=strconv.Atoi(a[3])
		h1,err4:=strconv.Atoi(a[4])
		if err != nil || err2 != nil || err3 != nil || err4 != nil {return "ERROR2"}
		x1,y1,b,h:= uint16(x1k),uint16(y1k),uint16(b1),uint16(h1)
		Restaurieren(x1,y1,b,h)
		return "OK"
		case "CLKO":
		if len (a) != 5 {return "ERROR1" }
		x1k,err:=strconv.Atoi(a[1])
		y1k,err2:=strconv.Atoi(a[2])
		b1,err3:=strconv.Atoi(a[3])
		h1,err4:=strconv.Atoi(a[4])
		if err != nil || err2 != nil || err3 != nil || err4 != nil {return "ERROR2"}
		x1,y1,b,h:= uint16(x1k),uint16(y1k),uint16(b1),uint16(h1)
		Clipboard_kopieren(x1,y1,b,h)
		return "OK"
		case "CLEI":
		if len (a) != 3 {return "ERROR1" }
		x1k,err:=strconv.Atoi(a[1])
		y1k,err2:=strconv.Atoi(a[2])
		if err != nil || err2 != nil {return "ERROR2"}
		x1,y1:= uint16(x1k),uint16(y1k)
		Clipboard_einfuegen(x1,y1)
		return "OK"
		case "CEMC":
		if len (a) != 6 {return "ERROR1" }
		x1k,err:=strconv.Atoi(a[1])
		y1k,err2:=strconv.Atoi(a[2])
		red,err3:=strconv.Atoi(a[3])
		green,err4:=strconv.Atoi(a[4])
		blue,err5:=strconv.Atoi(a[5])
		if err != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {return "ERROR2"}
		x1,y1,r,g,b:= uint16(x1k),uint16(y1k),uint8(red),uint8(green),uint8(blue)
		Clipboard_einfuegenMitColorKey(x1,y1,r,g,b)
		return "OK"
		case "UPAU":
		if len (a) != 1 {return "ERROR1" }
		UpdateAus()
		return "OK"
		case "UPAN":
		if len (a) != 1 {return "ERROR1" }
		UpdateAn()
		return "OK"
		case "TAL1":
		if len (a) != 1 {return "ERROR1" }
		taste,gedrueckt,tiefe:= TastaturLesen1()
		return fmt.Sprint(taste)+":"+fmt.Sprint(gedrueckt)+":"+fmt.Sprint(tiefe)
		case "TAZE":
		if len(a) != 3 {return "ERROR1" }
		ta,err:=strconv.Atoi(a[1])
		ti,err2:=strconv.Atoi(a[2])
		if err != nil || err2 != nil {return "ERROR2"}
		taste,tiefe:= uint16(ta),uint16(ti)
		zeichen:= Tastaturzeichen (taste,tiefe)
		return fmt.Sprint(zeichen)
		case "TPAN":
		if len (a) != 1 {return "ERROR1" }
		TastaturpufferAn()
		return "OK"
		case "TPAU":
		if len (a) != 1 {return "ERROR1" }
		TastaturpufferAus()
		return "OK"
		case "TPL1":
		if len (a) != 1 {return "ERROR1" }
		taste,gedrueckt,tiefe:= TastaturpufferLesen1 ()
		return fmt.Sprint(taste)+":"+fmt.Sprint(gedrueckt)+":"+fmt.Sprint(tiefe)
		case "MAL1":
		if len (a) != 1 {return "ERROR1" }
		taste,status,mausX,mausY:= MausLesen1 ()
		return fmt.Sprint(taste)+":"+fmt.Sprint(status)+":"+fmt.Sprint(mausX)+":"+fmt.Sprint(mausY)
		case "MPAN":
		if len (a) != 1 {return "ERROR1" }
		MauspufferAn()
		return "OK"
		case "MPAU":
		if len (a) != 1 {return "ERROR1" }
		MauspufferAus()
		return "OK"
		case "MPL1":
		if len (a) != 1 {return "ERROR1" }
		taste,status,mausX,mausY:= MauspufferLesen1 ()
		return fmt.Sprint(taste)+":"+fmt.Sprint(status)+":"+fmt.Sprint(mausX)+":"+fmt.Sprint(mausY)
		case "SPSO":
		if len (a) < 2 {return "ERROR1" }
		name:= a[1]
		for i:=2;i<len(a);i++ {name = name+":"+a[i]}
		SpieleSound(name)
		return "OK"
		case "GNTE": // GibNotenTempo
		if len (a) != 1 {return "ERROR1" }
		return fmt.Sprint(GibNotenTempo())
		case "SNTE": // SetzeNotentempo
		if len (a) != 2 {return "ERROR1" }
		t,err:=strconv.Atoi(a[1])
		if err != nil {return "ERROR2"}
		SetzeNotenTempo(uint8(t))
		return "OK"
		case "GKPA": // GibKlangparameter
		if len (a) != 1 {return "ERROR1" }
		r,b,k,s,pw:= GibKlangparameter()
		return fmt.Sprint(r)+":"+fmt.Sprint(b)+":"+fmt.Sprint(k)+":"+fmt.Sprint(s)+":"+fmt.Sprint(pw)
		case "SKPA": // SetzeKlangparameter
		if len (a) != 6 {return "ERROR1"}
		rate,err:=strconv.Atoi(a[1])
		bits,err2:=strconv.Atoi(a[2])
		kanaele,err3:=strconv.Atoi(a[3])
		signal,err4:=strconv.Atoi(a[4])
		pweite,err5:=strconv.ParseFloat(a[5],64)
		if err != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {return "ERROR2"}
		SetzeKlangparameter(uint32(rate),uint8(bits),uint8(kanaele),uint8(signal),pweite)
		return "OK"
		case "GHUE": // GibHuellkurve
		if len (a) != 1 {return "ERROR1"}
		an,ab,ha,au:=GibHuellkurve()
		return fmt.Sprint(an)+":"+fmt.Sprint(ab)+":"+fmt.Sprint(ha)+":"+fmt.Sprint(au)
		case "SHUE": // SetzeHuellkurve
		if len (a) != 5 {return "ERROR1"}
		an,err:=strconv.ParseFloat(a[1],64)
		ab,err2:=strconv.ParseFloat(a[2],64)
		halt,err3:=strconv.ParseFloat(a[3],64)
		aus,err4:=strconv.ParseFloat(a[4],64)
		if err != nil || err2 != nil || err3 != nil || err4 != nil {return "ERROR2"}
		SetzeHuellkurve(an, ab, halt,aus)
		return "OK"
		case "SPNO": // SpieleNote
		if len (a) < 4 {return "ERROR1"}
		tonname:= a[1]
		laenge,err:=strconv.ParseFloat(a[2],64)
		wartedauer,err2:=strconv.ParseFloat(a[3],64)
		if err != nil || err2 != nil {return "ERROR2"}
		SpieleNote(tonname,laenge,wartedauer)
		return "OK"
		case "SPAN":
		if len (a) != 1 {return "ERROR1" }
		setzeServerprotokoll(true)
		return "OK"
		case "SPAU":
		if len (a) != 1 {return "ERROR1" }
		setzeServerprotokoll(false)
		return "OK"
		default: //Sonst immer ...
		return "ERROR"
	}
}

		
func main () {
	var h, b, port uint16
	var ipadresse string
	if len(os.Args) != 5 {
		panic ("Startaufruf für den gfx-Server war falsch!")
	}
	// Damit die Funktion 'Tastaturzeichen' die richtigen Tastaturzeichen liefert:
	init_Tastatur_Deutsch () 
	// Nun wird nebenläufig der Gfx-Server gestartet!
	breite,err := strconv.Atoi(os.Args[1])
	if err != nil {breite = 640}
	if breite < 1 || breite > 1920 { breite = 1920 }
	b = uint16(breite)
	hoehe,err:= strconv.Atoi(os.Args[2])
	if err != nil {hoehe = 480}
	if hoehe < 1 || hoehe > 1200 { hoehe = 1200 }
	h = uint16(hoehe)
	portnummer,err:= strconv.Atoi(os.Args[3])
	if err != nil {portnummer = 55555}
	if portnummer < 0 || portnummer > 65535 { portnummer = 55555 }
	port = uint16(portnummer)
	ipadresse = os.Args[4]
	go starteGfxServer (ipadresse,port,f)
	// Nun wird das Grafikfenster geöffnet und mit dem Hauptprogramm in die Main-Loop
	// gegangen ...
	Fenster (b,h)  //terminiert nur, wenn man das Fenster schließt
	// Ist das Fenster geschlossen, so werden die Kanäle zum Programm geschlossen und 
	// 'serverläuft' auf False gesetzt. Damit ist das Programm beendet und der nebenläufig
	// gestartete Server ist damit auch tot.
	for serverLäuft { time.Sleep (1) }
}

	
