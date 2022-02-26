package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hamba/avro"
	"github.com/vmihailenco/msgpack/v5"
	yaml "gopkg.in/yaml.v2"

	"serialization/serialization/models"
)

type TestInstance struct {
	ID   int            `json:"id" yaml:"id" avro:"ID"`
	Frac float64        `json:"frac" yaml:"frac" avro:"Frac"`
	Name string         `json:"name" yaml:"name" avro:"Name"`
	Maps map[string]int `json:"maps" yaml:"maps" avro:"Maps"`
}

type Test struct {
	Tests []TestInstance `json:"tests" yaml:"tests" avro:"Tests"`
}

func printStartMessage(experimentName string) {
	fmt.Printf("Starting %s serialization experiment\n", experimentName)
}

type Experiment interface {
	ser(Test) []byte
	deser([]byte)
}

func run_timed(name string, exp Experiment, t Test) {
	fmt.Println("--------------------------------------------------")
	printStartMessage(name)
	start := time.Now()

	start_ser := time.Now()
	bytes := exp.ser(t)
	elapsed_ser := time.Since(start_ser)

	start_deser := time.Now()
	exp.deser(bytes)
	elapsed_deser := time.Since(start_deser)

	elapsed := time.Since(start)
	fmt.Printf("Experiment took %s\n", elapsed)
	fmt.Printf("Serialization took %s\n", elapsed_ser)
	fmt.Printf("Deserialization took %s\n", elapsed_deser)
	fmt.Printf("Total Bytes: %d\n", len(bytes))
}

type nativeExperiment struct {
	network *bytes.Buffer
	enc     *gob.Encoder
	dec     *gob.Decoder
}

func (nt nativeExperiment) ser(t Test) []byte {
	nt.enc.Encode(t)
	return nt.network.Bytes()
}

func (nt nativeExperiment) deser(b []byte) {
	var q Test
	nt.dec.Decode(&q)
}

type jsonExperiment struct{}

func (js jsonExperiment) ser(t Test) []byte {
	jsonBytes, _ := json.Marshal(t)
	return jsonBytes
}

func (js jsonExperiment) deser(b []byte) {
	var q Test
	json.Unmarshal(b, &q)
}

type xmlExperiment struct{}

func (xmle xmlExperiment) ser(t Test) []byte {
	xmlBytes, _ := xml.Marshal(t)
	return xmlBytes
}

func (xmle xmlExperiment) deser(b []byte) {
	var q Test
	xml.Unmarshal(b, &q)
}

type yamlExperiment struct{}

func (ya yamlExperiment) ser(t Test) []byte {
	yamlBytes, _ := yaml.Marshal(t)
	return yamlBytes
}

func (ya yamlExperiment) deser(b []byte) {
	var q Test
	yaml.Unmarshal(b, &q)
}

type msgpackExperiment struct{}

func (msg msgpackExperiment) ser(t Test) []byte {
	msgBytes, _ := msgpack.Marshal(t)
	return msgBytes
}

func (msg msgpackExperiment) deser(b []byte) {
	var q Test
	msgpack.Unmarshal(b, &q)
}

type avroExperiment struct {
	schema avro.Schema
}

func (avr avroExperiment) ser(t Test) []byte {
	avrBytes, _ := avro.Marshal(avr.schema, t)
	return avrBytes
}

func (avr avroExperiment) deser(b []byte) {
	var q Test
	avro.Unmarshal(avr.schema, b, &q)
}

type protoExperiment struct {
	test models.Test
}

func (prt protoExperiment) ser(t Test) []byte {
	prtBytes, _ := proto.Marshal(&prt.test)
	return prtBytes
}

func (prt protoExperiment) deser(b []byte) {
	var q models.Test
	proto.Unmarshal(b, &q)
}

func main() {
	var t Test

	file, _ := ioutil.ReadFile("json_init.json")
	_ = json.Unmarshal([]byte(file), &t)

	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	dec := gob.NewDecoder(&network)
	nt := nativeExperiment{network: &network, enc: enc, dec: dec}
	run_timed("Native", nt, t)

	var js jsonExperiment
	run_timed("JSON", js, t)

	// var xm xmlExperiment
	// run_timed("XML", xm, t)

	var ya yamlExperiment
	run_timed("YAML", ya, t)

	var msg msgpackExperiment
	run_timed("MSGPACK", msg, t)

	// schemaStr, _ := ioutil.ReadFile("schema.avsc")
	// schema, _ := avro.Parse(string(schemaStr))

	// avr := avroExperiment{schema: schema}
	// run_timed("AVRO", avr, t)

	var prtTest models.Test

	for _, test := range t.Tests {
		othrTestInstance := models.Test_TestInstance{
			Id: int32(test.ID), Frac: float32(test.Frac),
			Name: test.Name,
		}

		othrTestInstance.Maps = make(map[string]int32)

		for k, v := range test.Maps {
			othrTestInstance.Maps[k] = int32(v)
		}
		prtTest.Tests = append(prtTest.Tests, &othrTestInstance)
	}

	prt := protoExperiment{test: prtTest}
	run_timed("PROTO", prt, t)
}
