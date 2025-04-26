package Crypto

import "encoding/binary"

func HashDecrypt(pData []byte, nLength uint32, nKey uint32) uint32 {
	num := nKey ^ 347277256
	num1 := nKey ^ 2361332396
	num2 := nKey ^ 604215233
	num3 := nKey ^ 4089260480
	num4 := 0
	var num5 uint32 = 0
	var i int

	for i = 0; uint32(i) < (nLength >> 4); i++ {
		val1 := binary.LittleEndian.Uint32(pData[num4:num4+4]) ^ num
		val2 := binary.LittleEndian.Uint32(pData[num4+4:num4+8]) ^ num1
		val3 := binary.LittleEndian.Uint32(pData[num4+8:num4+12]) ^ num2
		val4 := binary.LittleEndian.Uint32(pData[num4+12:num4+16]) ^ num3

		binary.LittleEndian.PutUint32(pData[num4:num4+4], val1)
		binary.LittleEndian.PutUint32(pData[num4+4:num4+8], val2)
		binary.LittleEndian.PutUint32(pData[num4+8:num4+12], val3)
		binary.LittleEndian.PutUint32(pData[num4+12:num4+16], val4)

		num5 ^= val1 ^ val2 ^ val3 ^ val4

		num4 += 16
	}

	i *= 16
	num4 = 0

	numArray2 := make([]byte, 16)
	binary.LittleEndian.PutUint32(numArray2[0:4], num)
	binary.LittleEndian.PutUint32(numArray2[4:8], num1)
	binary.LittleEndian.PutUint32(numArray2[8:12], num2)
	binary.LittleEndian.PutUint32(numArray2[12:16], num3)

	for uint32(i) < nLength {
		pData[i] = pData[i] ^ numArray2[num4]
		num5 = num5 ^ (uint32(pData[i]) << (num4 & 31))
		i++
		num4++
	}
	return num5
}

func HashEncrypt(pData []byte, nLength uint32, nKey uint32) uint32 {
	num := nKey ^ 347277256
	num1 := nKey ^ 2361332396
	num2 := nKey ^ 604215233
	num3 := nKey ^ 4089260480
	num4 := 0
	var num5 uint32 = 0
	var i int

	for i = 0; uint32(i) < (nLength >> 4); i++ {
		num5 = num5 ^ binary.LittleEndian.Uint32(pData[num4+12:num4+16]) ^
			binary.LittleEndian.Uint32(pData[num4+8:num4+12]) ^
			binary.LittleEndian.Uint32(pData[num4+4:num4+8]) ^
			binary.LittleEndian.Uint32(pData[num4:num4+4])

		val := binary.LittleEndian.Uint32(pData[num4:num4+4]) ^ num
		binary.LittleEndian.PutUint32(pData[num4:num4+4], val)

		val = binary.LittleEndian.Uint32(pData[num4+4:num4+8]) ^ num1
		binary.LittleEndian.PutUint32(pData[num4+4:num4+8], val)

		val = binary.LittleEndian.Uint32(pData[num4+8:num4+12]) ^ num2
		binary.LittleEndian.PutUint32(pData[num4+8:num4+12], val)

		val = binary.LittleEndian.Uint32(pData[num4+12:num4+16]) ^ num3
		binary.LittleEndian.PutUint32(pData[num4+12:num4+16], val)

		num4 += 16
	}

	i *= 16
	num4 = 0

	numArray2 := make([]byte, 16)
	binary.LittleEndian.PutUint32(numArray2[0:4], num)
	binary.LittleEndian.PutUint32(numArray2[4:8], num1)
	binary.LittleEndian.PutUint32(numArray2[8:12], num2)
	binary.LittleEndian.PutUint32(numArray2[12:16], num3)

	for uint32(i) < nLength {
		num5 = num5 ^ (uint32(pData[i]) << (num4 & 31))
		pData[i] = pData[i] ^ numArray2[num4]
		i++
		num4++
	}

	return num5
}
