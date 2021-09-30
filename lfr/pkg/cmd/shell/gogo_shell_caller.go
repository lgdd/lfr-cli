package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/lgdd/liferay-cli/lfr/pkg/util/printutil"
	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"
)

// Telnet caller for Gogo Shell
var GogoShellCaller telnet.Caller = internalGogoShellCaller{}

type internalGogoShellCaller struct{}

func (caller internalGogoShellCaller) CallTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr

	go func(writer io.Writer, reader io.Reader) {
		var buffer [1]byte
		p := buffer[:]

		for {
			n, err := reader.Read(p)
			if n <= 0 && nil == err {
				continue
			} else if n <= 0 && nil != err {
				fmt.Println("Bye 👋")
				os.Exit(0)
			}

			_, err = oi.LongWrite(writer, p)

			if err != nil {
				printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
				os.Exit(1)
			}
		}
	}(stdout, r)

	var buffer bytes.Buffer
	var p []byte

	var crlfBuffer [2]byte = [2]byte{'\r', '\n'}
	crlf := crlfBuffer[:]

	scanner := bufio.NewScanner(stdin)
	scanner.Split(scannerSplitFunc)

	for scanner.Scan() {
		if "disconnect" == string(scanner.Bytes()) || "exit" == string(scanner.Bytes()) {
			buffer.Write([]byte("disconnect"))
			buffer.Write(crlf)
			buffer.Write([]byte("y"))
			buffer.Write(crlf)
		} else {
			buffer.Write(scanner.Bytes())
			buffer.Write(crlf)
		}

		p = buffer.Bytes()

		n, err := oi.LongWrite(w, p)
		if nil != err {
			break
		}
		if expected, actual := int64(len(p)), n; expected != actual {
			err := fmt.Errorf("Transmission problem: tried sending %d bytes, but actually only sent %d bytes.", expected, actual)
			fmt.Fprint(stderr, err.Error())
			return
		}

		buffer.Reset()
	}
	time.Sleep(3 * time.Millisecond)
}

func scannerSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		return 0, nil, nil
	}

	return bufio.ScanLines(data, atEOF)
}
