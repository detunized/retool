package shape

import "encoding/binary"

/* JS code:

function shapeStringHash(s) {
    let pairs = [];
    for (let i = 0; i < s.length; i += 2) {
        let c = s.charCodeAt(i) << 16;
        if (i + 1 < s.length) {
            c |= s.charCodeAt(i + 1);
        }
        pairs.push(c);
    }

    let h = 0;
    for (let i = 0; i < pairs.length; i++) {
        h = 0 | (h << 5) - h + pairs[i];
    }

    return h;
}
*/

func CalcStringHashBytes(bytes []byte) []byte {
	hash := make([]byte, 4)
	binary.LittleEndian.PutUint32(hash, CalcStringHashInt(bytes))

	return hash
}

func CalcStringHashInt(bytes []byte) uint32 {
	n := len(bytes)
	pairs := make([]uint32, 0, (n+1)/2)
	for i := 0; i < n; i += 2 {
		c := uint32(bytes[i]) << 16
		if i < n-1 {
			c |= uint32(bytes[i+1])
		}
		pairs = append(pairs, c)
	}

	h := int64(0)
	for _, p := range pairs {
		h = ((h << 5) - h + int64(p)) & 0xffff_ffff
	}

	return uint32(h)
}
