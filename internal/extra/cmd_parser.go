package extra

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"github.com/khafidprayoga/parking-app/internal/types"
	"os"
	"strconv"
	"strings"
)

func ParseImportCmd(filePath string) (cmdList []types.Socket, err error) {
	file, errOpen := os.Open(filePath)
	if errOpen != nil {
		err = fmt.Errorf("error on opening file: %s", errOpen.Error())
		return
	}

	var (
		scanner       = bufio.NewScanner(file)
		socketCommand = []types.Socket{}
	)

	for scanner.Scan() {
		if errScan := scanner.Err(); errScan != nil {
			err = fmt.Errorf("error on scanning file: %s", errScan.Error())
			return
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		// parse instruction set
		strCmd := strings.Split(strings.TrimSpace(line), " ")
		cmd := strCmd[0]
		args := strCmd[1:]

		allowedCommands := map[string]struct{}{
			types.CmdCreateStore: {},
			types.CmdPark:        {},
			types.CmdLeave:       {},
			types.CmdStatus:      {},
		}

		if _, ok := allowedCommands[cmd]; !ok {
			err = fmt.Errorf("invalid command: `%s`, this is not allowed", cmd)
			return
		}

		switch cmd {
		case types.CmdCreateStore:
			if len(args) == 0 {
				err = fmt.Errorf("lot capacity not specified")
				return
			}

			parkingLotCap := args[0]

			socketCommand = append(socketCommand, types.Socket{
				Command:    cmd,
				Data:       parkingLotCap,
				XRequestId: uuid.NewString(),
			})

		case types.CmdPark:
			policeNumber := strings.Join(args, "")
			if len(policeNumber) == 0 {
				//skipping invalid or malformed string
				continue
			}

			req := types.Socket{
				Command: cmd,
				Data: types.CarDTO{
					PoliceNumber: policeNumber,
				},
				XRequestId: uuid.NewString(),
			}
			socketCommand = append(socketCommand, req)
		case types.CmdLeave:
			// join string -1 before the hours parameter
			policeNumber := strings.Join(args[0:len(args)-1], "")

			// last is the hours count
			hours := args[len(args)-1]

			durationInHours, errParseDur := strconv.Atoi(hours)
			if errParseDur != nil {
				err = fmt.Errorf("error on parsing hours: %s at this instruction `%s` ", errParseDur.Error(), line)
				return
			}

			req := types.Socket{
				Command: cmd,
				Data: types.CarDTO{
					PoliceNumber: policeNumber,
					Hours:        durationInHours,
				},
				XRequestId: uuid.NewString(),
			}

			socketCommand = append(socketCommand, req)
		case types.CmdStatus:
			socketCommand = append(socketCommand, types.Socket{
				Command:    cmd,
				XRequestId: uuid.NewString(),
			})
		}
	}

	_ = file.Close()
	return socketCommand, nil
}
