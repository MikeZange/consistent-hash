package consistenthash

import (
	"errors"
	"github.com/OneOfOne/xxhash"
	"strconv"
	"strings"
)

type Hasher struct {
	VNC int
	Circle           []int
	Bucket           map[int]string
	Capacity         int
}

func New() *Hasher{
	return &Hasher{
		VNC: 255,
		Circle: []int{},
		Bucket: make(map[int]string),
		Capacity: 65536,
	}
}

func (ch *Hasher) Lookup(k string) (string, error){
	kh := ch.hash(k)
	i := ch.locateInCircle(kh)

	if i >= len(ch.Circle){
		i = 0
	}

	if i > len(ch.Circle){
		return "", errors.New("invalid index found for key")
	}

	return ch.Bucket[ch.Circle[i]], nil

}

func (ch *Hasher) AddNode(node string){
	ch.addNodeToCircle(node)
}

func (ch *Hasher) RemoveNode(node string){
	ch.removeNodeFromCircle(node)
}

func (ch *Hasher) addNodeToCircle(n string){

	for x := 0; x < ch.VNC; x++ {

		var sb strings.Builder

		sb.WriteString(n + "#" + strconv.Itoa(x))
		nh := ch.hash(sb.String())

		i := ch.locateInCircle(nh)

		if (len(ch.Circle) > i) && ch.Circle[i] == nh{
			continue
		}

		if i >= len(ch.Circle) {
			ch.Circle = append(ch.Circle, nh)
		} else{
			ch.Circle = append(ch.Circle, 0)
			copy(ch.Circle[i+1:], ch.Circle[i:])
			ch.Circle[i] = nh
		}

		ch.Bucket[nh] = n
	}
}

func (ch *Hasher) removeNodeFromCircle(n string){

	for x := 0; x < ch.VNC; x++ {

		var sb strings.Builder

		sb.WriteString(n + "#" + strconv.Itoa(x))
		nh := ch.hash(sb.String())

		i := ch.locateInCircle(nh)

		end := len(ch.Circle)-1
		copy(ch.Circle[i:], ch.Circle[i+1:])
		ch.Circle[end] = 0
		ch.Circle = ch.Circle[:end]

		_, ok := ch.Bucket[nh]

		if ok {
			delete(ch.Bucket, nh)
		}
	}
}

func (ch *Hasher) hash (k string) int {

	h := xxhash.Checksum32([]byte(k))

	return int(h) % ch.Capacity
}

func (ch *Hasher) locateInCircle(nh int) int{

	l := 0
	h := len(ch.Circle) - 1

	for l <= h {
		m := (l + h) / 2
		v := ch.Circle[m]

		if v == nh {
			return m
		}

		if v > nh {
			h = m - 1
		}

		if v < nh {
			l = m + 1
		}
	}

	return (l + h + 1) / 2
}
