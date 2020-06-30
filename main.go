package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var statusMsg = map[byte]string{
	0x0: "System Exclusive",
	0x1: "MIDI Time Code Qtr. Frame",
	0x2: "Song Position Pointer",
	0x3: "Song Select (Song #)",
	0x4: "Undefined (Reserved)",
	0x5: "Undefined (Reserved)",
	0x6: "Tune request",
	0x7: "End of SysEx (EOX)",
	0x8: "Timing clock",
	0x9: "Undefined (Reserved)",
	0xA: "Start",
	0xB: "Continue",
	0xC: "Stop",
	0xD: "Undefined (Reserved)",
	0xE: "Active Sensing",
	0xF: "System Reset",
}

func check(e error) {
	if e != nil {
		log.Fatalf("%s", e)
	}
}

func translateMidi(midi []byte) (string, bool) {
	if len(midi) != 3 {
		return "Wrong size message", false
	}
	status := midi[0]
	data1 := midi[1]
	data2 := midi[2]

	var s1, s2, d1, d2 = "Unset", "Unset", "Unset", "unset"

	s2 = fmt.Sprintf("Channel %d", status&0x0F)

	var clearOutput = false

	switch status & 0xF0 {
	case 0x80:
		d1 = fmt.Sprintf("Note Number %d", data1)
		d2 = fmt.Sprintf("Note Velocity %d", data2)
		s1 = "Note Off"
	case 0x90:
		d1 = fmt.Sprintf("Note Number %d", data1)
		d2 = fmt.Sprintf("Note Velocity %d", data2)
		s1 = "Note On"
	case 0xA0:
		d1 = fmt.Sprintf("Note Number %d", data1)
		d2 = fmt.Sprintf("Pressure %d", data2)
		s1 = "Polyphonic Aftertouch"
	case 0xB0:
		s1 = "Control/Mode Change"
	case 0xC0:
		s1 = "Program Change"
	case 0xD0:
		s1 = "Channel Aftertouch"
	case 0xE0:
		s1 = "Pitch Bend Change"
	case 0xF0:
		s1 = statusMsg[status&0x0f]
		s2 = ""
		clearOutput = true
	default:
		s1 = "Unknown Status Byte"
		clearOutput = true
	}

	return fmt.Sprintf("Midi Msg: %x | Status: %20s %20s | Data: %20s %20s",
		midi, s1, s2, d1, d2), clearOutput
}

func main() {

	f, err := os.Open("/dev/snd/midiC2D0")
	check(err)
	midiBuf := make([]byte, 3)
	for {
		f.Seek(3, os.SEEK_END)
		_, err := f.Read(midiBuf)
		check(err)

		midiInfo, clearOutput := translateMidi(midiBuf)
		var outputfmt string
		if clearOutput {
			outputfmt = "\r%s"
		} else {
			outputfmt = "\r%s\n"
		}

		fmt.Fprintf(os.Stdout, outputfmt, midiInfo)
		time.Sleep(100 * time.Millisecond)
	}

}
