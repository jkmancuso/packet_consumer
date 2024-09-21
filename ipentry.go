package main

import "time"

type ipEntry struct {
	Time   time.Time
	Src    ipInfo
	Dst    ipInfo
	Length int
}

// this is used to represent 1 connection essentially
type ipInfo struct {
	Ipv4 string
	DNS  string
}

func (i ipEntry) getTags() map[string]string {

	tags := make(map[string]string)
	tags["source_ip"] = i.Src.Ipv4
	tags["source_dns"] = i.Src.DNS
	tags["dest_ip"] = i.Dst.Ipv4
	tags["dest_dns"] = i.Dst.DNS

	return tags
}

func (i ipEntry) getFields() map[string]interface{} {

	return map[string]interface{}{"size": i.Length}
}

func (i ipEntry) getTime() time.Time {

	return i.Time
}

/*
SAMPLE
{"Time": "2024-09-21T17:48:58.283731821Z", "Src":{"Ipv4":"142.250.72.99","DNS":"www.gstatic.com"},"Dst":{"Ipv4":"192.168.7.249","DNS":""},"Length":1278}
*/
