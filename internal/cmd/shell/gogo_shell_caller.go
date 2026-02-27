package shell

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/lgdd/lfr-cli/pkg/util/logger"
)

// Telnet protocol constants (RFC 854)
const (
	telnetIAC  = 255
	telnetDO   = 253
	telnetDONT = 254
	telnetWILL = 251
	telnetWONT = 252
	telnetSB   = 250
	telnetSE   = 240

	optEcho         = 1
	optSuppressGA   = 3
	optTerminalType = 24
)

// connectGogoShell opens a raw TCP connection to the Gogo Shell, performs the
// telnet negotiation that Apache Felix/MINA requires, and then runs an
// interactive read-eval loop.  Without the negotiation the server waits ~60s
// for the handshake to time out before sending the first prompt.
func connectGogoShell(host string, port int) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Consume the server's IAC negotiation sequences and reply, then wait for
	// the first "g! " prompt before handing off to the interactive loop.
	if err := negotiateAndWaitForPrompt(conn, reader); err != nil {
		return err
	}

	// Goroutine: server → stdout (strips any residual IAC sequences)
	serverDone := make(chan error, 1)
	go func() {
		serverDone <- copyFromServer(reader)
	}()

	// Main goroutine: stdin → server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "disconnect" || text == "exit" {
			conn.Write([]byte("disconnect\r\ny\r\n"))
			break
		}
		if _, err := conn.Write(append([]byte(text), '\r', '\n')); err != nil {
			break
		}
	}

	// Wait for the server to close the connection (or for an error).
	if err := <-serverDone; err != nil {
		logger.Fatal(err.Error())
	}
	fmt.Println("Bye 👋")
	return nil
}

// negotiateAndWaitForPrompt handles the IAC exchange that Apache Felix/MINA
// initiates on every new connection, then blocks until the first "g! " prompt
// is received so the caller can start the interactive session immediately.
func negotiateAndWaitForPrompt(conn net.Conn, reader *bufio.Reader) error {
	var buf strings.Builder

	for {
		b, err := reader.ReadByte()
		if err != nil {
			return err
		}

		if b == telnetIAC {
			if err := handleIACCommand(conn, reader); err != nil {
				return err
			}
			continue
		}

		buf.WriteByte(b)
		os.Stdout.Write([]byte{b})

		if strings.HasSuffix(buf.String(), "g! ") {
			return nil
		}
	}
}

// handleIACCommand reads the rest of one IAC sequence and writes the
// appropriate reply so the server can complete its negotiation immediately.
func handleIACCommand(conn net.Conn, reader *bufio.Reader) error {
	cmd, err := reader.ReadByte()
	if err != nil {
		return err
	}

	switch cmd {
	case telnetDO:
		opt, err := reader.ReadByte()
		if err != nil {
			return err
		}
		if opt == optTerminalType {
			// Agree to send our terminal type; the server will follow up with SB.
			conn.Write([]byte{telnetIAC, telnetWILL, opt})
		} else {
			conn.Write([]byte{telnetIAC, telnetWONT, opt})
		}

	case telnetWILL:
		opt, err := reader.ReadByte()
		if err != nil {
			return err
		}
		if opt == optEcho || opt == optSuppressGA {
			conn.Write([]byte{telnetIAC, telnetDO, opt})
		} else {
			conn.Write([]byte{telnetIAC, telnetDONT, opt})
		}

	case telnetDONT, telnetWONT:
		_, err := reader.ReadByte() // consume the option byte
		return err

	case telnetSB:
		opt, err := reader.ReadByte()
		if err != nil {
			return err
		}
		// Read subnegotiation body until IAC SE
		for {
			c, err := reader.ReadByte()
			if err != nil {
				return err
			}
			if c == telnetIAC {
				next, err := reader.ReadByte()
				if err != nil {
					return err
				}
				if next == telnetSE {
					break
				}
			}
		}
		if opt == optTerminalType {
			// Send VT220 as our terminal type (matches what Liferay IDE uses)
			conn.Write([]byte{telnetIAC, telnetSB, optTerminalType, 0, 'V', 'T', '2', '2', '0', telnetIAC, telnetSE})
		}
	}

	return nil
}

// copyFromServer forwards data from the server to stdout, silently discarding
// any IAC sequences that arrive after the initial negotiation is complete.
func copyFromServer(reader *bufio.Reader) error {
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return err
		}

		if b == telnetIAC {
			cmd, err := reader.ReadByte()
			if err != nil {
				return err
			}
			switch cmd {
			case telnetDO, telnetDONT, telnetWILL, telnetWONT:
				if _, err := reader.ReadByte(); err != nil {
					return err
				}
			case telnetSB:
				for {
					c, err := reader.ReadByte()
					if err != nil {
						return err
					}
					if c == telnetIAC {
						next, err := reader.ReadByte()
						if err != nil {
							return err
						}
						if next == telnetSE {
							break
						}
					}
				}
			}
			continue
		}

		os.Stdout.Write([]byte{b})
	}
}
