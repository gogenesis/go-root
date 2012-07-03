package groot

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

type Buffer struct {
	order binary.ByteOrder // byte order of underlying data source
	data  []byte           // data source
	buf   *bytes.Buffer    // buffer for more efficient i/o from r
	klen  uint32           // to compute refs (used in read_class, read_object)
}

func NewBuffer(data []byte, order binary.ByteOrder, klen uint32) (b *Buffer, err error) {
	b = &Buffer{
		order: order,
		data:  data,
		klen:  klen,
	}
	b.buf = bytes.NewBuffer(b.data[:])
	return
}

func NewBufferFromKey(k *Key) (b *Buffer, err error) {
	buf, err := k.Buffer()
	if err != nil {
		return
	}
	return NewBuffer(buf, k.file.order, uint32(k.keysz))
}

func (b *Buffer) Pos() int {
	return len(b.data) - b.Len()
}

func (b *Buffer) Len() int {
	return len(b.Bytes())
}

func (b *Buffer) Bytes() []byte {
	return b.buf.Bytes()
}

func (b *Buffer) clone() *Buffer {
	bb, err := NewBuffer(b.data[:], b.order, b.klen)
	if err != nil {
		return nil
	}
	bb.read_nbytes(b.Pos())
	return bb
}

func (b *Buffer) rewind_nbytes(nbytes int) {
	idx := b.Pos()
	b.buf = bytes.NewBuffer(b.data[idx-nbytes:])
}

func (b *Buffer) read_nbytes(nbytes int) (o []byte) {
	o = make([]byte, nbytes)
	err := binary.Read(b.buf, b.order, o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntoi2() (o int16) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntoi4() (o int32) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntoi8() (o int64) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntobyte() (o byte) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntou2() (o uint16) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntou4() (o uint32) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntou8() (o uint64) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntof() (o float32) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) ntod() (o float64) {
	err := binary.Read(b.buf, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Buffer) read_array_F() (o []float32) {
	n := int(b.ntou4())
	o = make([]float32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntof()
	}
	return
}

func (b *Buffer) read_array_D() (o []float64) {
	n := int(b.ntou4())
	o = make([]float64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntod()
	}
	return
}

func (b *Buffer) read_array_S() (o []int16) {
	n := int(b.ntou4())
	o = make([]int16, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi2()
	}
	return
}

func (b *Buffer) read_array_I() (o []int32) {
	n := int(b.ntou4())
	o = make([]int32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi4()
	}
	return
}

func (b *Buffer) read_array_L() (o []int64) {
	n := int(b.ntou4())
	o = make([]int64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi8()
	}
	return
}

func (b *Buffer) read_array_C() (o []byte) {
	n := int(b.ntou4())
	o = make([]byte, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntobyte()
	}
	return
}

func (b *Buffer) read_static_array() (o []uint32) {
	n := int(b.ntou4())
	o = make([]uint32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntou4()
	}
	return
}

func (b *Buffer) read_fast_array_F(n int) (o []float32) {
	o = make([]float32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntof()
	}
	return
}

func (b *Buffer) read_fast_array_D(n int) (o []float64) {
	o = make([]float64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntod()
	}
	return
}

func (b *Buffer) read_fast_array_S(n int) (o []int16) {
	o = make([]int16, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi2()
	}
	return
}

func (b *Buffer) read_fast_array_I(n int) (o []int32) {
	o = make([]int32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi4()
	}
	return
}

func (b *Buffer) read_fast_array_L(n int) (o []int64) {
	o = make([]int64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi8()
	}
	return
}

func (b *Buffer) read_fast_array_UL(n int) (o []uint64) {
	o = make([]uint64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntou8()
	}
	return
}

func (b *Buffer) read_fast_array_C(n int) (o []byte) {
	o = make([]byte, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntobyte()
	}
	return
}

func (b *Buffer) read_fast_array_tstring(n int) (o []string) {
	o = make([]string, n)
	for i := 0; i < n; i++ {
		o[i] = b.read_tstring()
	}
	return
}

func (b *Buffer) read_fast_array(n int) (o []uint32) {
	o = make([]uint32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntou4()
	}
	return
}

func (b *Buffer) read_tstring() string {
	n := int(b.ntobyte())
	if n == 255 {
		// large string
		n = int(b.ntou4())
	}
	if n == 0 {
		return ""
	}
	v := b.ntobyte()
	if v == 0 {
		return ""
	}
	o := make([]byte, n)
	o[0] = v
	err := binary.Read(b.buf, b.order, o[1:])
	if err != nil {
		panic(err)
	}
	return string(o)
}

//FIXME
// readBasicPointerElem
// readBasicPointer

func (b *Buffer) read_string(max int) string {
	o := []byte{}
	n := 0
	var v byte
	for {
		v = b.ntobyte()
		if v == 0 {
			break
		}
		n += 1
		if max > 0 && n >= max {
			break
		}
		o = append(o, v)
	}
	return string(o)
}

func (b *Buffer) read_std_string() string {
	nwh := b.ntobyte()
	nchars := int32(nwh)
	if nwh == 255 {
		nchars = b.ntoi4()
	}
	if nchars < 0 {
		panic("groot.readStdString: negative char number")
	}
	return b.read_string(int(nchars))
}

func (b *Buffer) read_version() (vers uint16, pos, bcnt uint32) {

	bcnt = b.ntou4()
	if (int64(bcnt) & ^kByteCountMask) != 0 {
		bcnt = uint32(int64(bcnt) & ^kByteCountMask)
	} else {
		panic("groot.Buffer.read_version: too old file")
	}
	vers = b.ntou2()
	return
}

func (b *Buffer) read_object() (o Object) {
	spos := b.Pos()
	// before reading object, save start position
	startbuf := b.clone()

	clsname, bcnt, isref := b.read_class()
	dprintf(">>[class=%s] [bcnt=%v] [isref=%v]\n", clsname, bcnt, isref)
	if isref {
		bb := b.clone()
		dprintf("obj_offset: [%v] -> [%v] -> [%v]\n",
			bcnt, bcnt-kMapOffset, bcnt-kMapOffset-bb.klen)
		bb.rewind_nbytes(b.Pos() - int(bcnt - kMapOffset - bb.klen))
		ii := bb.ntou4()
		if (ii & kByteCountMask) != 0 {
			scls := bb.read_class_tag()
			if scls == "" {
				panic("groot.Buffer.read_object: read_class_tag did not find a class name")
			}
		} else {
			/* boo */
		}
		/*
		 // in principle at this point m_pos should be
		 //   m_buffer+startpos+sizeof(unsigned int)
		 // but enforce it anyway : 
		 m_pos = m_buffer+startpos+sizeof(unsigned int); 
		*/
		b = startbuf //FIXME ??
		b.read_nbytes(4)
	} else {
		if clsname == "" {
			o = nil
			// m_pos = m_buffer+startpos+bcnt+sizeof(unsigned int);
			b = startbuf
			b.read_nbytes(int(bcnt+4))
		} else {

			factory := Factory.Get(clsname)
			if factory == nil {
				dprintf("**err** no factory for class [%s]\n", clsname)
				return
			}

			vv := factory()
			o = vv.Interface().(Object)
			if vv, ok := vv.Interface().(ROOTStreamer); ok {
				err := vv.ROOTDecode(b)
				if err != nil {
					panic(err)
				} else {
					dprintf("--decoded[%s]--\n", o.Name())
				}
			} else {
				dprintf("**err** class [%s] does not satisfy the ROOTStreamer interface\n", clsname)
			}
			b.check_byte_count(0,bcnt, spos, clsname)
		}
	}
	return o
}

func (b *Buffer) read_class() (name string, bcnt uint32, isref bool) {

	//var bufvers = 0
	i := b.ntou4()
	dprintf("..first_int: %x (len=%d)\n", i, b.Len()/8)
	if i == kNullTag {
		/*empty*/
		dprintf("read_class: first_int is kNullTag\n")
	} else if (i & kByteCountMask) != 0 {
		//bufvers = 1
		dprintf("read_class: first_int & kByteCountMask\n")
		clstag := b.read_class_tag()
		if clstag == "" {
			panic("groot.Buffer.read_class: empty class tag")
		}
		name = clstag
		bcnt = uint32(int64(i) & ^kByteCountMask)
		dprintf("read_class: kNewClassTag: read class name='%s' bcnt=%d\n",
			name, bcnt)
	} else {
		dprintf("read_class: first_int %x ==> position toward object.\n", i)
		bcnt = uint32(i)
		isref = true
	}
	dprintf("--[cls=%s] [bcnt=%v] [isref=%v]\n", name, bcnt, isref)
	return
}

func (b *Buffer) read_class_tag() (clstag string) {
	tag := b.ntou4()

	tag_new_class := tag == kNewClassTag
	tag_class_mask := (int64(tag) & (^int64(kClassMask))) != 0

	dprintf("--tag:%v %x -> new_class=%v class_mask=%v\n", 
		tag, tag, 
		tag_new_class,
		tag_class_mask)

	if tag_new_class {
		clstag = b.read_string(80)
		dprintf("--class+tag: [%v] - kNewClassTag\n", clstag)
	} else if tag_class_mask {
		ref := uint32(int64(tag) & (^int64(kClassMask)))
		dprintf("--class-tag: [%v] & kClassMask -- ref=%d -- recurse\n", 
			clstag, ref)
		bb := b.clone()
		dprintf("cl_offset: [%v] -> [%v] -> [%v]\n",
			ref, ref-kMapOffset, ref-kMapOffset-bb.klen)
		bb.rewind_nbytes(b.Pos() - int(ref - kMapOffset - bb.klen))
		clstag = bb.read_class_tag()
		//printf("--class-tag: [%v] & kClassMask\n", clstag)
	} else {
		panic(fmt.Errorf("groot.read_class_tag: unknown class-tag [%v]", tag))
	}
	return
}

func (b *Buffer) read_tnamed() (name, title string) {
	spos := b.Pos()
	vers, pos, bcnt := b.read_version()
	id := b.ntou4()
	bits := b.ntou4()
	bits |= kIsOnHeap // by definition de-serialized object is on heap
	if (bits & kIsReferenced) == 0 {
		_ = b.read_nbytes(2)
	}
	name = b.read_tstring()
	title = b.read_tstring()
	printf("read_tnamed: vers=%v pos=%v bcnt=%v id=%v bits=%v name='%v' title='%v'\n",
		vers, pos, bcnt, id, bits, name, title)

	b.check_byte_count(pos,bcnt,spos,"TNamed")

	return
}

func (b *Buffer) read_elements() (elmts []Object) {
	name, bcnt, isref := b.read_class()
	printf("read_elements: name='%v' bcnt=%v isref=%v\n",
		name, bcnt, isref)
	elmts = b.read_obj_array()
	return elmts
}

func (b *Buffer) read_obj_array() (elmts []Object) {
	spos := b.Pos()
	vers, pos, bcnt := b.read_version()
	if vers > 2 {
		// skip version
		b.read_nbytes(2)
		// skip object bits and unique id
		b.read_nbytes(8)
	}
	name :=  b.read_tstring()

	nobjs := int(b.ntoi4())
	lbound := b.ntoi4()

	dprintf("read_obj_array: vers=%v pos=%v bcnt=%v name='%v' nobjs=%v lbound=%v\n",
		vers, pos, bcnt, name, nobjs, lbound)

	elmts = make([]Object, nobjs)
	for i := 0; i < nobjs; i++ {
		dprintf("read_obj_array: %d/%d...\n", i, nobjs)
		obj := b.read_object()
		dprintf("read_obj_array: %d/%d...[done]\n", i, nobjs)
		elmts[i] = obj
	}

	b.check_byte_count(pos, bcnt, spos, "TObjArray")
	return elmts
}

func (b *Buffer) read_attline() (color, style, width uint16) {
	spos := b.Pos()
	/*vers*/_, pos, bcnt := b.read_version()
	color = b.ntou2()
	style = b.ntou2()
	width = b.ntou2()
	b.check_byte_count(pos,bcnt, spos, "TAttLine")
	return
}

func (b *Buffer) read_attfill() (color, style uint16) {
	spos := b.Pos()
	/*vers*/_, pos, bcnt := b.read_version()
	color = b.ntou2()
	style = b.ntou2()
	b.check_byte_count(pos,bcnt, spos, "TAttFill")
	return
}

func (b *Buffer) read_attmarker() (color, style uint16, width float32) {
	spos := b.Pos()
	/*vers*/_, pos, bcnt := b.read_version()
	color = b.ntou2()
	style = b.ntou2()
	width = b.ntof()
	b.check_byte_count(pos,bcnt, spos,"TAttMarker")
	return
}

//FIXME
// readObjectAny
// readTList
// readTObjArray
// readTClonesArray
// readTCollection
// readTHashList
// readTNamed
// readTCanvas

func (b *Buffer) check_byte_count(start, count uint32, spos int, cls string) bool {
	if count == 0 {
		return true
	}

	lenbuf := uint64(start) + uint64(count) + uint64(unsafe.Sizeof(uint(0)))
	diff := uint64(b.Pos() - spos)

	if diff == lenbuf {
		return true
	}
	
	if diff < lenbuf {
		err := fmt.Errorf(
			"buffer.check_count: object of class [%s] read too few bytes (%d missing)",
			cls, lenbuf-diff)
		fmt.Printf("**error** %s\n", err.Error())
		//panic(err)

		dprintf("-->pos= %v\n", b.Pos())
		b.read_nbytes(int(lenbuf-diff))
		dprintf("-->pos= %v\n", b.Pos())
		return true
	}
	if diff > lenbuf {
		err := fmt.Errorf(
			"buffer.check_count: object of class [%s] read too many bytes (%d in excess)",
			cls, diff-lenbuf)
		fmt.Printf("**error** %s\n", err.Error())

		//panic(err)
		dprintf("-->pos= %v\n", b.Pos())
		b.rewind_nbytes(int(diff -lenbuf))
		dprintf("-->pos= %v\n", b.Pos())
		return true
	}
	return false
}
// EOF
