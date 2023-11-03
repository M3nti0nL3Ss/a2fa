package commands

import (
	"context"
	"fmt"
	"github.com/csyezheng/a2fa/oath"
	"github.com/spf13/cobra"
)

type generateCommand struct {
	r        *rootCommand
	name     string
	use      string
	commands []Commander

	// Flags
	hotp        bool
	totp        bool
	base32      bool
	hash        string
	valueLength int

	// Flags only for HOTP
	counter int64

	// Flags only for TOTP
	epoch    int64
	interval int64
}

func (c *generateCommand) Name() string {
	return c.name
}

func (c *generateCommand) Use() string {
	return c.use
}

func (c *generateCommand) Init(cd *Commandeer) error {
	cmd := cd.CobraCommand
	cmd.Short = "Generate one-time password from secret key"
	cmd.Long = "Generate one-time password from secret key"
	cmd.Flags().BoolVar(&c.hotp, "hotp", false, "use event-based HOTP mode (default is false)")
	cmd.Flags().BoolVar(&c.totp, "totp", true, "use use time-variant TOTP mode (default is true)")
	cmd.Flags().BoolVarP(&c.base32, "base32", "b", true, "use base32 encoding of KEY instead of hex, (default=true)")
	cmd.Flags().StringVarP(&c.hash, "hash", "H", "SHA1", "A cryptographic hash method H (SHA1, SHA256, SHA512, default is SHA1)")
	cmd.Flags().IntVarP(&c.valueLength, "length", "l", 6, "A HOTP value length d (6–10, default is 6, and 6–8 is recommended)")
	cmd.Flags().Int64VarP(&c.counter, "counter", "c", 0, "used for HOTP, A counter C, which counts the number of iterations")
	// epoch (T0) is the epoch as specified in seconds since the Unix epoch (e.g. if using Unix time, then T0 is 0)
	cmd.Flags().Int64VarP(&c.epoch, "epoch", "e", 0, "used for TOTP, epoch (T0) which is the Unix time from which to start counting time steps (default is 0),")
	// // interval (Tx) is the length of one time duration (e.g. 30 seconds).
	cmd.Flags().Int64VarP(&c.interval, "interval", "i", 30, "used for TOTP, an interval (Tx) which will be used to calculate the value of the counter CT (default is 30 seconds).")
	return nil
}

func (c *generateCommand) Args(ctx context.Context, cd *Commandeer, args []string) error {
	if err := cobra.ExactArgs(1)(cd.CobraCommand, args); err != nil {
		return err
	}
	return nil
}

func (c *generateCommand) PreRun(cd, runner *Commandeer) error {
	c.r = cd.Root.Command.(*rootCommand)
	return nil
}

func (c *generateCommand) Run(ctx context.Context, cd *Commandeer, args []string) error {
	if err := cobra.ExactArgs(1)(cd.CobraCommand, args); err != nil {
		return err
	}
	secretKey := args[0]
	otp := ""
	if c.hotp {
		hotp := oath.NewHOTP(c.base32, c.hash, c.counter, c.valueLength)
		otp = hotp.GeneratePassCode(secretKey)
	} else if c.totp {
		totp := oath.NewTOTP(c.base32, c.hash, c.valueLength, c.epoch, c.interval)
		otp = totp.GeneratePassCode(secretKey)
	} else {
		return fmt.Errorf("mode should be hotp or totp")
	}

	fmt.Println("Code: " + otp)
	return nil
}

func (c *generateCommand) Commands() []Commander {
	return c.commands
}

func newGenerateCommand() *generateCommand {
	generateCmd := &generateCommand{
		name: "generate",
		use:  "a2fa generate [flags] <secret key>",
	}
	return generateCmd
}