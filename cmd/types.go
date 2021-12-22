package cmd

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/thedonutfactory/donutbox/compiler"
)

type DonutBytecode struct {
	Version  int8
	Bytecode *compiler.Bytecode
}

func NewDonutByteCode() *DonutBytecode {
	return &DonutBytecode{
		Version: 1,
	}
}

func (key *DonutBytecode) write(filename string) error {
	return write(key, filename)
}

func (key *DonutBytecode) read(filename string) error {
	return read(key, filename)
}

type DonutTransaction struct {
	Version  int8
	Input    []int32
	FuncName string
}

func NewDonutTransaction(input []int32, funcName string) *DonutTransaction {
	return &DonutTransaction{
		Version:  1,
		Input:    input,
		FuncName: funcName,
	}
}

func (key *DonutTransaction) write(filename string) error {
	return write(key, filename)
}

func (key *DonutTransaction) read(filename string) error {
	return read(key, filename)
}

func write(key interface{}, filename string) error {
	file, _ := os.Create(filename)
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err := encoder.Encode(key)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func read(key interface{}, filename string) error {
	dataFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(key)
	if err != nil {
		fmt.Println(err)
		return err
	}
	dataFile.Close()
	return nil
}