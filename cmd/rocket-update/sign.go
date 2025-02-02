package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/mouuff/go-rocket-update/internal/constant"
	"github.com/mouuff/go-rocket-update/internal/crypto"
)

// Sign describes the sign subcommand
// this command is used to sign a folder
type Sign struct {
	flagSet *flag.FlagSet

	path string
	key  string
}

// Name gets the name of the command
func (cmd *Sign) Name() string {
	return "sign"
}

// Init initializes the command
func (cmd *Sign) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)

	cmd.flagSet.StringVar(&cmd.path, "path", "", "path to the package directory (required)")
	cmd.flagSet.StringVar(&cmd.key, "key", "", "path to the private key (required)")

	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *Sign) Run() error {

	log.Println("Reading private key...")
	privkeyBytes, err := ioutil.ReadFile(cmd.key)
	if err != nil {
		return err
	}
	privateKey, err := crypto.ParsePemPrivateKey(privkeyBytes)
	if err != nil {
		return err
	}
	log.Println("Computing signatures...")
	signatures, err := crypto.GetFolderSignatures(privateKey, cmd.path)
	if err != nil {
		return err
	}

	signaturesPath := filepath.Join(cmd.path, constant.SignatureRelPath)
	log.Println("Writing " + signaturesPath + " ...")
	err = crypto.WriteSignaturesToJSON(signaturesPath, signatures)
	if err != nil {
		return err
	}
	log.Println("Signed successfully! Don't forget to keep your private key in a safe place!")
	return nil
}
